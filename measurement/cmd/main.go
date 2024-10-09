package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/metalpoch/olt-blueprint/common/database"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/measurement/constants"
	"github.com/metalpoch/olt-blueprint/measurement/controller"
	"github.com/urfave/cli/v2"
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
	telegram := tracking.Telegram{
		BotID:  cfg.TelegramBotTokenID,
		ChatID: cfg.TelegramChatID,
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
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							err := controller.AddTemplate(db, telegram, &model.AddTemplate{
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
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							err := controller.UpdateTemplate(db, telegram, cCtx.Uint("id"), &model.AddTemplate{
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
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							err := controller.DeleteTemplate(db, telegram, cCtx.Uint("id"))
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
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							if err := controller.ShowAllTemplates(db, telegram, cCtx.Bool("csv")); err != nil {
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
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							if err := controller.AddDevice(db, telegram, &model.AddDevice{
								IP:        cCtx.String("ip"),
								Community: cCtx.String("community"),
								Template:  cCtx.Uint("template"),
							}); err != nil {
								log.Fatal(err)
							} else {
								fmt.Println("saved.")
							}
							return nil
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
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							err := controller.UpdateDevice(db, telegram, cCtx.Uint("id"), &model.AddDevice{
								IP:        cCtx.String("ip"),
								Community: cCtx.String("community"),
								Template:  cCtx.Uint("template-id"),
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
						Usage: "delete a device by id",
						Flags: []cli.Flag{
							&cli.UintFlag{Name: "id", Usage: "device id", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							err := controller.DeleteDevice(db, telegram, cCtx.Uint("id"))
							if err != nil {
								log.Fatal(err)
							}
							fmt.Println("deleted")
							return nil
						},
					},
					{
						Name:  "show",
						Usage: "show existing devices",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "csv", Usage: "show as csv"},
						},
						Action: func(cCtx *cli.Context) error {
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							if _, err := controller.ShowAllDevices(db, telegram, cCtx.Bool("csv")); err != nil {
								log.Fatal(err)
							}

							return nil
						},
					},
					{
						Name:  "show-interfaces",
						Usage: "show interface existing in devices",
						Flags: []cli.Flag{
							&cli.UintFlag{Name: "device", Usage: "device id", Required: true},
							&cli.BoolFlag{Name: "csv", Usage: "show as csv"},
						},
						Action: func(cCtx *cli.Context) error {
							db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
							if err := controller.ShowAllInterfaces(db, telegram, cCtx.Uint("device"), cCtx.Bool("csv")); err != nil {
								log.Fatal(err)
							}

							return nil
						},
					},
				},
			},
			{
				Name:  "traffic",
				Usage: "get the traffic from the devices and store into the database",
				Action: func(cCtx *cli.Context) error {
					for {
						db := database.Connect(cfg.DatabaseURI, cfg.IsProduction)
						sqlDB, err := db.DB()
						if err != nil {
							log.Fatal(err)
						}
						devices, err := controller.GetDeviceWithOIDRows(db, telegram)
						if err != nil {
							log.Fatalln(err)
						}

						go controller.Scan(db, telegram, devices)

						time.Sleep(5 * time.Minute)
						sqlDB.Close()
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
