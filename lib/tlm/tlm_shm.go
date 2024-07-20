package tlm

import (
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
	"sync"
)

type TelemetryData struct {
	data map[int][]byte
	mu   sync.RWMutex
}

var sharedTelemetryData = TelemetryData{
	data: make(map[int][]byte),
}

func TlmShmInit() {
	sharedTelemetryData.mu.Lock()
	defer sharedTelemetryData.mu.Unlock()

	for _, packet := range proc.GswConfig.TelemetryPackets {
		totalSize := 0
		for _, measurementName := range packet.Measurements {
			for _, measurement := range proc.GswConfig.Measurements {
				if measurement.Name == measurementName {
					totalSize += measurement.Size
				}
			}
		}
		sharedTelemetryData.data[packet.Port] = make([]byte, totalSize)
	}
}

func TlmShmWrite(port int, data []byte) error {
	sharedTelemetryData.mu.Lock()
	defer sharedTelemetryData.mu.Unlock()

	packetData, exists := sharedTelemetryData.data[port]
	if !exists {
		return fmt.Errorf("Port %d not found", port)
	}
	if len(data) != len(packetData) {
		return fmt.Errorf("Data size mismatch for port %d", port)
	}

	copy(packetData, data)
	return nil
}

func TlmShmRead(port int) ([]byte, error) {
	sharedTelemetryData.mu.RLock()
	defer sharedTelemetryData.mu.RUnlock()

	packetData, exists := sharedTelemetryData.data[port]
	if !exists {
		return nil, fmt.Errorf("port %d not found", port)
	}

	return append([]byte(nil), packetData...), nil
}
