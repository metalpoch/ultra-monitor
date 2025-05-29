package entity

import "time"

type Olt struct {
	ID          uint64    `db:"id"`
	IP          string    `db:"ip"`
	Community   string    `db:"community"`
	SysName     string    `db:"sys_name"`
	SysLocation string    `db:"sys_location"`
	IsAlive     bool      `db:"is_alive"`
	LastCheck   time.Time `db:"last_check"`
	CreatedAt   time.Time `db:"created_at"`
}
