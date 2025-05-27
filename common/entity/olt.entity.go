package entity

import "time"

type Olt struct {
	ID          uint64
	IP          string
	Community   string
	SysName     string
	SysLocation string
	IsAlive     bool
	TemplateID  uint
	LastCheck   time.Time
	CreatedAt   time.Time
}
type TrafficOlt struct {
	Date      time.Time
	MbpsIn    float64 `gorm:"column:mbps_in"`
	MbpsOut   float64 `gorm:"column:mbps_out"`
	Bandwidth float64 `gorm:"column:bandwidth_mbps_sec"`
	MBpsIn    float64 `gorm:"column:mbytes_in_sec"`
	MBpsOut   float64 `gorm:"column:mbytes_out_sec"`
}
