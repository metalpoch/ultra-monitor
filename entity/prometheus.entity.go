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
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}

type PrometheusPortStatus struct {
	Olts          uint32 `db:"olts"`
	Cards         uint32 `db:"cards"`
	GponActives   uint32 `db:"gpon_actives"`
	GponInactives uint32 `db:"gpon_inactives"`
	TotalGpon     uint32 `db:"total_gpon"`
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
