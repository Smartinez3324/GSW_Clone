package proc

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

// TelemetryPacketInfo Information about a telemetry packet
type TelemetryPacketInfo struct {
	name   string               // Name of the telemetry packet
	port   uint16               // Port number of the telemetry packet
	fields []TelemetryFieldInfo // Information about the fields in the telemetry packet
}

// TelemetryFieldInfo Information about a telemetry field in a telemetry packet
type TelemetryFieldInfo struct {
	name       string           // Name of the telemetry data
	size       uint16           // Size of the telemetry data in bytes
	padding    uint16           // Number of padding bytes
	endianness binary.ByteOrder // Endianness of the telemetry data
	signed     bool             // Whether the telemetry data is signed
	dataType   interface{}      // Data type of the telemetry data
}

func Parser(filename string) []TelemetryPacketInfo {
	file, _ := os.Open(filename)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return []TelemetryPacketInfo{}
}
