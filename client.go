package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/AarC10/GSW-V2/lib/ipc"
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

	ipcReader := ipc.CreateIpcShmHandler(proc.GswConfig.TelemetryPackets[0], false)
	defer ipcReader.Cleanup()

	for {
		data, err := ipcReader.Read()
		if err != nil {
			panic(err)
		}

		// Do something with data
		fmt.Println(BytesToInt32(data))
	}
}
