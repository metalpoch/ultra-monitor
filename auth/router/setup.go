package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"gorm.io/gorm"
)

func Setup(server *fiber.App, db *gorm.DB, secret []byte, telegram tracking.SmartModule) {
	newUserRouter(server, db, secret, telegram)
}
