package entity

import (
	"time"
)

type Fat struct {
	ID         uint `gorm:"primaryKey"`
	LocationID uint
	Splitter   uint8
	Latitude   float64 `gorm:"uniqueIndex:idx_unique_fat_location"`
	Longitude  float64 `gorm:"uniqueIndex:idx_unique_fat_location"`
	Address    string  `gorm:"uniqueIndex:idx_unique_fat_location"`
	Fat        string
	ODN        string
	Location   Location `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
