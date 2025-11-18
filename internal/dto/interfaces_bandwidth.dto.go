package dto

import "time"

type InterfaceBandwidth struct {
	OltVerbose string    `json:"olt_verbose"`
	Interface  string    `json:"interface"`
	Bandwidth  float64   `json:"bandwidth"`
	CreatedAt  time.Time `json:"created_at"`
}
