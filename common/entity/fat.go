package entity

import (
	"time"

	"gorm.io/gorm"
)

type Fat struct {
	ID          uint   `gorm:"primaryKey"`
	Fat         string `gorm:"uniqueIndex"`
	Splitter    uint8
	Address     string
	Latitude    float64
	Longitude   float64
	Interface   Interface `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	InterfaceID uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
