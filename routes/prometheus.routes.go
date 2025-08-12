package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewPrometheusRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis) {
	hdlr := &handler.PrometheusHandler{
		Usecase: usecase.NewPrometheusUsecase(db),
	}

	route := app.Group("/api/prometheus")
	route.Get("/status", hdlr.GetGponPortsStatus)
	route.Get("/status/region/:region", hdlr.GetGponPortsStatusByRegion)
	route.Get("/status/state/:state", hdlr.GetGponPortsStatusByState)
}
