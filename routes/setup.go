package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

func Setup(server *fiber.App, db *sqlx.DB, secret []byte) {
	newUserRoutes(server, db, secret)
}
