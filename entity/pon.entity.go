package entity

import "time"

type Pon struct { // Table
	ID        uint64    `db:"id"`
	OltID     uint64    `db:"olt_id"`
	IfIndex   uint64    `db:"if_index"`
	IfName    string    `db:"if_name"`
	IfDescr   string    `db:"if_descr"`
	IfAlias   string    `db:"if_alias"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
