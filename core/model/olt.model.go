package model

import "time"

type TrafficOlt struct {
	Date      time.Time
	MbpsIn    float64 `json:"mbps_in"`
	MbpsOut   float64 `json:"mbps_out"`
	Bandwidth float64 `json:"bandwidth_mbps_sec"`
	MBpsIn    float64 `json:"mbytes_in_sec"`
	MBpsOut   float64 `json:"mbytes_out_sec"`
}
