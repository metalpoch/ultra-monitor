package utils

import (
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
)

func FatResponse(fi *entity.FatInterface) *model.FatResponse {
	return &model.FatResponse{
		ID:        fi.FatID,
		ODN:       fi.Fat.ODN,
		Fat:       fi.Fat.Fat,
		Splitter:  fi.Fat.Splitter,
		Address:   fi.Fat.Address,
		Latitude:  fi.Fat.Latitude,
		Longitude: fi.Fat.Longitude,
		Interface: model.InterfaceLite{
			ID:        fi.Interface.ID,
			IfIndex:   fi.Interface.IfIndex,
			IfName:    fi.Interface.IfName,
			IfDescr:   fi.Interface.IfDescr,
			IfAlias:   fi.Interface.IfAlias,
			CreatedAt: fi.Interface.CreatedAt,
			UpdatedAt: fi.Interface.UpdatedAt,
		},
		Device: model.DeviceLite{
			ID:          fi.Interface.Device.ID,
			IP:          fi.Interface.Device.IP,
			SysName:     fi.Interface.Device.SysName,
			SysLocation: fi.Interface.Device.SysLocation,
			IsAlive:     fi.Interface.Device.IsAlive,
			LastCheck:   fi.Interface.Device.LastCheck,
			CreatedAt:   fi.Interface.Device.CreatedAt,
			UpdatedAt:   fi.Interface.Device.UpdatedAt,
		},
		Location: model.Location{
			ID:           fi.Fat.Location.ID,
			State:        fi.Fat.Location.State,
			County:       fi.Fat.Location.County,
			Municipality: fi.Fat.Location.Municipality,
		},
		CreatedAt: fi.Fat.CreatedAt,
		UpdatedAt: fi.Fat.UpdatedAt,
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
