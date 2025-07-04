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
	Shell        uint8     `json:"shell" csv:"pon_shell"`
	Card         uint8     `json:"card" csv:"pon_card"`
	Port         uint8     `json:"port" csv:"pon_port"`
	Actives      int       `db:"actives"`
	Inactive     int       `db:"inactive"`
	Others       int       `db:"others"`
	CreatedAt    time.Time `json:"created_at"`
}

type FatsOntStatusSummary struct {
	Day      time.Time `json:"day"`
	FatID    int       `json:"fat_id"`
	Actives  int       `json:"actives"`
	Inactive int       `json:"inactive"`
	Others   int       `json:"others"`
}
