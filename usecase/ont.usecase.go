package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/snmp"
	"github.com/metalpoch/ultra-monitor/repository"
)

type OntUsecase struct {
	repo repository.OntRepository
}

func NewOntUsecase(db *sqlx.DB) *OntUsecase {
	return &OntUsecase{repository.NewOntRepository(db)}
}

func (use *OntUsecase) SnmpSerialDespt(device dto.AllSerialDesptByPon) ([]dto.OntSerialsAndDespts, error) {
	snmpClient := snmp.NewSnmp(snmp.Config{
		IP:        device.IP,
		Community: device.Community,
		Timeout:   10 * time.Second,
		Retries:   3,
	})

	ontData, err := snmpClient.OntSerialsAndDespts(device.PonIdx)
	if err != nil {
		return nil, fmt.Errorf("SNMP query failed: %v", err)
	}

	var res []dto.OntSerialsAndDespts
	for idx, data := range ontData {
		res = append(res, dto.OntSerialsAndDespts{
			PonIdx:       device.PonIdx,
			OntIdx:       uint8(idx),
			Despt:        data.Despt,
			SerialNumber: data.SerialNumber,
		})
	}

	return res, nil
}

func (use *OntUsecase) SaveFromSNMP(device dto.CreateOntRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	snmpClient := snmp.NewSnmp(snmp.Config{
		IP:        device.IP,
		Community: device.Community,
		Timeout:   10 * time.Second,
		Retries:   3,
	})

	ontData, err := snmpClient.OntQuery(device.PonIdx, device.OntIdx)
	if err != nil {
		return fmt.Errorf("SNMP query failed: %v", err)
	}

	return use.repo.Create(ctx, entity.Ont{
		IP:          device.IP,
		OntIDX:      fmt.Sprintf("%d.%d", device.PonIdx, device.OntIdx),
		Serial:      ontData.SerialNumber,
		Despt:       ontData.Despt,
		LineProf:    ontData.LineProfName,
		Description: device.Description,
		Status:      ontData.ControlRunStatus == 1,
		LastCheck:   time.Now(),
	})
}

func (use *OntUsecase) GetByID(id int32) (dto.OntResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ont, err := use.repo.GetByID(ctx, id)
	if err != nil {
		return dto.OntResponse{}, err
	}

	return dto.OntResponse{
		ID:          ont.ID,
		IP:          ont.IP,
		OntIDX:      ont.OntIDX,
		Serial:      ont.Serial,
		Despt:       ont.Despt,
		LineProf:    ont.LineProf,
		Description: ont.Description,
		Enabled:     ont.Enabled,
		Status:      ont.Status,
		LastCheck:   ont.LastCheck.Format(time.RFC3339),
		CreatedAt:   ont.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (use *OntUsecase) GetAll() ([]dto.OntResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	onts, err := use.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.OntResponse
	for _, ont := range onts {
		responses = append(responses, dto.OntResponse{
			ID:          ont.ID,
			IP:          ont.IP,
			OntIDX:      ont.OntIDX,
			Serial:      ont.Serial,
			Despt:       ont.Despt,
			LineProf:    ont.LineProf,
			Description: ont.Description,
			Enabled:     ont.Enabled,
			Status:      ont.Status,
			LastCheck:   ont.LastCheck.Format(time.RFC3339),
			CreatedAt:   ont.CreatedAt.Format(time.RFC3339),
		})
	}

	return responses, nil
}

func (use *OntUsecase) Delete(id int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return use.repo.Delete(ctx, id)
}

func (use *OntUsecase) Enable(id int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return use.repo.Enable(ctx, id)
}

func (use *OntUsecase) Disable(id int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return use.repo.Disable(ctx, id)
}
