package ipc

import "github.com/AarC10/GSW-V2/proc"

type IpcServiceSide interface {
	Setup(packet proc.TelemetryPacket) error
	Cleanup()
	Write(data []byte) error
}

type IpcClientSide interface {
	Setup(packet proc.TelemetryPacket) error
	Cleanup()
	Read() ([]byte, error)
}
