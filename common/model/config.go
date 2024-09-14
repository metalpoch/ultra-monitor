package model

type Config struct {
	DatabaseURI string `json:"db_uri"`
	CacheURI    string `json:"cache_uri"`
	SecretKey   string `json:"secret_key"`
}
