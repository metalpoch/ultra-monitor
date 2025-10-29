package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
)

type Config struct {
	App        *fiber.App
	DB         *sqlx.DB
	Cache      *cache.Redis
	Secret     []byte
	Prometheus *prometheus.Prometheus
	WebAppDir  string
	ReportsDir string
	Enviroment string
}

func Init(cfg *Config) {
	if cfg.Enviroment == "production" {
		cfg.App.Use(favicon.New(favicon.Config{
			File: cfg.WebAppDir + "/favicon.svg",
			URL:  "/favicon.svg",
		}))
		NewWebRoutes(cfg.App, cfg.WebAppDir)
	}

	NewAuthRoutes(cfg.App, cfg.DB, cfg.Secret)
	NewFatRoutes(cfg.App, cfg.DB, cfg.Secret)
	NewReportRoutes(cfg.App, cfg.DB, cfg.Cache, cfg.ReportsDir, cfg.Secret)
	NewTrafficRoutes(cfg.App, cfg.DB, cfg.Cache, cfg.Prometheus, cfg.Secret)
	NewPrometheusRoutes(cfg.App, cfg.DB, cfg.Cache, cfg.Secret)
	NewOntRoutes(cfg.App, cfg.DB, cfg.Cache, cfg.Secret)
}
