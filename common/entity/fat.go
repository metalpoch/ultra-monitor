package entity

import (
	"time"
)

type Fat struct {
	ID          uint `gorm:"primaryKey"`
	LocationID  uint
	InterfaceID uint
	Splitter    uint8
	Fat         string    `gorm:"uniqueIndex"`
	Address     string    `gorm:"uniqueIndex:idx_unique_fat_location"`
	OND         string    `gorm:"uniqueIndex:idx_unique_fat_location"`
	Latitude    float64   `gorm:"uniqueIndex:idx_unique_fat_location"`
	Longitude   float64   `gorm:"uniqueIndex:idx_unique_fat_location"`
	Location    Location  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Interface   Interface `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
