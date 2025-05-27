package config

// Config represents the configuration structure for the application.
type Config struct {
	DatabaseURI string `json:"db_uri"`
	CacheURI    string `json:"cache_uri"`
	// SmartModuleTelegramURL string `json:"smart_module_telegram_url"`
	// SmartModuleOSMURL      string `json:"smart_module_osm_url"`
	// SecretKey              string `json:"secret_key"`
	// StaticReportDirectory  string `json:"static_report_directory"`
}
