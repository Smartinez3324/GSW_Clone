package main

import (
	"fmt"
	"github.com/AarC10/GSW-V2/lib/tlm"
	"github.com/AarC10/GSW-V2/proc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func calculateTimestamps(startLine int, packet proc.TelemetryPacket, rcvChan chan []byte) {

	for {
		// Get current timestamp in milliseconds
		data := <-rcvChan
		timestamp := time.Now().UnixNano() / int64(time.Millisecond)

		udpTimestampMeas, err := proc.FindMeasurementByName(proc.GswConfig.Measurements, "UdpSendTimestamp")
		if err != nil {
			fmt.Printf("\t\tMeasurement 'UdpSendTimestamp' not found: %v\n", err)
			continue
		}

		shmTimestampMeas, err := proc.FindMeasurementByName(proc.GswConfig.Measurements, "ShmSendTimestamp")
		if err != nil {
			fmt.Printf("\t\tMeasurement 'ShmSendTimestamp' not found: %v\n", err)
			continue
		}

		// Interpret each timestam as a uint64
		udpTimestamp := tlm.InterpretMeasurementValue(*udpTimestampMeas, data)

	}
}

func main() {
	_, err := proc.ParseConfig("data/test/benchmark.yaml")
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return
	}

	// Clear screen
	fmt.Print("\033[2J")

	// Hide the cursor
	fmt.Print("\033[?25l")

	for i, packet := range proc.GswConfig.TelemetryPackets {
		outChan := make(chan []byte)
		go proc.TelemetryPacketReader(packet, outChan)
		go calculateTimestamps(i, packet, outChan)
	}

	// Set up channel to catch interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Print("\033[2J")
	fmt.Print("\033[?25h")
}
