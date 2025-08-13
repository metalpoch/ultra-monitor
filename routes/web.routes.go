package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func NewWebRoutes(app *fiber.App, webAppDir string) {

	app.Get("/", static.New(webAppDir+"/index.html"))
	app.Get("/auth/login", static.New(webAppDir+"/auth/login"))
	app.Get("/_astro/*", static.New(webAppDir+"/_astro/"))
	app.Get("/404", static.New(webAppDir+"/404.html"))
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
}
