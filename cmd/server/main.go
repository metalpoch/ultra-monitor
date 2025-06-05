package main

import (
	"log"
	"os"
	"strconv"

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
var reportsDir string
var port string

func init() {
	reportsDir = os.Getenv("REPORTS_DIRECTORY")
	if reportsDir == "" {
		log.Fatal("error 'REPORTS_DIRECTORY' enviroment varables requried.")
	}

	if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(reportsDir, 0755); err != nil {
			log.Fatalf("failed to create reports directory: %v", err)
		}
	}

	dbURI := os.Getenv("POSTGRES_URI")
	if dbURI == "" {
		log.Fatal("error 'POSTGRES_URI' enviroment varables requried.")
	}

	jwtSecret = os.Getenv("AUTH_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("error 'JWT_SECRET' enviroment varables requried.")
	}

	port = os.Getenv("PORT")
	if port == "" {
		log.Fatal("error 'PORT' environment variable required.")
	}
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("error 'PORT' must be a valid number: %v", err)
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
		ServerHeader:    "KURBIL01",
		AppName:         "Ultra Monitor",
		CaseSensitive:   true,
		StrictRouting:   true,
	})

	app.Use(logger.New())
	app.Use(cors.New())

	routes.Init(app, db, redis, []byte(jwtSecret), reportsDir)

	app.Listen(":"+port, fiber.ListenConfig{
		EnablePrefork: true,
	})
}
