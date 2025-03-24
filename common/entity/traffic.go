package entity

import "time"

type Traffic struct {
	Interface   Interface `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	InterfaceID uint
	Date        time.Time
	Bandwidth   uint64
	In          uint64
	Out         uint64
}
