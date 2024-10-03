package usecase

import "github.com/metalpoch/olt-blueprint/common/model"

type DeviceUsecase interface {
	Add(device *model.AddDevice) error
	Check(device *model.Device) error
	GetAll() ([]*model.Device, error)
	GetDeviceWithOIDRows() ([]*model.DeviceWithOID, error)
	Update(id uint, device *model.AddDevice) error
	Delete(id uint) error
}
