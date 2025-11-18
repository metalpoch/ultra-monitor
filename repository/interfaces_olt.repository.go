package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type InterfaceOltRepository interface {
	GetAll(ctx context.Context) ([]entity.InterfacesDetailedOlt, error)
	Update(ctx context.Context, ip, oltVerbose string) error
}

type interfaceOltRepository struct {
	db *sqlx.DB
}

func NewInterfaceOltRepository(db *sqlx.DB) *interfaceOltRepository {
	return &interfaceOltRepository{db}
}

func (r *interfaceOltRepository) GetAll(ctx context.Context) ([]entity.InterfacesDetailedOlt, error) {
	var olts []entity.InterfacesDetailedOlt
	query := `SELECT DISTINCT st.region, st.state, st.sysname, st.ip, inf.olt_verbose
	FROM summary_traffic AS st
	LEFT JOIN  interfaces_olt AS inf ON st.ip = inf.olt_ip
	ORDER BY st.region, st.state, inf.olt_verbose, st.ip;`
	err := r.db.SelectContext(ctx, &olts, query)
	return olts, err
}

func (r *interfaceOltRepository) Update(ctx context.Context, ip, oltVerbose string) error {
	query := `INSERT INTO interfaces_olt (olt_verbose, olt_ip)
	VALUES ($1, $2)
	ON CONFLICT (olt_ip)
	DO UPDATE SET olt_verbose = EXCLUDED.olt_verbose;`
	_, err := r.db.ExecContext(ctx, query, oltVerbose, ip)
	return err
}
