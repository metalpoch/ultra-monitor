package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/controller"
	"github.com/metalpoch/ultra-monitor/internal/constants"
	"github.com/metalpoch/ultra-monitor/internal/database"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
)

var db *sqlx.DB
var rdb *redis.Client
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
		log.Fatal("error 'POSTGRES_URI' enviroment varables requried.")
	}

	var err error
	db, err = database.Connect(dbURI)
	if err != nil {
		log.Fatal(err)
	}

	opt, err := redis.ParseURL(cacheURI)
	if err != nil {
		log.Fatal(err)
	}

	rdb = redis.NewClient(opt)
}

func main() {
	app := &cli.App{
		Name:        constants.CLI_TITLE,
		Description: constants.CLI_DESCRIPTION,
		Authors:     []*cli.Author{{Name: "Keiber Urbila", Email: "keiberup.dev@gmail.com"}},
		Commands: []*cli.Command{
			{
				Name:  "device",
				Usage: "options for device settings",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a new device",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "ip", Usage: "device IP", Required: true},
							&cli.StringFlag{Name: "community", Usage: "device community", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewOltController(db).
								AddOlt(dto.NewOlt{
									IP:        cCtx.String("ip"),
									Community: cCtx.String("community"),
								})
							return err
						},
					},
					{
						Name:  "update",
						Usage: "update a device by id",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "id", Usage: "olt id", Required: true},
							&cli.StringFlag{Name: "ip", Usage: "new ip"},
							&cli.StringFlag{Name: "community", Usage: "new community"},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewOltController(db).
								Update(cCtx.Int("id"), dto.NewOlt{
									IP:        cCtx.String("ip"),
									Community: cCtx.String("community"),
								})
							return err
						},
					},
					{
						Name:  "delete",
						Usage: "delete a device by id",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "id", Usage: "olt id", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewOltController(db).Delete(cCtx.Int("id"))
							return err
						},
					},
					{
						Name:  "show",
						Usage: "show existing devices",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "csv", Usage: "show as csv"},
						},
						Action: func(cCtx *cli.Context) error {
							_, err := controller.NewOltController(db).
								ShowAllOlt(cCtx.Bool("csv"))

							return err
						},
					},
				},
			},
			{
				Name:  "traffic",
				Usage: "get the traffic from the devices and store into the database",
				Subcommands: []*cli.Command{
					{
						Name:  "pon",
						Usage: "Olt pon traffic scan",
						Flags: []cli.Flag{
							&cli.Uint64Flag{Name: "interval", Usage: "interval time in minute", Value: 5},
						},
						Action: func(cCtx *cli.Context) error {
							interval := cCtx.Uint64("interval")
							if interval < 5 {
								return fmt.Errorf("inverval time invalid")
							}

							ctrl := controller.NewMeasurementController(db)

							started := make(map[string]struct{})
							var mu sync.Mutex

							for {
								olts, err := controller.NewOltController(db).Usecase.Olts(constants.DEFAULT_PAGE, constants.DEFAULT_LIMIT)
								if err != nil {
									log.Fatalln(err)
								}
								mu.Lock()
								for _, olt := range olts {
									if _, ok := started[olt.IP]; !ok {
										started[olt.IP] = struct{}{}
										go func(device model.Olt) {
											for {
												ctrl.PonScan(device)
												time.Sleep(time.Duration(interval) * time.Minute)
											}
										}(olt)
									}
								}
								mu.Unlock()
								time.Sleep(1 * time.Minute)
							}
						},
					},
					{
						Name:  "ont",
						Usage: "Ont traffic scan",
						Flags: []cli.Flag{
							&cli.Uint64Flag{Name: "interval", Usage: "interval time in minute", Value: 60},
						},
						Action: func(cCtx *cli.Context) error {
							interval := cCtx.Uint64("interval")
							if interval < 30 {
								return fmt.Errorf("inverval time invalid")
							}
							ctrl := controller.NewMeasurementController(db)

							started := make(map[string]struct{})
							var mu sync.Mutex

							const redisKey = "ont:snmp:responses"
							const batchSize = 100
							go func() {
								for {
									// check length
									len, err := rdb.LLen(cCtx.Context, redisKey).Result() // check len in redis
									if err != nil {
										fmt.Println("Redis LLEN error:", err)
										time.Sleep(10 * time.Second)
										continue
									}
									if len >= batchSize {
										// get first 1000 items
										data, err := rdb.LRange(cCtx.Context, redisKey, 0, batchSize-1).Result()
										if err != nil {
											fmt.Println("Redis LRANGE error:", err)
											continue
										}
										if err := rdb.LTrim(cCtx.Context, redisKey, batchSize, -1).Err(); err != nil {
											fmt.Println("Redis LTRIM error:", err)
											continue
										}
										if err := ctrl.ProcessOntBatchData(data); err != nil {
											fmt.Println("Batch insert error:", err)
										}
									} else {
										time.Sleep(5 * time.Second)
									}
								}
							}()

							for {
								olts, err := controller.NewOltController(db).Usecase.Olts(constants.DEFAULT_PAGE, constants.DEFAULT_LIMIT)
								if err != nil {
									log.Fatalln(err)
								}
								for _, olt := range olts {
									pons, err := ctrl.GetPonsBySysname(olt.SysName)
									if err != nil {
										fmt.Println(err)
										continue
									}
									mu.Lock()
									for _, pon := range pons {
										unique_id := fmt.Sprintf("%s-%d", olt.IP, pon.IfIndex)
										if _, ok := started[unique_id]; !ok {
											started[unique_id] = struct{}{}
											go func(device model.Olt, ponID int32, idx int64) {
												for {

													data, err := ctrl.OntScan(device, ponID, idx)
													if err != nil {
														log.Println(err.Error())
														continue
													}
													for _, record := range data {
														jsonData, err := json.Marshal(record)
														if err == nil {
															if err := rdb.RPush(cCtx.Context, redisKey, jsonData).Err(); err != nil {
																fmt.Println("Redis RPUSH error:", err)
															}
														}
													}
													time.Sleep(time.Duration(interval) * time.Minute)
												}
											}(olt, pon.ID, pon.IfIndex)
										}
									}
									mu.Unlock()
								}
							}
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
