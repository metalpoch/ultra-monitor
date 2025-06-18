package usecase

import (
	"context"
	"database/sql"
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

type PonUsecase struct {
	repo  repository.PonRepository
	cache *cache.Redis
}

func NewPonUsecase(db *sqlx.DB, cache *cache.Redis) *PonUsecase {
	return &PonUsecase{repository.NewPonRepository(db), cache}
}

func (uc *PonUsecase) GetAllByDevice(sysname string) ([]model.Pon, error) {
	var interfaces []model.Pon
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.PonsByOLT(ctx, sysname)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		interfaces = append(interfaces, (model.Pon)(e))
	}

	return interfaces, nil
}

func (uc *PonUsecase) PonByOltAndPort(sysname string, shell, card, port int) (*model.Pon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ifname := fmt.Sprintf("GPON %d/%d/%d", shell, card, port)
	pon, err := uc.repo.PonByPort(ctx, sysname, ifname)
	if err != nil {
		return nil, err
	}
	return (*model.Pon)(&pon), nil
}

func (uc *PonUsecase) TrafficByState(state string, dates dto.RangeDate) ([]model.Traffic, error) {
	res, err := uc.repo.TrafficByState(context.Background(), state, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *PonUsecase) TrafficByMunicipaly(state, municipality string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByMunicipality(ctx, state, municipality, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *PonUsecase) TrafficByCounty(state, municipality, county string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByCounty(ctx, state, municipality, county, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *PonUsecase) TrafficByODN(state, municipality, county, odn string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByODN(ctx, state, municipality, county, odn, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}
func (uc *PonUsecase) TrafficByPon(sysname string, shell, card, port int, dates dto.RangeDate) ([]model.Traffic, error) {
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

func (uc *PonUsecase) GetTrafficPonForecast(futureDays int, dates dto.RangeDate) (*model.TrafficTrendForecast, error) {
	var traffic []model.TrafficSummary
	key := fmt.Sprintf("trafficPonForecast:%d:%d", dates.InitDate.Unix(), dates.EndDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := uc.repo.GetDailyAveragedHourlyMaxTrafficTrends(context.Background(), dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}

		for _, t := range res {
			traffic = append(traffic, (model.TrafficSummary)(t))
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
		mbytesIn = append(mbytesIn, t.MbytesInSec)
		mbytesOut = append(mbytesOut, t.MbytesOutSec)
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

	var forecast []model.TrafficSummary
	var lastDay time.Time
	if len(traffic) > 0 {
		lastDay = traffic[len(traffic)-1].Day
	} else {
		lastDay = dates.EndDate
	}
	for i := 1; i <= futureDays; i++ {
		forecast = append(forecast, model.TrafficSummary{
			Day:          lastDay.AddDate(0, 0, i),
			MbpsIn:       predMbpsIn[i-1],
			MbpsOut:      predMbpsOut[i-1],
			MbytesInSec:  predMbytesIn[i-1],
			MbytesOutSec: predMbytesOut[i-1],
		})
	}

	return &model.TrafficTrendForecast{
		Historical: traffic,
		Forecast:   forecast,
	}, nil
}

func (uc *PonUsecase) UpdateSummaryTraffic(dates dto.RangeDate) error {
	res, err := uc.repo.GetDailyAveragedHourlyMaxTrafficTrends(context.Background(), dates.InitDate, dates.EndDate)
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return sql.ErrNoRows
	}

	return uc.repo.UpsertSummaryTraffic(context.Background(), res)
}

func (uc *PonUsecase) GetTrafficSummary(dates dto.RangeDate) ([]model.TrafficTotalSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetTrafficSummary(ctx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.TrafficTotalSummary
	for _, t := range res {
		traffic = append(traffic, (model.TrafficTotalSummary)(t))
	}

	return traffic, err
}

func (uc *PonUsecase) GetTrafficStatesSummary(dates dto.RangeDate) ([]model.TrafficInfoSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetTrafficStatesSummary(ctx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.TrafficInfoSummary
	for _, t := range res {
		traffic = append(traffic, (model.TrafficInfoSummary)(t))
	}

	return traffic, err
}

func (uc *PonUsecase) GetTrafficMunicipalitySummary(state string, dates dto.RangeDate) ([]model.TrafficInfoSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetTrafficMunicipalitySummary(ctx, state, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.TrafficInfoSummary
	for _, t := range res {
		traffic = append(traffic, (model.TrafficInfoSummary)(t))
	}

	return traffic, err
}

func (uc *PonUsecase) GetTrafficCountySummary(state, municipality string, dates dto.RangeDate) ([]model.TrafficInfoSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetTrafficCountySummary(ctx, state, municipality, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.TrafficInfoSummary
	for _, t := range res {
		traffic = append(traffic, (model.TrafficInfoSummary)(t))
	}

	return traffic, err
}

func (uc *PonUsecase) GetTrafficOdnSummary(state, municipality, county string, dates dto.RangeDate) ([]model.TrafficInfoSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetTrafficOdnSummary(ctx, state, municipality, county, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.TrafficInfoSummary
	for _, t := range res {
		traffic = append(traffic, (model.TrafficInfoSummary)(t))
	}

	return traffic, err
}
