package api

import (
	"fmt"
	"strings"
)

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

// PrintBytes
func PrintBytes(size uint64) string {
	unit := ""
	value := float32(size)

	switch {
	case size >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case size >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case size >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case size >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case size >= BYTE:
		unit = "B"
	case size == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.2f", value)
	stringValue = strings.TrimSuffix(stringValue, ".00")
	return fmt.Sprintf("%s%s", stringValue, unit)
}
