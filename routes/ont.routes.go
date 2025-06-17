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

	route.Get("/status/ip/:ip", hdlr.OntStatusIPSummary)
	route.Get("/status/state", hdlr.OntStatusStateSummary)
	route.Get("/status/state/:state", hdlr.OntStatusByStateSummary)
	route.Get("/status/municipality/:state", hdlr.OntStatusMunicipalitySummary)
	route.Get("/status/county/:state/:municipality", hdlr.OntStatusCountySummary)
	route.Get("/status/odn/:state/:municipality/:county", hdlr.OntStatusOdnSummary)

	route.Get("/traffic/id/:ponID/:idx", hdlr.TrafficOnt)
	route.Get("/traffic/despt/:despt", hdlr.TrafficOntByDespt)

	route.Get("/forecast/status", hdlr.GetStatusSummaryForecast)
	route.Get("/forecast/status/:state", hdlr.OntStatusByStateForecast)
}
