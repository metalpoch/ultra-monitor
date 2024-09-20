package entity

import "time"

type Template struct {
	ID        uint
	Name      string
	OidBw     string
	OidIn     string
	OidOut    string
	PrefixBw  string
	PrefixIn  string
	PrefixOut string
	CreatedAt time.Time
	UpdatedAt time.Time
}
