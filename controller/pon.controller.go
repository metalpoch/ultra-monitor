package controller

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/constants"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type PonController struct {
	Usecase usecase.PonUsecase
}

func NewInterfaceController(db *sqlx.DB) *PonController {
	return &PonController{*usecase.NewPonUsecase(db)}
}

func (ctrl PonController) ShowAllInterfaces(sysname string, csv bool) error {
	interfaces, err := ctrl.Usecase.GetAllByDevice(sysname)
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
