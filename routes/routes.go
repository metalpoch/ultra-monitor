package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
)

func Init(app *fiber.App, db *sqlx.DB, cache *cache.Redis, secret []byte, reportsDir string, prometheus *prometheus.Prometheus) {
	NewAuthRoutes(app, db, secret)
	NewFatRoutes(app, db)
	NewReportRoutes(app, db, cache, reportsDir)
	NewTrafficRoutes(app, db, prometheus)
}
