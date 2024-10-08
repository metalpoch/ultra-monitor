package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Interface struct {
	ID        uint
	IfIndex   uint
	IfName    string
	IfDescr   string
	IfAlias   string
	Device    entity.Device
	DeviceID  uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InterfaceResponse struct {
	ID        uint
	IfIndex   uint
	IfName    string
	IfDescr   string
	IfAlias   string
	Device    DeviceLite
	Template  Template
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InterfaceLite struct {
	ID        uint      `json:"id"`
	IfIndex   uint      `json:"ifindex"`
	IfName    string    `json:"ifname"`
	IfDescr   string    `json:"ifDescr"`
	IfAlias   string    `json:"ifAlias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
