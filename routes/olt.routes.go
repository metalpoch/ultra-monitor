package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewOltRoutes(app *fiber.App, db *sqlx.DB) {
	hdlr := &handler.OltHandler{
		Usecase: *usecase.NewOltUsecase(db),
	}

	route := app.Group("/api/olt")
	route.Post("/", hdlr.Add)
	route.Get("/", hdlr.GetOlt)
	route.Get("/:id", hdlr.GetOlt)
	route.Put("/:id", hdlr.UpdateOne)
	route.Delete("/:id", hdlr.DeleteOne)
	route.Get("/location/:state", hdlr.GetOltsByState)
	route.Get("/location/:county", hdlr.GetOltsByCounty)
	route.Get("/location/:state:/municipality", hdlr.GetOltsByMunicipality)

}
