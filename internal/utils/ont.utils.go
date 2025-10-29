package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseOntIDX parses the ont_idx string to extract pon_idx and ont_idx
// Expected format: "pon_idx.ont_idx"
func ParseOntIDX(ontIDX string) (int64, uint8, error) {
	// Try dot-separated format: "pon_idx.ont_idx"
	if strings.Contains(ontIDX, ".") {
		parts := strings.Split(ontIDX, ".")
		if len(parts) == 2 {
			ponIdx, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return 0, 0, fmt.Errorf("invalid pon_idx: %v", err)
			}
			ontIdx, err := strconv.ParseUint(parts[1], 10, 8)
			if err != nil {
				return 0, 0, fmt.Errorf("invalid ont_idx: %v", err)
			}
			return ponIdx, uint8(ontIdx), nil
		}
	}

	// Try as a single number (assuming it's ont_idx and pon_idx is 0)
	ontIdx, err := strconv.ParseUint(ontIDX, 10, 8)
	if err == nil {
		return 0, uint8(ontIdx), nil
	}

	return 0, 0, fmt.Errorf("unable to parse ont_idx '%s', expected format: pon_idx.ont_idx", ontIDX)
}

