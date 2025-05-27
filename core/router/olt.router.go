package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/core/handler"
	"github.com/metalpoch/olt-blueprint/core/usecase"
	"gorm.io/gorm"
)

func newOltRouter(server *fiber.App, db *gorm.DB, cache cache.Redis) {
	hdlr := handler.OltHandler{
		Usecase: *usecase.NewTrafficUsecase(db, cache),
	}
	/*
	 *        server.Get("/traffic/interface/:id", hdlr.GetByInterface)
	 *        server.Get("/traffic/device/:id", hdlr.GetByDevice)
	 *        server.Get("/traffic/fat/:id", hdlr.GetByFat)
	 *        server.Get("/traffic/location/:state", hdlr.GetByState)
	 *        server.Get("/traffic/location/:state/:county", hdlr.GetByCounty)
	 *        server.Get("/traffic/location/:state/:county/:municipality", hdlr.GetByMunicipaly)
	 *        server.Get("/traffic/odn/:odn", hdlr.GetTrafficByODN)
	 *        server.Get(("/traffic/state/:month"), hdlr.GetTotalTrafficByState)
	 *        server.Get(("/traffic/state_n/:month/:n"), hdlr.GetTotalTrafficByState_N)
	 *        server.Get(("/traffic/odn_d/:month"), hdlr.GetTotalTrafficByOND)*/
}
