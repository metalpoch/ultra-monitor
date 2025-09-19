package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/middleware"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewAuthRoutes(app *fiber.App, db *sqlx.DB, secret []byte) {
	authUsecase := *usecase.NewUserUsecase(db, secret)
	hdlr := &handler.UserHandler{Usecase: authUsecase}

	route := app.Group("/api/auth")
	route.Post("/signin", hdlr.Login)
	route.Patch("/reset_password", middleware.ChangePassword(authUsecase, secret), hdlr.ChangePassword)

	route.Use(middleware.ValidateJWT(authUsecase, secret))

	route.Get("/", middleware.AdminAccess, hdlr.AllUsers)
	route.Get("/me", hdlr.GetOwn)
	route.Post("/signup", middleware.AdminAccess, hdlr.Create)
	route.Post("/:id", middleware.AdminAccess, hdlr.Enable)
	route.Delete("/:id", middleware.AdminAccess, hdlr.Disable)
}
