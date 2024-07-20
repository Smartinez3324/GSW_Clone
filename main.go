package main

import (
	"fmt"
	"github.com/AarC10/GSW-V2/lib/tlm"
	"github.com/AarC10/GSW-V2/proc"
)

func printTelemetryPackets() {
	fmt.Println("Telemetry Packets:")
	for _, packet := range proc.GswConfig.TelemetryPackets {
		fmt.Printf("\tName: %s\n\tPort: %d\n", packet.Name, packet.Port)
		if len(packet.Measurements) > 0 {
			fmt.Println("\tMeasurements:")
			for _, measurementName := range packet.Measurements {
				measurement, err := proc.FindMeasurementByName(proc.GswConfig.Measurements, measurementName)
				if err != nil {
					fmt.Printf("\t\tMeasurement '%s' not found: %v\n", measurementName, err)
					continue
				}
				fmt.Printf("\t\t%s\n", measurement.String())
			}
		} else {
			fmt.Println("\t\tNo measurements defined.")
		}
	}
}

func vcmInitialize() {
	_, err := proc.ParseConfig("data/config/backplane.yaml")
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return
	}

	printTelemetryPackets()
}

func decomInitialize() map[int]chan []byte {
	channelMap := make(map[int]chan []byte)

	for _, packet := range proc.GswConfig.TelemetryPackets {
		finalOutputChannel := make(chan []byte)
		channelMap[packet.Port] = finalOutputChannel

		go proc.PacketListener(packet, finalOutputChannel)
	}

	return channelMap
}

//func main() {
//	vcmInitialize()
//	channelMap := decomInitialize()
//
//	for _, channel := range channelMap {
//		go proc.TestReceiver(channel)
//	}
//
//	select {}
//}

func main() {
	proc.GswConfig = proc.Configuration{
		Name: "example",
		Measurements: []proc.Measurement{
			{Name: "measurement1", Size: 4},
			{Name: "measurement2", Size: 4},
		},
		TelemetryPackets: []proc.TelemetryPacket{
			{Name: "packet1", Port: 10000, Measurements: []string{"measurement1", "measurement2"}},
		},
	}

	tlm.TlmShmInit()

	// Write to shared memory
	err := tlm.TlmShmWrite(10000, []byte{1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		fmt.Println("Write error:", err)
	}

	// Read from shared memory
	data, err := tlm.TlmShmRead(10000)
	if err != nil {
		fmt.Println("Read error:", err)
	} else {
		fmt.Println("Read data:", data)
	}
}
