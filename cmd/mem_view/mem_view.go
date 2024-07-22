package main

import (
	"fmt"
	"github.com/AarC10/GSW-V2/lib/tlm"
	"github.com/AarC10/GSW-V2/lib/util"
	"github.com/AarC10/GSW-V2/proc"
	"strings"
	"time"
)

func main() {
	//_, err := proc.ParseConfig("data/config/backplane.yaml")
	_, err := proc.ParseConfig("data/test/good.yaml")
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return
	}

	outChan := make(chan []byte)
	for _, packet := range proc.GswConfig.TelemetryPackets {
		go proc.TelemetryPacketReader(packet, outChan)
	}

	for {
		data := <-outChan
		fmt.Print("\033[H\033[2J")

		var sb strings.Builder

		for _, packet := range proc.GswConfig.TelemetryPackets {
			// Print the measurement name, base-10 value and base-16 value. One for each line
			// Format: MeasurementName: Value (Base-10) [(Base-16)]

			offset := 0
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
		}
		fmt.Print(sb.String())
		time.Sleep(1 * time.Nanosecond)
	}
}
