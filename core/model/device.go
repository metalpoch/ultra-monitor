package model

import "time"

type Device struct {
	ID          uint
	IP          string
	Community   string
	SysName     string
	SysLocation string
	IsAlive     bool
	LastCheck   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
