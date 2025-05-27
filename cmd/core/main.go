package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/goccy/go-json"
	"github.com/metalpoch/olt-blueprint/common/database"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/core/router"
	"github.com/metalpoch/olt-blueprint/internal/config"
)

var cfg config.Config

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

	// db, err := sqlx.Connect("postgres", cfg.DatabaseURI)
	// if err != nil {
	// 	log.Panicln("Failed to connect to the database:", err)
	// }
	// if err := db.Ping(); err != nil {
	// 	log.Panicln("Failed to ping the database:", err)
	// }

	// 	queries := [7]string{
	// 	constants.CREATE_OLT_TABLE,
	// 	constants.CREATE_PON_TABLE,
	// 	constants.CREATE_PON_TRAFFIC_TABLE,
	// 	constants.CREATE_PON_TRAFFIC_TMP_TABLE,
	// 	constants.CREATE_INDEX_TRAFFIC_PONS_DATE,
	// 	constants.CREATE_INDEX_TMP_TRAFFIC_PONS_ID,
	// 	constants.CREATE_INDEX_TRAFFIC_PONS_ID_DATE,
	// }

	// for _, query := range queries {
	// 	if _, err := base_db.Exec(query); err != nil {
	// 		log.Printf("sql error `%s` on trying execute: %s", err.Error(), query)
	// 	}
	// }
}

func main() {
	telegram := tracking.SmartModule{
		URL: cfg.SmartModuleTelegramURL,
	}

	db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
	server := fiber.New(fiber.Config{
		StructValidator: &model.StructValidator{Validator: validator.New()},
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
	})
	cache := cache.NewCache(cfg.CacheURI)
	server.Use(logger.New())
	server.Use(cors.New())

	router.Setup(server, db, telegram, *cache)

	server.Listen(":3001")
}
