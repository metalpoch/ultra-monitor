package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/core/handler"
	"github.com/metalpoch/olt-blueprint/core/usecase"
	"gorm.io/gorm"
)

func newSummaryRouter(server *fiber.App, db *gorm.DB, cache cache.Redis) {
	hdlr := handler.OntHandler{
		Usecase: *usecase.NewOntUsecase(db, cache),
	}

	ont := server.Group("/ont")
	ont.Get("/summary/status", hdlr.OntStatus)
	ont.Get("/summary/status/:state", hdlr.OntStatusByState)
	ont.Get("/summary/status/:state/:odn", hdlr.OntStatusByOdn)
	ont.Get("/traffic/:interfaceID/:idx", hdlr.TrafficOnt)
	/*ont.Get("/traffic/state/all", hdlr.Traffic)*/
	/*ont.Get("/traffic/state/:state", hdlr.TrafficByState)*/
}
