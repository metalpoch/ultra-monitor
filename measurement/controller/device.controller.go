package controller

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/common/usecase"
	"github.com/metalpoch/olt-blueprint/measurement/internal/snmp"
	"gorm.io/gorm"
)

type deviceController struct {
	Usecase usecase.DeviceUsecase
}

func NewDeviceController(db *gorm.DB, telegram tracking.SmartModule) *deviceController {
	return &deviceController{Usecase: *usecase.NewDeviceUsecase(db, telegram)}
}

func (d deviceController) AddDevice(device *model.AddDevice) error {
	var isAlive bool
	info, err := snmp.NewSnmp(snmp.Config{
		IP:        device.IP,
		Community: device.Community,
		Timeout:   5 * time.Second,
	}).OltSysQuery()
	if err == nil {
		isAlive = true
	}

	newDevice := &model.Device{
		IP:          device.IP,
		SysName:     info.SysName,
		SysLocation: info.SysLocation,
		Community:   device.Community,
		IsAlive:     isAlive,
		TemplateID:  device.Template,
		LastCheck:   time.Now(),
	}

	return d.Usecase.Add(newDevice)
}

func (d deviceController) ShowAllDevices(db *gorm.DB, telegram tracking.SmartModule, csv bool) ([]model.Device, error) {
	devices, err := d.Usecase.GetAll()
	if err != nil {
		return nil, err
	}

	pretty := table.NewWriter()
	pretty.SetOutputMirror(os.Stdout)
	pretty.AppendHeader(table.Row{
		"ID",
		"IP",
		"Community",
		"SysName",
		"SysLocation",
		"Template ID",
		"Is Alive",
		"Last Check",
		"Created at",
		"Updated at",
	})

	for _, device := range devices {
		pretty.AppendRow(table.Row{
			device.ID,
			device.IP,
			device.Community,
			device.SysName,
			device.SysLocation,
			device.TemplateID,
			device.IsAlive,
			device.LastCheck.Local().Format(constants.FORMAT_DATE),
			device.CreatedAt.Local().Format(constants.FORMAT_DATE),
			device.UpdatedAt.Local().Format(constants.FORMAT_DATE),
		})
		pretty.AppendSeparator()
	}
	if csv {
		pretty.RenderCSV()
	} else {
		pretty.Render()
	}

	return nil, nil
}

func (d deviceController) GetDeviceWithOIDRows() ([]*model.DeviceWithOID, error) {
	return d.Usecase.GetDeviceWithOIDRows()
}

func (d deviceController) UpdateDevice(id uint, device *model.AddDevice) error {
	return d.Usecase.Update(id, device)
}

func (d deviceController) DeleteDevice(id uint64) error {
	return d.Usecase.Delete(id)
}
