package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/router"
	"github.com/metalpoch/olt-blueprint/common/database"
	"github.com/metalpoch/olt-blueprint/common/model"
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
	db := database.Connect(cfg.DatabaseURI)
	server := fiber.New()

	router.Setup(server, db, []byte(cfg.SecretKey))

	server.Listen(":3000")
}
