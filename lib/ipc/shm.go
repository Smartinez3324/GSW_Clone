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
	common       ShmCommon
	shmServer    *shmipc.Session
	bufferWriter shmipc.BufferWriter
	writeBuff    []byte
}

type ShmClientSide struct {
	common       ShmCommon
	shmClient    *shmipc.Session
	packetSize   int
	bufferReader shmipc.BufferReader
}

func setupCommonUds(common *ShmCommon, telemetryPacket proc.TelemetryPacket) error {
	// Setup Unix domain socket
	udsPath := filepath.Join(os.TempDir(), "gsw-service-", telemetryPacket.Name, "-", strconv.Itoa(telemetryPacket.Port), ".uds_sock")
	_ = syscall.Unlink(udsPath)
	unixListener, err := net.ListenUnix("unix", &net.UnixAddr{Name: udsPath, Net: "unix"})
	if err != nil {
		return fmt.Errorf("Creating Unix domain socket failed: %v", err)
	}
	common.unixListener = unixListener

	// Accept UDS
	common.conn, err = unixListener.Accept()
	if err != nil {
		return fmt.Errorf("Accepting Unix domain socket failed: %v", err)
	}

	return nil
}

func (shmHandler *ShmServiceSide) Setup(telemetryPacket proc.TelemetryPacket) error {
	err := setupCommonUds(&shmHandler.common, telemetryPacket)
	if err != nil {
		return err
	}

	// Create server session
	conf := shmipc.DefaultConfig()
	shmHandler.shmServer, err = shmipc.Server(shmHandler.common.conn, conf)
	if err != nil {
		return fmt.Errorf("IPC server creation failed: %v", err)
	}

	// Accept stream
	shmHandler.common.stream, err = shmHandler.shmServer.AcceptStream()
	if err != nil {
		return fmt.Errorf("Accept stream failed: %v", err)
	}

	// Create buffer writer
	shmHandler.bufferWriter = shmHandler.common.stream.BufferWriter()
	shmHandler.writeBuff, err = shmHandler.bufferWriter.Reserve(proc.GetPacketSize(telemetryPacket))
	if err != nil {
		return fmt.Errorf("Reserve buffer failed: %v", err)
	}

	return nil
}

func (shmHandler *ShmClientSide) Setup(packet proc.TelemetryPacket) error {
	err := setupCommonUds(&shmHandler.common, packet)
	if err != nil {
		return err
	}

	return nil
}

func (shmHandler *ShmServiceSide) Cleanup() {
	defer func(unixListener *net.UnixListener) {
		err := unixListener.Close()
		if err != nil {
			fmt.Printf("Closing Unix domain socket failed: %v\n", err)
		}
	}(shmHandler.common.unixListener)

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Closing connection failed: %v", err)
		}
	}(shmHandler.common.conn)

	defer func(shmServer *shmipc.Session) {
		err := shmServer.Close()
		if err != nil {
			fmt.Printf("Closing IPC server failed: %v\n", err)
		}
	}(shmHandler.shmServer)

	defer func(stream *shmipc.Stream) {
		err := stream.Close()
		if err != nil {
			fmt.Printf("Closing stream failed: %v\n", err)
		}
	}(shmHandler.common.stream)
}

func (shmHandler *ShmServiceSide) Write(data []byte) error {
	copy(shmHandler.writeBuff, data)
	return shmHandler.common.stream.Flush(false)
}

func (shmHandler *ShmClientSide) Read() ([]byte, error) {
	data := make([]byte, shmHandler.packetSize)
	n, err := shmHandler.common.stream.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Read failed: %v", err)
	}

	return data[:n], nil
}
