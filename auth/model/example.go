package model

type Example struct {
	Message string `bson:"message"`
}

type NewExample struct {
	Message string `json:"message" validate:"required"`
	Value   uint8  `json:"value" validate:"required,gte=18"`
}
