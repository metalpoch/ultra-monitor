package entity

import "time"

type Interface struct {
	ID        uint // IfIndex
	IfName    string
	IfDescr   string
	DeviceID  uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
