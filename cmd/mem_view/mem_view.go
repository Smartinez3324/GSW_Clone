package main

import (
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
	"time"
)

func main() {
	_, err := proc.ParseConfig("data/config/backplane.yaml")
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
		fmt.Println(data)
		time.Sleep(1 * time.Nanosecond)
	}
}
