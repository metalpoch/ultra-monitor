package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/metalpoch/olt-blueprint/common/database"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/measurement/controller"
	"github.com/metalpoch/olt-blueprint/measurement/internal/constants"
	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

const maxConcurrent = 20

var (
	cfg model.Config
	db  *gorm.DB
	rdb *redis.Client
)

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

	configCache, err := redis.ParseURL(cfg.CacheURI)
	if err != nil {
		log.Panicln(err)
	}

	db = database.Connect(cfg.DatabaseURI, cfg.IsProduction)
	rdb = redis.NewClient(configCache)
}

func main() {
	telegram := tracking.SmartModule{
		URL: cfg.SmartModuleTelegramURL,
	}

	app := &cli.App{
		Name:        constants.CLI_TITLE,
		Description: constants.CLI_DESCRIPTION,
		Authors:     []*cli.Author{{Name: "Keiber Urbila", Email: "keiberup.dev@gmail.com"}},
		Commands: []*cli.Command{
			{
				Name:  "template",
				Usage: "options for device templates",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a new template",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "name", Usage: "template name", Required: true},
							&cli.StringFlag{Name: "oid-bw", Usage: "bandwidth oid", Required: true},
							&cli.StringFlag{Name: "oid-in", Usage: "traffic in oid", Required: true},
							&cli.StringFlag{Name: "oid-out", Usage: "traffic out oid", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							ctrl := controller.NewTemplateController(db, telegram)
							err := ctrl.AddTemplate(&model.AddTemplate{
								Name:   cCtx.String("name"),
								OidBw:  cCtx.String("oid-bw"),
								OidIn:  cCtx.String("oid-in"),
								OidOut: cCtx.String("oid-out"),
							})
							if err != nil {
								log.Fatal(err)
							}
							fmt.Println("saved")
							return nil
						},
					},
					{
						Name:  "update",
						Usage: "update a template by id",
						Flags: []cli.Flag{
							&cli.UintFlag{Name: "id", Usage: "template id", Required: true},
							&cli.StringFlag{Name: "name", Usage: "new name of the template"},
							&cli.StringFlag{Name: "oid-bw", Usage: "new bandwidth oid"},
							&cli.StringFlag{Name: "oid-in", Usage: "new traffic in oid"},
							&cli.StringFlag{Name: "oid-out", Usage: "new traffic out oid"},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewTemplateController(db, telegram).
								UpdateTemplate(cCtx.Uint("id"), &model.AddTemplate{
									Name:   cCtx.String("name"),
									OidBw:  cCtx.String("oid-bw"),
									OidIn:  cCtx.String("oid-in"),
									OidOut: cCtx.String("oid-out"),
								})
							if err != nil {
								log.Fatal(err)
							}
							fmt.Println("updated")
							return nil
						},
					},
					{
						Name:  "delete",
						Usage: "delete a template by id",
						Flags: []cli.Flag{
							&cli.UintFlag{Name: "id", Usage: "template id", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewTemplateController(db, telegram).DeleteTemplate(cCtx.Uint("id"))
							if err != nil {
								log.Fatal(err)
							}
							fmt.Println("deleted")
							return nil
						},
					},
					{
						Name:  "show",
						Usage: "show existing template",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "csv", Usage: "show as csv"},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewTemplateController(db, telegram).ShowAllTemplates(cCtx.Bool("csv"))
							if err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "device",
				Usage: "options for device settings",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a new device",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "template", Usage: "template id", Required: true},
							&cli.StringFlag{Name: "ip", Usage: "device IP", Required: true},
							&cli.StringFlag{Name: "community", Usage: "device community", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewDeviceController(db, telegram).
								AddDevice(&model.AddDevice{
									IP:        cCtx.String("ip"),
									Community: cCtx.String("community"),
									Template:  cCtx.Uint("template"),
								})
							return err
						},
					},
					{
						Name:  "update",
						Usage: "update a device by id",
						Flags: []cli.Flag{
							&cli.UintFlag{Name: "id", Usage: "device id", Required: true},
							&cli.StringFlag{Name: "ip", Usage: "new ip"},
							&cli.StringFlag{Name: "community", Usage: "new community"},
							&cli.UintFlag{Name: "template-id", Usage: "new template id"},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewDeviceController(db, telegram).
								UpdateDevice(cCtx.Uint("id"), &model.AddDevice{
									IP:        cCtx.String("ip"),
									Community: cCtx.String("community"),
									Template:  cCtx.Uint("template-id"),
								})
							return err
						},
					},
					{
						Name:  "delete",
						Usage: "delete a device by id",
						Flags: []cli.Flag{
							&cli.Uint64Flag{Name: "id", Usage: "device id", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewDeviceController(db, telegram).DeleteDevice(cCtx.Uint64("id"))
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
							_, err := controller.NewDeviceController(db, telegram).
								ShowAllDevices(db, telegram, cCtx.Bool("csv"))

							return err
						},
					},
					{
						Name:  "show-interfaces",
						Usage: "show interface existing in devices",
						Flags: []cli.Flag{
							&cli.Uint64Flag{Name: "device", Usage: "device id", Required: true},
							&cli.BoolFlag{Name: "csv", Usage: "show as csv"},
						},
						Action: func(cCtx *cli.Context) error {
							err := controller.NewInterfaceController(db, telegram).
								ShowAllInterfaces(cCtx.Uint64("device"), cCtx.Bool("csv"))
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
						Name:  "olt",
						Usage: "Olt ports traffic scan",
						Flags: []cli.Flag{
							&cli.Uint64Flag{Name: "interval", Usage: "interval time in minute", Value: 5},
						},
						Action: func(cCtx *cli.Context) error {
							interval := cCtx.Uint64("interval")
							if interval < 5 {
								return fmt.Errorf("inverval time invalid")
							}

							ctrl := controller.NewMeasurementController(db, telegram)

							started := make(map[uint64]struct{})
							var mu sync.Mutex

							for {
								devices, err := ctrl.Device.GetDeviceWithOIDRows()
								if err != nil {
									log.Fatalln(err)
								}
								mu.Lock()
								for _, device := range devices {
									if _, ok := started[device.ID]; !ok {
										started[device.ID] = struct{}{}
										go func(device *model.DeviceWithOID) {
											for {
												log.Println("INIT")
												ctrl.OltScan(device)
												log.Println("END")
												time.Sleep(time.Duration(interval) * time.Minute)
											}
										}(device)
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
							ctrl := controller.NewMeasurementController(db, telegram)

							started := make(map[string]struct{})
							var mu sync.Mutex

							const redisKey = "ont:snmp:responses"
							const batchSize = 100
							go func() {
								for {
									// check length
									len, err := rdb.LLen(cCtx.Context, redisKey).Result() // check len in redis
									if err != nil {
										log.Println("Redis LLEN error:", err)
										time.Sleep(10 * time.Second)
										continue
									}
									if len >= batchSize {
										// get first 1000 items
										data, err := rdb.LRange(cCtx.Context, redisKey, 0, batchSize-1).Result()
										if err != nil {
											log.Println("Redis LRANGE error:", err)
											continue
										}
										if err := rdb.LTrim(cCtx.Context, redisKey, batchSize, -1).Err(); err != nil {
											log.Println("Redis LTRIM error:", err)
											continue
										}
										if err := ctrl.ProcessOntBatchData(data); err != nil {
											log.Println("Batch insert error:", err)
										}
									} else {
										time.Sleep(5 * time.Second)
									}
								}
							}()

							for {
								devices, err := ctrl.Device.GetDeviceWithOIDRows()
								if err != nil {
									log.Fatalln(err)
								}
								for _, device := range devices {
									pons, err := ctrl.Interface.GetAllByDevice(device.ID)
									if err != nil {
										log.Fatalln(err)
										continue
									}
									mu.Lock()
									for _, pon := range pons {
										unique_id := fmt.Sprintf("%d-%d", device.ID, pon.IfIndex)
										if _, ok := started[unique_id]; !ok {
											started[unique_id] = struct{}{}
											go func(device *model.DeviceWithOID, ponID, idx uint64) {
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
																log.Println("Redis RPUSH error:", err)
															}
														}
													}
													time.Sleep(time.Duration(interval) * time.Minute)
												}
											}(device, pon.ID, pon.IfIndex)
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
