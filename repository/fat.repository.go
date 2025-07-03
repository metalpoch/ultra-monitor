package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type FatRepository interface {
	AllFat(ctx context.Context, page, limit uint16) ([]entity.Fat, error)
	AddFat(ctx context.Context, fat *entity.Fat) error
	DeleteOne(ctx context.Context, id int32) error
	GetByID(ctx context.Context, id int32) (entity.Fat, error)
	GetTraffic(ctx context.Context, id int, initDate, endDate time.Time) ([]entity.Traffic, error)
	GetStates(ctx context.Context) ([]string, error)
	GetMunicipality(ctx context.Context, state string) ([]string, error)
	GetCounty(ctx context.Context, state, municipality string) ([]string, error)
	GetOdn(ctx context.Context, state, municipality, county string) ([]string, error)
	GetFatsByStates(ctx context.Context, state string) ([]entity.Fat, error)
	GetFatsByMunicipality(ctx context.Context, state, municipality string) ([]entity.Fat, error)
	GetFatsByCounty(ctx context.Context, state, municipality, county string) ([]entity.Fat, error)
	GetFatsBytOdn(ctx context.Context, state, municipality, county, odn string) ([]entity.Fat, error)
}

type fatRepository struct {
	db *sqlx.DB
}

func NewFatRepository(db *sqlx.DB) *fatRepository {
	return &fatRepository{db}
}

func (repo *fatRepository) AllFat(ctx context.Context, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats ORDER BY region, state, municipality, county LIMIT $1 OFFSET $2`
	err := repo.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (repo *fatRepository) AddFat(ctx context.Context, fat *entity.Fat) error {
	query := `
    INSERT INTO fats (
        fat, region, state, municipality, county, odn, olt_ip,
        pon_shell, pon_port, pon_card, latitude, longitude
    ) VALUES (
        :fat, :region, :state, :municipality, :county, :odn, :olt_ip,
        :pon_shell, :pon_port, :pon_card, :latitude, :longitude
    )
    ON CONFLICT (fat, state, municipality, county, olt_ip, odn, pon_shell, pon_card, pon_port)
    DO UPDATE SET
        region = EXCLUDED.region,
        latitude = EXCLUDED.latitude,
        longitude = EXCLUDED.longitude;`

	_, err := repo.db.NamedExecContext(ctx, query, fat)
	return err
}

func (repo *fatRepository) DeleteOne(ctx context.Context, id int32) error {
	query := `DELETE FROM fats WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *fatRepository) GetByID(ctx context.Context, id int32) (entity.Fat, error) {
	var res entity.Fat
	query := `SELECT * FROM fats WHERE id = $1`
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
	query := `SELECT DISTINCT odn FROM fats WHERE state = $1 AND municipality = $2 AND municipality = $3 ORDER BY odn;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (repo *fatRepository) GetFatsByStates(ctx context.Context, state string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 ORDER BY id;`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *fatRepository) GetFatsByMunicipality(ctx context.Context, state, municipality string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 ORDER BY id;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (repo *fatRepository) GetFatsByCounty(ctx context.Context, state, municipality, county string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY id;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (repo *fatRepository) GetFatsBytOdn(ctx context.Context, state, municipality, county, odn string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 AND odn = $4 ORDER BY id;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county, odn)
	return res, err
}
