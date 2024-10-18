package model

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	Validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

type Config struct {
	IsProduction   bool   `json:"is_production"`
	DatabaseURI    string `json:"db_uri"`
	SmartModuleURL string `json:"smart_module_url"`
	CacheURI       string `json:"cache_uri"`
	SecretKey      string `json:"secret_key"`
}
