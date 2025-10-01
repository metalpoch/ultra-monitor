package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type TrafficRepository interface {
	GetSnmpIndexByMunicipality(ctx context.Context, state, municipality string) ([]entity.OltIndex, error)
	GetSnmpIndexByCounty(ctx context.Context, state, municipality, county string) ([]entity.OltIndex, error)
	GetSnmpIndexByODN(ctx context.Context, state, municipality, odn string) ([]entity.OltIndex, error)
	SaveSummaryTraffic(ctx context.Context, trafficData []entity.SumaryTraffic) error
	GetTotalTrafficByIP(ctx context.Context, ip string, startTime, endTime time.Time) (*entity.TrafficSummary, error)
	GetTotalTrafficByState(ctx context.Context, state string, startTime, endTime time.Time) ([]entity.TrafficSummary, error)
	GetTotalTrafficByRegion(ctx context.Context, region string, startTime, endTime time.Time) ([]entity.TrafficSummary, error)
	GetTotalTraffic(ctx context.Context, startTime, endTime time.Time) ([]entity.TrafficSummary, error)
	GetTrafficGroupedByRegion(ctx context.Context, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error)
	GetTrafficGroupedByState(ctx context.Context, region string, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error)
	GetTrafficGroupedByIP(ctx context.Context, state string, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error)
}

type trafficRepository struct {
	db *sqlx.DB
}

func NewTrafficRepository(db *sqlx.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (r *trafficRepository) GetSnmpIndexByMunicipality(ctx context.Context, state, municipality string) ([]entity.OltIndex, error) {
	var res []entity.OltIndex
	query := `SELECT pd.ip, pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.state = $1 AND f.municipality = $2;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (r *trafficRepository) GetSnmpIndexByCounty(ctx context.Context, state, municipality, county string) ([]entity.OltIndex, error) {
	var res []entity.OltIndex
	query := `SELECT pd.ip, pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (r *trafficRepository) GetSnmpIndexByODN(ctx context.Context, state, municipality, odn string) ([]entity.OltIndex, error) {
	var res []entity.OltIndex
	query := `SELECT pd.ip, pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.state = $1 AND f.municipality = $2 AND f.odn = $3;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, odn)
	return res, err
}

func (r *trafficRepository) SaveSummaryTraffic(ctx context.Context, trafficData []entity.SumaryTraffic) error {
	query := `INSERT INTO summary_traffic (time, ip, state, region, bps_in, bps_out, bytes_in, bytes_out)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (ip, time)
		DO UPDATE SET
			state = EXCLUDED.state,
			region = EXCLUDED.region,
			bps_in = EXCLUDED.bps_in,
			bps_out = EXCLUDED.bps_out,
			bytes_in = EXCLUDED.bytes_in,
			bytes_out = EXCLUDED.bytes_out`

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, traffic := range trafficData {
		_, err := tx.ExecContext(ctx, query,
			traffic.Time,
			traffic.IP,
			traffic.State,
			traffic.Region,
			traffic.BpsIn,
			traffic.BpsOut,
			traffic.BytesIn,
			traffic.BytesOut)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *trafficRepository) GetTotalTrafficByIP(ctx context.Context, ip string, startTime, endTime time.Time) (*entity.TrafficSummary, error) {
	var res entity.TrafficSummary
	query := `SELECT
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE ip = $1 AND time BETWEEN $2 AND $3`

	err := r.db.GetContext(ctx, &res, query, ip, startTime, endTime)
	return &res, err
}

func (r *trafficRepository) GetTotalTrafficByState(ctx context.Context, state string, startTime, endTime time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `SELECT
		date_trunc('day', time) as time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE state = $1 AND time BETWEEN $2 AND $3
		GROUP BY date_trunc('day', time)`

	err := r.db.SelectContext(ctx, &res, query, state, startTime, endTime)
	return res, err
}

func (r *trafficRepository) GetTotalTrafficByRegion(ctx context.Context, region string, startTime, endTime time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `SELECT
		date_trunc('day', time) as time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE region = $1 AND time BETWEEN $2 AND $3
		GROUP BY date_trunc('day', time)`

	err := r.db.SelectContext(ctx, &res, query, region, startTime, endTime)
	return res, err
}

func (r *trafficRepository) GetTotalTraffic(ctx context.Context, startTime, endTime time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `SELECT
		date_trunc('day', time) as time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE time BETWEEN $1 AND $2
		GROUP BY date_trunc('day', time)`

	err := r.db.SelectContext(ctx, &res, query, startTime, endTime)
	return res, err
}

func (r *trafficRepository) GetTrafficGroupedByRegion(ctx context.Context, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error) {
	var rows []entity.TrafficByRegion
	query := `SELECT
		region,
		date_trunc('day', time) as time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE time BETWEEN $1 AND $2
		GROUP BY region, date_trunc('day', time)`

	err := r.db.SelectContext(ctx, &rows, query, startTime, endTime)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]entity.TrafficSummary)
	for _, row := range rows {
		traffic := entity.TrafficSummary{
			Time:          row.Time,
			TotalBpsIn:    row.TotalBpsIn,
			TotalBpsOut:   row.TotalBpsOut,
			TotalBytesIn:  row.TotalBytesIn,
			TotalBytesOut: row.TotalBytesOut,
		}
		result[row.Region] = append(result[row.Region], traffic)
	}

	return result, nil
}

func (r *trafficRepository) GetTrafficGroupedByState(ctx context.Context, region string, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error) {
	var rows []entity.TrafficByState
	query := `SELECT
		state,
		date_trunc('day', time) as time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE region = $1 AND time BETWEEN $2 AND $3
		GROUP BY state, date_trunc('day', time)`

	err := r.db.SelectContext(ctx, &rows, query, region, startTime, endTime)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]entity.TrafficSummary)
	for _, row := range rows {
		traffic := entity.TrafficSummary{
			Time:          row.Time,
			TotalBpsIn:    row.TotalBpsIn,
			TotalBpsOut:   row.TotalBpsOut,
			TotalBytesIn:  row.TotalBytesIn,
			TotalBytesOut: row.TotalBytesOut,
		}
		result[row.State] = append(result[row.State], traffic)
	}

	return result, nil
}

func (r *trafficRepository) GetTrafficGroupedByIP(ctx context.Context, state string, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error) {
	var rows []entity.TrafficByIP
	query := `SELECT
		ip,
		date_trunc('day', time) as time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE state = $1 AND time BETWEEN $2 AND $3
		GROUP BY ip, date_trunc('day', time)`

	err := r.db.SelectContext(ctx, &rows, query, state, startTime, endTime)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]entity.TrafficSummary)
	for _, row := range rows {
		traffic := entity.TrafficSummary{
			Time:          row.Time,
			TotalBpsIn:    row.TotalBpsIn,
			TotalBpsOut:   row.TotalBpsOut,
			TotalBytesIn:  row.TotalBytesIn,
			TotalBytesOut: row.TotalBytesOut,
		}
		result[row.IP] = append(result[row.IP], traffic)
	}

	return result, nil
}
