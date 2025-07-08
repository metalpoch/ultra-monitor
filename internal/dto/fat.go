package dto

import "time"

type Fat struct {
	Fat          string    `json:"fat" validate:"required"`
	Region       string    `json:"region" validate:"required"`
	State        string    `json:"state" validate:"required"`
	Municipality string    `json:"municipality" validate:"required"`
	County       string    `json:"county" validate:"required"`
	Odn          string    `json:"odn"`
	OltIP        string    `json:"olt_ip" validate:"required,ip_addr"`
	Shell        uint8     `json:"shell"`
	Card         uint8     `json:"card"`
	Port         uint8     `json:"port"`
	Actives      int       `json:"actives"`
	Inactive     int       `json:"inactive"`
	Others       int       `json:"others"`
	ReportDay    time.Time `json:"report_day" validate:"required"`
}
