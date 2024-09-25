package utils

import (
	"math"
)

func BytesToKbps(prevBytes, currBytes, diffDate int64) int32 {
	bps := math.Abs(float64(8*currBytes)-float64(8*prevBytes)) / float64(diffDate)
	return int32(math.Round(bps / 1000))
}
