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
	route.Get("/status", hdlr.OntStatus)
	route.Get("/status/:state", hdlr.OntStatusByState)
	route.Get("/status/:state/:odn", hdlr.OntStatusByOdn)
	route.Get("/status/by-oltp/:ip", hdlr.OntStatusByOltIP)
	route.Get("/status/by-sysname/:sysname", hdlr.OntStatusBySysname)
	route.Get("/traffic/:ponID/:idx", hdlr.TrafficOnt)
	route.Get("/forecast/status", hdlr.AllOntStatusForecast)
	route.Get("/forecast/status/:state", hdlr.OntStatusByStateForecast)
	route.Get("/forecast/status/:state/:odn", hdlr.OntStatusByODNForecast)
}
