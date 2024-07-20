package ipc

type ShmHandler struct {
}

func (shmHandler *ShmHandler) Setup() error {

	return nil
}

func (shmHandler *ShmHandler) Cleanup() error {

	return nil
}

func (shmHandler *ShmHandler) Write(id int, data []byte) error {

	return nil
}

func (shmHandler *ShmHandler) Read(id int) ([]byte, error) {

	return nil, nil
}
