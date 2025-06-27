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
	route.Get("/:ip", hdlr.GetOlt)
	route.Delete("/:ip", hdlr.DeleteOne)

	route.Get("/ip/", hdlr.GetAllIP)
	route.Get("/sysname", hdlr.GetAllSysname)

	route.Get("/location/:state", hdlr.GetOltsByState)
}
