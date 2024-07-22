package main

import (
	"encoding/binary"
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
)

func interpretMeasurementValue(measurement proc.Measurement, data []byte, offset int) (int, error) {
	if measurement.Endianness == "little" {

	}

	if measurement.Size == 1 {
		return int(data[offset]), nil
	} else if measurement.Size == 2 {
		return int(data[offset]) + (int(data[offset+1]) << 8), nil
	} else if measurement.Size == 4 {
		return int(data[offset]) + (int(data[offset+1]) << 8) + (int(data[offset+2]) << 16) + (int(data[offset+3]) << 24), nil
	} else {
		return 0, fmt.Errorf("Unsupported size for measurement: %d", measurement.Size)
	}
}

func main() {
	i := []byte{0x6, 0x4}
	fmt.Println(binary.LittleEndian.Uint32(i))
	fmt.Println(binary.BigEndian.Uint32(i))
	//_, err := proc.ParseConfig("data/config/backplane.yaml")
	//if err != nil {
	//	fmt.Printf("Error parsing YAML: %v\n", err)
	//	return
	//}
	//
	//outChan := make(chan []byte)
	//for _, packet := range proc.GswConfig.TelemetryPackets {
	//	go proc.TelemetryPacketReader(packet, outChan)
	//}
	//
	//for {
	//	data := <-outChan
	//	fmt.Print("\033[H\033[2J")
	//
	//	var sb strings.Builder
	//	var offset int
	//
	//	for _, packet := range proc.GswConfig.TelemetryPackets {
	//		// Print the measurement name, base-10 value and base-16 value. One for each line
	//		// Format: MeasurementName: Value (Base-10) [(Base-16)]
	//		for _, measurementName := range packet.Measurements {
	//			measurement, err := proc.FindMeasurementByName(proc.GswConfig.Measurements, measurementName)
	//			if err != nil {
	//				fmt.Printf("\t\tMeasurement '%s' not found: %v\n", measurementName, err)
	//				continue
	//			}
	//
	//			// Get the value of the measurement
	//			value, err := interpretMeasurementValue(*measurement, data, offset)
	//			if err != nil {
	//				fmt.Printf("\t\tError interpreting measurement value: %v\n", err)
	//				continue
	//			}
	//
	//			offset += measurement.Size
	//
	//			// Print the measurement name, base-10 value and base-16 value
	//			sb.WriteString(fmt.Sprintf("%s: %d (0x%X)\n", measurement.Name, value, value))
	//		}
	//	}
	//	fmt.Print(sb.String())
	//	time.Sleep(1 * time.Nanosecond)
	//}
}
