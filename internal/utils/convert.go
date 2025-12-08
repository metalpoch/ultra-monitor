package utils

import (
	"fmt"
	"strconv"
)

func ParseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func VolumeSuffix(value int64) string {
	if value > 1_000_000_000_000_000 {
		return fmt.Sprintf("%.2f PB", float64(value)/float64(1_000_000_000_000_000))
	} else if value > 1_000_000_000_000 {
		return fmt.Sprintf("%.2f TB", float64(value)/float64(1_000_000_000_000))
	} else if value > 1_000_000_000 {
		return fmt.Sprintf("%.2f GB", float64(value)/float64(1_000_000_000))
	} else if value > 1_000_000 {
		return fmt.Sprintf("%.2f MB", float64(value)/float64(1_000_000))
	} else if value > 1_000 {
		return fmt.Sprintf("%.2f KB", float64(value)/float64(1_000))
	} else {
		return fmt.Sprintf("%d Bytes", value)
	}
}
