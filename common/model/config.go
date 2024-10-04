package model

import "github.com/go-playground/validator/v10"

type StructValidator struct {
	Validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

type Config struct {
	IsProduction       bool   `json:"is_production"`
	DatabaseURI        string `json:"db_uri"`
	CacheURI           string `json:"cache_uri"`
	SecretKey          string `json:"secret_key"`
	TelegramChatID     string `json:"telegram_chat_id"`
	TelegramBotTokenID string `json:"telegram_bot_token_id"`
}
