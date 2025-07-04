package dto

import "time"

type Fat struct {
	Fat          string    `json:"fat"`
	Region       string    `json:"region"`
	State        string    `json:"state"`
	Municipality string    `json:"municipality"`
	County       string    `json:"county"`
	Odn          string    `json:"odn"`
	OltIP        string    `json:"olt_ip"`
	Shell        uint8     `json:"shell"`
	Card         uint8     `json:"card"`
	Port         uint8     `json:"port"`
	Actives      int       `json:"actives"`
	Inactive     int       `json:"inactive"`
	Others       int       `json:"others"`
	ReportDay    time.Time `json:"report_day"`
}
