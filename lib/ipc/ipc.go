package ipc

type IpcWriter interface {
	Write(data []byte) error
	Cleanup()
}

type IpcReader interface {
	Read() ([]byte, error)
	Cleanup()
}
