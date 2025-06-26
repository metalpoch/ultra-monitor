package entity

import "time"

type TrafficPon struct {
	PonID     int32     `db:"pon_id"`
	BpsIn     float64   `db:"bps_in"`
	BpsOut    float64   `db:"bps_out"`
	Bandwidth float64   `db:"bandwidth_mbps_sec"`
	BytesIn   float64   `db:"bytes_in_sec"`
	BytesOut  float64   `db:"bytes_out_sec"`
	Date      time.Time `db:"date"`
}

type TrafficSummary struct {
	Day          time.Time `db:"day"`
	OltIP        string    `db:"olt_ip"`
	MbpsIn       float64   `db:"mbps_in"`
	MbpsOut      float64   `db:"mbps_out"`
	MbytesInSec  float64   `db:"mbytes_in_sec"`
	MbytesOutSec float64   `db:"mbytes_out_sec"`
}

type TrafficTotalSummary struct {
	Day          time.Time `db:"day"`
	MbpsIn       float64   `db:"mbps_in"`
	MbpsOut      float64   `db:"mbps_out"`
	MbytesInSec  float64   `db:"mbytes_in_sec"`
	MbytesOutSec float64   `db:"mbytes_out_sec"`
}

type TrafficInfoSummary struct {
	Day          time.Time `db:"day"`
	Description  string    `db:"description"`
	MbpsIn       float64   `db:"mbps_in"`
	MbpsOut      float64   `db:"mbps_out"`
	MbytesInSec  float64   `db:"mbytes_in_sec"`
	MbytesOutSec float64   `db:"mbytes_out_sec"`
}

type TrafficOnt struct {
	Date             time.Time `db:"date"`
	Despt            string    `db:"despt"`
	SerialNumber     string    `db:"serial_number"`
	LineProfName     string    `db:"line_prof_name"`
	OltDistance      int32     `db:"olt_distance"`
	ControlMacCount  int8      `db:"control_mac_count"`
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
