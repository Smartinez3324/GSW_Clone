package tlm

import (
	"encoding/binary"
	"fmt"
	"github.com/AarC10/GSW-V2/proc"
	"math"
)

func byteSwap(data []byte, startIndex int, stopIndex int) {
	for i, j := startIndex, stopIndex; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func PacketEndianConverter(packet proc.TelemetryPacket, inChannel chan []byte, outChannel chan []byte) {
	byteIndicesToSwap := make([][]int, 0)

	startIndice := 0
	packetSize := 0
	for _, measurementName := range packet.Measurements {
		measurement, err := proc.FindMeasurementByName(proc.GswConfig.Measurements, measurementName)
		if err != nil {
			fmt.Printf("\t\tMeasurement '%s' not found: %v\n", measurementName, err)
			continue
		}

		if measurement.Endianness == "little" {
			byteIndicesToSwap = append(byteIndicesToSwap, []int{startIndice, startIndice + measurement.Size - 1})
		}

		startIndice += measurement.Size
		packetSize += measurement.Size
	}

	for {
		rcvData := <-inChannel
		data := make([]byte, packetSize)
		copy(data, rcvData)

		for _, byteIndices := range byteIndicesToSwap {
			byteSwap(data, byteIndices[0], byteIndices[1])
		}

		outChannel <- data
	}
}

func InterpretUnsignedInteger(data []byte, endianness string) interface{} {
	switch len(data) {
	case 1:
		return uint8(data[0])
	case 2:
		if endianness == "little" {
			return binary.LittleEndian.Uint16(data)
		} else {
			return binary.BigEndian.Uint16(data)
		}
	case 4:
		if endianness == "little" {
			return binary.LittleEndian.Uint32(data)
		} else {
			return binary.BigEndian.Uint32(data)
		}
	case 8:
		if endianness == "little" {
			return binary.LittleEndian.Uint64(data)
		} else {
			return binary.BigEndian.Uint64(data)
		}
	}

	var result uint64
	if endianness == "little" {
		for i := 0; i < len(data); i++ {
			result |= uint64(data[i]) << uint(i*8)
		}
	} else {
		for i := 0; i < len(data); i++ {
			result |= uint64(data[len(data)-1-i]) << uint(i*8)
		}
	}

	return result
}

func InterpretSignedInteger(data []byte, endianness string) interface{} {
	unsigned := InterpretUnsignedInteger(data, endianness)

	switch unsigned.(type) {
	case uint8:
		return int8(unsigned.(uint8))
	case uint16:
		return int16(unsigned.(uint16))
	case uint32:
		return int32(unsigned.(uint32))
	case uint64:
		return int64(unsigned.(uint64))
	}

	return nil
}

func InterpretFloat(data []byte, endianness string) interface{} {
	unsigned := InterpretUnsignedInteger(data, endianness)

	switch unsigned.(type) {
	case uint32:
		return math.Float32frombits(unsigned.(uint32))
	case uint64:
		return math.Float64frombits(unsigned.(uint64))
	}

	return 0
}

func InterpretMeasurementValue(measurement proc.Measurement, data []byte) interface{} {
	switch measurement.Type {
	case "int":
		if measurement.Unsigned {
			return InterpretUnsignedInteger(data, measurement.Endianness)
		}

		return InterpretSignedInteger(data, measurement.Endianness)
	case "float":
		return InterpretFloat(data, measurement.Endianness)
	}

	fmt.Printf("Unsupported type for measurement: %s\n", measurement.Type)
	return nil
}
