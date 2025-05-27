package model

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
