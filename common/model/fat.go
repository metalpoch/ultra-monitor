package model

import (
	"time"
)

type Fat struct {
	ID           uint
	LocationID   uint
	Splitter     uint8
	Latitude     float64
	Longitude    float64
	Address      string
	Fat          string
	ODN          string
	FatInterface uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type NewFat struct {
	ODN         string  `json:"odn" validate:"required"`
	Fat         string  `json:"fat" validate:"required"`
	Splitter    uint8   `json:"splitter" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	InterfaceID uint    `json:"Interface_id" validate:"required"`
}

type FatResponse struct {
	ID        uint `json:"id"`
	ODN       string
	Fat       string        `json:"fat"`
	Splitter  uint8         `json:"splitter"`
	Address   string        `json:"address"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Interface InterfaceLite `json:"interface"`
	Device    DeviceLite    `json:"device"`
	Location  Location      `json:"location"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
