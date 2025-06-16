package controller

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type OntController struct {
	Usecase usecase.OntUsecase
}

func NewOntController(db *sqlx.DB, cache *cache.Redis) *OntController {
	return &OntController{*usecase.NewOntUsecase(db, cache)}
}

func (ctrl *OntController) StatusSummary(date time.Time) error {
	return ctrl.Usecase.UpdateStatusSummary(
		dto.RangeDate{
			InitDate: time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()),
			EndDate:  time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), date.Location()),
		},
	)
}
