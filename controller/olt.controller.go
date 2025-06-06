package controller

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/constants"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type OltController struct {
	Usecase usecase.OltUsecase
}

func NewOltController(db *sqlx.DB) *OltController {
	return &OltController{Usecase: *usecase.NewOltUsecase(db)}
}

func (ctrl OltController) AddOlt(olt dto.NewOlt) error {
	return ctrl.Usecase.Add(olt)
}

func (ctrl OltController) ShowAllOlt(csv bool) ([]model.Olt, error) {
	olts, err := ctrl.Usecase.Olts()
	if err != nil {
		return nil, err
	}

	pretty := table.NewWriter()
	pretty.SetOutputMirror(os.Stdout)
	pretty.AppendHeader(table.Row{
		"IP",
		"Community",
		"SysName",
		"SysLocation",
		"Is Alive",
		"Last Check",
		"Created at",
	})

	for _, olt := range olts {
		pretty.AppendRow(table.Row{
			olt.IP,
			olt.Community,
			olt.SysName,
			olt.SysLocation,
			olt.IsAlive,
			olt.LastCheck.Local().Format(constants.FORMAT_DATE),
			olt.CreatedAt.Local().Format(constants.FORMAT_DATE),
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

func (d OltController) Delete(id string) error {
	return d.Usecase.Delete(id)
}
