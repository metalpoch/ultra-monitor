package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewAuthRoutes(app *fiber.App, db *sqlx.DB, secret []byte) {
	hdlr := &handler.UserHandler{
		Usecase: *usecase.NewUserUsecase(db, secret),
	}

	route := app.Group("/api/auth")
	route.Get("/", hdlr.GetOwn, middleware.ValidateJWT(secret))
	route.Post("/signin", hdlr.Login)
	route.Patch("/reset_password", hdlr.ChangePassword, middleware.ValidateJWT(secret))
	route.Post("/:p00", hdlr.Enable, middleware.ValidateJWT(secret), middleware.AdminAccess)
	route.Delete("/:p00", hdlr.Disable, middleware.ValidateJWT(secret), middleware.AdminAccess)
	route.Post("/signup", hdlr.Create, middleware.ValidateJWT(secret), middleware.AdminAccess)
}
