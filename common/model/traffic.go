package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Traffic struct {
	Interface   entity.Interface
	InterfaceID uint
	Date        time.Time
	Bandwidth   uint
	In          uint
	Out         uint
}

type TrafficResponse struct {
	Date      time.Time `json:"date"`
	Bandwidth uint      `json:"bandwidth_mbps"`
	In        uint      `json:"in_bps"`
	Out       uint      `json:"out_bps"`
}

type TranficRangeDate struct {
	InitDate time.Time `query:"init_date" validate:"required"`
	EndDate  time.Time `query:"end_date" validate:"required"`
}
