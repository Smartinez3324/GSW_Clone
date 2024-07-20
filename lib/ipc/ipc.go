package ipc

import "github.com/AarC10/GSW-V2/proc"

type IpcWriter interface {
	Setup(packet proc.TelemetryPacket) error
	Cleanup()
	Write(data []byte) error
}

type IpcReader interface {
	Setup(packet proc.TelemetryPacket) error
	Cleanup()
	Read() ([]byte, error)
}
