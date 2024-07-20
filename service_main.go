package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/AarC10/GSW-V2/lib/ipc"
	"github.com/AarC10/GSW-V2/proc"
	"time"
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

func Int32ToBytes(value int32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value) // Use binary.BigEndian if you need big-endian
	if err != nil {
		panic(err) // Handle error appropriately
	}
	return buf.Bytes()
}

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

	ipcWriter := ipc.CreateIpcShmWriter(proc.GswConfig.TelemetryPackets[0])
	defer ipcWriter.Cleanup()

	fmt.Println("Writing to shared memory...")
	// Increment a i32 value in the shared memory
	var i int32
	i = 0

	for {
		err := ipcWriter.Write(Int32ToBytes(i))
		if err != nil {
			fmt.Printf("Error writing to shared memory: %v\n", err)
		}
		fmt.Println(i)

		i++
		time.Sleep(1 * time.Second)
	}

}
