package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Device struct {
	ID          uint64
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
	ID          uint64
	IP          string
	SysName     string
	SysLocation string
	Community   string
	OidBw       string
	OidIn       string
	OidOut      string
}

type DeviceLite struct {
	ID          uint      `json:"id"`
	IP          string    `json:"ip"`
	SysName     string    `json:"sysname"`
	SysLocation string    `json:"syslocation"`
	IsAlive     bool      `json:"is_alive"`
	LastCheck   time.Time `json:"last_check"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeviceResponse struct {
	ID          uint      `json:"id"`
	IP          string    `json:"ip"`
	Community   string    `json:"community"`
	SysName     string    `json:"sysname"`
	SysLocation string    `json:"syslocation"`
	IsAlive     bool      `json:"is_alive"`
	Template    Template  `json:"template"`
	LastCheck   time.Time `json:"last_check"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
