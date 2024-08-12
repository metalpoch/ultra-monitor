package model

type ElementOLT struct {
	ID        string `json:"_id"`
	OLT       string `json:"olt"`       // index
	Interface string `json:"interface"` //index
	Slot      int8   `json:"slot"`
	Card      int8   `json:"card"`
	Port      int8   `json:"port"`
}
