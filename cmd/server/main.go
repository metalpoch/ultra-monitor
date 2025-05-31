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
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/database"
	"github.com/metalpoch/ultra-monitor/internal/validations"
	"github.com/metalpoch/ultra-monitor/routes"
)

var db *sqlx.DB
var redis *cache.Redis
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

	cacheURI := os.Getenv("REDIS_URI")
	if dbURI == "" {
		log.Fatal("error 'REDIS_URI' enviroment varables requried.")
	}

	var err error
	db, err = database.Connect(dbURI)
	if err != nil {
		log.Fatal(err)
	}

	redis = cache.NewCache(cacheURI)
}

func main() {
	app := fiber.New(fiber.Config{
		StructValidator: &validations.StructValidator{Validator: validator.New()},
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
		BodyLimit:       100 * 1024 * 1024, // 100 mb
	})

	app.Use(logger.New())
	app.Use(cors.New())

	routes.Init(app, db, redis, []byte(jwtSecret), "./data")

	app.Listen(":3000")
}
