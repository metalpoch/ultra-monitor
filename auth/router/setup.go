package router

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Setup(server *fiber.App, db *gorm.DB, secret []byte) {
	newUserRouter(server, db, secret)
}
