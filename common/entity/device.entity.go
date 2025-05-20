package entity

import "time"

type Device struct {
	ID          uint64 `gorm:"primaryKey"`
	IP          string `gorm:"uniqueIndex:idx_device"`
	Community   string `gorm:"uniqueIndex:idx_device"`
	SysName     string
	SysLocation string
	IsAlive     bool
	Template    Template `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TemplateID  uint
	LastCheck   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DeviceWithOID struct {
	ID          uint
	IP          string
	SysName     string
	SysLocation string
	Community   string
	IsAlive     bool
	Template    Template
	TemplateID  uint
	OidBw       string
	OidIn       string
	OidOut      string
	LastCheck   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
