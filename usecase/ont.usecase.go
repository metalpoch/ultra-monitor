package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/trend"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
	"github.com/redis/go-redis/v9"
)

type OntUsecase struct {
	repo  repository.OntRepository
	cache *cache.Redis
}

func NewOntUsecase(db *sqlx.DB, cache *cache.Redis) *OntUsecase {
	return &OntUsecase{repository.NewOntRepository(db), cache}
}

func (use *OntUsecase) AllOntStatus(dates dto.RangeDate) ([]model.OntStatusCounts, error) {
	var status []model.OntStatusCounts
	key := fmt.Sprintf("ontStatus:%d:%d", dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := use.repo.AllOntStatus(context.Background(), dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.OntStatusCounts)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil

	} else if err != nil {
		return nil, err
	}

	return status, err
}

func (uc *OntUsecase) OntStatusByState(state string, dates dto.RangeDate) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByState:%s:%d:%d", state, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.GetOntStatusByState(context.Background(), state, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.OntStatusCountsByState)(s))
		}
		err = uc.cache.InsertOne(context.Background(), key, 12*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil

	} else if err != nil {
		return nil, err
	}

	return status, err
}

func (use *OntUsecase) OntStatusByOdn(state, odn string, dates dto.RangeDate) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByODN:%s:%s:%d:%d", state, odn, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := use.repo.GetOntStatusByODN(context.Background(), state, odn, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.OntStatusCountsByState)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil

	} else if err != nil {
		return nil, err
	}

	return status, err
}

func (uc *OntUsecase) OntStatusByOltIP(ip string, dates dto.RangeDate) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByOltIP:%s:%d:%d", ip, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.GetOntStatusByOltIP(context.Background(), ip, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.OntStatusCountsByState)(s))
		}
		err = uc.cache.InsertOne(context.Background(), key, 12*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil
	} else if err != nil {
		return nil, err
	}
	return status, err
}

func (uc *OntUsecase) OntStatusBySysname(sysname string, dates dto.RangeDate) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusBySysname:%s:%d:%d", sysname, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.GetOntStatusBySysname(context.Background(), sysname, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.OntStatusCountsByState)(s))
		}
		err = uc.cache.InsertOne(context.Background(), key, 12*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil
	} else if err != nil {
		return nil, err
	}
	return status, err
}

func (use *OntUsecase) TrafficOnt(ponID uint64, idx string, dates dto.RangeDate) ([]model.TrafficOnt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var traffic []model.TrafficOnt
	res, err := use.repo.TrafficOnt(ctx, ponID, idx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}
	for _, s := range res {
		traffic = append(traffic, (model.TrafficOnt)(s))
	}

	return traffic, err
}

func (uc *OntUsecase) AllOntStatusForecast(dates dto.RangeDate, futureDays int) (*model.OntStatusForecast, error) {
	res, err := uc.AllOntStatus(dates)
	if err != nil {
		return nil, err
	}

	var historical []model.OntForecastBase
	for _, r := range res {
		historical = append(historical, model.OntForecastBase{
			Date:          r.Date,
			PonsCount:     r.PonsCount,
			ActiveCount:   r.ActiveCount,
			InactiveCount: r.InactiveCount,
			UnknownCount:  r.UnknownCount,
			TotalCount:    r.TotalCount,
		})
	}

	var actives, inactives, unknowns, totals []float64
	for _, s := range historical {
		actives = append(actives, float64(s.ActiveCount))
		inactives = append(inactives, float64(s.InactiveCount))
		unknowns = append(unknowns, float64(s.UnknownCount))
		totals = append(totals, float64(s.TotalCount))
	}

	predActives := trend.NewTrend(actives).Prediction(futureDays)
	predInactives := trend.NewTrend(inactives).Prediction(futureDays)
	predUnknowns := trend.NewTrend(unknowns).Prediction(futureDays)
	predTotals := trend.NewTrend(totals).Prediction(futureDays)

	var forecast []model.OntForecastBase
	var lastDay time.Time
	if len(historical) > 0 {
		lastDay = historical[len(historical)-1].Date
	} else {
		lastDay = dates.EndDate
	}
	for i := 1; i <= futureDays; i++ {
		forecast = append(forecast, model.OntForecastBase{
			Date:          lastDay.AddDate(0, 0, i),
			ActiveCount:   uint64(predActives[i-1]),
			InactiveCount: uint64(predInactives[i-1]),
			UnknownCount:  uint64(predUnknowns[i-1]),
			TotalCount:    uint64(predTotals[i-1]),
		})
	}

	return &model.OntStatusForecast{
		Historical: historical,
		Forecast:   forecast,
	}, nil
}

func (uc *OntUsecase) OntStatusByStateForecast(state string, dates dto.RangeDate, futureDays int) (*model.OntStatusForecast, error) {
	res, err := uc.OntStatusByState(state, dates)
	if err != nil {
		return nil, err
	}

	var historical []model.OntForecastBase
	for _, r := range res {
		historical = append(historical, model.OntForecastBase{
			Date:          r.Date,
			PonsCount:     r.PonsCount,
			ActiveCount:   r.ActiveCount,
			InactiveCount: r.InactiveCount,
			UnknownCount:  r.UnknownCount,
			TotalCount:    r.TotalCount,
		})
	}

	var actives, inactives, unknowns, totals []float64
	for _, s := range historical {
		actives = append(actives, float64(s.ActiveCount))
		inactives = append(inactives, float64(s.InactiveCount))
		unknowns = append(unknowns, float64(s.UnknownCount))
		totals = append(totals, float64(s.TotalCount))
	}

	predActives := trend.NewTrend(actives).Prediction(futureDays)
	predInactives := trend.NewTrend(inactives).Prediction(futureDays)
	predUnknowns := trend.NewTrend(unknowns).Prediction(futureDays)
	predTotals := trend.NewTrend(totals).Prediction(futureDays)

	var forecast []model.OntForecastBase
	var lastDay time.Time
	if len(historical) > 0 {
		lastDay = historical[len(historical)-1].Date
	} else {
		lastDay = dates.EndDate
	}
	for i := 1; i <= futureDays; i++ {
		forecast = append(forecast, model.OntForecastBase{
			Date:          lastDay.AddDate(0, 0, i),
			ActiveCount:   uint64(predActives[i-1]),
			InactiveCount: uint64(predInactives[i-1]),
			UnknownCount:  uint64(predUnknowns[i-1]),
			TotalCount:    uint64(predTotals[i-1]),
		})
	}

	return &model.OntStatusForecast{
		Historical: historical,
		Forecast:   forecast,
	}, nil
}

func (uc *OntUsecase) OntStatusByODNForecast(state, odn string, dates dto.RangeDate, futureDays int) (*model.OntStatusForecast, error) {
	res, err := uc.OntStatusByOdn(state, odn, dates)
	if err != nil {
		return nil, err
	}

	var historical []model.OntForecastBase
	for _, r := range res {
		historical = append(historical, model.OntForecastBase{
			Date:          r.Date,
			PonsCount:     r.PonsCount,
			ActiveCount:   r.ActiveCount,
			InactiveCount: r.InactiveCount,
			UnknownCount:  r.UnknownCount,
			TotalCount:    r.TotalCount,
		})
	}

	var actives, inactives, unknowns, totals []float64
	for _, s := range historical {
		actives = append(actives, float64(s.ActiveCount))
		inactives = append(inactives, float64(s.InactiveCount))
		unknowns = append(unknowns, float64(s.UnknownCount))
		totals = append(totals, float64(s.TotalCount))
	}

	predActives := trend.NewTrend(actives).Prediction(futureDays)
	predInactives := trend.NewTrend(inactives).Prediction(futureDays)
	predUnknowns := trend.NewTrend(unknowns).Prediction(futureDays)
	predTotals := trend.NewTrend(totals).Prediction(futureDays)

	var forecast []model.OntForecastBase
	var lastDay time.Time
	if len(historical) > 0 {
		lastDay = historical[len(historical)-1].Date
	} else {
		lastDay = dates.EndDate
	}
	for i := 1; i <= futureDays; i++ {
		forecast = append(forecast, model.OntForecastBase{
			Date:          lastDay.AddDate(0, 0, i),
			ActiveCount:   uint64(predActives[i-1]),
			InactiveCount: uint64(predInactives[i-1]),
			UnknownCount:  uint64(predUnknowns[i-1]),
			TotalCount:    uint64(predTotals[i-1]),
		})
	}

	return &model.OntStatusForecast{
		Historical: historical,
		Forecast:   forecast,
	}, nil
}

func (use *OntUsecase) TrafficOntByDespt(despt string, dates dto.RangeDate) ([]model.TrafficOnt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
