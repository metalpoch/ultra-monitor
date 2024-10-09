package entity

import "time"

type Template struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	OidBw     string `gorm:"default:1.3.6.1.2.1.2.2.1.5"`
	OidIn     string `gorm:"default:1.3.6.1.2.1.31.1.1.1.6"`
	OidOut    string `gorm:"default:1.3.6.1.2.1.31.1.1.1.10"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
