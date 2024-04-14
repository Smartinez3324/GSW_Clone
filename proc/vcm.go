package proc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type TelemetryPacket struct {
	Fields map[string]TelemetryPacketField
}

type TelemetryPacketField struct {
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

func ParseConfiguration(filename string) []TelemetryPacket {
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

	var packets []TelemetryPacket
	return packets
}
