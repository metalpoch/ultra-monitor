package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/core/handler"
	"github.com/metalpoch/olt-blueprint/core/usecase"
	"gorm.io/gorm"
)

func newSummaryRouter(server *fiber.App, db *gorm.DB, cache cache.Redis) {
	hdlr := handler.SummaryHandler{
		Usecase: *usecase.NewSummaryUsecase(db, cache),
	}

	summary := server.Group("/summary")
	summary.Get("/users", hdlr.UserStatus)
}
