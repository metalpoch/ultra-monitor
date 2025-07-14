package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type FatRepository interface {
	AllFat(ctx context.Context, page, limit uint16) ([]entity.FatResponse, error)
	AddFat(ctx context.Context, fat entity.Fat) (int64, error)
	AddFatsOntStatusSummary(ctx context.Context, summary entity.FatsOntStatusSummary) error
	DeleteOne(ctx context.Context, id int32) error
	GetByID(ctx context.Context, id int32) (entity.FatResponse, error)
	GetTraffic(ctx context.Context, id int, initDate, endDate time.Time) ([]entity.Traffic, error)
	GetStates(ctx context.Context) ([]string, error)
	GetMunicipality(ctx context.Context, state string) ([]string, error)
	GetCounty(ctx context.Context, state, municipality string) ([]string, error)
	GetOdn(ctx context.Context, state, municipality, county string) ([]string, error)
	GetFatsByStates(ctx context.Context, state string) ([]entity.FatResponse, error)
	GetFatsByMunicipality(ctx context.Context, state, municipality string) ([]entity.FatResponse, error)
	GetFatsByCounty(ctx context.Context, state, municipality, county string) ([]entity.FatResponse, error)
	GetFatsBytOdn(ctx context.Context, state, municipality, county, odn string) ([]entity.FatResponse, error)
}

type fatRepository struct {
	db *sqlx.DB
}

func NewFatRepository(db *sqlx.DB) *fatRepository {
	return &fatRepository{db}
}

func (repo *fatRepository) AllFat(ctx context.Context, page, limit uint16) ([]entity.FatResponse, error) {
	var res []entity.FatResponse
	offset := (page - 1) * limit
	query := `
        SELECT f.*, 
            COALESCE(s.actives,0) AS actives, 
            COALESCE(s.inactive,0) AS inactive, 
            COALESCE(s.others,0) AS others
        FROM fats f
        LEFT JOIN LATERAL (
            SELECT s.actives, s.inactive, s.others
            FROM fats_ont_status_summary s
            WHERE s.fat_id = f.id
            ORDER BY s.day DESC
            LIMIT 1
        ) s ON TRUE
        ORDER BY f.region, f.state, f.municipality, f.county
        LIMIT $1 OFFSET $2;`
	err := repo.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (repo *fatRepository) AddFat(ctx context.Context, fat entity.Fat) (int64, error) {
	query := `
    INSERT INTO fats (
        fat, region, state, municipality, county, odn, olt_ip,
        pon_shell, pon_port, pon_card
    ) VALUES (
        :fat, :region, :state, :municipality, :county, :odn, :olt_ip,
        :pon_shell, :pon_port, :pon_card
    )
    ON CONFLICT (fat, state, municipality, county, olt_ip, odn, pon_shell, pon_card, pon_port)
    DO UPDATE SET
        region = EXCLUDED.region
    RETURNING id;`

	var id int64
	stmt, err := repo.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &id, fat)
	return id, err
}

func (repo *fatRepository) AddFatsOntStatusSummary(ctx context.Context, summary entity.FatsOntStatusSummary) error {
	query := `
	INSERT INTO fats_ont_status_summary (day, fat_id, actives, inactives, others)
	VALUES (:day, :fat_id, :actives, :inactives, :others)
	ON CONFLICT (day, fat_id)
	DO UPDATE SET
	actives = EXCLUDED.actives,
	inactives = EXCLUDED.inactives,
	others = EXCLUDED.others;`
	_, err := repo.db.NamedExecContext(ctx, query, summary)
	return err
}

func (repo *fatRepository) DeleteOne(ctx context.Context, id int32) error {
	query := `DELETE FROM fats WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *fatRepository) GetByID(ctx context.Context, id int32) (entity.FatResponse, error) {
	var res entity.FatResponse
	query := `
        SELECT f.*, 
            COALESCE(s.actives,0) AS actives, 
            COALESCE(s.inactive,0) AS inactive, 
            COALESCE(s.others,0) AS others
        FROM fats f
        LEFT JOIN LATERAL (
            SELECT s.actives, s.inactive, s.others
            FROM fats_ont_status_summary s
            WHERE s.fat_id = f.id
            ORDER BY s.day DESC
            LIMIT 1
        ) s ON TRUE
        WHERE f.id = $1
        LIMIT 1;`
	err := repo.db.GetContext(ctx, &res, query, id)
	return res, err
}

func (repo *fatRepository) GetTraffic(ctx context.Context, id int, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic

	var fat entity.Fat
	err := repo.db.GetContext(ctx, &fat, `SELECT * FROM fats WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	ifName := fmt.Sprintf("GPON %d/%d/%d", fat.Shell, fat.Card, fat.Port)
	fmt.Println(fat.OltIP, ifName, initDate, endDate)
	query := `
    SELECT
        DATE_TRUNC('minute', tp.date) AS date,
        COALESCE(SUM(tp.bps_in),0) / 1000000 AS mbps_in,
        COALESCE(SUM(tp.bps_out),0) / 1000000 AS mbps_out,
        COALESCE(SUM(tp.bandwidth_mbps_sec),0) / 1000000 AS bandwidth_mbps_sec,
        COALESCE(SUM(tp.bytes_in),0) / 1000000 AS mbytes_in,
        COALESCE(SUM(tp.bytes_out),0) / 1000000 AS mbytes_out
    FROM traffic_pons tp
    JOIN pons p ON p.id = tp.pon_id
    WHERE p.olt_ip = $1 AND p.if_name = $2 AND tp.date BETWEEN $3 AND $4
    GROUP BY DATE_TRUNC('minute', tp.date)
    ORDER BY date;`
	err = repo.db.SelectContext(ctx, &res, query, fat.OltIP, ifName, initDate, endDate)
	return res, err
}

func (repo *fatRepository) GetStates(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT state FROM fats ORDER BY state;`
	err := repo.db.SelectContext(ctx, &res, query)
	return res, err
}

func (repo *fatRepository) GetMunicipality(ctx context.Context, state string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT municipality FROM fats WHERE state = $1 ORDER BY municipality;`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *fatRepository) GetCounty(ctx context.Context, state, municipality string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT county FROM fats WHERE state = $1 AND municipality = $2 ORDER BY county;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (repo *fatRepository) GetOdn(ctx context.Context, state, municipality, county string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY odn;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (repo *fatRepository) GetFatsByStates(ctx context.Context, state string) ([]entity.FatResponse, error) {
	var res []entity.FatResponse
	query := `
        SELECT f.*, 
            COALESCE(s.actives,0) AS actives, 
            COALESCE(s.inactive,0) AS inactive, 
            COALESCE(s.others,0) AS others
        FROM fats f
        LEFT JOIN LATERAL (
            SELECT s.actives, s.inactives, s.others
            FROM fats_ont_status_summary s
            WHERE s.fat_id = f.id
            ORDER BY s.day DESC
            LIMIT 1
        ) s ON TRUE
        WHERE f.state = $1
        ORDER BY f.state, f.municipality, f.county, f.odn;`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *fatRepository) GetFatsByMunicipality(ctx context.Context, state, municipality string) ([]entity.FatResponse, error) {
	var res []entity.FatResponse
	query := `
        SELECT f.*, 
            COALESCE(s.actives,0) AS actives, 
            COALESCE(s.inactive,0) AS inactive, 
            COALESCE(s.others,0) AS others
        FROM fats f
        LEFT JOIN LATERAL (
            SELECT s.actives, s.inactives, s.others
            FROM fats_ont_status_summary s
            WHERE s.fat_id = f.id
            ORDER BY s.day DESC
            LIMIT 1
        ) s ON TRUE
        WHERE f.state = $1 AND f.municipality = $2
        ORDER BY f.state, f.municipality, f.county, f.odn;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (repo *fatRepository) GetFatsByCounty(ctx context.Context, state, municipality, county string) ([]entity.FatResponse, error) {
	var res []entity.FatResponse
	query := `
        SELECT f.*, 
            COALESCE(s.actives,0) AS actives, 
            COALESCE(s.inactive,0) AS inactive, 
            COALESCE(s.others,0) AS others
        FROM fats f
        LEFT JOIN LATERAL (
            SELECT s.actives, s.inactives, s.others
            FROM fats_ont_status_summary s
            WHERE s.fat_id = f.id
            ORDER BY s.day DESC
            LIMIT 1
        ) s ON TRUE
        WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3
        ORDER BY f.state, f.municipality, f.county, f.odn;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (repo *fatRepository) GetFatsBytOdn(ctx context.Context, state, municipality, county, odn string) ([]entity.FatResponse, error) {
	var res []entity.FatResponse
	query := `
        SELECT f.*, 
            COALESCE(s.actives,0) AS actives, 
            COALESCE(s.inactive,0) AS inactive, 
            COALESCE(s.others,0) AS others
        FROM fats f
        LEFT JOIN LATERAL (
            SELECT s.actives, s.inactive, s.others
            FROM fats_ont_status_summary s
            WHERE s.fat_id = f.id
            ORDER BY s.day DESC
            LIMIT 1
        ) s ON TRUE
        WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3 AND f.odn = $4
        ORDER BY f.state, f.municipality, f.county, f.odn;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county, odn)
	return res, err
}
