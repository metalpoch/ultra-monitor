package entity

import "time"

type TrafficPon struct { // Table
	PonID     uint64    `db:"pon_id"`
	BpsIn     float64   `db:"bps_in"`
	BpsOut    float64   `db:"bps_out"`
	Bandwidth float64   `db:"bandwidth_mbps_sec"`
	BytesIn   float64   `db:"bytes_in_sec"`
	BytesOut  float64   `db:"bytes_out_sec"`
	Date      time.Time `db:"date"`
}

type TrafficOnt struct {
	Date             time.Time `db:"date"`
	Despt            string    `db:"despt"`
	SerialNumber     string    `db:"serial_number"`
	LineProfName     string    `db:"line_prof_name"`
	OltDistance      int64     `db:"olt_distance"`
	ControlMacCount  int64     `db:"control_mac_count"`
	ControlRunStatus int8      `db:"control_run_status"`
	MbpsIn           float64   `db:"mbps_in"`
	MbpsOut          float64   `db:"mbps_out"`
	MBpsIn           float64   `db:"mbytes_in_sec"`
	MBpsOut          float64   `db:"mbytes_out_sec"`
}

type Traffic struct {
	Date      time.Time `db:"date"`
	MbpsIn    float64   `db:"mbps_in"`
	MbpsOut   float64   `db:"mbps_out"`
	Bandwidth float64   `db:"bandwidth_mbps_sec"`
	MBpsIn    float64   `db:"mbytes_in_sec"`
	MBpsOut   float64   `db:"mbytes_out_sec"`
}
