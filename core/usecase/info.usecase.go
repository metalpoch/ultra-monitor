package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/core/repository"
	"github.com/metalpoch/olt-blueprint/core/utils"
	"gorm.io/gorm"
)

type InfoUsecase struct {
	repo     repository.InfoRepository
	telegram tracking.SmartModule
}

func NewInfoUsecase(db *gorm.DB, telegram tracking.SmartModule) *InfoUsecase {
	return &InfoUsecase{repository.NewInfoRepository(db), telegram}
}

func (use InfoUsecase) GetDevice(id uint64) (*model.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetDevice(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDevice - use.repo.GetDevice(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	return (*model.Device)(res), err
}

func (use InfoUsecase) GetDeviceByIP(ip string) (*model.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetDeviceByIP(ctx, ip)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceByIP - use.repo.GetDeviceByIP(ctx, %s)", ip),
			err,
		)
		return nil, err
	}

	return (*model.Device)(res), err
}

func (use InfoUsecase) GetDeviceBySysname(sysname string) (*model.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetDeviceBySysname(ctx, sysname)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceBySysname - use.repo.GetDeviceBySysname(ctx, %s)", sysname),
			err,
		)
		return nil, err
	}

	return (*model.Device)(res), err
}

func (use InfoUsecase) GetAllDevice() ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAllDevice(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetAllDevice - use.repo.GetAllDevice(ctx, %d)",
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use InfoUsecase) GetDeviceByState(state string) ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetDeviceByState(ctx, state)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceByState - use.repo.GetDeviceByState(ctx, %s)", state),
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use InfoUsecase) GetDeviceByCounty(state, county string) ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetDeviceByCounty(ctx, state, county)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceByCounty - use.repo.GetDeviceByCounty(ctx, %s, %s)", state, county),
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use InfoUsecase) GetDeviceByMunicipality(state, county, municipality string) ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetDeviceByMunicipality(ctx, state, county, municipality)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceByMunicipality - use.repo.GetDeviceByMunicipality(ctx, %s, %s, %s)", state, county, municipality),
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use InfoUsecase) GetInterface(id uint64) (*model.Interface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetInterface(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetInterface - use.repo.GetInterface(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	return (*model.Interface)(res), err
}

func (use InfoUsecase) GetInterfacesByDevice(id uint) ([]*model.InterfaceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetInterfacesByDevice(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetInterfacesByDevice - use.repo.GetInterface(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	var interfaces []*model.InterfaceLite
	for _, i := range res {
		interfaces = append(interfaces, &model.InterfaceLite{
			ID:        i.ID,
			IfIndex:   i.IfIndex,
			IfName:    i.IfName,
			IfDescr:   i.IfDescr,
			IfAlias:   i.IfAlias,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		})
	}

	return interfaces, err
}

func (use InfoUsecase) GetInterfacesByDeviceAndPorts(id uint, shell, card, port *uint8) ([]*model.InterfaceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pattern := fmt.Sprintf("%%GPON %d", *shell)
	if card != nil {
		pattern = fmt.Sprintf("%s/%d", pattern, *card)
	}
	if port != nil {
		pattern = fmt.Sprintf("%s/%d", pattern, *port)
	}
	pattern += "%%"

	res, err := use.repo.GetInterfacesByDeviceAndPorts(ctx, id, pattern)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetInterfacesByDeviceAndPorts - use.repo.GetInterfacesByDeviceAndPorts(ctx, %d, %s)", id, pattern),
			err,
		)
		return nil, err
	}

	var interfaces []*model.InterfaceLite
	for _, i := range res {
		interfaces = append(interfaces, &model.InterfaceLite{
			ID:        i.ID,
			IfIndex:   i.IfIndex,
			IfName:    i.IfName,
			IfDescr:   i.IfDescr,
			IfAlias:   i.IfAlias,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		})
	}

	return interfaces, err
}

func (use InfoUsecase) GetLocationStates() ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetLocationStates(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetLocationStates - use.repo.GetLocationStates(ctx)",
			err,
		)
		return nil, err
	}

	return res, err
}

func (use InfoUsecase) GetLocationCounties(state string) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetLocationCounties(ctx, state)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetLocationCounties - use.repo.GetLocationCounties(ctx, %s)", state),
			err,
		)
		return nil, err
	}

	return res, err
}

func (use InfoUsecase) GetLocationMunicipalities(state, county string) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetLocationMunicipalities(ctx, state, county)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetLocationMunicipalities - use.repo.GetLocationMunicipalities(ctx, %s, %s)", state, county),
			err,
		)
		return nil, err
	}

	return res, err
}

func (use InfoUsecase) GetODN(odn string) ([]*model.FatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	res, err := use.repo.GetODN(ctx, odn)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetODN - use.repo.GetODN(ctx, %d)",
			err,
		)
		return nil, err
	}

	var fats []*model.FatResponse

	for _, f := range res {
		fat, err := use.repo.GetFat(ctx, f.FatID)
		if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				"(feedUsecase).GetODN - use.repo.GetODN(ctx, %d)",
				err,
			)
			return nil, err
		}
		inter, err := use.repo.GetInterface(ctx, f.InterfaceID)
		if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				"(feedUsecase).GetODN - use.repo.GetODN(ctx, %d)",
				err,
			)
			return nil, err
		}
		device, err := use.repo.GetDevice(ctx, inter.DeviceID)
		if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				"(feedUsecase).GetODN - use.repo.GetODN(ctx, %d)",
				err,
			)
			return nil, err
		}
		location, err := use.repo.GetLocation(ctx, fat.LocationID)
		if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				"(feedUsecase).GetODN - use.repo.GetODN(ctx, %d)",
				err,
			)
			return nil, err
		}
		temp := utils.FatResponse2(inter, fat, device, location)

		fats = append(fats, temp)
	}

	return fats, err
}

func (use InfoUsecase) GetODNStates(state string) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetODNStates(ctx, state)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetODNStates - use.repo.GetODNStates(ctx, %d)",
			err,
		)
		return nil, err
	}
	return res, err
}

func (use InfoUsecase) GetODNStatesContries(state, country string) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetODNStatesContries(ctx, state, country)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetODNStatesContries - use.repo.GetODNStatesContries(ctx, %d)",
			err,
		)
		return nil, err
	}
	return res, err
}

func (use InfoUsecase) GetODNStatesContriesMunicipality(state, country, municipality string) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetODNStatesContriesMunicipality(ctx, state, country, municipality)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetODNStatesContriesMunicipality - use.repo.GetODNStatesContriesMunicipality(ctx, %d)",
			err,
		)
		return nil, err
	}
	return res, err
}

func (use InfoUsecase) GetODNDevice(id uint) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetODNDevice(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetODNDevice - use.repo.GetODNDevice(ctx, %d)",
			err,
		)
		return nil, err
	}
	return res, err
}

func (use InfoUsecase) GetODNDevicePort(id uint, shell *uint8, card, port *uint8) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pattern := fmt.Sprintf("%%GPON %d", *shell)
	if card != nil {
		pattern = fmt.Sprintf("%s/%d", pattern, *card)
	}
	if port != nil {
		pattern = fmt.Sprintf("%s/%d", pattern, *port)
	}
	pattern += "%%"

	res, err := use.repo.GetODNDevicePort(ctx, id, pattern)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetODNDevicePort - use.repo.GetODNDevicePort(ctx, %d, %s)", id, pattern),
			err,
		)
		return nil, err
	}

	return res, err
}
