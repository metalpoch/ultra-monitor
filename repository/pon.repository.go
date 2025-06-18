package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/repository/internal/model"
)

type PonRepository interface {
	PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error)
	PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error)
	TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByMunicipality(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByCounty(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByODN(ctx context.Context, state, municipality, county, odn string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByPon(ctx context.Context, sysname, ifname string, initDate, endDate time.Time) ([]entity.Traffic, error)
	GetDailyAveragedHourlyMaxTrafficTrends(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficSummary, error)

	UpsertSummaryTraffic(ctx context.Context, traffic []entity.TrafficSummary) error
	GetTrafficSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficTotalSummary, error)
	GetTrafficStatesSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
	GetTrafficMunicipalitySummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
	GetTrafficCountySummary(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
	GetTrafficOdnSummary(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
}

type ponRepository struct {
	db *sqlx.DB
}

func NewPonRepository(db *sqlx.DB) *ponRepository {
	return &ponRepository{db}
}

func (repo *ponRepository) PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error) {
	var res []entity.Pon
	query := `SELECT pons.* FROM pons JOIN olts ON olts.ip = pons.olt_ip WHERE olts.sys_name = $1`
	err := repo.db.SelectContext(ctx, &res, query, sysname)
	return res, err
}

func (repo *ponRepository) PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error) {
	var res entity.Pon
	query := `SELECT pons.* FROM pons JOIN olts ON olts.ip = pons.olt_ip WHERE olts.sys_name = $1 AND pons.if_name = $2`
	err := repo.db.GetContext(ctx, &res, query, sysname, port)
	return res, err
}

func (repo *ponRepository) TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in_sec,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out_sec
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND traffic_pons.date BETWEEN $2 AND $3
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByMunicipality(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in_sec,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out_sec
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND fats.municipality = $2 AND traffic_pons.date BETWEEN $3 AND $4
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByCounty(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in_sec,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out_sec
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND fats.municipality = $2 AND fats.county = $3 AND traffic_pons.date BETWEEN $4 AND $5
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, county, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByODN(ctx context.Context, state, municipality, county, odn string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in_sec,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out_sec
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND fats.municipality = $2 AND fats.county = $3 AND fats.odn = $4 AND traffic_pons.date BETWEEN $5 AND $6
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, county, odn, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
    FROM traffic_pons 
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    WHERE olts.sys_name = $1 AND traffic_pons.date BETWEEN $2 AND $3
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, sysname, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByPon(ctx context.Context, sysname, ifname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in_sec,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out_sec
    FROM traffic_pons 
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    WHERE olts.sys_name = $1 AND pons.if_name = $2 AND traffic_pons.date BETWEEN $3 AND $4
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`
	err := repo.db.SelectContext(ctx, &res, query, sysname, ifname, initDate, endDate)
	return res, err
}

func (repo *ponRepository) GetDailyAveragedHourlyMaxTrafficTrends(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficSummary, error) {
	trendsQuery := `
        SELECT
            day,
            pon_id,
            AVG(max_bps_in) / 1e6 AS mbps_in,
            AVG(max_bps_out) / 1e6 AS mbps_out,
            AVG(max_bytes_in_sec) / 1e6 AS mbytes_in_sec,
            AVG(max_bytes_out_sec) / 1e6 AS mbytes_out_sec
        FROM (
            SELECT
                DATE(date) AS day,
                pon_id,
                EXTRACT(HOUR FROM date) AS hour,
                MAX(bps_in) AS max_bps_in,
                MAX(bps_out) AS max_bps_out,
                MAX(bytes_in_sec) AS max_bytes_in_sec,
                MAX(bytes_out_sec) AS max_bytes_out_sec
            FROM traffic_pons
            WHERE date BETWEEN $1 AND $2
            GROUP BY day, pon_id, hour
        ) hourly_max
        GROUP BY day, pon_id
        ORDER BY day, pon_id;`

	var trends []model.TrafficTrend
	err := repo.db.SelectContext(ctx, &trends, trendsQuery, initDate, endDate)
	if err != nil {
		return nil, err
	}

	metaQuery := `
        SELECT DISTINCT p.id AS pon_id, f.id AS fat_id, p.olt_ip
        FROM pons p
        JOIN fats f ON f.olt_ip = p.olt_ip;`

	var metas []model.TrafficMeta
	err = repo.db.SelectContext(ctx, &metas, metaQuery)
	if err != nil {
		return nil, err
	}
	metaMap := make(map[int32]model.TrafficMeta)
	for _, m := range metas {
		metaMap[m.PonID] = m
	}

	type groupKey struct {
		Day   time.Time
		FatID int32
		OltIP string
	}
	grouped := make(map[groupKey]*entity.TrafficSummary)
	for _, t := range trends {
		meta, ok := metaMap[t.PonID]
		if !ok {
			continue
		}
		key := groupKey{
			Day:   t.Day,
			FatID: meta.FatID,
			OltIP: meta.OltIP,
		}
		if _, exists := grouped[key]; !exists {
			grouped[key] = &entity.TrafficSummary{
				Day:          t.Day,
				FatID:        meta.FatID,
				OltIP:        meta.OltIP,
				MbpsIn:       0,
				MbpsOut:      0,
				MbytesInSec:  0,
				MbytesOutSec: 0,
			}
		}
		grouped[key].MbpsIn += t.MbpsIn
		grouped[key].MbpsOut += t.MbpsOut
		grouped[key].MbytesInSec += t.MbytesInSec
		grouped[key].MbytesOutSec += t.MbytesOutSec
	}

	var res []entity.TrafficSummary
	for _, v := range grouped {
		res = append(res, *v)
	}
	return res, nil
}

func (repo *ponRepository) UpsertSummaryTraffic(ctx context.Context, counts []entity.TrafficSummary) error {
	const fieldCount = 7
	query := `
        INSERT INTO traffic_pons_summary (
            day, fat_id, olt_ip, mbps_in, mbps_out, mbytes_in_sec, mbytes_out_sec
        ) VALUES `
	valueStrings := make([]string, 0, len(counts))
	valueArgs := make([]interface{}, 0, len(counts)*fieldCount)

	for i, c := range counts {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			i*fieldCount+1, i*fieldCount+2, i*fieldCount+3, i*fieldCount+4, i*fieldCount+5, i*fieldCount+6, i*fieldCount+7))
		valueArgs = append(valueArgs,
			c.Day, c.FatID, c.OltIP, c.MbpsIn, c.MbpsOut, c.MbytesInSec, c.MbytesOutSec)
	}

	query += strings.Join(valueStrings, ", ")
	query += `
	    ON CONFLICT (day, fat_id, olt_ip) DO UPDATE SET
    	    mbps_in = EXCLUDED.mbps_in,
	        mbps_out = EXCLUDED.mbps_out,
	        mbytes_in_sec = EXCLUDED.mbytes_in_sec,
	        mbytes_out_sec = EXCLUDED.mbytes_out_sec`

	_, err := repo.db.ExecContext(ctx, query, valueArgs...)
	return err
}

func (repo *ponRepository) GetTrafficSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficTotalSummary, error) {
	var res []entity.TrafficTotalSummary
	query := `
    SELECT
        day,
        SUM(mbps_in) AS mbps_in,
        SUM(mbps_out) AS mbps_out,
        SUM(mbytes_in_sec) AS mbytes_in_sec,
        SUM(mbytes_out_sec) AS mbytes_out_sec
    FROM traffic_pons_summary
    WHERE day BETWEEN $1 AND $2
    GROUP BY day
    ORDER BY day;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ponRepository) GetTrafficStatesSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    SELECT
        traffic.day AS day,
        fats.state AS description,
        SUM(traffic.mbps_in) AS mbps_in,
        SUM(traffic.mbps_out) AS mbps_out,
        SUM(traffic.mbytes_in_sec) AS mbytes_in_sec,
        SUM(traffic.mbytes_out_sec) AS mbytes_out_sec
    FROM traffic_pons_summary AS traffic
    JOIN fats ON fats.olt_ip = traffic.olt_ip
    WHERE day BETWEEN $1 AND $2
    GROUP BY day, state
    ORDER BY state, day;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ponRepository) GetTrafficMunicipalitySummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    SELECT
        traffic.day AS day,
        fats.municipality AS description,
        SUM(traffic.mbps_in) AS mbps_in,
        SUM(traffic.mbps_out) AS mbps_out,
        SUM(traffic.mbytes_in_sec) AS mbytes_in_sec,
        SUM(traffic.mbytes_out_sec) AS mbytes_out_sec
    FROM traffic_pons_summary AS traffic
    JOIN fats ON fats.olt_ip = traffic.olt_ip
    WHERE fats.state = $1 AND day BETWEEN $2 AND $3
    GROUP BY day, municipality
    ORDER BY municipality, day;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ponRepository) GetTrafficCountySummary(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    SELECT
        traffic.day AS day,
        fats.county AS description,
        SUM(traffic.mbps_in) AS mbps_in,
        SUM(traffic.mbps_out) AS mbps_out,
        SUM(traffic.mbytes_in_sec) AS mbytes_in_sec,
        SUM(traffic.mbytes_out_sec) AS mbytes_out_sec
    FROM traffic_pons_summary AS traffic
    JOIN fats ON fats.olt_ip = traffic.olt_ip
    WHERE fats.state = $1 AND fats.municipality = $2 AND day BETWEEN $3 AND $4
    GROUP BY day, county
    ORDER BY county, day;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ponRepository) GetTrafficOdnSummary(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    SELECT
        traffic.day AS day,
        fats.odn AS description,
        SUM(traffic.mbps_in) AS mbps_in,
        SUM(traffic.mbps_out) AS mbps_out,
        SUM(traffic.mbytes_in_sec) AS mbytes_in_sec,
        SUM(traffic.mbytes_out_sec) AS mbytes_out_sec
    FROM traffic_pons_summary AS traffic
    JOIN fats ON fats.olt_ip = traffic.olt_ip
    WHERE fats.state = $1 AND fats.municipality = $2 AND fats.county = $3 AND day BETWEEN $4 AND $5
    GROUP BY day, odn
    ORDER BY odn, day;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, county, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}
