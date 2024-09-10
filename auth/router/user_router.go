package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/handler"
	"github.com/metalpoch/olt-blueprint/auth/repository"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func newUserRouter(server *fiber.App, client *mongo.Client) {

	hdlr := &handler.UserHandler{
		Usecase: usecase.NewUserUsecase(
			repository.NewUserRepository(client),
		),
	}
	server.Get("/user/Get", hdlr.Get)
	server.Post("/user/Post", hdlr.Create)
	server.Get("/user/Getvalue/:clave/:valor", hdlr.GetValue)

	server.Delete("/user/DeleteName/:p00", hdlr.DeleteName)

	server.Patch("/user/ChangeP", hdlr.ChangePassword)
	server.Post("/user/Login", hdlr.Login)
	server.Post("/user/ReadToken", hdlr.ReadToken)

}
