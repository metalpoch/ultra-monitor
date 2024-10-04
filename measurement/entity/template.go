package entity

import "time"

type Template struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	OidBw     string
	OidIn     string
	OidOut    string
	PrefixBw  string
	PrefixIn  string
	PrefixOut string
	CreatedAt time.Time
	UpdatedAt time.Time
}
