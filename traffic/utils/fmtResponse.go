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
		TemplateID:  device.TemplateID,
		LastCheck:   device.LastCheck,
		CreatedAt:   device.CreatedAt,
		UpdatedAt:   device.UpdatedAt,
	}
}
