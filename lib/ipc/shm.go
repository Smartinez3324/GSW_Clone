package ipc

import (
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
	"os"
	"path/filepath"
	"syscall"
)

type IpcShmCommon struct {
	file   *os.File
	data   []byte
	size   int
	packet proc.TelemetryPacket
}

type IpcShmWriter struct {
	common IpcShmCommon
}

type IpcShmReader struct {
	common IpcShmCommon
}

func CreateIpcShmWriter(packet proc.TelemetryPacket) *IpcShmWriter {
	shmHandler := &IpcShmWriter{
		common: IpcShmCommon{
			packet: packet,
			size:   proc.GetPacketSize(packet),
		},
	}

	filename := filepath.Join(os.TempDir(), fmt.Sprintf("gsw-service-%d.shm", packet.Port))
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	if err := file.Truncate(int64(shmHandler.common.size)); err != nil {
		panic(err)
	}

	shmHandler.common.file = file

	data, err := syscall.Mmap(int(file.Fd()), 0, shmHandler.common.size, syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	shmHandler.common.data = data

	return shmHandler
}

func CreateIpcShmReader(packet proc.TelemetryPacket) *IpcShmReader {
	shmHandler := &IpcShmReader{
		common: IpcShmCommon{
			packet: packet,
			size:   proc.GetPacketSize(packet),
		},
	}

	filename := filepath.Join(os.TempDir(), fmt.Sprintf("gsw-service-%d.shm", packet.Port))
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	shmHandler.common.file = file

	data, err := syscall.Mmap(int(file.Fd()), 0, shmHandler.common.size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	shmHandler.common.data = data

	return shmHandler
}

func (shmHandler *IpcShmWriter) Cleanup() {
	if shmHandler.common.data != nil {
		if err := syscall.Munmap(shmHandler.common.data); err != nil {
			panic(err)
		}
		shmHandler.common.data = nil
	}
	if shmHandler.common.file != nil {
		if err := shmHandler.common.file.Close(); err != nil {
			panic(err)
		}
		shmHandler.common.file = nil
	}
}

func (shmHandler *IpcShmReader) Cleanup() {
	if shmHandler.common.data != nil {
		if err := syscall.Munmap(shmHandler.common.data); err != nil {
			panic(err)
		}
		shmHandler.common.data = nil
	}
	if shmHandler.common.file != nil {
		if err := shmHandler.common.file.Close(); err != nil {
			panic(err)
		}
		shmHandler.common.file = nil
	}
}

func (shmHandler *IpcShmWriter) Write(data []byte) error {
	if len(data) > shmHandler.common.size {
		return fmt.Errorf("data size exceeds shared memory size")
	}

	copy(shmHandler.common.data, data)
	return nil
}

func (shmHandler *IpcShmReader) Read() ([]byte, error) {
	data := make([]byte, shmHandler.common.size)
	copy(data, shmHandler.common.data)
	return data, nil
}
