package model

import "time"

type Pon struct {
	ID        int32     `json:"id"`
	OltID     int32     `json:"olt_id"`
	IfIndex   int32     `json:"if_index"`
	IfName    string    `json:"if_name"`
	IfDescr   string    `json:"if_descr"`
	IfAlias   string    `json:"if_alias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
