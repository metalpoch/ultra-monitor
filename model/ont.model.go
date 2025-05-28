package model

import "time"

type OntStatusCounts struct {
	Date          time.Time `json:"date"`
	State         string    `json:"state"`
	PonsCount     uint64    `json:"pons_count"`
	ActiveCount   uint64    `json:"actives"`
	InactiveCount uint64    `json:"inactives"`
	UnknownCount  uint64    `json:"unknowns"`
	TotalCount    uint64    `json:"total"`
}

type OntStatusCountsByState struct {
	Date          time.Time `json:"date"`
	Sysname       string    `json:"sysname"`
	PonsCount     uint64    `json:"pons_count"`
	ActiveCount   uint64    `json:"actives"`
	InactiveCount uint64    `json:"inactives"`
	UnknownCount  uint64    `json:"unknowns"`
	TotalCount    uint64    `json:"total"`
}
