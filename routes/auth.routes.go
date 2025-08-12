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
	route.Post("/signin", hdlr.Login)
	route.Get("/", middleware.ValidateJWT(secret), hdlr.GetOwn)
	route.Patch("/reset_password", middleware.ValidateJWT(secret), hdlr.ChangePassword)
	route.Post("/signup", middleware.ValidateJWT(secret), middleware.AdminAccess, hdlr.Create)
	route.Post("/:id", middleware.ValidateJWT(secret), middleware.AdminAccess, hdlr.Enable)
	route.Delete("/:id", middleware.ValidateJWT(secret), middleware.AdminAccess, hdlr.Disable)
}
