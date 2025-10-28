package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewPrometheusRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis, secret []byte) {
	authUsecase := *usecase.NewUserUsecase(db, secret)

	hdlr := &handler.PrometheusHandler{
		Usecase: usecase.NewPrometheusUsecase(db),
	}

	route := app.Group("/api/prometheus")

	route.Use(middleware.ValidateJWT(authUsecase, secret))

	route.Get("/pons/ip/:ip", hdlr.GetDeviceByIP)
	route.Get("/status", hdlr.GetGponPortsStatus)
	route.Get("/status/region/:region", hdlr.GetGponPortsStatusByRegion)
	route.Get("/status/state/:state", hdlr.GetGponPortsStatusByState)
}
