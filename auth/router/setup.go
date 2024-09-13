package router

import (
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(server *fiber.App, client *mongo.Client) {
	newUserRouter(server, client)
}
