package utils

import "github.com/metalpoch/olt-blueprint/measurement/constants"

func BytesToBbps(prev, curr, diffDate uint) uint {
	var bps uint
	if prev >= curr {
		bps = (8 * ((curr + constants.OVERFLOW_COUNTER64 + 1) - prev)) / diffDate
	} else {
		bps = (8 * (curr - prev)) / diffDate
	}

	return bps
}
