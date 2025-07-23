package entity

import (
	"time"
)

type InfoDevice struct {
	ID           int32     `db:"id"`
	IP           string    `db:"ip"`
	Region       string    `db:"region"`
	State        string    `db:"state"`
	Municipality string    `db:"municipality"`
	County       string    `db:"county"`
	Odn          string    `db:"odn"`
	Fat          string    `db:"fat"`
	Shell        uint8     `db:"pon_shell"`
	Card         uint8     `db:"pon_card"`
	Port         uint8     `db:"pon_port"`
	CreatedAt    time.Time `db:"created_at"`
}
