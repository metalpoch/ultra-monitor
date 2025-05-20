package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Interface struct {
	ID        uint64
	IfIndex   uint64
	IfName    string
	IfDescr   string
	IfAlias   string
	Device    entity.Device
	DeviceID  uint64
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
	ID        uint64    `json:"id"`
	IfIndex   uint64    `json:"ifindex"`
	IfName    string    `json:"ifname"`
	IfDescr   string    `json:"ifDescr"`
	IfAlias   string    `json:"ifAlias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GponPort struct {
	Shell *uint8 `query:"shell" validate:"required,lte=0"`
	Card  *uint8 `query:"card" validate:"required_with=Port"`
	Port  *uint8 `query:"port"`
}
