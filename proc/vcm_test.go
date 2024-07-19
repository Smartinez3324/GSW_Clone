package proc

import "testing"

func TestParseConfigBadFile(test *testing.T) {
	_, err := ParseConfig("non-existing-file123")
	if err == nil {
		test.Errorf("Expected error, got nil")
	}
}

func TestParseConfigBadYaml(test *testing.T) {
	_, err := ParseConfig("data/test/bad.yaml")
	if err == nil {
		test.Errorf("Expected error, got nil")
	}
}

func TestParseConfig(test *testing.T) {
	config, err := ParseConfig("data/test/good.yaml")
	if err != nil {
		test.Errorf("Expected nil, got %v", err)
	}

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
