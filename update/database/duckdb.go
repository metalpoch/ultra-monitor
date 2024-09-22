package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/metalpoch/olt-blueprint/update/constants"
)

func initDB(db *sql.DB, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch name {
	case constants.DATABASE_DEVICE:
		querys := [3]string{constants.SQL_DEVICE_SEQUENCES, constants.SQL_CREATE_TABLE_TEMPLATE, constants.SQL_CREATE_TABLE_DEVICE}
		for _, q := range querys {
			if _, err := db.ExecContext(ctx, q); err != nil {
				fmt.Println(q, err.Error())
				return err
			}
		}

	default:
		tables := [3]string{constants.SQL_CREATE_TABLE_INTERFACE, constants.SQL_CREATE_TABLE_MEASUREMENT, constants.SQL_CREATE_TABLE_TMP_MEASUREMENT}
		for _, table := range tables {
			if _, err := db.ExecContext(ctx, table); err != nil {
				return err
			}
		}
	}

	return nil
}

func Duckdb(database string) (*sql.DB, error) {
	db, err := sql.Open("duckdb", fmt.Sprintf("./data/%s.db", database))
	if err != nil {
		return nil, err
	}

	initDB(db, database)

	return db, nil
}
