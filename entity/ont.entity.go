package entity

import "time"

type OntSummaryStatusCounts struct {
	Day           time.Time `db:"day"`
	FatID         int32     `db:"fat_id"`
	OltIP         string    `db:"olt_ip"`
	PonsCount     uint64    `db:"ports_pon"`
	ActiveCount   uint64    `db:"actives"`
	InactiveCount uint64    `db:"inactives"`
	UnknownCount  uint64    `db:"unknowns"`
}

type OntSummaryStatus struct {
	Day           time.Time `db:"day"`
	PonsCount     uint64    `db:"ports_pon"`
	ActiveCount   uint64    `db:"actives"`
	InactiveCount uint64    `db:"inactives"`
	UnknownCount  uint64    `db:"unknowns"`
}

type GetStatusSummary struct {
	Day           time.Time `db:"day"`
	Description   string    `db:"description"`
	PonsCount     uint64    `db:"ports_pon"`
	ActiveCount   uint64    `db:"actives"`
	InactiveCount uint64    `db:"inactives"`
	UnknownCount  uint64    `db:"unknowns"`
}
