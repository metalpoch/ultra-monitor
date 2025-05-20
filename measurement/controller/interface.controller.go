package controller

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/common/usecase"
	"gorm.io/gorm"
)

type interfaceController struct {
	Usecase usecase.InterfaceUsecase
}

func NewInterfaceController(db *gorm.DB, telegram tracking.SmartModule) *interfaceController {
	return &interfaceController{*usecase.NewInterfaceUsecase(db, telegram)}
}

func (i interfaceController) ShowAllInterfaces(deviceID uint64, csv bool) error {
	interfaces, err := i.Usecase.GetAllByDevice(deviceID)
	if err != nil {
		return err
	}

	pretty := table.NewWriter()
	pretty.SetOutputMirror(os.Stdout)
	pretty.AppendHeader(table.Row{
		"ID",
		"IfIndex",
		"IfName",
		"IfDescr",
		"IfAlias",
		"Created at",
		"Updated at",
	})

	for _, i := range interfaces {
		pretty.AppendRow(table.Row{
			i.ID,
			i.IfIndex,
			i.IfName,
			i.IfDescr,
			i.IfAlias,
			i.CreatedAt.Local().Format(constants.FORMAT_DATE),
			i.UpdatedAt.Local().Format(constants.FORMAT_DATE),
		})
		pretty.AppendSeparator()
	}
	if csv {
		pretty.RenderCSV()
	} else {
		pretty.Render()
	}

	return nil
}
