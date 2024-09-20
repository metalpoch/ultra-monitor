package main

import (
	"fmt"
	"log"
	"os"

	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/controller"
	"github.com/metalpoch/olt-blueprint/update/database"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/urfave/cli/v2"
)

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
							db, err := database.Duckdb(constants.DATABASE_DEVICE)
							if err != nil {
								log.Fatalln(err)
							}
							defer db.Close()

							err = controller.AddTemplate(db, &model.AddTemplate{
								Name:      cCtx.String("name"),
								OidBw:     cCtx.String("oid-bw"),
								OidIn:     cCtx.String("oid-in"),
								OidOut:    cCtx.String("oid-out"),
								PrefixBw:  cCtx.String("prefix-bw"),
								PrefixIn:  cCtx.String("prefix-in"),
								PrefixOut: cCtx.String("prefix-out"),
							})
							if err != nil {
								log.Fatalln(err)
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
							db, err := database.Duckdb(constants.DATABASE_DEVICE)
							if err != nil {
								log.Fatalln(err)
							}
							defer db.Close()

							if err := controller.ShowAllTemplates(db, cCtx.Bool("csv")); err != nil {
								log.Fatalln(err)
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
							&cli.UintFlag{Name: "template", Usage: "template id", Required: true},
							&cli.StringFlag{Name: "ip", Usage: "device IP", Required: true},
							&cli.StringFlag{Name: "community", Usage: "device community", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							db, err := database.Duckdb(constants.DATABASE_DEVICE)
							if err != nil {
								log.Fatalln(err)
							}
							defer db.Close()

							if err := controller.AddDevice(db, &model.AddDevice{
								IP:         cCtx.String("ip"),
								Community:  cCtx.String("community"),
								TemplateID: cCtx.Uint("template"),
							}); err != nil {
								log.Fatalln(err)
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
							db, err := database.Duckdb(constants.DATABASE_DEVICE)
							if err != nil {
								log.Fatalln(err)
							}
							defer db.Close()

							if err := controller.ShowAllDevices(db, cCtx.Bool("csv")); err != nil {
								log.Fatalln(err)
							}

							return nil
						},
					},
				},
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "complete a task on the list",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
