package database

import (
	"log"

	authEntity "github.com/metalpoch/olt-blueprint/auth/entity"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(uri string, isProduction bool) *gorm.DB {
	var typeLogger logger.Interface
	if uri == "" {
		log.Fatalln("Set your 'db_uri' on config.json")
	}

	if isProduction {
		typeLogger = logger.Default.LogMode(logger.Silent)
	} else {
		typeLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{Logger: typeLogger})
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.AutoMigrate(
		authEntity.User{},
		entity.Template{},
		entity.Device{},
		entity.Interface{},
		entity.Traffic{},
		entity.Measurement{},
	); err != nil {
		log.Fatalln(err)
	}

	return db
}
