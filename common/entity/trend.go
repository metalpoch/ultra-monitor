package entity

import "time"

type Trend struct {
	Device    Device `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DeviceID  uint
	Date      time.Time
	Bandwidth uint
	In        uint
	Out       uint
}
