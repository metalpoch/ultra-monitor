package model

import "time"

type UserStatusCounts struct {
	Hour          time.Time `json:"date"`
	ActiveCount   uint64    `json:"active"`
	InactiveCount uint64    `json:"inactive"`
	UnknownCount  uint64    `json:"unknown"`
	TotalCount    uint64    `json:"total"`
}

type UserStatusByState struct {
	States   []string  `query:"states" validate:"required"`
	InitDate time.Time `query:"init_date" validate:"required"`
	EndDate  time.Time `query:"end_date" validate:"required"`
}

type TrafficResponse struct {
	Date      time.Time
	In        float64
	Out       float64
	Bandwidth float64
	BytesIn   float64
	BytesOut  float64
}
