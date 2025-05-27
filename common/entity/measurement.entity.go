package entity

import (
	"time"
)

// Measurement OLT
type Measurement struct {
	Interface   Interface `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	InterfaceID uint64    `gorm:"uniqueIndex"`
	Date        time.Time
	Bandwidth   uint64
	In          uint64
	Out         uint64
}

type MeasurementOnt struct {
	Interface        Interface `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	InterfaceID      uint64
	Idx              uint64
	Despt            string
	SerialNumber     string
	LineProfName     string
	OltDistance      int64
	ControlMacCount  int64
	ControlRunStatus int8
	BytesIn          uint64
	BytesOut         uint64
	Date             time.Time
}
