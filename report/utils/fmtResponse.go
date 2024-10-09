package utils

import "github.com/metalpoch/olt-blueprint/common/model"

func FatResponse(f *model.Fat) model.FatResponse {
	return model.FatResponse{
		ID:        f.ID,
		Fat:       f.Fat,
		Splitter:  f.Splitter,
		Address:   f.Address,
		Latitude:  f.Latitude,
		Longitude: f.Longitude,
		Interface: model.InterfaceLite{
			ID:        f.Interface.ID,
			IfIndex:   f.Interface.IfIndex,
			IfName:    f.Interface.IfName,
			IfDescr:   f.Interface.IfDescr,
			IfAlias:   f.Interface.IfAlias,
			CreatedAt: f.Interface.CreatedAt,
			UpdatedAt: f.Interface.UpdatedAt,
		},
		Device: model.DeviceLite{
			ID:          f.Interface.Device.ID,
			IP:          f.Interface.Device.IP,
			SysName:     f.Interface.Device.SysName,
			SysLocation: f.Interface.Device.SysLocation,
			IsAlive:     f.Interface.Device.IsAlive,
			LastCheck:   f.Interface.Device.LastCheck,
			CreatedAt:   f.Interface.Device.CreatedAt,
			UpdatedAt:   f.Interface.Device.UpdatedAt,
		},
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func ReportResponse(rp *model.Report) model.ReportResponse {
	return model.ReportResponse{
		Category:         rp.Category,
		OriginalFilename: rp.OriginalFilename,
		Filepath:         rp.Filepath,
		User: model.UserLite{
			ID:       rp.User.ID,
			Email:    rp.User.Email,
			Fullname: rp.User.Fullname,
		},
		CreatedAt: rp.CreatedAt,
		UpdatedAt: rp.UpdatedAt,
		DeletedAt: rp.DeletedAt.Time,
	}
}
