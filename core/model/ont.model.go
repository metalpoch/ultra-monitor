package model

import "time"

type OntStatusCounts struct {
	Hour          time.Time `json:"date"`
	State         string    `json:"state"`
	PonsCount     uint64    `json:"pons_count"`
	ActiveCount   uint64    `json:"active"`
	InactiveCount uint64    `json:"inactive"`
	UnknownCount  uint64    `json:"unknown"`
	TotalCount    uint64    `json:"total"`
}

type OntStatusCountsByState struct {
	Hour          time.Time `json:"date"`
	Sysname       string    `json:"sysname"`
	PonsCount     uint64    `json:"pons_count"`
	ActiveCount   uint64    `json:"active"`
	InactiveCount uint64    `json:"inactive"`
	UnknownCount  uint64    `json:"unknown"`
	TotalCount    uint64    `json:"total"`
}

type OntTraffic struct {
	Hour             time.Time `json:"total"`
	Despt            string    `json:"despt"`
	SerialNumber     string    `json:"serial_number"`
	LineProfName     string    `json:"line_prof_name"`
	OltDistance      int64     `json:"olt_distance"`
	ControlMacCount  int64     `json:"control_mac_count"`
	ControlRunStatus int8      `json:"control_run_status"`
	MbpsIn           float64   `json:"mbps_in"`
	MbpsOut          float64   `json:"mbps_out"`
	MBpsIn           float64   `json:"mbytes_in_sec"`
	MBpsOut          float64   `json:"mbytes_out_sec"`
}
