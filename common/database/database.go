package database

import (
	"log"

	authEntity "github.com/metalpoch/olt-blueprint/auth/entity"
	updateEntity "github.com/metalpoch/olt-blueprint/update/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(uri string) *gorm.DB {
	if uri == "" {
		log.Fatalln("Set your 'db_uri' on config.json")
	}
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)

	}
	if err := db.AutoMigrate(
		authEntity.User{},
		updateEntity.Template{},
		updateEntity.Device{},
		updateEntity.Interface{},
		updateEntity.Traffic{},
		updateEntity.Measurement{},
	); err != nil {
		log.Fatalln(err)
	}

	return db
}
