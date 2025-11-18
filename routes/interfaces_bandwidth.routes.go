package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/mongodb"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewInterfacesBandwidthRoutes(app *fiber.App, db *sqlx.DB, mongo *mongodb.MongoDB, secret []byte) {
	authUsecase := *usecase.NewUserUsecase(db, secret)

	hdlr := &handler.InterfacesBandwidthHandler{
		Usecase: usecase.NewInterfaceBandwidthUsecase(db, mongo),
	}

	route := app.Group("/api/interfaces-bandwidth")

	route.Use(middleware.ValidateJWT(authUsecase, secret))
	route.Use(middleware.AdminAccess)

	route.Get("/", hdlr.GetAll)
}

