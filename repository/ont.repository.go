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

type OntRepository interface {
	UpdateStatusSummary(ctx context.Context, counts []entity.OntSummaryStatusCounts) error
	GetDailyAveragedHourlyStatusSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.OntSummaryStatusCounts, error)
	GetStatusSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.OntSummaryStatus, error)
	GetStatusIPSummary(ctx context.Context, ip string, initDate, endDate time.Time) ([]entity.OntSummaryStatus, error)
	GetStatusStateSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.GetStatusSummary, error)
	GetStatusByStateSummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntSummaryStatus, error)
	GetStatusMunicipalitySummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.GetStatusSummary, error)
	GetStatusCountySummary(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.GetStatusSummary, error)
	GetStatusOdnSummary(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.GetStatusSummary, error)
	TrafficOnt(ctx context.Context, ponID int, idx int64, initDate, endDate time.Time) ([]entity.TrafficOnt, error)
	TrafficOntByDespt(ctx context.Context, despt string, initDate, endDate time.Time) ([]entity.TrafficOnt, error)
}

type ontRepository struct {
	db *sqlx.DB
}

func NewOntRepository(db *sqlx.DB) *ontRepository {
	return &ontRepository{db}
}

func (repo *ontRepository) GetDailyAveragedHourlyStatusSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.OntSummaryStatusCounts, error) {
	var counts []model.StatusCounts
	countsQuery := `
        SELECT
            DATE_TRUNC('day', m.date) AS day,
            m.pon_id,
            p.olt_ip,
            COUNT(*) FILTER (WHERE m.control_run_status = 1) AS actives,
            COUNT(*) FILTER (WHERE m.control_run_status = 2) AS inactives,
            COUNT(*) FILTER (WHERE m.control_run_status NOT IN (1, 2)) AS unknowns
        FROM measurement_onts m
        JOIN pons p ON p.id = m.pon_id
        WHERE m.date BETWEEN $1 AND $2
        GROUP BY day, m.pon_id, p.olt_ip
        ORDER BY day, m.pon_id, p.olt_ip;
    `
	err := repo.db.SelectContext(ctx, &counts, countsQuery, initDate, endDate)
	if err != nil {
		return nil, err
	}

	var metas []model.FatMeta
	metaQuery := `
        SELECT p.id AS pon_id, f.id AS fat_id, p.olt_ip
        FROM pons p
        JOIN fats f ON f.olt_ip = p.olt_ip
    `
	err = repo.db.SelectContext(ctx, &metas, metaQuery)
	if err != nil {
		return nil, err
	}
	metaMap := make(map[int32]model.FatMeta)
	for _, m := range metas {
		metaMap[m.PonID] = m
	}

	type groupKey struct {
		Day   time.Time
		FatID int32
		OltIP string
	}
	grouped := make(map[groupKey]*entity.OntSummaryStatusCounts)

	for _, c := range counts {
		meta, ok := metaMap[c.PonID]
		if !ok {
			continue
		}
		key := groupKey{
			Day:   c.Day,
			FatID: meta.FatID,
			OltIP: c.OltIP,
		}
		if _, exists := grouped[key]; !exists {
			grouped[key] = &entity.OntSummaryStatusCounts{
				Day:           c.Day,
				FatID:         meta.FatID,
				OltIP:         c.OltIP,
				PonsCount:     0,
				ActiveCount:   0,
				InactiveCount: 0,
				UnknownCount:  0,
			}
		}
		grouped[key].PonsCount++
		grouped[key].ActiveCount += c.Actives
		grouped[key].InactiveCount += c.Inactives
		grouped[key].UnknownCount += c.Unknowns
	}

	var res []entity.OntSummaryStatusCounts
	for _, v := range grouped {
		res = append(res, *v)
	}
	return res, nil
}

func (repo *ontRepository) UpdateStatusSummary(ctx context.Context, counts []entity.OntSummaryStatusCounts) error {
	const fieldCount = 7
	query := `
    INSERT INTO ont_summary_status_count (
        day, fat_id, olt_ip, ports_pon, actives, inactives, unknowns
    ) VALUES `
	valueStrings := make([]string, 0, len(counts))
	valueArgs := make([]interface{}, 0, len(counts)*fieldCount)

	for i, c := range counts {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			i*fieldCount+1, i*fieldCount+2, i*fieldCount+3, i*fieldCount+4, i*fieldCount+5, i*fieldCount+6, i*fieldCount+7))
		valueArgs = append(valueArgs,
			c.Day, c.FatID, c.OltIP, c.PonsCount, c.ActiveCount, c.InactiveCount, c.UnknownCount)
	}

	query += strings.Join(valueStrings, ", ")
	query += `
    ON CONFLICT (day, fat_id, olt_ip) DO UPDATE SET
        ports_pon = EXCLUDED.ports_pon,
        actives = EXCLUDED.actives,
        inactives = EXCLUDED.inactives,
        unknowns = EXCLUDED.unknowns`

	_, err := repo.db.ExecContext(ctx, query, valueArgs...)
	return err
}

func (repo *ontRepository) GetStatusSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.OntSummaryStatus, error) {
	var res []entity.OntSummaryStatus
	query := `
    SELECT
        day,
        SUM(ports_pon) AS ports_pon,
        SUM(actives) AS actives,
        SUM(inactives) AS inactives,
        SUM(unknowns) AS unknowns
    FROM ont_summary_status_count AS ont
    WHERE day BETWEEN $1 AND $2
    GROUP BY day
    ORDER BY day;`
	err := repo.db.SelectContext(ctx, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ontRepository) GetStatusIPSummary(ctx context.Context, ip string, initDate, endDate time.Time) ([]entity.OntSummaryStatus, error) {
	var res []entity.OntSummaryStatus
	query := `
    SELECT
        day,
        SUM(ports_pon) AS ports_pon,
        SUM(actives) AS actives,
        SUM(inactives) AS inactives,
        SUM(unknowns) AS unknowns
    FROM ont_summary_status_count
    WHERE olt_ip = $1 AND day BETWEEN $2 AND $3
    GROUP BY day
    ORDER BY day;`
	err := repo.db.SelectContext(ctx, &res, query, ip, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ontRepository) GetStatusStateSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.GetStatusSummary, error) {
	var res []entity.GetStatusSummary
	query := `
    SELECT
        ont.day AS day,
        fats.state AS description,
        SUM(ont.ports_pon) AS ports_pon,
        SUM(ont.actives) AS actives,
        SUM(ont.inactives) AS inactives,
        SUM(ont.unknowns) AS unknowns
    FROM ont_summary_status_count AS ont
    JOIN fats ON fats.id = ont.fat_id
    WHERE day BETWEEN $1 AND $2
    GROUP BY day, state
    ORDER BY state, day;`
	err := repo.db.SelectContext(ctx, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ontRepository) GetStatusByStateSummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntSummaryStatus, error) {
	var res []entity.OntSummaryStatus
	query := `
    SELECT
        day,
        SUM(ports_pon) AS ports_pon,
        SUM(actives) AS actives,
        SUM(inactives) AS inactives,
        SUM(unknowns) AS unknowns
    FROM ont_summary_status_count
    JOIN fats ON fats.id = fat_id
    WHERE fats.state = $1 AND day BETWEEN $2 AND $3
    GROUP BY day
    ORDER BY day;`
	err := repo.db.SelectContext(ctx, &res, query, state, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ontRepository) GetStatusMunicipalitySummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.GetStatusSummary, error) {
	var res []entity.GetStatusSummary
	query := `
    SELECT
        ont.day AS day,
        fats.municipality AS description,
        SUM(ont.ports_pon) AS ports_pon,
        SUM(ont.actives) AS actives,
        SUM(ont.inactives) AS inactives,
        SUM(ont.unknowns) AS unknowns
    FROM ont_summary_status_count AS ont
    JOIN fats ON fats.id = ont.fat_id
    WHERE fats.state = $1 AND day BETWEEN $2 AND $3
    GROUP BY day, municipality
    ORDER BY municipality, day;`
	err := repo.db.SelectContext(ctx, &res, query, state, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ontRepository) GetStatusCountySummary(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.GetStatusSummary, error) {
	var res []entity.GetStatusSummary
	query := `
    SELECT
        ont.day AS day,
        fats.county AS description,
        SUM(ont.ports_pon) AS ports_pon,
        SUM(ont.actives) AS actives,
        SUM(ont.inactives) AS inactives,
        SUM(ont.unknowns) AS unknowns
    FROM ont_summary_status_count AS ont
    JOIN fats ON fats.id = ont.fat_id
    WHERE fats.state = $1 AND fats.municipality = $2 AND day BETWEEN $3 AND $4
    GROUP BY day, county
    ORDER BY county, day;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ontRepository) GetStatusOdnSummary(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.GetStatusSummary, error) {
	var res []entity.GetStatusSummary
	query := `
    SELECT
        ont.day AS day,
        fats.odn AS description,
        SUM(ont.ports_pon) AS ports_pon,
        SUM(ont.actives) AS actives,
        SUM(ont.inactives) AS inactives,
        SUM(ont.unknowns) AS unknowns
    FROM ont_summary_status_count AS ont
    JOIN fats ON fats.id = ont.fat_id
    WHERE fats.state = $1 AND fats.municipality = $2 AND fats.county = $3 AND day BETWEEN $4 AND $5
    GROUP BY day, odn
    ORDER BY odn, day;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ontRepository) TrafficOnt(ctx context.Context, ponID int, idx int64, initDate, endDate time.Time) ([]entity.TrafficOnt, error) {
	var res []entity.TrafficOnt
	query := `
	SELECT
        date,
        despt,
        serial_number,
        line_prof_name,
        olt_distance,
        control_mac_count,
        control_run_status,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_in - prev_bytes_in) * 8) / (time_diff * 1000000)
        END AS mbps_in,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_out - prev_bytes_out) * 8) / (time_diff * 1000000)
        END AS mbps_out,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) / (time_diff * 1000000)
            ELSE (curr_bytes_in - prev_bytes_in) / (time_diff * 1000000)
        END AS mbytes_in_sec,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) / (time_diff * 1000000)
            ELSE (curr_bytes_out - prev_bytes_out) / (time_diff * 1000000)
        END AS mbytes_out_sec
    FROM (
        SELECT
            date,
            despt,
            serial_number,
            line_prof_name,
            olt_distance,
            control_mac_count,
            control_run_status,
            bytes_in_count AS prev_bytes_in,
            bytes_out_count AS prev_bytes_out,
            LEAD(bytes_in_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_in,
            LEAD(bytes_out_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_out,
            EXTRACT(EPOCH FROM (LEAD(date) OVER (PARTITION BY pon_id ORDER BY date) - date)) AS time_diff
        FROM measurement_onts
        WHERE pon_id = $1 AND idx = $2 AND bytes_in_count > 0 AND bytes_out_count > 0 AND date BETWEEN $3 AND $4
        ORDER BY date
    ) AS subquery
    WHERE curr_bytes_in IS NOT NULL
      AND curr_bytes_out IS NOT NULL
      AND time_diff IS NOT NULL;`

	err := repo.db.SelectContext(ctx, &res, query, ponID, idx, initDate, endDate)
	return res, err
}

func (repo *ontRepository) TrafficOntByDespt(ctx context.Context, despt string, initDate, endDate time.Time) ([]entity.TrafficOnt, error) {
	var res []entity.TrafficOnt
	query := `SELECT
        date,
        despt,
        serial_number,
        line_prof_name,
        olt_distance,
        control_mac_count,
        control_run_status,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_in - prev_bytes_in) * 8) / (time_diff * 1000000)
        END AS mbps_in,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_out - prev_bytes_out) * 8) / (time_diff * 1000000)
        END AS mbps_out,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) / (time_diff * 1000000)
            ELSE (curr_bytes_in - prev_bytes_in) / (time_diff * 1000000)
        END AS mbytes_in_sec,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) / (time_diff * 1000000)
            ELSE (curr_bytes_out - prev_bytes_out) / (time_diff * 1000000)
        END AS mbytes_out_sec
    FROM (
        SELECT
            date,
            despt,
            serial_number,
            line_prof_name,
            olt_distance,
            control_mac_count,
            control_run_status,
            bytes_in_count AS prev_bytes_in,
            bytes_out_count AS prev_bytes_out,
            LEAD(bytes_in_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_in,
            LEAD(bytes_out_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_out,
            EXTRACT(EPOCH FROM (LEAD(date) OVER (PARTITION BY pon_id ORDER BY date) - date)) AS time_diff
        FROM measurement_onts
        WHERE despt = $1 AND bytes_in_count > 0 AND bytes_out_count > 0 AND date BETWEEN $2 AND $3
        ORDER BY date
    ) AS subquery
    WHERE curr_bytes_in IS NOT NULL
      AND curr_bytes_out IS NOT NULL
      AND time_diff IS NOT NULL;`

	err := repo.db.SelectContext(ctx, &res, query, despt, initDate, endDate)
	return res, err
}
