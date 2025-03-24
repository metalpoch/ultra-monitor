package utils

import (
	"fmt"

	"github.com/metalpoch/olt-blueprint/measurement/constants"
)

func BytesToBbps(prev, curr, bandwidth, diffDate uint64) uint64 {
	maxPossible := bandwidth + (bandwidth / 10) // +10% de tolerancia
	var delta uint64
	if prev > curr {
		fmt.Println("ACTIVO! curr =", curr, "| prev =", prev, "| bandwidth =", bandwidth, "| maxPossible =", maxPossible)
		delta = (curr + constants.OVERFLOW_COUNTER64 + 1) - prev
	} else {
		delta = curr - prev
	}

	bps := (8 * delta) / diffDate

	if bps > maxPossible {
		return bandwidth
	}

	return bps
}
