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
	route.Post("/", hdlr.AddFat)
	route.Post("/upload-csv", hdlr.UpdateFats)
	route.Delete("/:id", hdlr.DeleteOne)
	route.Get("/", hdlr.GetAll)
	route.Get("/:id", hdlr.GetByID)
	route.Get("/fat/:fat", hdlr.GetByFat)
	route.Get("/odn/:state/:odn", hdlr.GetByOdn)
	route.Get("/location/:state", hdlr.GetOdnStates)
	route.Get("/location/:state/:county", hdlr.GetOdnCounty)
	route.Get("/location/:state/:county/:municipality", hdlr.GetOdnMunicipality)
	route.Get("/olt/:oltIP", hdlr.GetOdnByOlt)
	route.Get("/olt/:oltIP/:shell/:card/:port", hdlr.GetOdnOltPort)
}
