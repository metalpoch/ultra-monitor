package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewPonRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis) {
	hdlr := &handler.PonHandler{
		Usecase: usecase.NewPonUsecase(db, cache),
	}

	route := app.Group("/api/pon")

	route.Get("/traffic/:sysname", hdlr.TrafficOlt)
	route.Get("/traffic/:sysname/:shell/:card/:port", hdlr.TrafficPon)

	route.Get("/traffic/summary", hdlr.GetTrafficSummary)
	route.Get("/traffic/summary/states", hdlr.GetTrafficStatesSummary)
	route.Get("/traffic/summary/municipality/:state", hdlr.GetTrafficMunicipalitySummary)
	route.Get("/traffic/summary/county/:state/:municipality", hdlr.GetTrafficCountySummary)
	route.Get("/traffic/summary/odn/:state/:municipality/:county", hdlr.GetTrafficOdnSummary)

	route.Get("/traffic/forecast", hdlr.GetTrafficPonForecast)

	route.Get("/:sysname", hdlr.GetAllByDevice)
	route.Get("/:sysname/:shell/:card/:port", hdlr.GetByOltAndPort)
}
