package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
	"github.com/metalpoch/ultra-monitor/repository"
	"github.com/redis/go-redis/v9"
)

type TrafficUsecase struct {
	repo       repository.TrafficRepository
	cache      *cache.Redis
	prometheus prometheus.Prometheus
}

func NewTrafficUsecase(db *sqlx.DB, cache *cache.Redis, prometheus *prometheus.Prometheus) *TrafficUsecase {
	return &TrafficUsecase{repository.NewTrafficRepository(db), cache, *prometheus}
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
	var result []dto.Traffic

	keyCache := fmt.Sprintf("total-%d-%d", initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &result); err == nil {
		return result, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffic, err := use.prometheus.TrafficTotalByField(context.Background(), "", "", initDate, finalDate)
	if err != nil {
		return nil, err
	}

	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, result)

	return result, nil
}

func (use *TrafficUsecase) Region(region string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	var result []dto.Traffic

	keyCache := fmt.Sprintf("region-%s-%d-%d", region, initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &result); err == nil {
		return result, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffic, err := use.prometheus.TrafficTotalByField(context.Background(), "region", region, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, result)

	return result, nil
}

func (use *TrafficUsecase) State(state string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	var result []dto.Traffic

	keyCache := fmt.Sprintf("states-%s-%d-%d", state, initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &result); err == nil {
		return result, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffic, err := use.prometheus.TrafficTotalByField(context.Background(), "state", state, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, result)

	return result, nil
}

func (use *TrafficUsecase) Regions(initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	results := make(dto.TrafficByLabel)

	keyCache := fmt.Sprintf("regions-%d-%d", initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &results); err == nil {
		return results, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffics, err := use.prometheus.TrafficGroupedByField(context.Background(), "", "", "region", initDate, finalDate)
	if err != nil {
		return nil, err
	}

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, results)

	return results, nil
}

func (use *TrafficUsecase) StatesByRegion(region string, initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	results := make(dto.TrafficByLabel)

	keyCache := fmt.Sprintf("statesByRegion-%s-%d-%d", region, initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &results); err == nil {
		return results, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffics, err := use.prometheus.TrafficGroupedByField(context.Background(), "region", region, "state", initDate, finalDate)
	if err != nil {
		return nil, err
	}

	if len(traffics) < 1 {
		return results, nil
	}

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, results)

	return results, nil
}

func (use *TrafficUsecase) SysnameByState(state string, initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	results := make(dto.TrafficByLabel)

	keyCache := fmt.Sprintf("oltsByState-%s-%d-%d", state, initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &results); err == nil {
		return results, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffics, err := use.prometheus.SysnameByState(context.Background(), state, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	if len(traffics) < 1 {
		return results, nil
	}

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, results)

	return results, nil
}

func (use *TrafficUsecase) GponStats(ip string, initDate, finalDate time.Time) ([]dto.GponStats, error) {
	stats, err := use.prometheus.GponStatsByInstance(context.Background(), ip, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.GponStats
	for _, s := range stats {
		result = append(result, (dto.GponStats)(s))
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

func (use *TrafficUsecase) ByMunicipality(state, municipality string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetSnmpIndexByMunicipality(ctx, state, municipality)
	if err != nil {
		return nil, err
	}

	instancesMap := make(map[string]string)
	for _, r := range res {
		_, ok := instancesMap[r.IP]
		if ok {
			instancesMap[r.IP] += "|"
		}
		instancesMap[r.IP] += r.Idx
	}
	accum := make(map[time.Time]prometheus.Traffic)
	for ip, indexes := range instancesMap {
		traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, indexes, initDate, finalDate)
		if err != nil {
			return nil, err
		}

		for _, t := range traffic {
			key := t.Time.Truncate(15 * time.Minute)

			if data, ok := accum[key]; ok {
				data.BpsIn += t.BpsIn
				data.BpsOut += t.BpsOut
				data.BytesIn += t.BytesIn
				data.BytesOut += t.BytesOut
				data.Bandwidth += t.Bandwidth
				accum[key] = data
			} else {
				cloned := *t
				cloned.Time = key
				accum[key] = cloned
			}
		}
	}
	var result []dto.Traffic
	for _, val := range accum {
		result = append(result, (dto.Traffic)(val))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (use *TrafficUsecase) ByCounty(state, municipality, county string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetSnmpIndexByCounty(ctx, state, municipality, county)
	if err != nil {
		return nil, err
	}

	instancesMap := make(map[string]string)
	for _, r := range res {
		_, ok := instancesMap[r.IP]
		if ok {
			instancesMap[r.IP] += "|"
		}
		instancesMap[r.IP] += r.Idx
	}
	accum := make(map[time.Time]prometheus.Traffic)
	for ip, indexes := range instancesMap {
		traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, indexes, initDate, finalDate)
		if err != nil {
			return nil, err
		}

		for _, t := range traffic {
			key := t.Time.Truncate(15 * time.Minute)

			if data, ok := accum[key]; ok {
				data.BpsIn += t.BpsIn
				data.BpsOut += t.BpsOut
				data.BytesIn += t.BytesIn
				data.BytesOut += t.BytesOut
				data.Bandwidth += t.Bandwidth
				accum[key] = data
			} else {
				cloned := *t
				cloned.Time = key
				accum[key] = cloned
			}
		}
	}
	var result []dto.Traffic
	for _, val := range accum {
		result = append(result, (dto.Traffic)(val))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil

}

func (use *TrafficUsecase) ByODN(state, municipality, odn string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetSnmpIndexByODN(ctx, state, municipality, odn)
	if err != nil {
		return nil, err
	}

	instancesMap := make(map[string]string)
	for _, r := range res {
		_, ok := instancesMap[r.IP]
		if ok {
			instancesMap[r.IP] += "|"
		}
		instancesMap[r.IP] += r.Idx
	}
	accum := make(map[time.Time]prometheus.Traffic)
	for ip, indexes := range instancesMap {
		traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, indexes, initDate, finalDate)
		if err != nil {
			return nil, err
		}

		for _, t := range traffic {
			key := t.Time.Truncate(15 * time.Minute)

			if data, ok := accum[key]; ok {
				data.BpsIn += t.BpsIn
				data.BpsOut += t.BpsOut
				data.BytesIn += t.BytesIn
				data.BytesOut += t.BytesOut
				data.Bandwidth += t.Bandwidth
				accum[key] = data
			} else {
				cloned := *t
				cloned.Time = key
				accum[key] = cloned
			}
		}
	}
	var result []dto.Traffic
	for _, val := range accum {
		result = append(result, (dto.Traffic)(val))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (use *TrafficUsecase) ByIdx(ip, idx string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, idx, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	return result, nil
}
