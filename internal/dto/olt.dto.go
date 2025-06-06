package dto

type NewOlt struct {
	IP        string `json:"ip" validate:"required,ip_addr"`
	Community string `json:"community" validate:"required"`
}

type OltIP struct {
	IP string `url:"ip" validate:"required,ip_addr"`
}
