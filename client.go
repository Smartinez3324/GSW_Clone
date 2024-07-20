package main

import (
	"fmt"
	"github.com/AarC10/GSW-V2/lib/tlm"
	"github.com/AarC10/GSW-V2/proc"
)

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

	tlmPacketService, err := tlm.TlmClientInit(proc.GswConfig.TelemetryPackets[0])
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Starting telemetry packet service")
	for {
		buff, err := tlmPacketService.Read()
		if err != nil {
			fmt.Println("Read error:", err)
		}

		fmt.Println("Received data:", buff)
	}

}
