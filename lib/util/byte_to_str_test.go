package util

import (
	"testing"
)

func TestBase2String(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		groupSize int
		expected  string
	}{
		{"empty byte slice", []byte{}, 4, ""},
		{"single byte", []byte{0x01}, 8, "0b00000001"},
		{"multiple bytes", []byte{0x01, 0x02}, 8, "0b00000001 0b00000010"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base2String(tt.input, tt.groupSize)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBase16String(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		groupSize int
		expected  string
	}{
		{"empty byte slice", []byte{}, 2, ""},
		{"single byte", []byte{0x01}, 2, "0x01"},
		{"multiple bytes", []byte{0x01, 0x02}, 2, "0x01 0x02"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base16String(tt.input, tt.groupSize)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBase2StringNoHeader(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		groupSize int
		expected  string
	}{
		{"empty byte slice", []byte{}, 4, ""},
		{"single byte", []byte{0x01}, 8, "00000001"},
		{"multiple bytes", []byte{0x01, 0x02}, 8, "00000001 00000010"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base2StringNoHeader(tt.input, tt.groupSize)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBase16StringNoHeader(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		groupSize int
		expected  string
	}{
		{"empty byte slice", []byte{}, 2, ""},
		{"single byte", []byte{0x01}, 2, "01"},
		{"multiple bytes", []byte{0x01, 0x02, 0x03, 0x04}, 2, "01 02 03 04"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Base16StringNoHeader(tt.input, tt.groupSize)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestUnsupportedBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic")
		}
	}()

	formatBytes([]byte{0x01}, 0, 2, true)
}
