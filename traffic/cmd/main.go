package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/goccy/go-json"
	"github.com/metalpoch/olt-blueprint/common/database"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/traffic/router"
)

var cfg model.Config

func init() {
	fileJSON := os.Getenv("CONFIG_JSON")
	if fileJSON == "" {
		log.Panicln("CONFIG_JSON env is required")
	}

	f, err := os.ReadFile(fileJSON)
	if err != nil {
		log.Panicln(err)
	}

	json.Unmarshal([]byte(f), &cfg)
}

func main() {
	db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
	server := fiber.New(fiber.Config{
		StructValidator: &model.StructValidator{Validator: validator.New()},
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
	})

	server.Use(logger.New())

	router.Setup(server, db, []byte(cfg.SecretKey))

	server.Listen(":3000")

}
