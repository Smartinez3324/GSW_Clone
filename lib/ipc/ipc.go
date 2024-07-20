package ipc

type IpcHandler interface {
	Setup() error
	Cleanup() error
	Write(id int, data []byte) error
	Read(id int) ([]byte, error)
}
