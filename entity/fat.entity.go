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
	Card         uint8     `db:"pon_card"`
	Port         uint8     `db:"pon_port"`
	CreatedAt    time.Time `db:"created_at"`
}

type FatResponse struct {
	ID           int32     `db:"id"`
	Fat          string    `db:"fat"`
	Region       string    `db:"region"`
	State        string    `db:"state"`
	Municipality string    `db:"municipality"`
	County       string    `db:"county"`
	Odn          string    `db:"odn"`
	OltIP        string    `db:"olt_ip"`
	Shell        uint8     `db:"pon_shell"`
	Card         uint8     `db:"pon_card"`
	Port         uint8     `db:"pon_port"`
	Actives      int       `db:"actives"`
	Inactives    int       `db:"inactivew"`
	Others       int       `db:"others"`
	CreatedAt    time.Time `db:"created_at"`
}

type FatsOntStatusSummary struct {
	Day       time.Time `db:"day"`
	FatID     int64     `db:"fat_id"`
	Actives   int       `db:"actives"`
	Inactives int       `db:"inactives"`
	Others    int       `db:"others"`
}
