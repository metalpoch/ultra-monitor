package entity

import "time"

type Device struct {
	ID         uint
	IP         string
	Sysname    string
	Community  string
	TemplateID uint
	IsAlive    bool
	LastCheck  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
