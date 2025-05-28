package dto

type NewOlt struct {
	IP        string `json:"ip" validate:"ipv4,required"`
	Community string `json:"community" validate:"required"`
}
