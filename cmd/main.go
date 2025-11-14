package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/database"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/mongodb"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
	"github.com/metalpoch/ultra-monitor/internal/validations"
	"github.com/metalpoch/ultra-monitor/routes"
	"github.com/metalpoch/ultra-monitor/usecase"
)

var db *sqlx.DB
var redis *cache.Redis
var prometheusClient prometheus.Prometheus
var mongoDB *mongodb.MongoDB
var jwtSecret string
var reportsDir string
var port string
var webAppDir string
var allowOrigin string
var enviroment string
var mongoURI string

func init() {
	enviroment = os.Getenv("ENVIROMENT")
	allowOrigin = "*"

	if enviroment == "production" {
		allowOrigin = os.Getenv("CORS_ALLOW_ORIGIN")
		if allowOrigin == "" {
			log.Fatal("error 'CORS_ALLOW_ORIGIN' enviroment varables requried.")
		}

		webAppDir = os.Getenv("WEB_APP_DIRECTORY")
		if webAppDir == "" {
			log.Fatal("error 'WEB_APP_DIRECTORY' enviroment varables requried.")
		}
	}

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
	if cacheURI == "" {
		log.Fatal("error 'REDIS_URI' enviroment varables requried.")
	}

	prometheusURL := os.Getenv("PROMETHEUS_URL")
	if prometheusURL == "" {
		log.Fatal("error 'PROMETHEUS_URL' enviroment varables requried.")
	}

	// MongoDB connection variables
	mongoURI = os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("error 'MONGODB_URI' enviroment varables requried.")
	}

	// Initialize MongoDB connection
	mongoDB = mongodb.NewMongoDB(mongoURI)

	db = database.Connect(dbURI)
	redis = cache.NewCache(cacheURI)
	prometheusClient = prometheus.NewPrometheusClient(prometheusURL)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Uso: <app> [ server | scan | traffic ]")
		return
	}

	switch os.Args[1] {
	case "server":

		app := fiber.New(fiber.Config{
			StructValidator: validations.NewValidator(),
			JSONEncoder:     json.Marshal,
			JSONDecoder:     json.Unmarshal,
			BodyLimit:       100 * 1024 * 1024, // 100 mb
			ServerHeader:    "KURBIL01",
			AppName:         "Gestor Ultra",
			StrictRouting:   true,
		})

		app.Use(logger.New())
		app.Use(limiter.New(limiter.Config{
			Max:        60,
			Expiration: 1 * time.Minute,
			// },
			LimitReached: func(c fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
				})
			},
		}))
		app.Use(cors.New(cors.Config{
			AllowOrigins: []string{allowOrigin},
			AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		}))

		routes.Init(&routes.Config{
			App:        app,
			DB:         db,
			Cache:      redis,
			Secret:     []byte(jwtSecret),
			Prometheus: &prometheusClient,
			WebAppDir:  webAppDir,
			ReportsDir: reportsDir,
			Enviroment: enviroment,
		})

		app.Listen("localhost"+":"+port, fiber.ListenConfig{
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
				Status: device.IfOperStatus,
			}); err != nil {
				log.Fatalf("error upserting batch: %v", err)
			}
		}

		log.Printf("Escaneo completado. Dispositivos totales: %d | ifName inválidos: %d", totalScanned, invalidFormat)
		log.Println("Total de dispositivos con shell en 0:", totalShellInZero)
		log.Println("Total de shells procesados:", totalShells)

	case "traffic":
		// Check if date parameter is provided
		if len(os.Args) < 3 {
			log.Fatal("Uso: <app> traffic YYYYMMDD (ejemplo: 20250930)")
		}

		dateParam := os.Args[2]

		// Parse date in YYYYMMDD format
		if len(dateParam) != 8 {
			log.Fatal("Formato de fecha inválido. Use YYYYMMDD (ejemplo: 20250930)")
		}

		year, err := strconv.Atoi(dateParam[0:4])
		if err != nil {
			log.Fatalf("Año inválido: %v", err)
		}

		month, err := strconv.Atoi(dateParam[4:6])
		if err != nil {
			log.Fatalf("Mes inválido: %v", err)
		}

		day, err := strconv.Atoi(dateParam[6:8])
		if err != nil {
			log.Fatalf("Día inválido: %v", err)
		}

		// Calculate date range from 00:00 to 23:59 of the specified date
		location := time.Now().Location()
		initDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
		finalDate := time.Date(year, time.Month(month), day, 23, 59, 59, 0, location)

		// Use the traffic usecase to update summary traffic
		trafficUsecase := usecase.NewTrafficUsecase(db, redis, &prometheusClient)
		if err := trafficUsecase.UpdateSummaryTraffic(initDate, finalDate); err != nil {
			log.Fatalf("Error updating summary traffic: %v", err)
		}

	case "traffic-ont":
		// Check if community parameter is provided
		if len(os.Args) < 3 {
			log.Fatal("Uso: <app> traffic-ont <community> (ejemplo: <app> traffic-ont public)")
		}

		community := os.Args[2]

		// Use the ont usecase to update traffic for all ONTs with individual intervals
		ontUsecase := usecase.NewOntUsecase(db, redis)

		// Create context with cancellation for graceful shutdown
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle SIGINT (Ctrl+C) for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		log.Println("Starting distributed ONT traffic monitoring with individual intervals...")
		log.Println("Press Ctrl+C to stop")

		// Get all enabled ONTs
		onts, err := ontUsecase.GetAll()
		if err != nil {
			log.Fatalf("Error getting ONTs: %v", err)
		}

		// Filter only enabled ONTs
		var enabledONTs []dto.OntResponse
		for _, ont := range onts {
			if ont.Enabled {
				enabledONTs = append(enabledONTs, ont)
			}
		}

		if len(enabledONTs) == 0 {
			log.Fatal("No enabled ONTs found")
		}

		log.Printf("Monitoring %d enabled ONTs with individual intervals", len(enabledONTs))

		// Each ONT will be queried every minute, but start times are staggered
		// to distribute load evenly throughout the minute
		intervalPerONT := 1 * time.Minute

		log.Printf("Each ONT will be queried every %v", intervalPerONT)

		// Create individual tickers for each ONT
		type ontTicker struct {
			ontID  int32
			ticker *time.Ticker
		}

		var ontTickers []ontTicker

		// Calculate staggered start times to distribute load
		staggerInterval := 1 * time.Minute / time.Duration(len(enabledONTs))
		staggerInterval = max(staggerInterval, 10*time.Second) // Minimum stagger

		log.Printf("Stagger interval between ONTs: %v", staggerInterval)

		// Start tickers with staggered intervals
		for i, ont := range enabledONTs {
			// Calculate initial delay to distribute start times
			initialDelay := time.Duration(i) * staggerInterval

			ticker := time.NewTicker(intervalPerONT)
			ontTickers = append(ontTickers, ontTicker{
				ontID:  ont.ID,
				ticker: ticker,
			})

			// Start goroutine for this ONT
			go func(ontID int32, delay time.Duration, ticker *time.Ticker) {
				// Initial delay to stagger starts
				time.Sleep(delay)

				// Run first execution immediately after delay
				log.Printf("Starting monitoring for ONT %d (first query)", ontID)
				if err := ontUsecase.UpdateTrafficForONT(ctx, ontID, community); err != nil {
					log.Printf("Error updating ONT %d: %v", ontID, err)
				}

				// Continue with 1-minute ticker
				for {
					select {
					case <-ticker.C:
						log.Printf("Querying ONT %d", ontID)
						if err := ontUsecase.UpdateTrafficForONT(ctx, ontID, community); err != nil {
							log.Printf("Error updating ONT %d: %v", ontID, err)
						}
					case <-ctx.Done():
						return
					}
				}
			}(ont.ID, initialDelay, ticker)
		}

		log.Println("All ONT monitors started successfully")

		// Refresh ticker to check for new ONTs periodically
		refreshTicker := time.NewTicker(5 * time.Minute)
		defer refreshTicker.Stop()

		// Track current ONTs to detect changes
		currentONTs := make(map[int32]bool)
		for _, ont := range enabledONTs {
			currentONTs[ont.ID] = true
		}

		// Function to refresh ONT list
		refreshONTs := func() {
			log.Println("Checking for ONT list changes...")

			// Get current ONTs from database
			currentONTsList, err := ontUsecase.GetAll()
			if err != nil {
				log.Printf("Error getting current ONTs: %v", err)
				return
			}

			// Filter enabled ONTs
			var newEnabledONTs []dto.OntResponse
			for _, ont := range currentONTsList {
				if ont.Enabled {
					newEnabledONTs = append(newEnabledONTs, ont)
				}
			}

			// Check for new ONTs
			newONTs := make(map[int32]dto.OntResponse)
			for _, ont := range newEnabledONTs {
				if !currentONTs[ont.ID] {
					newONTs[ont.ID] = ont
				}
			}

			// Check for removed ONTs
			removedONTs := make(map[int32]bool)
			for ontID := range currentONTs {
				found := false
				for _, ont := range newEnabledONTs {
					if ont.ID == ontID {
						found = true
						break
					}
				}
				if !found {
					removedONTs[ontID] = true
				}
			}

			// Add new ONTs
			if len(newONTs) > 0 {
				log.Printf("Found %d new ONTs, adding to monitoring...", len(newONTs))

				// Calculate stagger interval for new ONTs
				staggerInterval := 1 * time.Minute / time.Duration(len(newEnabledONTs))
				staggerInterval = max(staggerInterval, 10*time.Second) // Minimum stagger

				// Start monitoring for new ONTs
				for _, ont := range newONTs {
					// Calculate initial delay based on current time
					secondsInMinute := time.Now().Second()
					initialDelay := time.Duration(60-secondsInMinute) * time.Second

					ticker := time.NewTicker(intervalPerONT)
					ontTickers = append(ontTickers, ontTicker{
						ontID:  ont.ID,
						ticker: ticker,
					})

					// Start goroutine for new ONT
					go func(ontID int32, delay time.Duration, ticker *time.Ticker) {
						// Initial delay to align with minute boundary
						time.Sleep(delay)

						log.Printf("Starting monitoring for new ONT %d", ontID)
						if err := ontUsecase.UpdateTrafficForONT(ctx, ontID, community); err != nil {
							log.Printf("Error updating new ONT %d: %v", ontID, err)
						}

						// Continue with 1-minute ticker
						for {
							select {
							case <-ticker.C:
								if err := ontUsecase.UpdateTrafficForONT(ctx, ontID, community); err != nil {
									log.Printf("Error updating ONT %d: %v", ontID, err)
								}
							case <-ctx.Done():
								return
							}
						}
					}(ont.ID, initialDelay, ticker)

					// Add to current ONTs tracking
					currentONTs[ont.ID] = true
				}
			}

			// Remove ONTs that are no longer enabled
			if len(removedONTs) > 0 {
				log.Printf("Found %d ONTs to remove from monitoring...", len(removedONTs))

				// Stop tickers for removed ONTs
				var remainingTickers []ontTicker
				for _, ot := range ontTickers {
					if removedONTs[ot.ontID] {
						ot.ticker.Stop()
						log.Printf("Stopped monitoring for ONT %d", ot.ontID)
					} else {
						remainingTickers = append(remainingTickers, ot)
					}
				}
				ontTickers = remainingTickers

				// Remove from current ONTs tracking
				for ontID := range removedONTs {
					delete(currentONTs, ontID)
				}
			}

			if len(newONTs) == 0 && len(removedONTs) == 0 {
				log.Println("No changes in ONT list")
			}
		}

		// Run first refresh immediately
		go refreshONTs()

		// Main loop for refresh and signal handling
		for {
			select {
			case <-refreshTicker.C:
				refreshONTs()
			case <-sigChan:
				log.Println("Received interrupt signal, shutting down...")
				cancel()

				// Stop all tickers
				for _, ot := range ontTickers {
					ot.ticker.Stop()
				}

				log.Println("All ONT monitors stopped")
				return
			}
		}

case "interface-bandwidth":
		// Use the interface bandwidth usecase to update interface bandwidth data
		interfaceBandwidthUsecase := usecase.NewInterfaceBandwidthUsecase(db, mongoDB)
		if err := interfaceBandwidthUsecase.UpdateInterfaceBandwidth(context.Background()); err != nil {
			log.Fatalf("Error updating interface bandwidth: %v", err)
		}

		log.Printf("Interface bandwidth data updated successfully for date range: yesterday 00:00 to today 23:59")
	default:
		fmt.Println("Comando no reconocido. Usa 'server', 'scan', 'traffic', 'traffic-ont' o 'interface-bandwidth'")
	}
}
