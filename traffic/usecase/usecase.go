package usecase

import "github.com/metalpoch/olt-blueprint/common/model"

type FeedUsecase interface {
	GetDevice(id uint) (*model.Device, error)
	GetAllDevice() ([]*model.DeviceLite, error)
	GetInterface(id uint) (*model.Interface, error)
	GetInterfacesByDevice(id uint) ([]*model.InterfaceWithoutDevice, error)
}
