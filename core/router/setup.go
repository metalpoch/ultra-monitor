package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(server *fiber.App, db *gorm.DB, telegram tracking.SmartModule, cache *redis.Client) {
	newTrafficRouter(server, db, telegram, cache)
	newInfoRouter(server, db, telegram)
}
