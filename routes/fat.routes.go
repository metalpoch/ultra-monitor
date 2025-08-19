package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewFatRoutes(app *fiber.App, db *sqlx.DB, secret []byte) {
	authUsecase := *usecase.NewUserUsecase(db, secret)

	hdlr := &handler.FatHandler{
		Usecase: usecase.NewFatUsecase(db),
	}

	route := app.Group("/api/fat")

	route.Use(middleware.ValidateJWT(authUsecase, secret))

	// Base
	route.Get("/", hdlr.GetAll)
	route.Post("/", hdlr.UpsertFats)
	route.Get("/id/:id", hdlr.GetByID)
	route.Get("/ip/:ip", hdlr.GetAllByIP)

	// Locations
	route.Get("/regions/all", hdlr.GetRegions)
	route.Get("/regions/:region", hdlr.GetStates)
	route.Get("/regions/:region/:state", hdlr.GetMunicipalities)
	route.Get("/regions/:region/:state/:municipality", hdlr.GetCounties)
	route.Get("/regions/:region/:state/:municipality/:county", hdlr.GetODN)
	route.Get("/regions/:region/:state/:municipality/:county/:odn", hdlr.GetFat)

	// Fat info by location
	route.Get("/location/:state", hdlr.FindByStates)
	route.Get("/location/:state/:municipality", hdlr.FindByMunicipality)
	route.Get("/location/:state/:municipality/:county", hdlr.FindByCounty)

	// Trend status
	route.Get("/trend/status", hdlr.GetAllFatStatus)
	route.Get("/trend/status/:region", hdlr.GetAllFatStatusByRegion)
	route.Get("/trend/status/state/:state", hdlr.GetAllFatStatusByState)
	route.Get("/trend/status/state/:state/:municipality", hdlr.GetAllFatStatusByMunicipality)
	route.Get("/trend/status/state/:state/:municipality/:county", hdlr.GetAllFatStatusByCounty)
	route.Get("/trend/status/state/:state/:municipality/:county/:odn", hdlr.GetAllFatStatusByODN)
	route.Get("/trend/status/state/:state/:municipality/:county/:odn/:fat", hdlr.GetAllFatStatusByFat)
}
