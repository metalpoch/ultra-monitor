package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/config"
	"github.com/metalpoch/ultra-monitor/internal/database"
	"github.com/metalpoch/ultra-monitor/routes"
)

var db *sqlx.DB
var jwtSecret string

func init() {
	dbURI := os.Getenv("POSTGRES_URI")
	if dbURI == "" {
		log.Fatal("error 'POSTGRES_URI' enviroment varables requried.")
	}

	jwtSecret = os.Getenv("AUTH_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("error 'JWT_SECRET' enviroment varables requried.")
	}

	var err error
	db, err = database.Connect(dbURI)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &config.StructValidator{Validator: validator.New()},
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
	})

	app.Use(logger.New())
	app.Use(cors.New())
	routes.Users(app, db, []byte(jwtSecret))
	app.Listen(":3001")
}
