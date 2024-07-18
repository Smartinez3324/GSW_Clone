package proc

import "fmt"

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

func PacketListener(packet TelemetryPacket) {
	packetSize := getPacketSize(packet)
	fmt.Printf("Packet size for port %d: %d\n", packet.Port, packetSize)

	for {
	}
}
