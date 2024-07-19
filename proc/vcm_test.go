package proc

import (
	"testing"
)

const TEST_DATA_DIR = "../data/test/"

func compareMeasurements(expected Measurement, actual Measurement, test *testing.T) {
	if expected != actual {
		test.Errorf("Expected:\n\tName: %s\n\tSize: %d\n\tType: %s\n\tUnsigned: %t\n\tEndianness: %s\nGot:\n\tName: %s\n\tSize: %d\n\tType: %s\n\tUnsigned: %t\n\tEndianness: %s", expected.Name, expected.Size, expected.Type, expected.Unsigned, expected.Endianness, actual.Name, actual.Size, actual.Type, actual.Unsigned, actual.Endianness)
	}
}

func compareTelemetryPackets(expected TelemetryPacket, actual TelemetryPacket, test *testing.T) {
	if expected.Name != actual.Name {
		test.Errorf("Expected %s, got %s for telemetry packet name", expected.Name, actual.Name)
	}

	if expected.Port != actual.Port {
		test.Errorf("Expected %d, got %d for telemetry packet port", expected.Port, actual.Port)
	}

	if len(expected.Measurements) != len(actual.Measurements) {
		test.Errorf("Expected %d measurements, got %d", len(expected.Measurements), len(actual.Measurements))
	}

	for i := range expected.Measurements {
		if expected.Measurements[i] != actual.Measurements[i] {
			test.Errorf("Expected %s, got %s for measurement name %d", expected.Measurements[i], actual.Measurements[i], i)
		}
	}
}

func TestParseConfigBadFile(test *testing.T) {
	_, err := ParseConfig("non-existing-file123")
	if err == nil {
		test.Errorf("Expected error, got nil")
	}
}

func TestParseConfig(test *testing.T) {
	config, err := ParseConfig(TEST_DATA_DIR + "good.yaml")
	if err != nil {
		test.Errorf("Expected nil, got %v", err)
	}

	if config.Name != "vcm_test" {
		test.Errorf("Expected vcm_test, got %s", config.Name)
	}

	if len(config.Measurements) != 5 {
		test.Errorf("Expected 5 measurements, got %d", len(config.Measurements))
	}

	if len(config.TelemetryPackets) != 2 {
		test.Errorf("Expected 2 telemetry packets, got %d", len(config.TelemetryPackets))
	}

	// Check each measurement
	compareMeasurements(Measurement{Name: "Default", Size: 4, Type: "int", Unsigned: false, Endianness: "big"}, config.Measurements[0], test)
	compareMeasurements(Measurement{Name: "BigEndian", Size: 4, Type: "int", Unsigned: false, Endianness: "big"}, config.Measurements[1], test)
	compareMeasurements(Measurement{Name: "LittleEndian", Size: 4, Type: "int", Unsigned: false, Endianness: "little"}, config.Measurements[2], test)
	compareMeasurements(Measurement{Name: "Unsigned", Size: 4, Type: "int", Unsigned: true, Endianness: "big"}, config.Measurements[3], test)
	compareMeasurements(Measurement{Name: "SixteenBit", Size: 2, Type: "int", Unsigned: false, Endianness: "big"}, config.Measurements[4], test)

	// Check each telemetry packet
	compareTelemetryPackets(TelemetryPacket{Name: "Default", Port: 0, Measurements: []string{"Default", "Unsigned", "SixteenBit"}}, config.TelemetryPackets[0], test)
	compareTelemetryPackets(TelemetryPacket{Name: "Endian", Port: 1, Measurements: []string{"BigEndian", "LittleEndian"}}, config.TelemetryPackets[1], test)
}

func TestFindMeasurementByName(test *testing.T) {

}

func TestFindMeasurementByNameMissing(test *testing.T) {
}

func TestMeasurementToStringNoType(test *testing.T) {

}

func TestMeasurementToStringHasType(test *testing.T) {
}

func TestMeasurementToStringSigned(test *testing.T) {

}
