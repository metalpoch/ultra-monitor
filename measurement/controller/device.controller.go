package controller

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/common/usecase"
	"gorm.io/gorm"
)

type deviceController struct {
	Usecase usecase.DeviceUsecase
}

func newDeviceController(db *gorm.DB, telegram tracking.SmartModule) *deviceController {
	return &deviceController{Usecase: *usecase.NewDeviceUsecase(db, telegram)}
}

func AddDevice(db *gorm.DB, telegram tracking.SmartModule, device *model.AddDevice) error {
	return newDeviceController(db, telegram).Usecase.Add(device)
}

func ShowAllDevices(db *gorm.DB, telegram tracking.SmartModule, csv bool) ([]model.Device, error) {
	devices, err := newDeviceController(db, telegram).Usecase.GetAll()
	if err != nil {
		return nil, err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
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
		t.AppendRow(table.Row{
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
		t.AppendSeparator()
	}
	if csv {
		t.RenderCSV()
	} else {
		t.Render()
	}

	return nil, nil
}

func GetDeviceWithOIDRows(db *gorm.DB, telegram tracking.SmartModule) ([]*model.DeviceWithOID, error) {
	return newDeviceController(db, telegram).Usecase.GetDeviceWithOIDRows()
}

func UpdateDevice(db *gorm.DB, telegram tracking.SmartModule, id uint, device *model.AddDevice) error {
	return newDeviceController(db, telegram).Usecase.Update(id, device)
}

func DeleteDevice(db *gorm.DB, telegram tracking.SmartModule, id uint) error {
	return newDeviceController(db, telegram).Usecase.Delete(id)
}
