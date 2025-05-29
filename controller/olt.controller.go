package controller

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/internal/constants"
	"github.com/metalpoch/olt-blueprint/internal/dto"
	"github.com/metalpoch/olt-blueprint/model"
	"github.com/metalpoch/olt-blueprint/usecase"
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

func (ctrl OltController) ShowAllDevices(csv bool) ([]model.Olt, error) {
	olts, err := ctrl.Usecase.Olts(constants.DEFAULT_PAGE, constants.DEFAULT_LIMIT)
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
		"Is Alive",
		"Last Check",
		"Created at",
	})

	for _, olt := range olts {
		pretty.AppendRow(table.Row{
			olt.ID,
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

func (d OltController) Update(id uint64, olt dto.NewOlt) error {
	return d.Usecase.Update(id, olt)
}

func (d OltController) Delete(id uint64) error {
	return d.Usecase.Delete(id)
}
