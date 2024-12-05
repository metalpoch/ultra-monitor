package utils

import (
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
)

func DeviceResponse(device *model.Device) model.DeviceResponse {
	return model.DeviceResponse{
		ID:          device.ID,
		IP:          device.IP,
		Community:   device.Community,
		SysName:     device.SysName,
		SysLocation: device.SysLocation,
		IsAlive:     device.IsAlive,
		Template:    model.Template(device.Template),
		LastCheck:   device.LastCheck,
		CreatedAt:   device.CreatedAt,
		UpdatedAt:   device.UpdatedAt,
	}
}

func InterfaceResponse(i *model.Interface) model.InterfaceResponse {
	return model.InterfaceResponse{
		ID:      i.ID,
		IfIndex: i.IfIndex,
		IfName:  i.IfName,
		IfDescr: i.IfDescr,
		IfAlias: i.IfAlias,
		Device: model.DeviceLite{
			ID:          i.Device.ID,
			IP:          i.Device.IP,
			SysName:     i.Device.SysName,
			SysLocation: i.Device.SysLocation,
			IsAlive:     i.Device.IsAlive,
			LastCheck:   i.Device.LastCheck,
			CreatedAt:   i.Device.CreatedAt,
			UpdatedAt:   i.Device.UpdatedAt,
		},
		Template:  model.Template(i.Device.Template),
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

// func FatResponse2(i *entity.Interface, f *entity.Fat, d *entity.Device, l *entity.Location) *model.FatResponse {
// 	return &model.FatResponse{
// 		ID:        f.ID,
// 		ODN:       f.ODN,
// 		Fat:       f.Fat,
// 		Splitter:  f.Splitter,
// 		Address:   f.Address,
// 		Latitude:  f.Latitude,
// 		Longitude: f.Longitude,
// 		Interface: model.InterfaceLite{
// 			ID:        i.ID,
// 			IfIndex:   i.IfIndex,
// 			IfName:    i.IfName,
// 			IfDescr:   i.IfDescr,
// 			IfAlias:   i.IfAlias,
// 			CreatedAt: i.CreatedAt,
// 			UpdatedAt: i.UpdatedAt,
// 		},
// 		Device: model.DeviceLite{
// 			ID:          d.ID,
// 			IP:          d.IP,
// 			SysName:     d.SysName,
// 			SysLocation: d.SysLocation,
// 			IsAlive:     d.IsAlive,
// 			LastCheck:   d.LastCheck,
// 			CreatedAt:   d.CreatedAt,
// 			UpdatedAt:   d.UpdatedAt,
// 		},
// 		Location: model.Location{
// 			ID:           l.ID,
// 			State:        l.State,
// 			County:       l.County,
// 			Municipality: l.Municipality,
// 		},
// 		CreatedAt: f.CreatedAt,
// 		UpdatedAt: f.UpdatedAt,
// 	}
// }

func FatResponse21(fi *entity.FatInterface, f *entity.Fat) *model.Fat {
	return &model.Fat{
		ID:           f.ID,
		ODN:          f.ODN,
		Fat:          f.Fat,
		Splitter:     f.Splitter,
		Address:      f.Address,
		Latitude:     f.Latitude,
		Longitude:    f.Longitude,
		LocationID:   f.LocationID,
		FatInterface: fi.ID,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}
}
