package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewTrafficRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis) {
	hdlr := &handler.TrafficHandler{
		Usecase: usecase.NewTrafficUsecase(db, cache),
	}

	route := app.Group("/api/pon/traffic")

	route.Get("/all", hdlr.GetTotalTraffic)
	route.Get("/location/:state", hdlr.TrafficByState)
	route.Get("/location/:state/:county", hdlr.TrafficByCounty)
	route.Get("/location/:state/:county/:municipality", hdlr.TrafficByMunicipaly)
	route.Get("/odn/:state/:odn", hdlr.TrafficByODN)
	route.Get("/pon/:sysname/:shell/:card/:port", hdlr.TrafficPon)
}
