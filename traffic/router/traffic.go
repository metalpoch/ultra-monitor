package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/common/usecase"
	"github.com/metalpoch/olt-blueprint/traffic/handler"
	"gorm.io/gorm"
)

func newTrafficRouter(server *fiber.App, db *gorm.DB, telegram tracking.Telegram) {
	hdlr := handler.TrafficHandler{
		Usecase: *usecase.NewTrafficUsecase(db, telegram),
	}

	server.Get("/traffic/interface/:id", hdlr.GetByInterface)
	server.Get("/traffic/device/:id", hdlr.GetByDevice)
}
