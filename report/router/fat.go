package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/openstreetmap"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/report/handler"
	"github.com/metalpoch/olt-blueprint/report/usecase"
	"gorm.io/gorm"
)

func newFatRouter(server *fiber.App, db *gorm.DB, telegram tracking.SmartModule, openstreetmap openstreetmap.OSM) {
	hdlr := handler.FatHandler{
		Usecase: *usecase.NewFatUsecase(db, telegram, openstreetmap),
	}
	server.Post("/fat/", hdlr.Add)
	server.Get("/fat/", hdlr.GetAll)
	server.Delete("/fat/:id", hdlr.Delete)
	server.Get("/fat/:id", hdlr.Get)
}
