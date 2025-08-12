package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewTrafficRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis, prometheus *prometheus.Prometheus) {
	hdlr := &handler.TrafficHandler{Usecase: usecase.NewTrafficUsecase(db, cache, prometheus)}

	route := app.Group("/api/traffic")
	route.Get("/info", hdlr.DeviceLocation)
	route.Get("/info/instance/:ip", hdlr.InfoInstance)

	// Total
	route.Get("/total", hdlr.Total)
	route.Get("/region/:region", hdlr.Region)
	route.Get("/state/:state", hdlr.State)
	route.Get("/instances", hdlr.GroupIP)
	route.Get("/instance/:ip/:index", hdlr.IndexAndIP)

	// Details
	route.Get("/regions", hdlr.Regions)
	route.Get("/states/:region", hdlr.StatesByRegion)
	route.Get("/sysname/:state", hdlr.SysnameByState)
}
