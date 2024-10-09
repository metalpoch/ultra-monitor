package model

import "time"

type Template struct {
	ID        uint
	Name      string
	OidBw     string
	OidIn     string
	OidOut    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AddTemplate struct {
	Name   string `json:"name"`
	OidBw  string `json:"oid_bw"`
	OidIn  string `json:"oid_in"`
	OidOut string `json:"oid_out"`
}
