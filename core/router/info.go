package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/core/handler"
	"github.com/metalpoch/olt-blueprint/core/usecase"
	"gorm.io/gorm"
)

func newInfoRouter(server *fiber.App, db *gorm.DB, telegram tracking.SmartModule) {
	hdlr := handler.InfoHandler{
		Usecase: *usecase.NewInfoUsecase(db, telegram),
	}

	server.Get("/info/device/", hdlr.GetAllDevice)
	server.Get("/info/device/ip/:ip", hdlr.GetDeviceByIP)
	server.Get("/info/device/:id", hdlr.GetDevice)
	server.Get("/info/device/location/:state", hdlr.GetDeviceByState)
	server.Get("/info/device/location/:state/:county", hdlr.GetDeviceByCounty)
	server.Get("/info/device/location/:state/:county/:municipality", hdlr.GetDeviceByMunicipality)

	server.Get("/info/interface/:id", hdlr.GetInterface)
	server.Get("/info/interface/device/:id", hdlr.GetInterfacesByDevice)
	server.Get("/info/interface/device/:id/find", hdlr.GetInterfacesByDeviceAndPorts)

	server.Get("/info/location/state", hdlr.GetLocationStates)
	server.Get("/info/location/:state", hdlr.GetLocationCounties)
	server.Get("/info/location/:state/:county", hdlr.GetLocationMunicipalities)
}
