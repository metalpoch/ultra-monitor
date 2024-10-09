package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Fat struct {
	ID          uint
	Fat         string
	Splitter    uint8
	Address     string
	Latitude    float64
	Longitude   float64
	Interface   entity.Interface
	InterfaceID uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewFat struct {
	Fat         string  `json:"fat" validate:"required"`
	Splitter    uint8   `json:"splitter" validate:"required"`
	Address     string  `json:"address" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	InterfaceID uint    `json:"Interface_id" validate:"required"`
}

type FatResponse struct {
	ID        uint          `json:"id"`
	Fat       string        `json:"fat"`
	Splitter  uint8         `json:"splitter"`
	Address   string        `json:"address"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Interface InterfaceLite `json:"interface"`
	Device    DeviceLite    `json:"device"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
