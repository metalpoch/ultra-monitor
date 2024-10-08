package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/report/handler"
	"github.com/metalpoch/olt-blueprint/report/usecase"
	"gorm.io/gorm"
)

func newFatRouter(server *fiber.App, db *gorm.DB, telegram tracking.Telegram) {
	hdlr := handler.FatHandler{
		Usecase: *usecase.NewFatUsecase(db, telegram),
	}

	server.Post("/fat/", hdlr.Add)
	server.Get("/fat/:id", hdlr.Get)
}
