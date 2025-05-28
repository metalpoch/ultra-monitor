package model

import (
	"time"
)

type MeasurementPon struct {
	PonID     uint64    `json:"pon_id"`
	Bandwidth uint64    `json:"bandwidth"`
	In        uint64    `json:"bytes_in_count"`
	Out       uint64    `json:"bytes_out_count"`
	Date      time.Time `json:"date"`
}

type MeasurementOnt struct {
	PonID            uint64    `json:"pon_id"`
	Idx              uint64    `json:"idx"`
	Despt            string    `json:"despt"`
	SerialNumber     string    `json:"serial_number"`
	LineProfName     string    `json:"line_prof_name"`
	OltDistance      int64     `json:"olt_distance"`
	ControlMacCount  int64     `json:"control_mac_count"`
	ControlRunStatus int8      `json:"control_run_status"`
	BytesIn          uint64    `json:"bytes_in_count"`
	BytesOut         uint64    `json:"bytes_out_count"`
	Date             time.Time `json:"date"`
}
