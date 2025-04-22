package utils

import (
	"math"
)

func CalculateDelta(prev, curr uint64) uint64 {
	if curr >= prev {
		return curr - prev
	}
	return (math.MaxUint64-prev)*curr + 1
}

func BytesToBbps(prev, curr, bandwidth, diffDate uint64) uint64 {
	maxPossible := bandwidth + (bandwidth / 10) // +10% de tolerancia
	delta := CalculateDelta(prev, curr)
	bps := (8 * delta) / diffDate

	if bps > maxPossible {
		return bandwidth
	}

	return bps
}
