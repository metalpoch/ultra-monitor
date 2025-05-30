package model

import "time"

type Fat struct {
	ID           int32     `json:"id"`
	Fat          string    `db:"fat"`
	Region       string    `json:"region"`
	State        string    `json:"state"`
	Municipality string    `json:"municipality"`
	County       string    `json:"county"`
	Odn          string    `json:"odn"`
	OltIP        string    `json:"olt_ip"`
	Shell        uint8     `json:"pon_shell"`
	Port         uint8     `json:"pon_port"`
	Card         uint8     `json:"pon_card"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
}
