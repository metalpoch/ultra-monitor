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
	Date      time.Time `json:"total"`
	In        float64   `json:"mbps_in"`
	Out       float64   `json:"mbps_out"`
	Bandwidth float64   `json:"bandwidth_mbps"`
	BytesIn   float64   `json:"mbytes_in_sec"`
	BytesOut  float64   `json:"mbytes_out_sec"`
}
