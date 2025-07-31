package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/database"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
	"github.com/metalpoch/ultra-monitor/internal/validations"
	"github.com/metalpoch/ultra-monitor/routes"
	"github.com/metalpoch/ultra-monitor/usecase"
)

var db *sqlx.DB
var redis *cache.Redis
var prometheusClient prometheus.Prometheus
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

	prometheusURL := os.Getenv("PROMETHEUS_URL")
	if prometheusURL == "" {
		log.Fatal("error 'PROMETHEUS_URL' enviroment varables requried.")
	}

	db = database.Connect(dbURI)
	redis = cache.NewCache(cacheURI)
	prometheusClient = prometheus.NewPrometheusClient(prometheusURL)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Uso: <app> [server | scan ]")
		return
	}

	switch os.Args[1] {
	case "server":

		app := fiber.New(fiber.Config{
			StructValidator: &validations.StructValidator{Validator: validator.New()},
			JSONEncoder:     json.Marshal,
			JSONDecoder:     json.Unmarshal,
			BodyLimit:       100 * 1024 * 1024, // 100 mb
			ServerHeader:    "KURBIL01",
			AppName:         "Gestor Ultra",
			CaseSensitive:   true,
			StrictRouting:   true,
		})

		app.Use(logger.New())
		app.Use(cors.New())

		routes.Init(app, db, redis, []byte(jwtSecret), reportsDir, &prometheusClient)

		app.Listen(":"+port, fiber.ListenConfig{
			EnablePrefork: true,
		})

	case "scan":
		prometheusDevices, err := prometheusClient.DeviceScan(context.Background())
		if err != nil {
			log.Fatalf("error scanning Prometheus devices: %v", err)
		}

		prometheusUsecase := usecase.NewPrometheusUsecase(db)
		invalidFormat := 0
		totalScanned := 0
		totalShellInZero := 0
		totalShells := 0

		re := regexp.MustCompile(`GPON (\d+)/(\d+)/(\d+)`)

		for _, device := range prometheusDevices {
			totalScanned++
			matches := re.FindStringSubmatch(device.IfName)
			if len(matches) != 4 {
				log.Printf("ifName no válido (esperado GPON x/y/z) en %s: %s", device.IP, device.IfName)
				invalidFormat++
				continue
			}

			shell, _ := strconv.Atoi(matches[1])
			card, _ := strconv.Atoi(matches[2])
			port, _ := strconv.Atoi(matches[3])
			if shell == 0 {
				totalShellInZero++
			}
			totalShells++
			if err := prometheusUsecase.Upsert(context.Background(), dto.Prometheus{
				Region: device.Region,
				State:  device.State,
				IP:     device.IP,
				IDX:    device.IfIndex,
				Shell:  uint8(shell),
				Card:   uint8(card),
				Port:   uint8(port),
			}); err != nil {
				log.Fatalf("error upserting batch: %v", err)
			}
		}

		log.Printf("Escaneo completado. Dispositivos totales: %d | ifName inválidos: %d", totalScanned, invalidFormat)
		log.Println("Total de dispositivos con shell en 0:", totalShellInZero)
		log.Println("Total de shells procesados:", totalShells)

	default:
		fmt.Println("Comando no reconocido. Usa 'server' o 'scan'")

	}
}
