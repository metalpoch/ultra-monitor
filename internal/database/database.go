package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
		SQL_TABLE_USERS,
		SQL_TABLE_REPORT,
		SQL_TABLE_OLT,
		SQL_TABLE_PON,
		SQL_TABLE_MEASUREMENT_PON,
		SQL_TABLE_TRAFFIC_PON,
		SQL_TABLE_TRAFFIC_PON_SUMMARY,
		SQL_TABLE_MEASUREMENT_ONT,
		SQL_TABLE_FAT,
		SQL_TABLE_ONT_SUMMARY_STATUS_COUNTS,
	}
	indexQueries := []string{
		SQL_INDEX_FAT_STATE,
		SQL_INDEX_FAT_COUNTY,
		SQL_INDEX_FAT_MUNICIPALITY,
		SQL_INDEX_FAT_ODN,
		SQL_INDEX_FAT_OLT,
		SQL_INDEX_MEASUREMENT_ONT_DATE,
		SQL_INDEX_MEASUREMENT_ONT_IDX,
		SQL_INDEX_MEASUREMENT_ONT_DESPT,
		SQL_INDEX_MEASUREMENT_ONT_PON_ID,
		SQL_INDEX_REPORT_CATEGORY,
		SQL_INDEX_REPORT_USER_ID,
		SQL_INDEX_USERS_USERNAME,
		SQL_INDEX_MEASUREMENT_ONT_PON_STATS,
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
