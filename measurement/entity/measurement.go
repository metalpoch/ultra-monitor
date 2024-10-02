package entity

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Measurement struct {
	Interface   entity.Interface `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	InterfaceID uint             `gorm:"uniqueIndex"`
	Date        time.Time
	Bandwidth   uint
	In          uint
	Out         uint
}
