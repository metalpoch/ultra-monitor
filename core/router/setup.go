package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"gorm.io/gorm"
)

func Setup(server *fiber.App, db *gorm.DB, telegram tracking.SmartModule, cache cache.Redis) {
	newTrafficRouter(server, db, telegram, cache)
	newInfoRouter(server, db, telegram)
}
