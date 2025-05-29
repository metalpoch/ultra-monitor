package model

import "time"

type Tao struct {
	ID           uint64    `json:"id"`
	Tao          string    `db:"tao"`
	Region       string    `json:"region"`
	State        string    `json:"state"`
	Municipality string    `json:"municipality"`
	County       string    `json:"county"`
	Sector       string    `json:"sector"`
	Odn          string    `json:"odn"`
	OltIP        string    `json:"olt_ip"`
	Shell        uint8     `json:"pon_shell"`
	Port         uint8     `json:"pon_port"`
	Card         uint8     `json:"pon_card"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
}
