package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

func Connect(uri string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return nil, err
	}

	if err := createTablesAndIndexes(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTablesAndIndexes(db *sqlx.DB) error {
	tableQueries := []string{
		constants.SQL_TABLE_USERS,
		constants.SQL_TABLE_REPORT,
		constants.SQL_TABLE_OLT,
		constants.SQL_TABLE_PON,
		constants.SQL_TABLE_MEASUREMENT_PON,
		constants.SQL_TABLE_TRAFFIC_PON,
		constants.SQL_TABLE_MEASUREMENT_ONT,
		constants.SQL_TABLE_FAT,
	}
	indexQueries := []string{
		constants.SQL_INDEX_FAT_STATE,
		constants.SQL_INDEX_FAT_COUNTY,
		constants.SQL_INDEX_FAT_MUNICIPALITY,
		constants.SQL_INDEX_FAT_ODN,
		constants.SQL_INDEX_FAT_OLT,
		constants.SQL_INDEX_MEASUREMENT_PON_DATE,
		constants.SQL_INDEX_MEASUREMENT_PON_PON_ID,
		constants.SQL_INDEX_MEASUREMENT_ONT_DATE,
		constants.SQL_INDEX_MEASUREMENT_ONT_IDX,
		constants.SQL_INDEX_MEASUREMENT_ONT_DESPT,
		constants.SQL_INDEX_MEASUREMENT_ONT_PON_ID,
		constants.SQL_INDEX_REPORT_CATEGORY,
		constants.SQL_INDEX_REPORT_USER_ID,
		constants.SQL_INDEX_USERS_USERNAME,
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
