package tlm

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// MeasurementInfo represents the metadata of a measurement.
type MeasurementInfo struct {
	Name      string
	Size      int
	Locations []LocationInfo
}

// LocationInfo represents where a measurement is located within a packet.
type LocationInfo struct {
	PacketIndex uint32
	Offset      int
}

// TelemetryShm represents shared memory for telemetry.
type TelemetryShm struct {
	// Add necessary fields and methods for managing shared memory
}

// VCM represents the Vehicle Configuration Manager.
type VCM struct {
	NumPackets int
	Packets    []PacketInfo
}

func (vcm *VCM) GetInfo(name string) (*MeasurementInfo, error) {
	// Retrieve measurement info by name
	return nil, nil
}

// PacketInfo represents a telemetry packet.
type PacketInfo struct {
	Size int
}

// TelemetryViewer is responsible for reading telemetry measurements.
type TelemetryViewer struct {
	updateMode    UpdateMode
	packetIDs     []uint32
	packetSizes   map[uint32]int
	packetBuffers map[uint32][]byte
	shm           *TelemetryShm
	vcm           *VCM
	checkAll      bool
	mu            sync.Mutex
}

type UpdateMode int

const (
	StandardUpdate UpdateMode = iota
	BlockingUpdate
	NonBlockingUpdate
)

// NewTelemetryViewer creates a new TelemetryViewer instance.
func NewTelemetryViewer() *TelemetryViewer {
	return &TelemetryViewer{
		updateMode:    StandardUpdate,
		packetSizes:   make(map[uint32]int),
		packetBuffers: make(map[uint32][]byte),
	}
}

// Init initializes the TelemetryViewer with VCM and shared memory.
func (tv *TelemetryViewer) Init(vcm *VCM, shm *TelemetryShm) error {
	tv.vcm = vcm
	tv.shm = shm

	tv.packetIDs = make([]uint32, tv.vcm.NumPackets)
	tv.packetSizes = make(map[uint32]int)
	tv.packetBuffers = make(map[uint32][]byte)

	return nil
}

// RemoveAll removes all tracked packets.
func (tv *TelemetryViewer) RemoveAll() {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	tv.checkAll = false
	tv.packetIDs = nil
	tv.packetSizes = make(map[uint32]int)
	tv.packetBuffers = make(map[uint32][]byte)
}

// AddAll adds all packets to the tracker.
func (tv *TelemetryViewer) AddAll() error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	tv.checkAll = true

	for id := 0; id < tv.vcm.NumPackets; id++ {
		if err := tv.addPacket(uint32(id)); err != nil {
			return err
		}
	}

	return nil
}

// AddMeasurement adds a measurement to the tracker.
func (tv *TelemetryViewer) AddMeasurement(info *MeasurementInfo) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	for _, loc := range info.Locations {
		if err := tv.addPacket(loc.PacketIndex); err != nil {
			return err
		}
	}

	return nil
}

// AddPacket adds a packet to the tracker by ID.
func (tv *TelemetryViewer) AddPacket(packetID uint32) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	return tv.addPacket(packetID)
}

func (tv *TelemetryViewer) addPacket(packetID uint32) error {
	if packetID < uint32(tv.vcm.NumPackets) {
		if _, exists := tv.packetSizes[packetID]; !exists {
			tv.packetSizes[packetID] = tv.vcm.Packets[packetID].Size
			tv.packetBuffers[packetID] = make([]byte, tv.packetSizes[packetID])
			tv.packetIDs = append(tv.packetIDs, packetID)
		}
		return nil
	}
	return fmt.Errorf("invalid packet id: %d", packetID)
}

// SetUpdateMode sets the update mode.
func (tv *TelemetryViewer) SetUpdateMode(mode UpdateMode) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	tv.updateMode = mode
	// Set the corresponding read mode in shared memory (not implemented here)
}

// Update updates the telemetry data.
func (tv *TelemetryViewer) Update(timeout time.Duration) error {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	// Lock shared memory and read the packets
	// Copy packets from shared memory to local buffers
	// Unlock shared memory

	return nil
}

// LatestData returns the latest data for a measurement.
func (tv *TelemetryViewer) LatestData(meas *MeasurementInfo) ([]byte, error) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	var bestLoc *LocationInfo
	best := int64(time.Now().UnixNano())

	for _, loc := range meas.Locations {
		// Implement logic to find the best location based on update value
		// bestLoc = &loc // Example assignment
	}

	if bestLoc == nil {
		return nil, errors.New("measurement does not exist anywhere")
	}

	return tv.packetBuffers[bestLoc.PacketIndex][bestLoc.Offset:], nil
}
