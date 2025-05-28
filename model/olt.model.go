package model

import "time"

type Olt struct {
	ID          uint64    `json:"id"`
	IP          string    `json:"ip"`
	Community   string    `json:"community"`
	SysName     string    `json:"sys_name"`
	SysLocation string    `json:"sys_location"`
	IsAlive     bool      `json:"is_alive"`
	LastCheck   time.Time `json:"last_check"`
	CreatedAt   time.Time `json:"created_at"`
}
