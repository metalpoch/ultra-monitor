package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewTrafficRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis, prometheus *prometheus.Prometheus, secret []byte) {
	authUsecase := *usecase.NewUserUsecase(db, secret)

	hdlr := &handler.TrafficHandler{Usecase: usecase.NewTrafficUsecase(db, cache, prometheus)}

	route := app.Group("/api/traffic")

	route.Use(middleware.ValidateJWT(authUsecase, secret))

	route.Get("/info", hdlr.DeviceLocation)
	route.Get("/info/instance/:ip", hdlr.InfoInstance)

	// Total
	route.Get("/total", hdlr.GetNationalTraffic)
	route.Get("/region/:region", hdlr.GetRegionalTraffic)
	route.Get("/state/:state", hdlr.GetStateTraffic)
	route.Get("/olt/:ip", hdlr.GetOLTByIPTraffic)

	// Details
	route.Get("/regions", hdlr.GetTrafficByRegions)
	route.Get("/states/:region", hdlr.GetTrafficByStates)
	route.Get("/sysname/:state", hdlr.GetTrafficByIPs)

	// stats
	route.Get("/stats/region/:region", hdlr.RegionStats)
	route.Get("/stats/state/:state", hdlr.StateStats)
	route.Get("/stats/ip/:ip", hdlr.GponStats)

	// Basic
	route.Get("/basic/criteria/:criteria/:value", hdlr.GetTrafficByCriteria)
	route.Get("/basic/instances", hdlr.GroupIP)
	route.Get("/basic/index/:ip/:idx", hdlr.ByIdx)

	// using fats table
	route.Get("/basic/municipality/:state/:municipality", hdlr.ByMunicipality)
	route.Get("/basic/county/:state/:municipality/:county", hdlr.ByCounty)
	route.Get("/basic/odn/:state/:municipality/:odn", hdlr.ByOdn)

	// Trend prediction routes
	route.Get("/trend/national", hdlr.GetNationalTrend)
	route.Get("/trend/region/:region", hdlr.GetRegionalTrend)
	route.Get("/trend/state/:state", hdlr.GetStateTrend)
	route.Get("/trend/olt/:ip", hdlr.GetOLTTrend)

	// Location hierarchy
	route.Get("/hierarchy", hdlr.GetLocationHierarchy)

	// ONT traffic
	route.Get("/ont/:id", hdlr.GetOntTraffic)
}
