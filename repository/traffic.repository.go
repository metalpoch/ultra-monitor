package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/dto"
)

type TrafficRepository interface {
	GetSnmpIndexByMunicipality(ctx context.Context, state, municipality string) ([]entity.OltIndex, error)
	GetSnmpIndexByCounty(ctx context.Context, state, municipality, county string) ([]entity.OltIndex, error)
	GetSnmpIndexByODN(ctx context.Context, state, municipality, odn string) ([]entity.OltIndex, error)
	SaveSummaryTraffic(ctx context.Context, trafficData []entity.SumaryTraffic) error
	GetTotalTrafficByIP(ctx context.Context, ip string, startTime, endTime time.Time) ([]entity.TrafficSummary, error)
	GetTotalTrafficByState(ctx context.Context, state string, startTime, endTime time.Time) ([]entity.TrafficSummary, error)
	GetTotalTrafficByRegion(ctx context.Context, region string, startTime, endTime time.Time) ([]entity.TrafficSummary, error)
	GetTotalTraffic(ctx context.Context, startTime, endTime time.Time) ([]entity.TrafficSummary, error)
	GetTrafficGroupedByRegion(ctx context.Context, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error)
	GetTrafficGroupedByState(ctx context.Context, region string, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error)
	GetTrafficGroupedByIP(ctx context.Context, state string, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error)
	GetLocationHierarchy(ctx context.Context, initDate, finalDate time.Time) (*dto.LocationHierarchy, error)
	GetOntTraffic(ctx context.Context, ontID int32, startTime, endTime time.Time) ([]entity.OntTraffic, error)
	GetRealBandwidthByIP(ctx context.Context, ip string) (float64, error)
	GetRealBandwidthByState(ctx context.Context, state string, initDate, finalDate time.Time) (float64, error)
	GetRealBandwidthByRegion(ctx context.Context, region string, initDate, finalDate time.Time) (float64, error)
	GetSwitchByIP(ctx context.Context, ip string) (string, error)
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
	query := `INSERT INTO summary_traffic (time, ip, state, region, sysname, bps_in, bps_out, bytes_in, bytes_out, volume_in, volume_out)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (ip, time)
		DO UPDATE SET
			state = EXCLUDED.state,
			region = EXCLUDED.region,
			sysname = EXCLUDED.sysname,
			bps_in = EXCLUDED.bps_in,
			bps_out = EXCLUDED.bps_out,
			bytes_in = EXCLUDED.bytes_in,
			bytes_out = EXCLUDED.bytes_out,
			volume_in = EXCLUDED.volume_in,
			volume_out = EXCLUDED.volume_out`

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
			traffic.Sysname,
			traffic.BpsIn,
			traffic.BpsOut,
			traffic.BytesIn,
			traffic.BytesOut,
			traffic.VolumeIn,
			traffic.VolumeOut)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *trafficRepository) GetTotalTrafficByIP(ctx context.Context, ip string, startTime, endTime time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `SELECT
		date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas' AS time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out,
		SUM(volume_in) as total_volume_in,
		SUM(volume_out) as total_volume_out
		FROM summary_traffic
		WHERE ip = $1 AND time BETWEEN $2 AND $3
		GROUP BY date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas'
		ORDER BY time`

	err := r.db.SelectContext(ctx, &res, query, ip, startTime, endTime)
	return res, err
}

func (r *trafficRepository) GetTotalTrafficByState(ctx context.Context, state string, startTime, endTime time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `SELECT
		date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas' AS time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out,
		SUM(volume_in) as total_volume_in,
		SUM(volume_out) as total_volume_out
		FROM summary_traffic
		WHERE state = $1 AND time BETWEEN $2 AND $3
		GROUP BY date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas'
		ORDER BY time`

	err := r.db.SelectContext(ctx, &res, query, state, startTime, endTime)
	return res, err
}

func (r *trafficRepository) GetTotalTrafficByRegion(ctx context.Context, region string, startTime, endTime time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `SELECT
		date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas' AS time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out,
		SUM(volume_in) as total_volume_in,
		SUM(volume_out) as total_volume_out
		FROM summary_traffic
		WHERE region = $1 AND time BETWEEN $2 AND $3
		GROUP BY date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas'
		ORDER BY time`

	err := r.db.SelectContext(ctx, &res, query, region, startTime, endTime)
	return res, err
}

func (r *trafficRepository) GetTotalTraffic(ctx context.Context, startTime, endTime time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `SELECT
		date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas' AS time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out,
		SUM(volume_in) as total_volume_in,
		SUM(volume_out) as total_volume_out
		FROM summary_traffic
		WHERE time BETWEEN $1 AND $2
		GROUP BY date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas'
		ORDER BY time`

	err := r.db.SelectContext(ctx, &res, query, startTime, endTime)
	return res, err
}

func (r *trafficRepository) GetTrafficGroupedByRegion(ctx context.Context, startTime, endTime time.Time) (map[string][]entity.TrafficSummary, error) {
	var rows []entity.TrafficByRegion
	query := `SELECT
		region,
		date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas' AS time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE time BETWEEN $1 AND $2
		GROUP BY region, date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas'
		ORDER BY time`

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
		date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas' AS time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE region = $1 AND time BETWEEN $2 AND $3
		GROUP BY state, date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas'
		ORDER BY time`

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
	var rows []entity.TrafficBySysname
	query := `SELECT
		sysname,
		date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas' AS time,
		SUM(bps_in) as total_bps_in,
		SUM(bps_out) as total_bps_out,
		SUM(bytes_in) as total_bytes_in,
		SUM(bytes_out) as total_bytes_out
		FROM summary_traffic
		WHERE state = $1 AND time BETWEEN $2 AND $3
		GROUP BY sysname, date_trunc('day', time AT TIME ZONE 'America/Caracas') AT TIME ZONE 'America/Caracas'
		ORDER BY sysname, time`

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
		result[row.Sysname] = append(result[row.Sysname], traffic)
	}

	return result, nil
}

func (r *trafficRepository) GetLocationHierarchy(ctx context.Context, initDate, finalDate time.Time) (*dto.LocationHierarchy, error) {
	initDate = time.Date(initDate.Year(), initDate.Month(), initDate.Day(), 0, 0, 0, 0, initDate.Location())
	finalDate = time.Date(finalDate.Year(), finalDate.Month(), finalDate.Day(), 23, 59, 59, 0, finalDate.Location())
	hierarchy := &dto.LocationHierarchy{
		Regions: []string{},
		States:  make(map[string][]string),
		Olts:    make(map[string][]dto.OltInfo),
	}

	// Get unique regions
	var regions []string
	queryRegions := `SELECT DISTINCT region FROM summary_traffic WHERE time BETWEEN $1 AND $2 ORDER BY region`
	err := r.db.SelectContext(ctx, &regions, queryRegions, initDate, finalDate)
	if err != nil {
		return nil, err
	}
	hierarchy.Regions = regions

	// Get unique states by region
	var stateRows []struct {
		Region string `db:"region"`
		State  string `db:"state"`
	}
	queryStates := `SELECT DISTINCT region, state FROM summary_traffic WHERE time BETWEEN $1 AND $2 ORDER BY region, state`
	err = r.db.SelectContext(ctx, &stateRows, queryStates, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	for _, row := range stateRows {
		hierarchy.States[row.Region] = append(hierarchy.States[row.Region], row.State)
	}

	// Get OLTs by state
	var oltRows []struct {
		State   string `db:"state"`
		IP      string `db:"ip"`
		SysName string `db:"sysname"`
	}
	queryOlts := `SELECT DISTINCT state, ip, sysname FROM summary_traffic WHERE time BETWEEN $1 AND $2 ORDER BY state, ip`
	err = r.db.SelectContext(ctx, &oltRows, queryOlts, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	for _, row := range oltRows {
		oltInfo := dto.OltInfo{
			IP:      row.IP,
			SysName: row.SysName,
		}
		hierarchy.Olts[row.State] = append(hierarchy.Olts[row.State], oltInfo)
	}

	return hierarchy, nil
}

func (r *trafficRepository) GetRealBandwidthByIP(ctx context.Context, ip string) (float64, error) {
	var bandwidth float64
	query := `
		SELECT COALESCE(SUM(ib.bandwidth), 0) as total_bandwidth
		FROM interfaces_bandwidth ib
		JOIN interfaces_olt io ON ib.olt_verbose = io.olt_verbose
		WHERE io.olt_ip = $1
	`
	err := r.db.GetContext(ctx, &bandwidth, query, ip)
	return bandwidth, err
}

func (r *trafficRepository) GetRealBandwidthByState(ctx context.Context, state string, initDate, finalDate time.Time) (float64, error) {
	var bandwidth float64
	query := `
		SELECT COALESCE(SUM(ib.bandwidth), 0) as total_bandwidth
		FROM interfaces_bandwidth ib
		JOIN interfaces_olt io ON ib.olt_verbose = io.olt_verbose
		WHERE io.olt_ip IN (
			SELECT DISTINCT ip FROM summary_traffic
			WHERE state = $1 AND time BETWEEN $2 AND $3
		)
	`
	err := r.db.GetContext(ctx, &bandwidth, query, state, initDate, finalDate)
	return bandwidth, err
}

func (r *trafficRepository) GetRealBandwidthByRegion(ctx context.Context, region string, initDate, finalDate time.Time) (float64, error) {
	var bandwidth float64
	query := `
		SELECT COALESCE(SUM(ib.bandwidth), 0) as total_bandwidth
		FROM interfaces_bandwidth ib
		JOIN interfaces_olt io ON ib.olt_verbose = io.olt_verbose
		WHERE io.olt_ip IN (
			SELECT DISTINCT ip FROM summary_traffic
			WHERE region = $1 AND time BETWEEN $2 AND $3
		)
	`
	err := r.db.GetContext(ctx, &bandwidth, query, region, initDate, finalDate)
	return bandwidth, err
}

func (r *trafficRepository) GetOntTraffic(ctx context.Context, id int32, initDate, finalDate time.Time) ([]entity.OntTraffic, error) {
	var traffic []entity.OntTraffic
	query := `SELECT * FROM onts_traffic WHERE ont_id = $1 AND time BETWEEN $2 AND $3 AND (bps_in <= 2.49e9 OR bps_out <= 2.49e9) ORDER BY time`
	err := r.db.SelectContext(ctx, &traffic, query, id, initDate, finalDate)
	return traffic, err
}

func (r *trafficRepository) GetSwitchByIP(ctx context.Context, ip string) (string, error) {
	var interfaceName string
	query := `
		SELECT ib.interface
		FROM interfaces_bandwidth ib
		JOIN interfaces_olt io ON ib.olt_verbose = io.olt_verbose
		WHERE io.olt_ip = $1
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &interfaceName, query, ip)
	if err != nil {
		return "", err
	}

	// Extract switch prefix from interface name
	return extractSwitchFromInterface(interfaceName), nil
}


