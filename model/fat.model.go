package model

import "time"

type Fat struct {
	ID           int32     `json:"id"`
	Fat          string    `json:"fat" csv:"fat"`
	Region       string    `json:"region" csv:"region"`
	State        string    `json:"state" csv:"state"`
	Municipality string    `json:"municipality" csv:"municipality"`
	County       string    `json:"county" csv:"county"`
	Odn          string    `json:"odn" csv:"odn"`
	OltIP        string    `json:"olt_ip" csv:"olt_ip"`
	Shell        uint8     `json:"pon_shell" csv:"pon_shell"`
	Port         uint8     `json:"pon_port" csv:"pon_port"`
	Card         uint8     `json:"pon_card" csv:"pon_card"`
	Latitude     float64   `json:"latitude" csv:"latitude"`
	Longitude    float64   `json:"longitude" csv:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
}
