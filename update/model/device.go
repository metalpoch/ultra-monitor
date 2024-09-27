package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
)

type Device struct {
	ID          uint
	IP          string
	Community   string
	SysName     string
	SysLocation string
	IsAlive     bool
	Template    entity.Template
	TemplateID  uint
	LastCheck   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AddDevice struct {
	IP        string
	Community string
	Template  uint
}

type DeviceWithOID struct {
	ID          uint
	IP          string
	SysName     string
	SysLocation string
	Community   string
	IsAlive     bool
	Template    entity.Template
	TemplateID  uint
	OidBw       string
	OidIn       string
	OidOut      string
	LastCheck   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
