package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/handler"
	"github.com/metalpoch/olt-blueprint/auth/middleware"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
)

func newUserRouter(server *fiber.App, db *sql.DB, secret []byte) {
	hdlr := &handler.UserHandler{
		Usecase: *usecase.NewUserUsecase(db, secret),
	}

	server.Post("/api/auth/signin", hdlr.Login)
	server.Get("/api/auth/user/profile", middleware.ValidateJWT(secret), hdlr.GetOwn)
	server.Patch("/api/auth/user/reset_password", middleware.ValidateJWT(secret), hdlr.ChangePassword)

	// Admin routes
	server.Post("/api/auth/signup", middleware.ValidateJWT(secret), middleware.AdminAccess, hdlr.Create)
	server.Get("/api/auth/user/all", middleware.ValidateJWT(secret), middleware.AdminAccess, hdlr.GetAll)
	server.Delete("/api/auth/user/:id", middleware.ValidateJWT(secret), middleware.AdminAccess, hdlr.DeleteUser)
}
