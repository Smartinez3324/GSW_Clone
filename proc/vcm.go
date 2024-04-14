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

// TelemetryPacketInfo Information about a telemetry packet
type telemetryPacketConfig struct {
	Port   uint16
	Fields []telemetryFieldInfoConfig
}

// telemetryFieldInfoConfig Configuration structure meant for unmarshalling JSON
type telemetryFieldInfoConfig struct {
	Name   string `json:"name"`
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

	info := new(telemetryPacketConfig)

	err := json.Unmarshal(jsonStr, &info)
	if err != nil {
		fmt.Println("Error unmarshalling JSON")
		return nil
	}

	fmt.Println(info.Port)
	for _, field := range info.Fields {
		fmt.Println(field.Type)
		fmt.Println(field.Endian)
	}

	return nil
}
