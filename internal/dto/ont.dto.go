package dto

type AllSerialDesptByPon struct {
	IP        string `json:"ip" validate:"required,ipv4"`
	Community string `json:"community" validate:"required"`
	PonIdx    int64  `json:"pon_idx" validate:"required"`
}

type OntSerialsAndDespts struct {
	PonIdx       int64  `json:"pon_idx"`
	OntIdx       uint8  `json:"ont_idx"`
	Despt        string `json:"despt"`
	SerialNumber string `json:"serial"`
}

type CreateOntRequest struct {
	IP          string `json:"ip" validate:"required,ipv4"`
	Community   string `json:"community" validate:"required"`
	PonIdx      int64  `json:"pon_idx" validate:"required"`
	OntIdx      uint8  `json:"ont_idx" validate:"required"`
	Description string `json:"description"`
}

type OntResponse struct {
	ID          int32  `json:"id"`
	IP          string `json:"ip"`
	OntIDX      string `json:"ont_idx"`
	Serial      string `json:"serial"`
	Despt       string `json:"despt"`
	LineProf    string `json:"line_prof"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Status      bool   `json:"status"`
	LastCheck   string `json:"last_check"`
	CreatedAt   string `json:"created_at"`
}

type IDRequest struct {
	ID int32 `json:"id" validate:"required"`
}
