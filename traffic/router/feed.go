package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/traffic/handler"
	"github.com/metalpoch/olt-blueprint/traffic/usecase"
	"gorm.io/gorm"
)

func newFeedRouter(server *fiber.App, db *gorm.DB, telegram tracking.Telegram) {
	hdlr := handler.FeedHandler{
		Usecase: *usecase.NewFeedUsecase(db, telegram),
	}

	server.Get("/feed/device/", hdlr.GetAllDevice)
	server.Get("/feed/device/:id", hdlr.GetDevice)
	server.Get("/feed/interface/:id", hdlr.GetInterface)
	server.Get("/feed/interface/device/:id", hdlr.GetInterfacesByDevice)
	server.Get("/feed/location/state", hdlr.GetLocationStates)
	server.Get("/feed/location/:state", hdlr.GetLocationCounties)
	server.Get("/feed/location/:state/:county", hdlr.GetLocationMunicipalities)
}
