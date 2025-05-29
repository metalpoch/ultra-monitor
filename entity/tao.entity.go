package entity

import (
	"time"
)

type Tao struct {
	ID           uint64    `db:"id"`
	Tao          string    `db:"tao"`
	Region       string    `db:"region"`
	State        string    `db:"state"`
	Municipality string    `db:"municipality"`
	County       string    `db:"county"`
	Sector       string    `db:"sector"`
	Odn          string    `db:"odn"`
	OltIP        string    `db:"olt_ip"`
	Shell        uint8     `db:"pon_shell"`
	Port         uint8     `db:"pon_port"`
	Card         uint8     `db:"pon_card"`
	Latitude     float64   `db:"latitude"`
	Longitude    float64   `db:"longitude"`
	CreatedAt    time.Time `db:"created_at"`
}
