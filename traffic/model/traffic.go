package model

import "time"

type RangeDate struct {
	InitDate time.Time `query:"init_date" validate:"required"`
	EndDate  time.Time `query:"end_date" validate:"required"`
}

type Traffic struct {
	Date      time.Time `json:"date"`
	Bandwidth uint      `json:"bandwidth_mbps"`
	In        uint      `json:"in_bps"`
	Out       uint      `json:"out_bps"`
}
