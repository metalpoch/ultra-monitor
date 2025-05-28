package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/handler"
	"github.com/metalpoch/olt-blueprint/middleware"
	"github.com/metalpoch/olt-blueprint/usecase"
)

func newUserRoutes(server *fiber.App, db *sqlx.DB, secret []byte) {
	hdlr := &handler.UserHandler{
		Usecase: *usecase.NewUserUsecase(db, secret),
	}

	server.Post("/login", hdlr.Login)
	server.Get("/user/profile", hdlr.GetOwn, middleware.ValidateJWT(secret))
	server.Patch("/user/reset_password", hdlr.ChangePassword, middleware.ValidateJWT(secret))

	// for develop
	server.Post("/signup", hdlr.Create)

	// Admin routes
	// server.Post("/api/auth/signup", hdlr.Create, middleware.ValidateJWT(secret), middleware.AdminAccess)
	// server.Get("/user/all", hdlr.GetAll, middleware.ValidateJWT(secret), middleware.AdminAccess)
	// server.Delete("/user/:id", hdlr.DeleteUser, middleware.ValidateJWT(secret), middleware.AdminAccess)
}
