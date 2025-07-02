package model

import "time"

type TrafficPon struct {
	PonID     int32     `json:"pon_id"`
	BpsIn     float64   `json:"bps_in"`
	BpsOut    float64   `json:"bps_out"`
	Bandwidth float64   `json:"bandwidth_mbps_sec"`
	BytesIn   float64   `json:"bytes_in"`
	BytesOut  float64   `json:"bytes_out"`
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
	MbytesIn         float64   `json:"mbytes_in"`
	MbytesOut        float64   `json:"mbytes_out"`
}

type Traffic struct {
	Date      time.Time `json:"date"`
	MbpsIn    float64   `json:"mbps_in"`
	MbpsOut   float64   `json:"mbps_out"`
	Bandwidth float64   `json:"bandwidth_mbps_sec"`
	MbytesIn  float64   `json:"mbytes_in"`
	MbytesOut float64   `json:"mbytes_out"`
}

type TrafficTrendForecast struct {
	Historical []TrafficSummary `json:"historical"`
	Forecast   []TrafficSummary `json:"forecast"`
}

type TrafficSummary struct {
	Day       time.Time `json:"day"`
	OltIP     string    `json:"olt_ip"`
	MbpsIn    float64   `json:"mbps_in"`
	MbpsOut   float64   `json:"mbps_out"`
	MbytesIn  float64   `json:"mbytes_in"`
	MbytesOut float64   `json:"mbytes_out"`
}

type TrafficTotalSummary struct {
	Day       time.Time `json:"day"`
	MbpsIn    float64   `json:"mbps_in"`
	MbpsOut   float64   `json:"mbps_out"`
	MbytesIn  float64   `json:"mbytes_in"`
	MbytesOut float64   `json:"mbytes_out"`
}

type TrafficInfoSummary struct {
	Day         time.Time `json:"day"`
	Description string    `json:"description"`
	MbpsIn      float64   `json:"mbps_in"`
	MbpsOut     float64   `json:"mbps_out"`
	MbytesIn    float64   `json:"mbytes_in"`
	MbytesOut   float64   `json:"mbytes_out"`
}
