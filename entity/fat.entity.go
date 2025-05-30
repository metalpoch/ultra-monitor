package entity

import (
	"time"
)

type Fat struct {
	ID           int32     `db:"id"`
	Fat          string    `db:"fat"`
	Region       string    `db:"region"`
	State        string    `db:"state"`
	Municipality string    `db:"municipality"`
	County       string    `db:"county"`
	Odn          string    `db:"odn"`
	OltIP        string    `db:"olt_ip"`
	Shell        uint8     `db:"pon_shell"`
	Port         uint8     `db:"pon_port"`
	Card         uint8     `db:"pon_card"`
	Latitude     float64   `db:"latitude"`
	Longitude    float64   `db:"longitude"`
	CreatedAt    time.Time `db:"created_at"`
}
