package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/entity"
)

type trafficRepository struct {
	db *sql.DB
}

func NewTrafficRepository(db *sql.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (repo trafficRepository) UpsertInterface(ctx context.Context, element *entity.Interface) (uint, error) {
	var id uint
	if err := repo.db.QueryRowContext(ctx,
		constants.SQL_UPSERT_INTERFACE,
		element.ID,
		element.IfName,
		element.IfDescr,
		element.IfAlias,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo trafficRepository) GetAllInterfaces(ctx context.Context) ([]*entity.Interface, error) {
	rows, err := repo.db.QueryContext(ctx, constants.SQL_SELECT_ALL_INTERFACES)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	interfaces := []*entity.Interface{}
	for rows.Next() {
		i := new(entity.Interface)
		rows.Scan(
			&i.ID,
			&i.IfName,
			&i.IfDescr,
			&i.IfAlias,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		interfaces = append(interfaces, i)
	}

	return interfaces, nil
}

func (repo trafficRepository) GetInterface(ctx context.Context, id uint) (*entity.Interface, error) {
	i := new(entity.Interface)
	err := repo.db.QueryRowContext(ctx, constants.SQL_SELECT_INTERFACE, id).Scan(
		&i.ID,
		&i.IfName,
		&i.IfDescr,
		&i.IfAlias,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (repo trafficRepository) SaveMeasurement(ctx context.Context, tmp bool, measurement entity.Measurement) error {
	query := constants.SQL_INSERT_MEASUREMENT
	if tmp {
		query = constants.SQL_INSERT_TMP_MEASUREMENT
	}

	result, err := repo.db.ExecContext(ctx,
		query,
		measurement.Date,
		measurement.Bw,
		measurement.In,
		measurement.Out,
		measurement.InterfaceID,
	)

	if err != nil {
		return err
	}

	qty, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if qty < 1 {
		return errors.New("no row was inserted")
	}

	return nil
}

func (repo trafficRepository) GetTmpMeasurement(ctx context.Context, measurement entity.Measurement) (*entity.Measurement, error) {
	m := new(entity.Measurement)
	err := repo.db.QueryRowContext(ctx, constants.SQL_SELECT_TMP_MEASUREMENT).Scan(
		&m.Date,
		&m.Bw,
		&m.In,
		&m.Out,
		&m.InterfaceID,
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}
