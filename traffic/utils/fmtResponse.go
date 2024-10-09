package utils

import "github.com/metalpoch/olt-blueprint/common/model"

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
