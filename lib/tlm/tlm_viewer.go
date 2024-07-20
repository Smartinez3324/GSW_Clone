package tlm

import (
	"fmt"
	"github.com/AarC10/GSW-V2/lib/ipc"
	"github.com/AarC10/GSW-V2/proc"
)

func byteSwap(data []byte, startIndex int, stopIndex int) {
	for i, j := startIndex, stopIndex; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func ReadTelemetryPacket(packet proc.TelemetryPacket, outChannel chan []byte) {
	ipcReader, err := ipc.CreateIpcShmHandler(packet, false)
	if err != nil {
		fmt.Println("Error creating IPC handler: %v\n", err)
		return
	}
	defer ipcReader.Cleanup()

	lastUpdate := ipcReader.LastUpdate()
	for {
		latestUpdate := ipcReader.LastUpdate()
		if lastUpdate != latestUpdate {
			data, err := ipcReader.Read()
			if err != nil {
				fmt.Println("Error reading from shared memory: %v\n", err)
				continue
			}
			lastUpdate = latestUpdate
			outChannel <- data
		}
	}
}

func EndiannessConverter(packet proc.TelemetryPacket, inChannel chan []byte, outChannel chan []byte) {
	byteIndicesToSwap := make([][]int, 0)

	startIndice := 0
	packetSize := 0
	for _, measurementName := range packet.Measurements {
		measurement, err := proc.FindMeasurementByName(proc.GswConfig.Measurements, measurementName)
		if err != nil {
			fmt.Printf("\t\tMeasurement '%s' not found: %v\n", measurementName, err)
			continue
		}

		if measurement.Endianness == "little" {
			byteIndicesToSwap = append(byteIndicesToSwap, []int{startIndice, startIndice + measurement.Size - 1})
		}

		startIndice += measurement.Size
		packetSize += measurement.Size
	}

	for {
		rcvData := <-inChannel
		data := make([]byte, packetSize)
		copy(data, rcvData)

		for _, byteIndices := range byteIndicesToSwap {
			byteSwap(data, byteIndices[0], byteIndices[1])
		}

		outChannel <- data
	}
}
