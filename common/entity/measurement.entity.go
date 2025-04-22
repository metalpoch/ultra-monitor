package entity

import (
	"time"
)

type Measurement struct {
	Interface   Interface `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	InterfaceID uint      `gorm:"uniqueIndex"`
	Date        time.Time
	Bandwidth   uint64
	In          uint64
	Out         uint64
}
