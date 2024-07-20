package main

import (
	"bytes"
	"encoding/binary"
	"github.com/AarC10/GSW-V2/lib/tlm"
	"github.com/AarC10/GSW-V2/proc"
)

func BytesToInt32(data []byte) int32 {
	var i int32
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &i)
	if err != nil {
		panic(err)
	}
	return i
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

	outChan := make(chan []byte)
	go tlm.ReadTelemetryPacket(proc.GswConfig.TelemetryPackets[0], outChan)

	for {
		data := <-outChan
		i := BytesToInt32(data)
		println(i)
	}
}
