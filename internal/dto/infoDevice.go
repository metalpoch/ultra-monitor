package dto

type InfoDevice struct {
	Fat          string `json:"fat" validate:"required"`
	Region       string `json:"region" validate:"required"`
	State        string `json:"state" validate:"required"`
	Municipality string `json:"municipality" validate:"required"`
	County       string `json:"county" validate:"required"`
	Odn          string `json:"odn"`
	IP           string `json:"ip" validate:"required,ip_addr"`
	Shell        uint8  `json:"shell"`
	Card         uint8  `json:"card"`
	Port         uint8  `json:"port"`
}
