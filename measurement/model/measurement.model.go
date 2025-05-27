package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type MeasurementOlt struct {
	Interface   entity.Interface
	InterfaceID uint64
	Date        time.Time
	Bandwidth   uint64
	In          uint64
	Out         uint64
}

type MeasurementOnt struct {
	/*Interface        entity.Interface*/
	InterfaceID      uint64    `json:"interface_id"`
	Idx              uint64    `json:"idx"`
	Despt            string    `json:"despt"`
	SerialNumber     string    `json:"serial_number"`
	LineProfName     string    `json:"line_prof_name"`
	OltDistance      int64     `json:"olt_distance"`
	ControlMacCount  int64     `json:"mac_count"`
	ControlRunStatus int8      `json:"control_run_status"`
	BytesIn          uint64    `json:"bytes_in"`
	BytesOut         uint64    `json:"bytes_out"`
	Date             time.Time `json:"date"`
}
