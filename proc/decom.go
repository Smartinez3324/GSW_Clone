package proc

import (
	"fmt"
	"github.com/AarC10/GSW-V2/lib/ipc"
	"net"
)

func PacketListener(packet TelemetryPacket) {
	shmWriter, _ := ipc.CreateIpcShmHandler(packet, true)
	if shmWriter == nil {
		fmt.Printf("Failed to create shared memory writer\n")
		return
	}
	defer shmWriter.Cleanup()

	packetSize := GetPacketSize(packet)
	fmt.Printf("Packet size for port %d: %d\n", packet.Port, packetSize)

	// Listen over UDP
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", packet.Port))
	if err != nil {
		fmt.Printf("Error resolving UDP address: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Error listening on UDP: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Listening on port %d for telemetry packet...\n", packet.Port)

	// Receive data
	buffer := make([]byte, packetSize)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		if n == packetSize {
			err := shmWriter.Write(buffer)
			if err != nil {
				fmt.Printf("Error writing to shared memory: %v\n", err)
			}
		} else {
			fmt.Printf("Received packet of incorrect size. Expected: %d, Received: %d\n", packetSize, n)
		}
	}
}
