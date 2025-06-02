package model

import "time"

type Pon struct {
	ID        int32     `json:"id"`
	OltIP     string    `json:"olt_ip"`
	IfIndex   int64     `json:"if_index"`
	IfName    string    `json:"if_name"`
	IfDescr   string    `json:"if_descr"`
	IfAlias   string    `json:"if_alias"`
	CreatedAt time.Time `json:"created_at"`
}
