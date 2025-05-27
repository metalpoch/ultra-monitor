package entity

import "time"

type MeasurementOlt struct {
	PonID     uint64    `db:"pon_id"`
	Bandwidth uint64    `db:"bandwidth"`
	In        uint64    `db:"bytes_in_count"`
	Out       uint64    `db:"bytes_out_count"`
	Date      time.Time `db:"date"`
}

type MeasurementOnt struct {
	PonID            uint64    `db:"pon_id"`
	Idx              uint64    `db:"idx"`
	Despt            string    `db:"despt"`
	SerialNumber     string    `db:"serial_number"`
	LineProfName     string    `db:"line_prof_name"`
	OltDistance      int64     `db:"olt_distance"`
	ControlMacCount  int64     `db:"control_mac_count"`
	ControlRunStatus int8      `db:"control_run_status"`
	BytesIn          uint64    `db:"bytes_in"`
	BytesOut         uint64    `db:"bytes_out"`
	Date             time.Time `db:"date"`
}
