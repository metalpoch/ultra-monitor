package model

type Device struct {
	IP        string `json:"ip"`
	Community string `json:"community"`
	Sysname   string `json:"sysname"`
}
