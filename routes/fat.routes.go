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
	route.Post("/", hdlr.UpsertFats)
	route.Delete("/:id", hdlr.DeleteOne)
	route.Get("/:id", hdlr.GetByID)

	// Fat info by location
	route.Get("/location/:state", hdlr.FindByStates)
	route.Get("/location/:state/:municipality", hdlr.FindByMunicipality)
	route.Get("/location/:state/:municipality/:county", hdlr.FindByCounty)
	route.Get("/location/:state/:municipality/:county/:odn", hdlr.FindBytOdn)
}
