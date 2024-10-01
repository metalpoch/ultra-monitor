package database

import (
	"log"

	authEntity "github.com/metalpoch/olt-blueprint/auth/entity"
	measurementEntity "github.com/metalpoch/olt-blueprint/measurement/entity"
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
		measurementEntity.Template{},
		measurementEntity.Device{},
		measurementEntity.Interface{},
		measurementEntity.Traffic{},
		measurementEntity.Measurement{},
	); err != nil {
		log.Fatalln(err)
	}

	return db
}
