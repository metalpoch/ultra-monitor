package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(uri string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		log.Fatal(err)
	}

	if err := createTablesAndIndexes(db); err != nil {
		log.Fatal(err)
	}

	return db
}

func createTablesAndIndexes(db *sqlx.DB) error {
	tableQueries := []string{
		SQL_TABLE_USERS,
		SQL_TABLE_REPORT,
		SQL_TABLE_FAT,
		SQL_TABLE_FAT_STATUS,
		SQL_TABLE_PROMETHEUS_DEVICES,
		SQL_TABLE_SUMMARY_TRAFFIC,
		SQL_TABLE_ONT,
		SQL_TABLE_ONT_TRAFFIC,
	}
	indexQueries := []string{
		SQL_INDEX_REPORT_CATEGORY,
		SQL_INDEX_REPORT_USER_ID,
		SQL_INDEX_USERS_USERNAME,
	}

	for _, q := range tableQueries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("Error creating table: %v\nQuery: %s", err, q)
			return err
		}
	}
	for _, q := range indexQueries {
		if _, err := db.Exec(q); err != nil {
			log.Printf("Error creating index: %v\nQuery: %s", err, q)
			return err
		}
	}
	return nil
}
