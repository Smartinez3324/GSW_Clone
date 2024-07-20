package ipc

import (
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
	"github.com/cloudwego/shmipc-go"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

type ShmCommon struct {
	unixListener *net.UnixListener
	conn         net.Conn
	stream       *shmipc.Stream
}

type ShmServiceSide struct {
	common    ShmCommon
	shmServer *shmipc.Session
}

type ShmClientSide struct {
	common    ShmCommon
	shmClient *shmipc.Session
}

func (shmHandler *ShmServiceSide) Setup(telemetryPacket proc.TelemetryPacket) error {
	// Setup Unix domain socket
	udsPath := filepath.Join(os.TempDir(), "gsw-service-", telemetryPacket.Name, "-", strconv.Itoa(telemetryPacket.Port), ".uds_sock")
	_ = syscall.Unlink(udsPath)
	unixListener, err := net.ListenUnix("unix", &net.UnixAddr{Name: udsPath, Net: "unix"})
	if err != nil {
		return fmt.Errorf("Creating Unix domain socket failed: %v", err)
	}

	// Accept UDS
	conn, err := unixListener.Accept()
	if err != nil {
		return fmt.Errorf("Accepting Unix domain socket failed: %v", err)
	}
	shmHandler.common.unixListener = unixListener

	// Create server session
	conf := shmipc.DefaultConfig()
	shmServer, err := shmipc.Server(conn, conf)
	if err != nil {
		return fmt.Errorf("IPC server creation failed: %v", err)
	}
	shmHandler.common.conn = conn
	shmHandler.shmServer = shmServer

	// Accept stream
	stream, err := shmServer.AcceptStream()
	if err != nil {
		return fmt.Errorf("Accept stream failed: %v", err)
	}
	shmHandler.common.stream = stream

	return nil
}

func (shmHandler *ShmServiceSide) Cleanup() error {
	shmHandler.common.unixListener.Close()
	shmHandler.common.conn.Close()
	shmHandler.shmServer.Close()
	shmHandler.common.stream.Close()
	return nil
}

func (shmHandler *ShmServiceSide) Write(data []byte) error {

	return nil
}

func (shmHandler *ShmClientSide) Read() ([]byte, error) {
	return nil, nil
}
