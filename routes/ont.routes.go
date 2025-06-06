package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewOntRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis) {
	hdlr := &handler.OntHandler{
		Usecase: *usecase.NewOntUsecase(db, cache),
	}

	route := app.Group("/api/ont")

	route.Get("/status/:state", hdlr.OntStatusByState)
	route.Get("/status/:state/:municipality/:county/:odn", hdlr.OntStatusByOdn)
	route.Get("/status/by-olt/:ip", hdlr.OntStatusByOltIP)
	route.Get("/status/by-sysname/:sysname", hdlr.OntStatusBySysname)

	route.Get("/traffic/by-id/:ponID/:idx", hdlr.TrafficOnt)
	route.Get("/traffic/by-despt/:despt", hdlr.TrafficOntByDespt)

	route.Get("/forecast/status/:state", hdlr.OntStatusByStateForecast)
	route.Get("/forecast/status/:state/:municipality/:county/:odn", hdlr.OntStatusByODNForecast)
}
