package ipc

import (
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
	"os"
	"path/filepath"
	"syscall"
)

type IpcShmHandler struct {
	file   *os.File
	data   []byte
	size   int
	packet proc.TelemetryPacket
	mode   int // 0 for reader, 1 for writer
}

const (
	modeReader = iota
	modeWriter
)

func CreateIpcShmHandler(packet proc.TelemetryPacket, isWriter bool) *IpcShmHandler {
	handler := &IpcShmHandler{
		packet: packet,
		size:   proc.GetPacketSize(packet),
		mode:   modeReader,
	}

	filename := filepath.Join(os.TempDir(), fmt.Sprintf("gsw-service-%d.shm", packet.Port))

	if isWriter {
		handler.mode = modeWriter
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}

		if err := file.Truncate(int64(handler.size)); err != nil {
			panic(err)
		}

		handler.file = file

		data, err := syscall.Mmap(int(file.Fd()), 0, handler.size, syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			panic(err)
		}

		handler.data = data
	} else {
		file, err := os.OpenFile(filename, os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}

		handler.file = file

		data, err := syscall.Mmap(int(file.Fd()), 0, handler.size, syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			panic(err)
		}

		handler.data = data
	}

	return handler
}

func (handler *IpcShmHandler) Cleanup() {
	if handler.data != nil {
		if err := syscall.Munmap(handler.data); err != nil {
			panic(err)
		}
		handler.data = nil
	}
	if handler.file != nil {
		if err := handler.file.Close(); err != nil {
			panic(err)
		}
		handler.file = nil
	}
}

func (handler *IpcShmHandler) Write(data []byte) error {
	if handler.mode != modeWriter {
		return fmt.Errorf("Handler is in reader mode")
	}
	if len(data) > handler.size {
		return fmt.Errorf("Data size exceeds shared memory size")
	}

	copy(handler.data, data)
	return nil
}

func (handler *IpcShmHandler) Read() ([]byte, error) {
	if handler.mode != modeReader {
		return nil, fmt.Errorf("Handler is in writer mode")
	}
	data := make([]byte, handler.size)
	copy(data, handler.data)
	return data, nil
}
