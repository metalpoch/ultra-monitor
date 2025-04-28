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
	BytesIn     uint64
	BytesOut    uint64
}

type TrafficResponse struct {
	Date      time.Time `json:"date"`
	Bandwidth uint64    `json:"bandwidth_bps"`
	In        uint64    `json:"in_bps"`
	Out       uint64    `json:"out_bps"`
	BytesIn   uint64    `json:"byte_in"`
	BytesOut  uint64    `json:"bytes_out"`
}

type TranficRangeDate struct {
	InitDate time.Time `query:"init_date" validate:"required"`
	EndDate  time.Time `query:"end_date" validate:"required"`
}

type TrafficState struct {
	State     string `json:"state"`
	Bandwidth uint64 `json:"bandwidth_bps"`
	In        uint64 `json:"in_bps"`
	Out       uint64 `json:"out_bps"`
	BytesIn   uint64 `json:"bytes_in"`
	BytesOut  uint64 `json:"bytes_out"`
}

type TrafficODN struct {
	ODN       string `json:"odn"`
	Bandwidth uint64 `json:"bandwidth_bps"`
	In        uint64 `json:"in_bps"`
	Out       uint64 `json:"out_bps"`
	BytesIn   uint64 `json:"bytes_in"`
	BytesOut  uint64 `json:"bytes_out"`
}
