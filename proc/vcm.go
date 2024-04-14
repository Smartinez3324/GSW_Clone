package proc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type TelemetryPacket struct {
	Fields map[string]TelemetryPacketField
	Port   uint16
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
	var config map[string]telemetryConfiguration

	err := json.Unmarshal(jsonStr, &config)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Convert to internal data structure
	var packets []TelemetryPacket
	convertConfiguration(config)

	return packets
}

func convertConfiguration(config map[string]telemetryConfiguration) []TelemetryPacket {
	// Allocate a slice of TelemetryPackets
	packets := make([]TelemetryPacket, 0)

	fmt.Println("Processing", len(config), "ports")
	for key, value := range config {
		fmt.Println("Processing port", key)
		port, _ := strconv.Atoi(key)
		packets = append(packets, TelemetryPacket{
			Port: uint16(port),
		})

		fmt.Println(value)
	}

	fmt.Println(packets)

	return packets
}
