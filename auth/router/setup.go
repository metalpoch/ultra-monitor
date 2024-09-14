package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
)

func Setup(server *fiber.App, db *sql.DB) {
	newUserRouter(server, db)
}
