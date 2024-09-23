package database

import (
	"fmt"
	"log"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDB() {
	sc := constants.SC
	dbN := constants.DB

	DB, _ := gorm.Open(postgres.Open(sc), &gorm.Config{})

	createDataBaseCommand := fmt.Sprintf("Create DATABASE %s", dbN)
	DB.Exec(createDataBaseCommand)

}

func DBConnection(uri string) *gorm.DB {
	if uri == "" {
		log.Fatal("Set your 'db_uri' on config.json")
	}
	DB, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return DB
}
