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
		formatString = fmt.Sprintf("%%0%db", groupSize)
	case 16:
		prefix = "0x"
		formatString = fmt.Sprintf("%%0%dX", groupSize)
	default:
		panic("Unsupported base")
	}

	bytesForNumWritten := 0
	for _, byteValue := range b {
		if bytesForNumWritten >= groupSize {
			result.WriteString(" ")
			bytesForNumWritten = 0

		}

		if bytesForNumWritten == 0 && includePrefix {
			result.WriteString(prefix)
		}

		bytesForNumWritten, _ = result.WriteString(fmt.Sprintf(formatString, byteValue))
	}

	return result.String()
}
