package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/metalpoch/olt-blueprint/update/controller"
	"github.com/metalpoch/olt-blueprint/update/database"
	"github.com/metalpoch/olt-blueprint/update/model"
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
	if len(os.Args) == 1 {
		fmt.Println("error: You need to specify an action [add-device|get-device|scan]")
		os.Exit(0)
	}

	args := os.Args[1:]
	action := args[0]

	switch action {
	case "add-device":
		if len(args) != 3 {
			fmt.Println("error: it is necessary to specify IP and community, in that same order")
			return
		}
		db := database.Connect(cfg.DatabaseURI)
		controller.AddDevice(args[1], args[2], db)

	case "scan":
		db := database.Connect(cfg.DatabaseURI)
		devices := controller.GetDevices(db)
		controller.TrafficUpdate(devices, db)
	}
}
