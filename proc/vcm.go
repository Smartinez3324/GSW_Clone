package proc

import (
	"encoding/binary"
)

// TelemetryPacketInfo Information about a telemetry packet
type TelemetryPacketInfo struct {
	name   string               // Name of the telemetry packet
	port   uint16               // Port number of the telemetry packet
	fields []TelemetryFieldInfo // Information about the fields in the telemetry packet
}

// TelemetryFieldInfo Information about a telemetry field in a telemetry packet
type TelemetryFieldInfo struct {
	name       string           // Name of the telemetry data
	size       uint16           // Size of the telemetry data in bytes
	padding    uint16           // Number of padding bytes
	endianness binary.ByteOrder // Endianess of the telemetry data
	signed     bool             // Whether the telemetry data is signed
}
