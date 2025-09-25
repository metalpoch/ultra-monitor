package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type FatRepository interface {
	AllInfo(ctx context.Context, find, value string, page, limit uint16) ([]entity.FatInfoStatus, error)
	UpsertFats(ctx context.Context, fats []entity.UpsertFat) (int64, error)
	FindByID(ctx context.Context, id int32) (entity.FatInfoStatus, error)
	GetAllByIP(ctx context.Context, ip string) ([]entity.FatInfoStatus, error)
	GetFieldsOptions(ctx context.Context, field string) ([]string, error)

	FindByStates(ctx context.Context, state string, page, limit uint16) ([]entity.FatInfoStatus, error)
	FindByMunicipality(ctx context.Context, state, municipality string, page, limit uint16) ([]entity.FatInfoStatus, error)
	FindByCounty(ctx context.Context, state, municipality, county string, page, limit uint16) ([]entity.FatInfoStatus, error)
	FindByOdn(ctx context.Context, state, municipality, county, odn string, page, limit uint16) ([]entity.FatInfoStatus, error)

	GetRegions(ctx context.Context) ([]string, error)
	GetStates(ctx context.Context) ([]string, error)
	GetMunicipalities(ctx context.Context, state string) ([]string, error)
	GetCounties(ctx context.Context, state, municipality string) ([]string, error)
	GetODN(ctx context.Context, state, municipality, county string) ([]string, error)

	GetAllFatStatus(ctx context.Context) ([]entity.FatStatusSummary, error)
	GetAllFatStatusByRegion(ctx context.Context, region string) ([]entity.FatStatusSummary, error)
	GetAllFatStatusByState(ctx context.Context, state string) ([]entity.FatStatusSummary, error)
	GetAllFatStatusByMunicipality(ctx context.Context, state, municipality string) ([]entity.FatStatusSummary, error)
	GetAllFatStatusByCounty(ctx context.Context, state, municipality, county string) ([]entity.FatStatusSummary, error)
	GetAllFatStatusByODN(ctx context.Context, state, municipality, county, odn string) ([]entity.FatStatusSummary, error)
	GetAllFatStatusByFat(ctx context.Context, state, municipality, county, odn, fat string) ([]entity.FatStatusSummary, error)

	GetFatStatusStateByRegion(ctx context.Context, region string) ([]entity.LastFatStatus, error)
	GetFatStatusOltByState(ctx context.Context, state string) ([]entity.LastFatStatus, error)
	GetFatStatusGponByOlt(ctx context.Context, olt string) ([]entity.LastFatStatus, error)
}

type fatRepository struct {
	db *sqlx.DB
}

func NewFatRepository(db *sqlx.DB) *fatRepository {
	return &fatRepository{db}
}

func (r *fatRepository) AllInfo(ctx context.Context, field, value string, page, limit uint16) ([]entity.FatInfoStatus, error) {
	var findField string
	var res []entity.FatInfoStatus

	if field != "" || value != "" {
		findField = fmt.Sprintf("f.%s = '%s' AND ", field, value)
	}

	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT
		f.*,
		fs.date,
		fs.actives,
		fs.provisioned_offline,
		fs.cut_off,
		fs.in_progress
	FROM fats AS f
	INNER JOIN fat_status AS fs ON fs.fats_id = f.id
	WHERE %s fs.date = (SELECT MAX(date) FROM fat_status)
	ORDER BY f.region, f.state, f.municipality, f.county
	LIMIT $1 OFFSET $2;`, findField)

	err := r.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (r *fatRepository) UpsertFats(ctx context.Context, fats []entity.UpsertFat) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryFats := `
        INSERT INTO fats (ip, region, state, municipality, county, odn, fat, shell, card, port)
        VALUES (:ip, :region, :state, :municipality, :county, :odn, :fat, :shell, :card, :port)
        ON CONFLICT (ip, region, state, municipality, county, odn, fat, shell, card, port) DO NOTHING
        RETURNING id;`

	queryGetFatID := `
        SELECT id FROM fats
        WHERE
			ip = :ip
			AND	region = :region
			AND	state = :state
			AND municipality = :municipality
			AND county = :county
			AND odn = :odn
			AND fat = :fat
			AND	shell = :shell
			AND card = :card
			AND port = :port;`

	stmtInsert, err := tx.PrepareNamedContext(ctx, queryFats)
	if err != nil {
		return 0, err
	}
	defer stmtInsert.Close()

	stmtGetID, err := tx.PrepareNamedContext(ctx, queryGetFatID)
	if err != nil {
		return 0, err
	}
	defer stmtGetID.Close()

	fatStatusesMap := make(map[string]entity.FatStatus)
	var totalProcessed int64

	for _, fat := range fats {
		var id int32

		err := stmtInsert.GetContext(ctx, &id, fat)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				err = stmtGetID.GetContext(ctx, &id, fat)
				if err != nil {
					return totalProcessed, err
				}
			} else {
				return totalProcessed, err
			}
		}

		key := fmt.Sprintf("%d-%s", id, fat.Date.Format("2006-01-02"))
		fatStatusesMap[key] = entity.FatStatus{
			FatsID:             id,
			Date:               fat.Date,
			Actives:            fat.Actives,
			ProvisionedOffline: fat.ProvisionedOffline,
			CutOff:             fat.CutOff,
			InProgress:         fat.InProgress,
		}
		totalProcessed++
	}

	if len(fatStatusesMap) == 0 {
		return 0, tx.Commit()
	}

	fatStatusesDeduplicated := make([]entity.FatStatus, 0, len(fatStatusesMap))
	for _, status := range fatStatusesMap {
		fatStatusesDeduplicated = append(fatStatusesDeduplicated, status)
	}

	queryFatStatus := `
	INSERT INTO fat_status (fats_id, date, actives, provisioned_offline, cut_off, in_progress)
	VALUES (:fats_id, :date, :actives, :provisioned_offline, :cut_off, :in_progress)
	ON CONFLICT (fats_id, date) DO UPDATE SET
		actives = EXCLUDED.actives,
		provisioned_offline = EXCLUDED.provisioned_offline,
		cut_off = EXCLUDED.cut_off,
		in_progress = EXCLUDED.in_progress;`

	batchSize := 10000
	for i := 0; i < len(fatStatusesDeduplicated); i += batchSize {
		end := i + batchSize
		if end > len(fatStatusesDeduplicated) {
			end = len(fatStatusesDeduplicated)
		}

		batch := fatStatusesDeduplicated[i:end]
		_, err = tx.NamedExecContext(ctx, queryFatStatus, batch)
		if err != nil {
			return totalProcessed, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return totalProcessed, err
	}

	return totalProcessed, nil
}

func (r *fatRepository) FindByID(ctx context.Context, id int32) (entity.FatInfoStatus, error) {
	var fat entity.FatInfoStatus
	query := `
	SELECT
		f.*,
		fs.date,
		fs.actives,
		fs.provisioned_offline,
		fs.cut_off,
		fs.in_progress
	FROM fats AS f
	INNER JOIN fat_status AS fs ON fs.fats_id = f.id
	WHERE f.id = $1 AND fs.date = (SELECT MAX(date) FROM fat_status);
	`
	err := r.db.GetContext(ctx, &fat, query, id)
	if err != nil {
		return entity.FatInfoStatus{}, err
	}
	return fat, nil
}

func (r *fatRepository) GetAllByIP(ctx context.Context, ip string) ([]entity.FatInfoStatus, error) {
	var fats []entity.FatInfoStatus
	query := `
	SELECT
		f.*,
		fs.date,
		fs.actives,
		fs.provisioned_offline,
		fs.cut_off,
		fs.in_progress
	FROM fats AS f
	INNER JOIN fat_status AS fs ON fs.fats_id = f.id
	WHERE f.ip = $1 AND fs.date = (SELECT MAX(date) FROM fat_status)
	ORDER BY f.odn, f.fat, f.shell, f.card, f.port;`

	err := r.db.SelectContext(ctx, &fats, query, ip)
	return fats, err
}

func (r *fatRepository) GetRegions(ctx context.Context) ([]string, error) {
	var res []string
	err := r.db.SelectContext(ctx, &res, `SELECT DISTINCT region FROM fats;`)
	return res, err
}

func (r *fatRepository) GetStates(ctx context.Context) ([]string, error) {
	var res []string
	err := r.db.SelectContext(ctx, &res, "SELECT DISTINCT state FROM fats ORDER BY state ASC;")
	return res, err
}

func (r *fatRepository) GetMunicipalities(ctx context.Context, state string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT municipality FROM fats WHERE state = $1 ORDER BY municipality ASC;`
	err := r.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (r *fatRepository) GetCounties(ctx context.Context, state, municipality string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT county FROM fats WHERE state = $1 AND municipality = $2 ORDER BY county ASC;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (r *fatRepository) GetODN(ctx context.Context, state, municipality, county string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY odn ASC;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (r *fatRepository) FindByStates(ctx context.Context, state string, page, limit uint16) ([]entity.FatInfoStatus, error) {
	var res []entity.FatInfoStatus
	offset := (page - 1) * limit
	query := `
	SELECT
		f.*,
		fs.date,
		fs.actives,
		fs.provisioned_offline,
		fs.cut_off,
		fs.in_progress
	FROM fats AS f
	INNER JOIN fat_status AS fs ON fs.fats_id = f.id
	WHERE f.state = $1 AND fs.date = (SELECT MAX(date) FROM fat_status)
	ORDER BY f.region, f.state, f.municipality, f.county
	LIMIT $2 OFFSET $3;
	`
	err := r.db.SelectContext(ctx, &res, query, state, limit, offset)
	return res, err
}

func (r *fatRepository) FindByMunicipality(ctx context.Context, state, municipality string, page, limit uint16) ([]entity.FatInfoStatus, error) {
	var res []entity.FatInfoStatus
	offset := (page - 1) * limit
	query := `
	SELECT
		f.*,
		fs.date,
		fs.actives,
		fs.provisioned_offline,
		fs.cut_off,
		fs.in_progress
	FROM fats AS f
	INNER JOIN fat_status AS fs ON fs.fats_id = f.id
	WHERE f.state = $1 AND f.municipality = $2 AND fs.date = (SELECT MAX(date) FROM fat_status)
	ORDER BY f.region, f.state, f.municipality, f.county
	LIMIT $3 OFFSET $4;
	`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, limit, offset)
	return res, err
}

func (r *fatRepository) FindByCounty(ctx context.Context, state, municipality, county string, page, limit uint16) ([]entity.FatInfoStatus, error) {
	var res []entity.FatInfoStatus
	offset := (page - 1) * limit
	query := `
	SELECT
		f.*,
		fs.date,
		fs.actives,
		fs.provisioned_offline,
		fs.cut_off,
		fs.in_progress
	FROM fats AS f
	INNER JOIN fat_status AS fs ON fs.fats_id = f.id
	WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3 AND fs.date = (SELECT MAX(date) FROM fat_status)
	ORDER BY f.region, f.state, f.municipality, f.county
	LIMIT $4 OFFSET $5;
	`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, limit, offset)
	return res, err
}

func (r *fatRepository) FindByOdn(ctx context.Context, state, municipality, county, odn string, page, limit uint16) ([]entity.FatInfoStatus, error) {
	var res []entity.FatInfoStatus
	offset := (page - 1) * limit
	query := `
	SELECT
		f.*,
		fs.date,
		fs.actives,
		fs.provisioned_offline,
		fs.cut_off,
		fs.in_progress
	FROM fats AS f
	INNER JOIN fat_status AS fs ON fs.fats_id = f.id
	WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3 AND f.odn = $4 AND fs.date = (SELECT MAX(date) FROM fat_status)
	ORDER BY f.region, f.state, f.municipality, f.county
	LIMIT $5 OFFSET $6;
	`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, odn, limit, offset)
	return res, err
}

func (r *fatRepository) GetAllFatStatus(ctx context.Context) ([]entity.FatStatusSummary, error) {
	var res []entity.FatStatusSummary

	query := `
	SELECT
		date,
		SUM(actives) AS actives,
		SUM(provisioned_offline) AS provisioned_offline,
		SUM(cut_off) AS cut_off,
		SUM(in_progress) AS in_progress
	FROM fat_status
	GROUP BY date
	ORDER BY date ASC;`
	err := r.db.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *fatRepository) GetAllFatStatusByRegion(ctx context.Context, region string) ([]entity.FatStatusSummary, error) {
	var res []entity.FatStatusSummary

	states := constants.STATES_BY_REGION[region]
	query := `
        SELECT
            fs.date,
            SUM(fs.actives) AS actives,
            SUM(fs.provisioned_offline) AS provisioned_offline,
            SUM(fs.cut_off) AS cut_off,
            SUM(fs.in_progress) AS in_progress
        FROM fat_status AS fs
        INNER JOIN fats AS f ON f.id = fs.fats_id
        WHERE f.state = ANY($1)
        GROUP BY fs.date
        ORDER BY fs.date DESC;`
	err := r.db.SelectContext(ctx, &res, query, pq.Array(states))
	return res, err
}

func (r *fatRepository) GetAllFatStatusByState(ctx context.Context, state string) ([]entity.FatStatusSummary, error) {
	var res []entity.FatStatusSummary
	query := `
	SELECT
		fs.date,
		SUM(fs.actives) AS actives,
		SUM(fs.provisioned_offline) AS provisioned_offline,
		SUM(fs.cut_off) AS cut_off,
		SUM(fs.in_progress) AS in_progress
	FROM fat_status AS fs
	INNER JOIN fats AS f ON f.id = fs.fats_id
	WHERE f.state = $1
	GROUP BY fs.date
	ORDER BY fs.date DESC;`
	err := r.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (r *fatRepository) GetAllFatStatusByMunicipality(ctx context.Context, state, municipality string) ([]entity.FatStatusSummary, error) {
	var res []entity.FatStatusSummary
	query := `
	SELECT
		fs.date,
		SUM(fs.actives) AS actives,
		SUM(fs.provisioned_offline) AS provisioned_offline,
		SUM(fs.cut_off) AS cut_off,
		SUM(fs.in_progress) AS in_progress
	FROM fat_status AS fs
	INNER JOIN fats AS f ON f.id = fs.fats_id
	WHERE f.state = $1 AND f.municipality = $2
	GROUP BY fs.date
	ORDER BY fs.date DESC;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (r *fatRepository) GetAllFatStatusByCounty(ctx context.Context, state, municipality, county string) ([]entity.FatStatusSummary, error) {
	var res []entity.FatStatusSummary
	query := `
	SELECT
		fs.date,
		SUM(fs.actives) AS actives,
		SUM(fs.provisioned_offline) AS provisioned_offline,
		SUM(fs.cut_off) AS cut_off,
		SUM(fs.in_progress) AS in_progress
	FROM fat_status AS fs
	INNER JOIN fats AS f ON f.id = fs.fats_id
	WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3
	GROUP BY fs.date
	ORDER BY fs.date DESC;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (r *fatRepository) GetAllFatStatusByODN(ctx context.Context, state, municipality, county, odn string) ([]entity.FatStatusSummary, error) {
	var res []entity.FatStatusSummary
	query := `
	SELECT
		fs.date,
		SUM(fs.actives) AS actives,
		SUM(fs.provisioned_offline) AS provisioned_offline,
		SUM(fs.cut_off) AS cut_off,
		SUM(fs.in_progress) AS in_progress
	FROM fat_status AS fs
	INNER JOIN fats AS f ON f.id = fs.fats_id
	WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3 AND f.odn = $4
	GROUP BY fs.date
	ORDER BY fs.date DESC;
	`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, odn)
	return res, err
}

func (r *fatRepository) GetAllFatStatusByFat(ctx context.Context, state, municipality, county, odn, fat string) ([]entity.FatStatusSummary, error) {
	var res []entity.FatStatusSummary
	query := `
	SELECT
		fs.date,
		SUM(fs.actives) AS actives,
		SUM(fs.provisioned_offline) AS provisioned_offline,
		SUM(fs.cut_off) AS cut_off,
		SUM(fs.in_progress) AS in_progress
	FROM fat_status AS fs
	INNER JOIN fats AS f ON f.id = fs.fats_id
	WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3 AND f.odn = $4 AND f.fat = $5
	GROUP BY fs.date
	ORDER BY fs.date DESC;
	`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, odn, fat)
	return res, err
}

func (r *fatRepository) GetFatStatusStateByRegion(ctx context.Context, region string) ([]entity.LastFatStatus, error) {
	var res []entity.LastFatStatus
	states := constants.STATES_BY_REGION[region]
	query := `
        SELECT
    f.state AS name,
    SUM(fs.actives) AS actives,
    SUM(fs.provisioned_offline) AS provisioned_offline,
    SUM(fs.cut_off) AS cut_off,
    SUM(fs.in_progress) AS in_progress
FROM fat_status AS fs
INNER JOIN fats AS f ON f.id = fs.fats_id
INNER JOIN (
    SELECT f.state, MAX(fs.date) AS max_date
    FROM fat_status AS fs
    INNER JOIN fats AS f ON f.id = fs.fats_id
    WHERE f.state = ANY($1)
    GROUP BY f.state
) AS latest ON latest.state = f.state AND latest.max_date = fs.date
WHERE f.state = ANY($1)
GROUP BY name;`
	err := r.db.SelectContext(ctx, &res, query, pq.Array(states))
	return res, err
}

func (r *fatRepository) GetFatStatusOltByState(ctx context.Context, state string) ([]entity.LastFatStatus, error) {
	var res []entity.LastFatStatus
	query := `
	SELECT
		f.ip AS name,
		SUM(fs.actives) AS actives,
		SUM(fs.provisioned_offline) AS provisioned_offline,
		SUM(fs.cut_off) AS cut_off,
		SUM(fs.in_progress) AS in_progress
	FROM fat_status AS fs
	INNER JOIN fats AS f ON f.id = fs.fats_id
	INNER JOIN (
		SELECT f.ip, MAX(fs.date) AS max_date
		FROM fat_status AS fs
		INNER JOIN fats AS f ON f.id = fs.fats_id
		WHERE f.state = $1
		GROUP BY f.ip
	) AS latest ON latest.ip = f.ip AND latest.max_date = fs.date
	WHERE f.state = $1
	GROUP BY name;
	`
	err := r.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (r *fatRepository) GetFatStatusGponByOlt(ctx context.Context, ip string) ([]entity.LastFatStatus, error) {
	var res []entity.LastFatStatus
	query := `
	SELECT
		CONCAT('GPON ', f.shell, '/', f.card, '/', f.port) AS name,
		SUM(fs.actives) AS actives,
		SUM(fs.provisioned_offline) AS provisioned_offline,
		SUM(fs.cut_off) AS cut_off,
		SUM(fs.in_progress) AS in_progress
	FROM fat_status AS fs
	INNER JOIN fats AS f ON f.id = fs.fats_id
	INNER JOIN (
		SELECT f.shell, f.card, f.port, MAX(fs.date) AS max_date
		FROM fat_status AS fs
		INNER JOIN fats AS f ON f.id = fs.fats_id
		WHERE f.ip = $1
		GROUP BY f.shell, f.card, f.port
	) AS latest ON latest.shell = f.shell AND latest.card = f.card AND latest.port = f.port AND latest.max_date = fs.date
	WHERE f.ip = $1
	GROUP BY f.shell, f.card, f.port;
	`
	err := r.db.SelectContext(ctx, &res, query, ip)
	return res, err
}

func (r *fatRepository) GetFieldsOptions(ctx context.Context, field string) ([]string, error) {
	var res []string

	validFields := map[string]bool{
		"ip":           true,
		"region":       true,
		"state":        true,
		"municipality": true,
		"county":       true,
		"odn":          true,
	}
	if !validFields[field] {
		return nil, fmt.Errorf("invalid field: %s", field)
	}

	query := fmt.Sprintf(`SELECT DISTINCT %s FROM fats AS f ORDER BY %s;`, field)

	err := r.db.SelectContext(ctx, &res, query)
	return res, err
}
