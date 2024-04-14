package proc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// TelemetryPacketInfo Information about a telemetry packet
type TelemetryPacketInfo struct {
	Port   uint16               // Port number of the telemetry packet
	Fields []TelemetryFieldInfo // Information about the fields in the telemetry packet
}

// TelemetryFieldInfo Information about a telemetry field in a telemetry packet
type TelemetryFieldInfo struct {
	Name   string
	Type   interface{}
	Endian string
}

type telemetryConfiguration struct {
	Fields map[string]telemetryConfigurationField `json:"fields"`
}

type telemetryConfigurationField struct {
	Type   string `json:"type"`
	Endian string `json:"endian"`
}

func ParseConfiguration(filename string) []TelemetryPacketInfo {
	file, _ := os.Open(filename)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(file)

	jsonStr, _ := io.ReadAll(file)

	var config map[string]telemetryConfiguration

	err := json.Unmarshal(jsonStr, &config)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	fmt.Println(config)

	return nil
}
