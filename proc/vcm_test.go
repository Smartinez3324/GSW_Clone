package proc

import (
	"testing"
)

const TEST_DATA_DIR = "../data/test/"

func compareMeasurements(expected Measurement, actual Measurement, test *testing.T) {
	if expected != actual {
		test.Errorf("Expected:, \tName: %s, \tSize: %d, \tType: %s, \tUnsigned: %t, \tEndianness: %s, Got:, \tName: %s, \tSize: %d, \tType: %s, \tUnsigned: %t, \tEndianness: %s", expected.Name, expected.Size, expected.Type, expected.Unsigned, expected.Endianness, actual.Name, actual.Size, actual.Type, actual.Unsigned, actual.Endianness)
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

func CompareMeasurementString(expected string, actual string, test *testing.T) {
	if expected != actual {
		test.Errorf("\nExp: %s\nGot: %s", expected, actual)
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

	compareMeasurements(Measurement{Name: "Default", Size: 4, Type: "int", Unsigned: false, Endianness: "big"}, config.Measurements[0], test)
	compareMeasurements(Measurement{Name: "BigEndian", Size: 4, Type: "int", Unsigned: false, Endianness: "big"}, config.Measurements[1], test)
	compareMeasurements(Measurement{Name: "LittleEndian", Size: 4, Type: "int", Unsigned: false, Endianness: "little"}, config.Measurements[2], test)
	compareMeasurements(Measurement{Name: "Unsigned", Size: 4, Type: "int", Unsigned: true, Endianness: "big"}, config.Measurements[3], test)
	compareMeasurements(Measurement{Name: "SixteenBit", Size: 2, Type: "int", Unsigned: false, Endianness: "big"}, config.Measurements[4], test)

	compareTelemetryPackets(TelemetryPacket{Name: "Default", Port: 0, Measurements: []string{"Default", "Unsigned", "SixteenBit"}}, config.TelemetryPackets[0], test)
	compareTelemetryPackets(TelemetryPacket{Name: "Endian", Port: 1, Measurements: []string{"BigEndian", "LittleEndian"}}, config.TelemetryPackets[1], test)
}

func TestFindMeasurementByName(test *testing.T) {
	config, _ := ParseConfig(TEST_DATA_DIR + "good.yaml")
	measurement, err := FindMeasurementByName(config.Measurements, "Default")
	if err != nil {
		test.Errorf("Expected nil, got %v", err)
	}

	compareMeasurements(Measurement{Name: "Default", Size: 4, Type: "int", Unsigned: false, Endianness: "big"}, *measurement, test)

	measurement, err = FindMeasurementByName(config.Measurements, "Missing")
	if err == nil {
		test.Errorf("Expected error, got nil")
	}
}

func TestMeasurementToString(test *testing.T) {
	bigSigned := Measurement{Name: "Test", Size: 4, Type: "int", Unsigned: false, Endianness: "big"}
	expected := "Name: Test, Size: 4, Type: int, Signed, Endianness: big"
	CompareMeasurementString(expected, bigSigned.String(), test)

	littleUnsigned := Measurement{Name: "Test", Size: 4, Type: "int", Unsigned: true, Endianness: "little"}
	expected = "Name: Test, Size: 4, Type: int, Unsigned, Endianness: little"
	CompareMeasurementString(expected, littleUnsigned.String(), test)

	noType := Measurement{Name: "Test", Size: 4, Unsigned: true, Endianness: "little"}
	expected = "Name: Test, Size: 4, Unsigned, Endianness: little"
	CompareMeasurementString(expected, noType.String(), test)
}
