package model

import "time"

type TrafficPon struct {
	PonID     int32     `json:"pon_id"`
	BpsIn     float64   `json:"bps_in"`
	BpsOut    float64   `json:"bps_out"`
	Bandwidth float64   `json:"bandwidth_mbps_sec"`
	BytesIn   float64   `json:"bytes_in_sec"`
	BytesOut  float64   `json:"bytes_out_sec"`
	Date      time.Time `json:"date"`
}

type TrafficOnt struct {
	Date             time.Time `json:"date"`
	Despt            string    `json:"despt"`
	SerialNumber     string    `json:"serial_number"`
	LineProfName     string    `json:"line_prof_name"`
	OltDistance      int32     `json:"olt_distance"`
	ControlMacCount  int8      `json:"control_mac_count"`
	ControlRunStatus int8      `json:"control_run_status"`
	MbpsIn           float64   `json:"mbps_in"`
	MbpsOut          float64   `json:"mbps_out"`
	MBpsIn           float64   `json:"mbytes_in_sec"`
	MBpsOut          float64   `json:"mbytes_out_sec"`
}

type Traffic struct {
	Date      time.Time `json:"date"`
	MbpsIn    float64   `json:"mbps_in"`
	MbpsOut   float64   `json:"mbps_out"`
	Bandwidth float64   `json:"bandwidth_mbps_sec"`
	MBpsIn    float64   `json:"mbytes_in_sec"`
	MBpsOut   float64   `json:"mbytes_out_sec"`
}

type TrafficTrendForecast struct {
	Historical []TrafficSummary `json:"historical"`
	Forecast   []TrafficSummary `json:"forecast"`
}

type TrafficSummary struct {
	Day          time.Time `json:"day"`
	FatID        int32     `json:"fat_id"`
	OltIP        string    `json:"olt_ip"`
	MbpsIn       float64   `json:"mbps_in"`
	MbpsOut      float64   `json:"mbps_out"`
	MbytesInSec  float64   `json:"mbytes_in_sec"`
	MbytesOutSec float64   `json:"mbytes_out_sec"`
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
