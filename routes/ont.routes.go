package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewOntRoutes(app *fiber.App, db *sqlx.DB, secret []byte) {
	authUsecase := *usecase.NewUserUsecase(db, secret)

	hdlr := &handler.OntHandler{
		Usecase: usecase.NewOntUsecase(db),
	}

	route := app.Group("/api/ont")

	route.Use(middleware.ValidateJWT(authUsecase, secret))
	route.Use(middleware.AdminAccess)

	route.Post("/snmp/serial-despt", hdlr.SnmpSerialDespt)
	route.Get("/", hdlr.GetAll)
	route.Post("/", hdlr.SaveOnt)
	route.Get("/:id", hdlr.GetByID)
	route.Delete("/:id", hdlr.Delete)
	route.Patch("/:id/enable", hdlr.Enable)
	route.Patch("/:id/disable", hdlr.Disable)
}
