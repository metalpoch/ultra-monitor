package model

import "time"

type Olt struct {
	IP           string    `json:"ip"`
	Community    string    `json:"community"`
	SysName      string    `json:"sys_name"`
	SysLocation  string    `json:"sys_location"`
	State        string    `json:"state"`
	Municipality string    `json:"municipality"`
	County       string    `json:"county"`
	Odn          string    `json:"odn"`
	Fat          string    `json:"fat"`
	Shell        uint8     `json:"pon_shell"`
	Card         uint8     `json:"pon_card"`
	Port         uint8     `json:"pon_port"`
	IsAlive      bool      `json:"is_alive"`
	LastCheck    time.Time `json:"last_check"`
	CreatedAt    time.Time `json:"created_at"`
}
