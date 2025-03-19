package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/handler"
	"github.com/metalpoch/olt-blueprint/auth/middleware"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"gorm.io/gorm"
)

func newUserRouter(server *fiber.App, db *gorm.DB, secret []byte, telegram tracking.SmartModule) {
	hdlr := &handler.UserHandler{
		Usecase: *usecase.NewUserUsecase(db, secret, telegram),
	}

	server.Post("/api/auth/login", hdlr.Login)
	server.Get("/api/auth/user/profile", hdlr.GetOwn, middleware.ValidateJWT(secret))
	server.Patch("/api/auth/user/reset_password", hdlr.ChangePassword, middleware.ValidateJWT(secret))

	// for develop
	server.Post("/api/auth/signup", hdlr.Create)

	// Admin routes
	// server.Post("/api/auth/signup", hdlr.Create, middleware.ValidateJWT(secret), middleware.AdminAccess)
	server.Get("/api/auth/user/all", hdlr.GetAll, middleware.ValidateJWT(secret), middleware.AdminAccess)
	server.Delete("/api/auth/user/:id", hdlr.DeleteUser, middleware.ValidateJWT(secret), middleware.AdminAccess)
}
