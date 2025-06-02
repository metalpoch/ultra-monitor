package entity

import "time"

type Pon struct {
	ID        int32     `db:"id"`
	OltIP     string    `db:"olt_ip"`
	IfIndex   int64     `db:"if_index"`
	IfName    string    `db:"if_name"`
	IfDescr   string    `db:"if_descr"`
	IfAlias   string    `db:"if_alias"`
	CreatedAt time.Time `db:"created_at"`
}
