package entity

import "time"

type UserStatusCounts struct {
	Hour          time.Time
	ActiveCount   uint64
	InactiveCount uint64
	UnknownCount  uint64
	TotalCount    uint64
}
