package util

import (
	"fmt"
	"strings"
)

func BytesToString(b []byte, base int, groupSize ...int) string {
	var result strings.Builder
	var prefix string
	var formatString string

	defaultGroupSize := 2
	if base == 2 {
		defaultGroupSize = 4
	}
	size := defaultGroupSize
	if len(groupSize) > 0 {
		size = groupSize[0]
	}

	includePrefix := true
	if len(groupSize) > 1 {
		includePrefix = groupSize[1] != 0
	}

	switch base {
	case 2:
		prefix = "0b"
		formatString = "%08b"
	case 16:
		prefix = "0x"
		formatString = "%02X"
	default:
		panic("unsupported base")
	}

	for i, byteValue := range b {
		if i > 0 && i%size == 0 {
			result.WriteString(" ")
		}
		if includePrefix && i%size == 0 {
			result.WriteString(prefix)
		}
		result.WriteString(fmt.Sprintf(formatString, byteValue))
	}

	return result.String()
}
