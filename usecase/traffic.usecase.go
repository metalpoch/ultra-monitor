package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
	"github.com/metalpoch/ultra-monitor/repository"
)

type TrafficUsecase struct {
	repo       repository.TrafficRepository
	prometheus prometheus.Prometheus
}

func NewTrafficUsecase(db *sqlx.DB, prometheus *prometheus.Prometheus) *TrafficUsecase {
	return &TrafficUsecase{repository.NewTrafficRepository(db), *prometheus}
}

func (use *TrafficUsecase) DeviceLocation() ([]dto.DeviceLocation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	devices, err := use.prometheus.DeviceLocations(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.DeviceLocation
	for _, d := range devices {
		res = append(res, (dto.DeviceLocation)(d))
	}

	sort.SliceStable(res, func(i, j int) bool {
		if res[i].Region != res[j].Region {
			return res[i].Region < res[j].Region
		}

		if res[i].State != res[j].State {
			return res[i].State < res[j].State
		}
		return res[i].SysName < res[j].SysName
	})

	return res, nil
}

func (use *TrafficUsecase) InfoInstance(ip string) ([]dto.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	devices, err := use.prometheus.InstanceScan(ctx, ip)
	if err != nil {
		return nil, err
	}

	var res []dto.InfoDevice
	for _, d := range devices {
		res = append(res, (dto.InfoDevice)(d))
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].IfName < res[j].IfName
	})

	return res, nil
}

func (use *TrafficUsecase) Total(initDate, finalDate time.Time) ([]dto.Traffic, error) {
	traffic, err := use.prometheus.TrafficTotal(context.Background(), initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	return result, nil
}

func (use *TrafficUsecase) Regions(initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	traffics, err := use.prometheus.TrafficGroupedByField(context.Background(), "", "", "region", initDate, finalDate)
	if err != nil {
		return nil, err
	}

	results := make(dto.TrafficByLabel)

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}
	return results, nil
}

func (use *TrafficUsecase) StatesByRegion(region string, initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	traffics, err := use.prometheus.TrafficGroupedByField(context.Background(), "region", region, "state", initDate, finalDate)
	if err != nil {
		return nil, err
	}

	results := make(dto.TrafficByLabel)

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}
	return results, nil
}

func (use *TrafficUsecase) Region(region string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	traffic, err := use.prometheus.TrafficByRegion(context.Background(), region, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	return result, nil
}

func (use *TrafficUsecase) States(state string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	traffic, err := use.prometheus.TrafficByState(context.Background(), state, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	return result, nil
}

func (use *TrafficUsecase) GroupIP(ips []string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	traffic, err := use.prometheus.TrafficGroupInstance(context.Background(), ips, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	return result, nil
}

func (use *TrafficUsecase) IndexAndIP(ip, index string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, index, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	return result, nil
}
