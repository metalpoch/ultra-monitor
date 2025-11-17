package entity

import "time"

type InterfaceBandwidth struct {
	OltVerbose string    `db:"olt_verbose"`
	Interface  string    `db:"interface"`
	Bandwidth  float64   `db:"bandwidth"`
	CreatedAt  time.Time `db:"created_at"`
}

type VerboseOltID struct {
	OltVerbose string `db:"olt_verbose"`
	OltIP      string `db:"olt_ip"`
}
