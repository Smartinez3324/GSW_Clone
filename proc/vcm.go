package proc

import (
	"bufio"
	"fmt"
	"os"
)

// TelemetryPacketInfo Information about a telemetry packet
type TelemetryPacketInfo struct {
	Name   string               // Name of the telemetry packet
	Port   uint16               // Port number of the telemetry packet
	Fields []TelemetryFieldInfo // Information about the fields in the telemetry packet
}

// TelemetryFieldInfo Information about a telemetry field in a telemetry packet
type TelemetryFieldInfo struct {
	Type   interface{}
	Endian string
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
