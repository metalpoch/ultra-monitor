package entity

import "time"

type InterfaceBandwidth struct {
	Interface string    `db:"interface"`
	Olt       string    `db:"olt"`
	Bandwidth float64   `db:"bandwidth"`
	CreatedAt time.Time `db:"created_at"`
}

