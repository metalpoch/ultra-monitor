package controller

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/constants"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type PonController struct {
	Usecase usecase.PonUsecase
}

func NewPonController(db *sqlx.DB, cache *cache.Redis) *PonController {
	return &PonController{*usecase.NewPonUsecase(db, cache)}
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

func (ctrl PonController) SummaryTraffic(date time.Time) error {
	return ctrl.Usecase.UpdateSummaryTraffic(
		dto.RangeDate{
			InitDate: time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()),
			EndDate:  time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), date.Location()),
		},
	)
}
