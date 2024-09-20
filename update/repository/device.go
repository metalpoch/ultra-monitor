package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
)

type deviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *deviceRepository {
	return &deviceRepository{db}
}

func (repo deviceRepository) Add(ctx context.Context, device *entity.Device) (uint, error) {
	var id uint
	q := `INSERT INTO device
		(ip, sysname, community, template_id, is_alive, last_check, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id;`
	now := time.Now()
	if err := repo.db.QueryRowContext(
		ctx,
		q,
		device.IP,
		device.Sysname,
		device.Community,
		device.TemplateID,
		device.IsAlive,
		device.LastCheck,
		now,
		now,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo deviceRepository) GetAll(ctx context.Context) ([]*entity.Device, error) {
	devices := []*entity.Device{}
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM device;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := new(entity.Device)
		rows.Scan(
			&d.ID,
			&d.IP,
			&d.Sysname,
			&d.Community,
			&d.TemplateID,
			&d.IsAlive,
			&d.LastCheck,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		devices = append(devices, d)
	}

	return devices, nil
}
