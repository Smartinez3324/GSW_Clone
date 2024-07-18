package proc

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Configuration struct {
	Name             string            `yaml:"name"`
	Measurements     []Measurement     `yaml:"measurements"`
	TelemetryPackets []TelemetryPacket `yaml:"telemetry_packets"`
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
var Cfg Configuration

func ParseYAML(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	// Set default values for measurements if not specified
	for i := range Cfg.Measurements {
		if Cfg.Measurements[i].Name == "" {
			return nil, fmt.Errorf("Measurement name missing")
		}

		if Cfg.Measurements[i].Endianness == "" {
			Cfg.Measurements[i].Endianness = "big" // Default to big endian
		} else if Cfg.Measurements[i].Endianness != "little" && Cfg.Measurements[i].Endianness != "big" {
			return nil, fmt.Errorf("Endianess not specified as big or little got %s", Cfg.Measurements[i].Endianness)
		}
	}

	return &Cfg, nil
}

func FindMeasurementByName(measurements []Measurement, name string) (*Measurement, error) {
	for _, meas := range measurements {
		if meas.Name == name {
			return &meas, nil
		}
	}
	return nil, fmt.Errorf("measurement '%s' not found", name)
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
