package model

import "time"

type Fat struct {
	ID           int32     `json:"id"`
	IP           string    `json:"ip"`
	Region       string    `json:"region"`
	State        string    `json:"state"`
	Municipality string    `json:"municipality"`
	County       string    `json:"county"`
	Odn          string    `json:"odn"`
	Fat          string    `json:"fat"`
	Shell        uint8     `json:"shell"`
	Card         uint8     `json:"card"`
	Port         uint8     `json:"port"`
	CreatedAt    time.Time `json:"created_at"`
}
