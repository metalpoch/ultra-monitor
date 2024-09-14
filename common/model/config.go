package model

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	Validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

type Config struct {
	DatabaseURI string `json:"db_uri"`
	CacheURI    string `json:"cache_uri"`
	SecretKey   string `json:"secret_key"`
}
