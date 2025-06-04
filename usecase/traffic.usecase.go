package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/trend"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
	"github.com/redis/go-redis/v9"
)

type TrafficUsecase struct {
	repo  repository.TrafficRepository
	cache *cache.Redis
}

func NewTrafficUsecase(db *sqlx.DB, cache *cache.Redis) *TrafficUsecase {
	return &TrafficUsecase{repository.NewTrafficRepository(db), cache}
}

func (uc *TrafficUsecase) GetTotalTraffic(dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TotalTraffic(ctx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByState(state string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByState(ctx, state, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByCounty(state, county string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByCounty(ctx, state, county, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByMunicipaly(state, county, municipality string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByMunicipality(ctx, state, county, municipality, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByODN(state, odn string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByODN(ctx, state, odn, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}
func (uc *TrafficUsecase) TrafficByPon(sysname string, shell, card, port int, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ifname := fmt.Sprintf("GPON %d/%d/%d", shell, card, port)
	res, err := uc.repo.TrafficByPon(ctx, sysname, ifname, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) GetDailyAveragedHourlyMaxTrafficTrends(futureDays int, dates dto.RangeDate) (*model.TrafficTrendForecast, error) {
	var traffic []model.TrafficTrend
	key := fmt.Sprintf("averagedHourlyMaxTraffic:%d:%d", dates.InitDate.Unix(), dates.EndDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := uc.repo.GetDailyAveragedHourlyMaxTrafficTrends(context.Background(), dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}

		for _, t := range res {
			traffic = append(traffic, (model.TrafficTrend)(t))
		}

		if err := uc.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var mbpsIn, mbpsOut, mbytesIn, mbytesOut []float64
	for _, t := range traffic {
		mbpsIn = append(mbpsIn, t.MbpsIn)
		mbpsOut = append(mbpsOut, t.MbpsOut)
		mbytesIn = append(mbytesIn, t.MBpsIn)
		mbytesOut = append(mbytesOut, t.MBpsOut)
	}

	var (
		predMbpsIn, predMbpsOut, predMbytesIn, predMbytesOut []float64
		wg                                                   sync.WaitGroup
	)
	wg.Add(4)
	go func() {
		defer wg.Done()
		predMbpsIn = trend.NewTrend(mbpsIn).Prediction(futureDays)
	}()
	go func() {
		defer wg.Done()
		predMbpsOut = trend.NewTrend(mbpsOut).Prediction(futureDays)
	}()
	go func() {
		defer wg.Done()
		predMbytesIn = trend.NewTrend(mbytesIn).Prediction(futureDays)
	}()
	go func() {
		defer wg.Done()
		predMbytesOut = trend.NewTrend(mbytesOut).Prediction(futureDays)
	}()
	wg.Wait()

	var forecast []model.TrafficTrend
	var lastDay time.Time
	if len(traffic) > 0 {
		lastDay = traffic[len(traffic)-1].Day
	} else {
		lastDay = dates.EndDate
	}
	for i := 1; i <= futureDays; i++ {
		forecast = append(forecast, model.TrafficTrend{
			Day:     lastDay.AddDate(0, 0, i),
			MbpsIn:  predMbpsIn[i-1],
			MbpsOut: predMbpsOut[i-1],
			MBpsIn:  predMbytesIn[i-1],
			MBpsOut: predMbytesOut[i-1],
		})
	}

	return &model.TrafficTrendForecast{
		Historical: traffic,
		Forecast:   forecast,
	}, nil
}
