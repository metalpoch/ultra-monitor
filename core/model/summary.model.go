package model

import "time"

type UserStatusQuery struct {
	InitDate time.Time `query:"init_date" validate:"required"`
	EndDate  time.Time `query:"end_date" validate:"required"`
}

type UserStatusCounts struct {
	Hour          time.Time `json:"date"`
	ActiveCount   uint64    `json:"active"`
	InactiveCount uint64    `json:"inactive"`
	UnknownCount  uint64    `json:"unknown"`
	TotalCount    uint64    `json:"total"`
}
