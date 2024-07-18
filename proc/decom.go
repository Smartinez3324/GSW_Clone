package proc

import (
	"fmt"
	"net"
)

func getPacketSize(packet TelemetryPacket) int {
	size := 0
	for _, measurementName := range packet.Measurements {
		measurement, err := FindMeasurementByName(GswConfig.Measurements, measurementName)
		if err != nil {
			fmt.Printf("\t\tMeasurement '%s' not found: %v\n", measurementName, err)
			continue
		}
		size += measurement.Size
	}
	return size
}

func PacketListener(packet TelemetryPacket, channel chan []byte) {
	packetSize := getPacketSize(packet)
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
			channel <- buffer[:n] // Send data over channel
		} else {
			fmt.Printf("Received packet of incorrect size. Expected: %d, Received: %d\n", packetSize, n)
		}
	}
}

func ReceiverTest(channel chan []byte) {
	for {
		data := <-channel
		fmt.Printf("Received data: %v\n", data)
	}
}
