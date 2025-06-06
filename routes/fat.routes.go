package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewFatRoutes(app *fiber.App, db *sqlx.DB) {
	hdlr := &handler.FatHandler{
		Usecase: usecase.NewFatUsecase(db),
	}

	route := app.Group("/api/fat")

	// Base
	route.Get("/", hdlr.GetAll)
	route.Post("/", hdlr.AddFat)
	route.Delete("/:id", hdlr.DeleteOne)
	route.Get("/:id", hdlr.GetByID)

	// Locations
	route.Get("/info/location/", hdlr.GetStates)
	route.Get("/info/location/:state/", hdlr.GetMunicipality)
	route.Get("/info/location/:state/:municipality/", hdlr.GetCounty)

	// Fat info by location
	route.Get("/location/:state", hdlr.GetFatsByStates)
	route.Get("/location/:state/:municipality", hdlr.GetFatsByMunicipality)
	route.Get("/location/:state/:municipality/:county", hdlr.GetFatsByCounty)
}
