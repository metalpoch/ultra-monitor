package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewInterfacesOltRoutes(app *fiber.App, db *sqlx.DB, secret []byte) {
	authUsecase := *usecase.NewUserUsecase(db, secret)

	hdlr := &handler.InterfacesOltHandler{
		Usecase: usecase.NewInterfaceOltUsecase(db),
	}

	route := app.Group("/api/interfaces-olt")

	route.Use(middleware.ValidateJWT(authUsecase, secret))
	route.Use(middleware.AdminAccess)

	route.Get("/", hdlr.GetAll)
	route.Patch("/", hdlr.Update)
}

