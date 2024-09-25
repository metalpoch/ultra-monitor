package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
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
