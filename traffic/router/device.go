package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/common/usecase"
	"github.com/metalpoch/olt-blueprint/traffic/handler"
	"gorm.io/gorm"
)

func newDeviceRouter(server *fiber.App, db *gorm.DB, telegram tracking.Telegram) {
	hdlr := handler.DeviceHandler{
		Usecase: *usecase.NewDeviceUsecase(db, telegram),
	}

	server.Get("/device/", hdlr.GetAll)
	server.Get("/device/:id", hdlr.GetByID)
}
