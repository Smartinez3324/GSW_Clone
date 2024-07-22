package util

import (
	"fmt"
	"strings"
)

func Base2String(b []byte, groupSize int) string {
	return formatBytes(b, 2, groupSize, true)
}

func Base16String(b []byte, groupSize int) string {
	return formatBytes(b, 16, groupSize, true)
}

func Base2StringNoHeader(b []byte, groupSize int) string {
	return formatBytes(b, 2, groupSize, false)
}

func Base16StringNoHeader(b []byte, groupSize int) string {
	return formatBytes(b, 16, groupSize, false)
}

func formatBytes(b []byte, base int, groupSize int, includePrefix bool) string {
	var result strings.Builder
	var prefix string
	var formatString string

	switch base {
	case 2:
		prefix = "0b"
		formatString = "%08b"
	case 16:
		prefix = "0x"
		formatString = "%02X"
	default:
		panic("Unsupported base")
	}

	for i, byteValue := range b {
		if i > 0 && i%groupSize == 0 {
			result.WriteString(" ")
		}
		if includePrefix && i%groupSize == 0 {
			result.WriteString(prefix)
		}
		result.WriteString(fmt.Sprintf(formatString, byteValue))
	}

	return result.String()
}
