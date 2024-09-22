package entity

import "time"

type Interface struct {
	ID        uint // IfIndex
	IfName    string
	IfDescr   string
	IfAlias   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
