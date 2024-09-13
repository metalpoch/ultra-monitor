package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(uri string) *sql.DB {
	if uri == "" {
		log.Fatal("Set your 'db_uri' on config.json")
	}

	client, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}

	return client
}
