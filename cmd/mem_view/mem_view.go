package main

import (
	"fmt"
	"github.com/AarC10/GSW-V2/lib/tlm"
	"github.com/AarC10/GSW-V2/lib/util"
	"github.com/AarC10/GSW-V2/proc"
	"strings"
)

func printTelemetryPacket(startLine int, packet proc.TelemetryPacket, rcvChan chan []byte) {
	for {
		var sb strings.Builder
		offset := 0

		data := <-rcvChan

		// Print the measurement name, base-10 value and base-16 value. One for each line
		// Format: MeasurementName: Value (Base-10) [(Base-16)]
		sb.WriteString(fmt.Sprintf("\033[%d;0H", startLine))
		for _, measurementName := range packet.Measurements {
			measurement, err := proc.FindMeasurementByName(proc.GswConfig.Measurements, measurementName)
			if err != nil {
				fmt.Printf("\t\tMeasurement '%s' not found: %v\n", measurementName, err)
				continue
			}

			value := tlm.InterpretMeasurementValue(*measurement, data[offset:offset+measurement.Size])
			if err != nil {
				fmt.Printf("\t\tError interpreting measurement value: %v\n", err)
				continue
			}

			sb.WriteString(fmt.Sprintf("%s: %v [%s]\n", measurementName, value, util.Base16String(data[offset:offset+measurement.Size], 1)))
			offset += measurement.Size
		}

		fmt.Print(sb.String())
	}
}

func main() {
	//_, err := proc.ParseConfig("data/config/backplane.yaml")
	_, err := proc.ParseConfig("data/test/good.yaml")
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return
	}

	// Clear screen
	fmt.Print("\033[2J")

	outChan := make(chan []byte)
	startLine := 0
	for _, packet := range proc.GswConfig.TelemetryPackets {
		go proc.TelemetryPacketReader(packet, outChan)
		go printTelemetryPacket(startLine, packet, outChan)
		startLine += len(packet.Measurements) + 1
	}

	select {}
}
