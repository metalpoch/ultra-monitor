package entity

import "time"

type Prometheus struct {
	Region    string    `db:"region"`
	State     string    `db:"state"`
	IP        string    `db:"ip"`
	IDX       int64     `db:"idx"`
	Shell     uint8     `db:"shell"`
	Card      uint8     `db:"card"`
	Port      uint8     `db:"port"`
	CreatedAt time.Time `db:"created_at"`
}

type PrometheusUpsert struct {
	Region string `db:"region"`
	State  string `db:"state"`
	IP     string `db:"ip"`
	IDX    int64  `db:"idx"`
	Shell  uint8  `db:"shell"`
	Card   uint8  `db:"card"`
	Port   uint8  `db:"port"`
	Status int8   `db:"status"`
}
