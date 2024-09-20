package model

import "time"

type Template struct {
	ID        uint
	Name      string
	OidBw     string
	OidIn     string
	OidOut    string
	PrefixBw  string
	PrefixIn  string
	PrefixOut string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AddTemplate struct {
	Name      string `json:"name"`
	OidBw     string `json:"oid_bw"`
	OidIn     string `json:"oid_in"`
	OidOut    string `json:"oid_out"`
	PrefixBw  string `json:"prefix_bw"`
	PrefixIn  string `json:"prefix_in"`
	PrefixOut string `json:"prefix_out"`
}
