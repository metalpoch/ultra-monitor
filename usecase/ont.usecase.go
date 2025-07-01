package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/trend"
	"github.com/metalpoch/ultra-monitor/internal/utils"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
)

type OntUsecase struct {
	repo  repository.OntRepository
	cache *cache.Redis
}

func NewOntUsecase(db *sqlx.DB, cache *cache.Redis) *OntUsecase {
	return &OntUsecase{repository.NewOntRepository(db), cache}
}

func (uc *OntUsecase) GetStatusIPSummary(ip string, dates dto.RangeDate) ([]model.OntSummaryStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetStatusIPSummary(ctx, ip, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var summary []model.OntSummaryStatus
	for _, r := range res {
		summary = append(summary, (model.OntSummaryStatus)(r))
	}
	return summary, nil
}

func (uc *OntUsecase) GetStatusSysnameSummary(sysname string, dates dto.RangeDate) ([]model.OntSummaryStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetStatusSysnameSummary(ctx, sysname, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var summary []model.OntSummaryStatus
	for _, r := range res {
		summary = append(summary, (model.OntSummaryStatus)(r))
	}
	return summary, nil
}

func (uc *OntUsecase) GetStatusStateSummary(dates dto.RangeDate) ([]model.GetStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetStatusStateSummary(ctx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var summary []model.GetStatusSummary
	for _, r := range res {
		summary = append(summary, (model.GetStatusSummary)(r))
	}
	return summary, nil
}

func (uc *OntUsecase) GetStatusByStateSummary(state string, dates dto.RangeDate) ([]model.OntSummaryStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetStatusByStateSummary(ctx, state, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var summary []model.OntSummaryStatus
	for _, r := range res {
		summary = append(summary, (model.OntSummaryStatus)(r))
	}
	return summary, nil
}

func (uc *OntUsecase) GetStatusMunicipalitySummary(state string, dates dto.RangeDate) ([]model.GetStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetStatusMunicipalitySummary(ctx, state, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var summary []model.GetStatusSummary
	for _, r := range res {
		summary = append(summary, (model.GetStatusSummary)(r))
	}
	return summary, nil
}

func (uc *OntUsecase) GetStatusCountySummary(state, municipality string, dates dto.RangeDate) ([]model.GetStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetStatusCountySummary(ctx, state, municipality, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var summary []model.GetStatusSummary
	for _, r := range res {
		summary = append(summary, (model.GetStatusSummary)(r))
	}
	return summary, nil
}

func (uc *OntUsecase) GetStatusOdnSummary(state, municipality, county string, dates dto.RangeDate) ([]model.GetStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetStatusOdnSummary(ctx, state, municipality, county, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var summary []model.GetStatusSummary
	for _, r := range res {
		summary = append(summary, (model.GetStatusSummary)(r))
	}
	return summary, nil
}

func (uc *OntUsecase) GetStatusSummaryForecast(futureDays int) (*model.OntStatusForecast, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	initDate, endDate := utils.DateRangeFromYear()
	res, err := uc.repo.GetStatusSummary(ctx, initDate, endDate)
	if err != nil {
		return nil, err
	}

	var historical []model.OntSummaryStatus
	for _, r := range res {
		historical = append(historical, (model.OntSummaryStatus)(r))
	}

	var actives, inactives, unknowns []float64
	for _, s := range historical {
		actives = append(actives, float64(s.ActiveCount))
		inactives = append(inactives, float64(s.InactiveCount))
		unknowns = append(unknowns, float64(s.UnknownCount))
	}

	predActives := trend.NewTrend(actives).Prediction(futureDays)
	predInactives := trend.NewTrend(inactives).Prediction(futureDays)
	predUnknowns := trend.NewTrend(unknowns).Prediction(futureDays)

	var forecast []model.OntSummaryStatus
	var lastDay time.Time
	if len(historical) > 0 {
		lastDay = historical[len(historical)-1].Day
	} else {
		lastDay = endDate
	}
	for i := 1; i <= futureDays; i++ {
		forecast = append(forecast, model.OntSummaryStatus{
			Day:           lastDay.AddDate(0, 0, i),
			ActiveCount:   uint64(predActives[i-1]),
			InactiveCount: uint64(predInactives[i-1]),
			UnknownCount:  uint64(predUnknowns[i-1]),
		})
	}

	return &model.OntStatusForecast{
		Historical: historical,
		Forecast:   forecast,
	}, nil
}

func (uc *OntUsecase) GetStatusByStateForecast(state string, futureDays int) (*model.OntStatusForecast, error) {
	initDate, endDate := utils.DateRangeFromYear()
	res, err := uc.GetStatusByStateSummary(state, dto.RangeDate{InitDate: initDate, EndDate: endDate})
	if err != nil {
		return nil, err
	}

	var historical []model.OntSummaryStatus
	for _, r := range res {
		historical = append(historical, (model.OntSummaryStatus)(r))
	}

	var actives, inactives, unknowns []float64
	for _, s := range historical {
		actives = append(actives, float64(s.ActiveCount))
		inactives = append(inactives, float64(s.InactiveCount))
		unknowns = append(unknowns, float64(s.UnknownCount))
	}

	predActives := trend.NewTrend(actives).Prediction(futureDays)
	predInactives := trend.NewTrend(inactives).Prediction(futureDays)
	predUnknowns := trend.NewTrend(unknowns).Prediction(futureDays)

	var forecast []model.OntSummaryStatus
	var lastDay time.Time
	if len(historical) > 0 {
		lastDay = historical[len(historical)-1].Day
	} else {
		lastDay = endDate
	}
	for i := 1; i <= futureDays; i++ {
		forecast = append(forecast, model.OntSummaryStatus{
			Day:           lastDay.AddDate(0, 0, i),
			ActiveCount:   uint64(predActives[i-1]),
			InactiveCount: uint64(predInactives[i-1]),
			UnknownCount:  uint64(predUnknowns[i-1]),
		})
	}

	return &model.OntStatusForecast{
		Historical: historical,
		Forecast:   forecast,
	}, nil
}

func (use *OntUsecase) TrafficOnt(ponID int, idx int64, dates dto.RangeDate) ([]model.TrafficOnt, error) {
	if !utils.IsDateRangeWithin7Days(dates.InitDate, dates.EndDate) {
		return nil, fmt.Errorf("the date range invalor or cannot be greater than 7 days")
	}

	var traffic []model.TrafficOnt
	res, err := use.repo.TrafficOnt(context.Background(), ponID, idx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}
	for _, s := range res {
		traffic = append(traffic, (model.TrafficOnt)(s))
	}

	return traffic, err
}

func (use *OntUsecase) TrafficOntByDespt(despt string, dates dto.RangeDate) ([]model.TrafficOnt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var traffic []model.TrafficOnt
	res, err := use.repo.TrafficOntByDespt(ctx, despt, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}
	for _, s := range res {
		traffic = append(traffic, (model.TrafficOnt)(s))
	}

	return traffic, err
}

func (uc *OntUsecase) UpdateStatusSummary(dates dto.RangeDate) error {
	res, err := uc.repo.GetDailyAveragedHourlyStatusSummary(context.Background(), dates.InitDate, dates.EndDate)
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return sql.ErrNoRows
	}

	return uc.repo.UpdateStatusSummary(context.Background(), res)
}
