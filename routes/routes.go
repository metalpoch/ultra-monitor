package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
)

func Init(app *fiber.App, db *sqlx.DB, cache *cache.Redis, secret []byte) {
	NewAuthRoutes(app, db, secret)
	NewFatRoutes(app, db)
	NewOltRoutes(app, db)
	NewOntRoutes(app, db, cache)
	NewPonRoutes(app, db)
	NewTrafficRoutes(app, db, cache)
	// NewReportRoutes(app, db)
}
