package ipc

import (
	"github.com/AarC10/GSW-V2/proc"
)

type IpcShmCommon struct {
}

type IpcShmWriter struct {
	common IpcShmCommon
}

type IpcShmReader struct {
	common IpcShmCommon
}

func (shmHandler *IpcShmWriter) Setup(telemetryPacket proc.TelemetryPacket) error {

	return nil
}

func (shmHandler *IpcShmWriter) Cleanup() {

}

func (shmHandler *IpcShmReader) Setup(packet proc.TelemetryPacket) error {
	return nil
}

func (shmHandler *IpcShmReader) Cleanup() {

}

func (shmHandler *IpcShmWriter) Write(data []byte) error {

	return nil
}

func (shmHandler *IpcShmReader) Read() ([]byte, error) {
	return nil, nil
}
