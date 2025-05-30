package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewPonRoutes(app *fiber.App, db *sqlx.DB) {
	hdlr := &handler.PonHandler{
		Usecase: usecase.NewPonUsecase(db),
	}

	route := app.Group("/api/pon")
	route.Get("/:sysname", hdlr.GetAllByDevice)
	route.Get("/:sysname/:shell/:card/:port", hdlr.GetByOltAndPort)
}
