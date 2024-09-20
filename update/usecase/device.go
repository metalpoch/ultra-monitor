package usecase

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/pkg/snmp"
	"github.com/metalpoch/olt-blueprint/update/repository"
)

type deviceUsecase struct {
	repo repository.DeviceRepository
}

func NewDeviceUsecase(db *sql.DB) *deviceUsecase {
	return &deviceUsecase{repository.NewDeviceRepository(db)}
}

func (use deviceUsecase) Add(device *model.AddDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var isAlive bool
	sysname, err := snmp.GetSysname(device.IP, device.Community)
	if err != nil {
		log.Println("snmp error on try get the sysname:", err)
		isAlive = false
	} else {
		isAlive = true
	}

	now := time.Now()
	use.repo.Add(ctx, &entity.Device{
		ID:         0,
		IP:         device.IP,
		Sysname:    sysname,
		Community:  device.Community,
		TemplateID: device.TemplateID,
		IsAlive:    isAlive,
		LastCheck:  now,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	return nil

}

func (use deviceUsecase) GetAll() ([]model.Device, error) {
	devices := []model.Device{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)

	// Gestionar errores (con Axios por ejemplo)
	// ...

	for _, e := range res {
		devices = append(devices, model.Device{
			ID:         e.ID,
			IP:         e.IP,
			Sysname:    e.Sysname,
			Community:  e.Community,
			TemplateID: e.TemplateID,
			IsAlive:    e.IsAlive,
			LastCheck:  e.LastCheck,
			CreatedAt:  e.CreatedAt,
			UpdatedAt:  e.UpdatedAt,
		})
	}

	return devices, err
}
