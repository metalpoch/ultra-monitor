package model

import "time"

type OntSummaryStatusCounts struct {
	Day           time.Time `json:"day"`
	FatID         int32     `json:"fat_id"`
	OltIP         string    `json:"olt_ip"`
	PonsCount     uint64    `json:"ports_pon"`
	ActiveCount   uint64    `json:"actives"`
	InactiveCount uint64    `json:"inactives"`
	UnknownCount  uint64    `json:"unknowns"`
}

type OntSummaryStatus struct {
	Day           time.Time `json:"day"`
	PonsCount     uint64    `json:"ports_pon"`
	ActiveCount   uint64    `json:"actives"`
	InactiveCount uint64    `json:"inactives"`
	UnknownCount  uint64    `json:"unknowns"`
}

type GetStatusSummary struct {
	Day           time.Time `json:"day"`
	Description   string    `json:"description"`
	PonsCount     uint64    `json:"ports_pon"`
	ActiveCount   uint64    `json:"actives"`
	InactiveCount uint64    `json:"inactives"`
	UnknownCount  uint64    `json:"unknowns"`
}

type OntStatusForecast struct {
	Historical []OntSummaryStatus `json:"historical"`
	Forecast   []OntSummaryStatus `json:"forecast"`
}
