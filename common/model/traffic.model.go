package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Traffic struct {
	Interface   entity.Interface
	InterfaceID uint64
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

type TrafficRangeDate struct {
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

type TrafficOnt struct {
	InterfaceID      uint64    `json:"interface_id"`
	Idx              uint64    `json:"idx"`
	Despt            *string   `json:"desp"`
	SerialNumber     *string   `json:"serial_number"`
	LineProfName     *string   `json:"line_prof_name"`
	ControlMacCount  *int64    `json:"mac_count"`
	OltDistance      *int64    `json:"olt_distance"`
	ControlRunStatus *int8     `json:"control_run_status"`
	InBps            *uint64   `json:"in_bps"`
	OutBps           *uint64   `json:"out_bps"`
	BytesIn          *uint64   `json:"bytes_in"`
	BytesOut         *uint64   `json:"bytes_out"`
	Date             time.Time `json:"date"`
}
