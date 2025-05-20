package entity

import "time"

type Interface struct {
	ID        uint64 `gorm:"primaryKey"`
	IfIndex   uint64 `gorm:"uniqueIndex:idx_interface"`
	IfName    string
	IfDescr   string
	IfAlias   string
	Device    Device `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DeviceID  uint64 `gorm:"uniqueIndex:idx_interface"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
