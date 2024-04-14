package proc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type TelemetryPacket struct {
	Fields map[string]TelemetryPacketField `json:"fields"`
}

type TelemetryPacketField struct {
	Type   string `json:"type"`
	Endian string `json:"endian"`
}

func ParseConfiguration(filename string) map[string]TelemetryPacket {
	// Put file data into memory
	file, _ := os.Open(filename)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(file)

	jsonStr, _ := io.ReadAll(file)

	// Convert JSON into map of telemetryConfigurations
	var config map[string]TelemetryPacket

	err := json.Unmarshal(jsonStr, &config)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return config
}
