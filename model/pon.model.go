package model

import "time"

type Pon struct {
	ID        uint64    `json:"id"`
	OltID     uint64    `json:"olt_id"`
	IfIndex   uint64    `json:"if_index"`
	IfName    string    `json:"if_name"`
	IfDescr   string    `json:"if_descr"`
	IfAlias   string    `json:"if_alias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
