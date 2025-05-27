package entity

import "time"

type OntStatusCounts struct {
	Date          time.Time `db:"date"`
	State         string    `db:"state"`
	PonsCount     uint64    `db:"pons_count"`
	ActiveCount   uint64    `db:"actives"`
	InactiveCount uint64    `db:"inactives"`
	UnknownCount  uint64    `db:"unknowns"`
	TotalCount    uint64    `db:"total"`
}

type OntStatusCountsByState struct {
	Date          time.Time `db:"date"`
	Sysname       string    `db:"sysname"`
	PonsCount     uint64    `db:"pons_count"`
	ActiveCount   uint64    `db:"actives"`
	InactiveCount uint64    `db:"inactives"`
	UnknownCount  uint64    `db:"unknowns"`
	TotalCount    uint64    `db:"total"`
}
