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
	route.Get("/:sysname", hdlr.GetAllByDevice)
	route.Get("/:sysname/:shell/:card/:port", hdlr.GetByOltAndPort)
	route.Get("/traffic/location/:state", hdlr.TrafficByState)
	route.Get("/traffic/location/:state/:municipality", hdlr.TrafficByMunicipaly)
	route.Get("/traffic/location/:state/:municipality/:county/:", hdlr.TrafficByCounty)
	route.Get("/traffic/odn/:state/:odn", hdlr.TrafficByODN)
	route.Get("/traffic/pon/:sysname/:shell/:card/:port", hdlr.TrafficPon)
	route.Get("/traffic/summary", hdlr.GetTrafficSummary)
	route.Get("/traffic/forecast", hdlr.GetTrafficPonForecast)
}
