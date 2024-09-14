package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func createTableIfNotExist(db *sql.DB) {
	usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id 					serial NOT NULL,
        fullname 			text NOT NULL,
        email 				text NOT NULL UNIQUE,
        password			text NOT NULL,
		change_password 	boolean DEFAULT true,
		is_admin			boolean DEFAULT false,
		is_disabled			boolean DEFAULT false,
        created_at			timestamp DEFAULT now(),
        updated_at			timestamp DEFAULT now(),
        CONSTRAINT pk_users PRIMARY KEY(id)
    );`

	_, err := db.Exec(usersTable)
	if err != nil {
		log.Fatalf("error creating users table: %v\n", err)
	}

}

func Connect(uri string) *sql.DB {
	if uri == "" {
		log.Fatal("Set your 'db_uri' on config.json")
	}

	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}

	createTableIfNotExist(db)

	return db
}
