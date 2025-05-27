package entity

import "time"

type UserStatusCounts struct {
	Hour          time.Time
	State         string
	PonsCount     uint64
	ActiveCount   uint64
	InactiveCount uint64
	UnknownCount  uint64
	TotalCount    uint64
}

type OntStatusCountsByState struct {
	Hour          time.Time
	Sysname       string
	PonsCount     uint64
	ActiveCount   uint64
	InactiveCount uint64
	UnknownCount  uint64
	TotalCount    uint64
}

type OntTraffic struct {
	Hour             time.Time
	Despt            string
	SerialNumber     string
	LineProfName     string
	OltDistance      int64
	ControlMacCount  int64
	ControlRunStatus int8
	MbpsIn           float64 `gorm:"column:mbps_in"`
	MbpsOut          float64 `gorm:"column:mbps_out"`
	MBpsIn           float64 `gorm:"column:mbytes_in_sec"`
	MBpsOut          float64 `gorm:"column:mbytes_out_sec"`
}
