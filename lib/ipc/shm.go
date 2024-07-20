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

// CreateIpcShmHandler initializes a shared memory handler as a reader or writer.
// Returns an instance of IpcShmHandler and an error if any occurs.
func CreateIpcShmHandler(packet proc.TelemetryPacket, isWriter bool) (*IpcShmHandler, error) {
	handler := &IpcShmHandler{
		packet: packet,
		size:   proc.GetPacketSize(packet),
		mode:   modeReader,
	}

	// Use /dev/shm for shared memory
	filename := filepath.Join("/dev/shm", fmt.Sprintf("gsw-service-%d.shm", packet.Port))

	if isWriter {
		handler.mode = modeWriter
		file, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %v", err)
		}

		if err := file.Truncate(int64(handler.size)); err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to truncate file: %v", err)
		}

		handler.file = file

		data, err := syscall.Mmap(int(file.Fd()), 0, handler.size, syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to memory map file: %v", err)
		}

		handler.data = data
	} else {
		file, err := os.OpenFile(filename, os.O_RDWR, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}

		handler.file = file

		data, err := syscall.Mmap(int(file.Fd()), 0, handler.size, syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to memory map file: %v", err)
		}

		handler.data = data
	}

	return handler, nil
}

// Cleanup releases resources used by the handler.
func (handler *IpcShmHandler) Cleanup() {
	if handler.data != nil {
		if err := syscall.Munmap(handler.data); err != nil {
			fmt.Printf("Failed to unmap memory: %v\n", err)
		}
		handler.data = nil
	}
	if handler.file != nil {
		if err := handler.file.Close(); err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
		handler.file = nil
	}
}

// Write writes data to the shared memory.
func (handler *IpcShmHandler) Write(data []byte) error {
	if handler.mode != modeWriter {
		return fmt.Errorf("handler is in reader mode")
	}
	if len(data) > handler.size {
		return fmt.Errorf("data size exceeds shared memory size")
	}

	copy(handler.data, data)
	return nil
}

// Read reads data from the shared memory.
func (handler *IpcShmHandler) Read() ([]byte, error) {
	if handler.mode != modeReader {
		return nil, fmt.Errorf("handler is in writer mode")
	}
	data := make([]byte, handler.size)
	copy(data, handler.data)
	return data, nil
}
