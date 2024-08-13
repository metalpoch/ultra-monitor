package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/handler"
	"github.com/metalpoch/olt-blueprint/auth/repository"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func newExampleRouter(server *fiber.App, client *mongo.Client) {
	hdlr := &handler.ExampleHandler{
		Usecase: usecase.NewExampleUsecase(
			repository.NewExampleRepository(client),
		),
	}
	server.Get("/example/:id", hdlr.Get)
	server.Post("/example", hdlr.Create)

}
