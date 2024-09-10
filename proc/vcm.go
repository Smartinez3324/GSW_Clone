package proc

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Name             string                 `yaml:"name"`
	Measurements     map[string]Measurement `yaml:"measurements"`
	TelemetryPackets []TelemetryPacket      `yaml:"telemetry_packets"`
}

type Measurement struct {
	Name       string `yaml:"name"`
	Size       int    `yaml:"size"`
	Type       string `yaml:"type,omitempty"`
	Unsigned   bool   `yaml:"unsigned,omitempty"`
	Endianness string `yaml:"endianness,omitempty"`
}

type TelemetryPacket struct {
	Name         string   `yaml:"name"`
	Port         int      `yaml:"port"`
	Measurements []string `yaml:"measurements"`
}

// TODO: Make global safer
var GswConfig Configuration

func ParseConfig(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	// Unmarshalling doesn't seem to lead to errors with bad data. Better to check result config
	_ = yaml.Unmarshal(data, &GswConfig)
	if GswConfig.Name == "" {
		return nil, fmt.Errorf("Error parsing YAML. No config name.")
	}

	if len(GswConfig.Measurements) == 0 {
		return nil, fmt.Errorf("Error parsing YAML. No measurements.")
	}

	if len(GswConfig.TelemetryPackets) == 0 {
		return nil, fmt.Errorf("Error parsing YAML. No telemetry packets.")
	}

	// Set default values for measurements if not specified
	for k, _ := range GswConfig.Measurements {
		// TODO: More strict checks of configuration and input handling
		if GswConfig.Measurements[k].Name == "" {
			return nil, fmt.Errorf("Measurement name missing")
		}

		if GswConfig.Measurements[k].Endianness == "" {
			entry := GswConfig.Measurements[k] // Workaround to avoid UnaddressableFieldAssign
			entry.Endianness = "big"           // Default to big endian
			GswConfig.Measurements[k] = entry
		} else if GswConfig.Measurements[k].Endianness != "little" && GswConfig.Measurements[k].Endianness != "big" {
			return nil, fmt.Errorf("Endianess not specified as big or little got %s", GswConfig.Measurements[k].Endianness)
		}
	}

	return &GswConfig, nil
}

func GetPacketSize(packet TelemetryPacket) int {
	size := 0
	for _, measurementName := range packet.Measurements {
		measurement, ok := GswConfig.Measurements[measurementName]
		if !ok {
			fmt.Printf("\t\tMeasurement '%s' not found\n", measurementName)
			continue
		}
		size += measurement.Size
	}
	return size
}

func (m Measurement) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s, Size: %d", m.Name, m.Size))
	if m.Type != "" {
		sb.WriteString(fmt.Sprintf(", Type: %s", m.Type))
	}

	if m.Unsigned {
		sb.WriteString(", Unsigned")
	} else {
		sb.WriteString(", Signed")
	}
	sb.WriteString(fmt.Sprintf(", Endianness: %s", m.Endianness))
	return sb.String()
}
