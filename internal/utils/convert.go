package utils

import "strconv"

func ParseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}
