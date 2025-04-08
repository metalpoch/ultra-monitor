package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Traffic struct {
	Interface   entity.Interface
	InterfaceID uint
	Date        time.Time
	Bandwidth   uint64
	In          uint64
	Out         uint64
}

type TrafficResponse struct {
	Date      time.Time `json:"date"`
	Bandwidth uint64    `json:"bandwidth_bps"`
	In        uint64    `json:"in_bps"`
	Out       uint64    `json:"out_bps"`
}

type TranficRangeDate struct {
	InitDate time.Time `query:"init_date" validate:"required"`
	EndDate  time.Time `query:"end_date" validate:"required"`
}

type TrafficState struct {
	State     string `json:"state"`
	In        uint64 `json:"in_bps"`
	Out       uint64 `json:"out_bps"`
	Bandwidth uint64 `json:"bandwidth_bps"`
}

type TrafficStateN struct {
	State     string  `json:"state"`
	In        float64 `json:"in_bps"`
	Out       float64 `json:"out_bps"`
	Bandwidth float64 `json:"bandwidth_bps"`
}
