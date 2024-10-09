package model

import "time"

type Template struct {
	ID        uint
	Name      string
	OidBw     string
	OidIn     string
	OidOut    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
