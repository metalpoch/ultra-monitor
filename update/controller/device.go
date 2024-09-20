package controller

import (
	"database/sql"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/usecase"
)

type deviceHandler struct {
	Usecase usecase.DeviceUsecase
}

func newDeviceHandler(db *sql.DB) *deviceHandler {
	return &deviceHandler{Usecase: *usecase.NewDeviceUsecase(db)}
}

func AddDevice(db *sql.DB, device *model.AddDevice) error {
	return newDeviceHandler(db).Usecase.Add(device)
}

func ShowAllDevices(db *sql.DB, csv bool) error {
	devices, err := newDeviceHandler(db).Usecase.GetAll()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"ID",
		"IP",
		"Community",
		"Sysname",
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
			device.Sysname,
			device.TemplateID,
			device.IsAlive,
			device.LastCheck.Local().Format("2006-01-02 15:04:05"),
			device.CreatedAt.Local().Format("2006-01-02 15:04:05"),
			device.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
		})
		t.AppendSeparator()
	}
	if csv {
		t.RenderCSV()
	} else {
		t.Render()
	}

	return nil
}
