package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type FatRepository interface {
	AllInfo(ctx context.Context, page, limit uint16) ([]entity.Fat, error)
	UpsertFats(ctx context.Context, fats []entity.UpsertFat) (int64, error)
	DeleteOne(ctx context.Context, id int32) error
	FindByID(ctx context.Context, id int32) (entity.Fat, error)
	FindByStates(ctx context.Context, state string, page, limit uint16) ([]entity.Fat, error)
	FindByMunicipality(ctx context.Context, state, municipality string, page, limit uint16) ([]entity.Fat, error)
	FindByCounty(ctx context.Context, state, municipality, county string, page, limit uint16) ([]entity.Fat, error)
	FindBytOdn(ctx context.Context, state, municipality, county, odn string, page, limit uint16) ([]entity.Fat, error)
}

type fatRepository struct {
	db *sqlx.DB
}

func NewFatRepository(db *sqlx.DB) *fatRepository {
	return &fatRepository{db}
}

func (r *fatRepository) AllInfo(ctx context.Context, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats ORDER BY region, state, municipality, county LIMIT $1 OFFSET $2;`
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

		// Creamos una clave única para el mapa basada en los campos de conflicto.
		// `strconv.Itoa` convierte el int a string.
		key := fmt.Sprintf("%d-%s", id, fat.Date.Format("2006-01-02"))

		// Si la clave ya existe, tomamos el valor más reciente o el que prefieras.
		// Aquí, simplemente sobrescribimos el valor existente.
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

	// Convertimos el mapa de vuelta a un slice para la inserción masiva.
	fatStatusesDeduplicated := make([]entity.FatStatus, 0, len(fatStatusesMap))
	for _, status := range fatStatusesMap {
		fatStatusesDeduplicated = append(fatStatusesDeduplicated, status)
	}

	// Dividimos el slice deduplicado en lotes para evitar el error de 65535 parámetros.
	batchSize := 10000

	queryFatStatus := `
        INSERT INTO fat_status (fats_id, date, actives, provisioned_offline, cut_off, in_progress)
        VALUES (:fats_id, :date, :actives, :provisioned_offline, :cut_off, :in_progress)
        ON CONFLICT (fats_id, date) DO UPDATE SET
            actives = EXCLUDED.actives,
            provisioned_offline = EXCLUDED.provisioned_offline,
            cut_off = EXCLUDED.cut_off,
            in_progress = EXCLUDED.in_progress;`

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

// func (r *fatRepository) UpsertFats(ctx context.Context, fats []entity.UpsertFat) (int64, error) {
// 	tx, err := r.db.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	queryFats := `
//         INSERT INTO fats (ip, region, state, municipality, county, odn, fat, shell, card, port)
//         VALUES (:ip, :region, :state, :municipality, :county, :odn, :fat, :shell, :card, :port)
//         ON CONFLICT (ip, region, state, municipality, county, odn, fat, shell, card, port) DO NOTHING
//         RETURNING id;`

// 	queryGetFatID := `
//         SELECT id FROM fats
//         WHERE ip = :ip AND fat = :fat AND shell = :shell AND card = :card AND port = :port;`

// 	stmtInsert, err := tx.PrepareNamedContext(ctx, queryFats)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer stmtInsert.Close()

// 	stmtGetID, err := tx.PrepareNamedContext(ctx, queryGetFatID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer stmtGetID.Close()

// 	var fatStatuses []entity.FatStatus
// 	var totalProcessed int64

// 	for _, fat := range fats {
// 		var id int32

// 		err := stmtInsert.GetContext(ctx, &id, fat)
// 		if err != nil {
// 			if err.Error() == "sql: no rows in result set" {
// 				err = stmtGetID.GetContext(ctx, &id, fat)
// 				if err != nil {
// 					return totalProcessed, err
// 				}
// 			} else {
// 				return totalProcessed, err
// 			}
// 		}

// 		fatStatuses = append(fatStatuses, entity.FatStatus{
// 			FatsID:             id,
// 			Date:               fat.Date,
// 			Actives:            fat.Actives,
// 			ProvisionedOffline: fat.ProvisionedOffline,
// 			CutOff:             fat.CutOff,
// 			InProgress:         fat.InProgress,
// 		})
// 		totalProcessed++
// 	}

// 	if len(fatStatuses) > 0 {
// 		queryFatStatus := `
//             INSERT INTO fat_status (fats_id, date, actives, provisioned_offline, cut_off, in_progress)
//             VALUES (:fats_id, :date, :actives, :provisioned_offline, :cut_off, :in_progress)
//             ON CONFLICT (fats_id, date) DO UPDATE SET
//                 actives = EXCLUDED.actives,
//                 provisioned_offline = EXCLUDED.provisioned_offline,
//                 cut_off = EXCLUDED.cut_off,
//                 in_progress = EXCLUDED.in_progress;`

// 		_, err = tx.NamedExecContext(ctx, queryFatStatus, fatStatuses)
// 		if err != nil {
// 			return totalProcessed, err
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return totalProcessed, err
// 	}

// 	return totalProcessed, nil
// }

func (r *fatRepository) DeleteOne(ctx context.Context, id int32) error {
	query := `DELETE FROM fats WHERE id = $1;`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *fatRepository) FindByID(ctx context.Context, id int32) (entity.Fat, error) {
	var fat entity.Fat
	query := `SELECT * FROM fats WHERE id = $1;`
	err := r.db.GetContext(ctx, &fat, query, id)
	if err != nil {
		return entity.Fat{}, err
	}
	return fat, nil
}

func (r *fatRepository) FindByStates(ctx context.Context, state string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 ORDER BY region, state, municipality, county LIMIT $2 OFFSET $3;`
	err := r.db.SelectContext(ctx, &res, query, state, limit, offset)
	return res, err
}

func (r *fatRepository) FindByMunicipality(ctx context.Context, state, municipality string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 ORDER BY region, state, municipality, county LIMIT $3 OFFSET $4;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, limit, offset)
	return res, err
}

func (r *fatRepository) FindByCounty(ctx context.Context, state, municipality, county string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY region, state, municipality, county LIMIT $4 OFFSET $5;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, limit, offset)
	return res, err
}

func (r *fatRepository) FindBytOdn(ctx context.Context, state, municipality, county, odn string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 AND odn = $4 ORDER BY region, state, municipality, county LIMIT $5 OFFSET $6;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, odn, limit, offset)
	return res, err
}
