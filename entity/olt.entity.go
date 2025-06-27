package entity

import "time"

type Olt struct {
	IP           string    `db:"ip"`
	Community    string    `db:"community"`
	SysName      string    `db:"sys_name"`
	SysLocation  string    `db:"sys_location"`
	State        string    `db:"state"`
	Municipality string    `db:"municipality"`
	County       string    `db:"county"`
	Odn          string    `db:"odn"`
	Fat          string    `db:"fat"`
	Shell        uint8     `db:"pon_shell"`
	Card         uint8     `db:"pon_card"`
	Port         uint8     `db:"pon_port"`
	IsAlive      bool      `db:"is_alive"`
	LastCheck    time.Time `db:"last_check"`
	CreatedAt    time.Time `db:"created_at"`
}
