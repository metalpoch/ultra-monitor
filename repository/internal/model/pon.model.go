package model

import "time"

type TrafficTrend struct {
	Day          time.Time `db:"day"`
	PonID        int32     `db:"pon_id"`
	MbpsIn       float64   `db:"mbps_in"`
	MbpsOut      float64   `db:"mbps_out"`
	MbytesInSec  float64   `db:"mbytes_in_sec"`
	MbytesOutSec float64   `db:"mbytes_out_sec"`
}

type TrafficMeta struct {
	PonID int32  `db:"pon_id"`
	FatID int32  `db:"fat_id"`
	OltIP string `db:"olt_ip"`
}
