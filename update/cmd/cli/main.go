package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/metalpoch/olt-blueprint/common/database"
	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/update/controller"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/urfave/cli/v2"
)

var cfg commonModel.Config

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
	app := &cli.App{
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
							&cli.StringFlag{Name: "prefix-bw", Usage: "bandwidth prefix", Value: "octe"},
							&cli.StringFlag{Name: "prefix-in", Usage: "traffic in SI prefix", Value: "octe"},
							&cli.StringFlag{Name: "prefix-out", Usage: "traffic out SI prefix", Value: "octe"},
						},
						Action: func(cCtx *cli.Context) error {
							db := database.Connect(cfg.DatabaseURI)
							err := controller.AddTemplate(db, &model.AddTemplate{
								Name:      cCtx.String("name"),
								OidBw:     cCtx.String("oid-bw"),
								OidIn:     cCtx.String("oid-in"),
								OidOut:    cCtx.String("oid-out"),
								PrefixBw:  cCtx.String("prefix-bw"),
								PrefixIn:  cCtx.String("prefix-in"),
								PrefixOut: cCtx.String("prefix-out"),
							})
							if err != nil {
								log.Fatal(err)
							}
							fmt.Println("saved.")

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
							db := database.Connect(cfg.DatabaseURI)
							if err := controller.ShowAllTemplates(db, cCtx.Bool("csv")); err != nil {
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
							db := database.Connect(cfg.DatabaseURI)
							if err := controller.AddDevice(db, &model.AddDevice{
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
						Name:  "show",
						Usage: "show existing devices",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "csv", Usage: "show as csv"},
						},
						Action: func(cCtx *cli.Context) error {
							db := database.Connect(cfg.DatabaseURI)
							if _, err := controller.ShowAllDevices(db, cCtx.Bool("csv")); err != nil {
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
							db := database.Connect(cfg.DatabaseURI)
							if err := controller.ShowAllInterfaces(db, cCtx.Uint("device"), cCtx.Bool("csv")); err != nil {
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
						db := database.Connect(cfg.DatabaseURI)
						devices, err := controller.GetDeviceWithOIDRows(db)
						if err != nil {
							log.Fatalln(err)
						}

						go controller.Scan(db, devices)

						time.Sleep(20 * time.Second)
						// time.Sleep(5 * time.Minute)
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
