package model

import "time"

type StatusCounts struct {
	Day       time.Time `db:"day"`
	PonID     int32     `db:"pon_id"`
	OltIP     string    `db:"olt_ip"`
	Actives   uint64    `db:"actives"`
	Inactives uint64    `db:"inactives"`
	Unknowns  uint64    `db:"unknowns"`
}

type FatMeta struct {
	PonID int32  `db:"pon_id"`
	FatID int32  `db:"fat_id"`
	OltIP string `db:"olt_ip"`
}
