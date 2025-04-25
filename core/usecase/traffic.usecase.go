package usecase

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/core/model"
	"github.com/metalpoch/olt-blueprint/core/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TrafficUsecase struct {
	repo     repository.TrafficRepository
	repoinfo repository.InfoRepository
	telegram tracking.SmartModule
	cache    cache.Redis
}

func NewTrafficUsecase(db *gorm.DB, telegram tracking.SmartModule, cache cache.Redis) *TrafficUsecase {
	return &TrafficUsecase{repository.NewTrafficRepository(db), repository.NewInfoRepository(db), telegram, cache}
}

func (use TrafficUsecase) GetTrafficByInterface(id uint, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByInterface(ctx, id, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByInterface - use.repo.GetTrafficByInterface(ctx, %d, %v)", id, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTrafficByDevice(id uint, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByDevice(ctx, id, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByDevice - use.repo.GetTrafficByDevice(ctx, %d, %v)", id, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTrafficByFat(id uint, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByFat(ctx, id, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByFat - use.repo.GetTrafficByFat(ctx, %d, %v)", id, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTrafficByLocationID(id uint, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByLocationID(ctx, id, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByLocationID - use.repo.GetTrafficByLocationID(ctx, %d, %v)", id, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTrafficByState(state string, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByState(ctx, state, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByState - use.repo.GetTrafficByState(ctx, %s, %v)", state, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTrafficByCounty(state, county string, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByCounty(ctx, state, county, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByCounty - use.repo.GetTrafficByCounty(ctx, %s, %s, %v)", state, county, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTrafficByMunicipality(state, county, municipality string, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByMunicipality(ctx, state, county, municipality, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByMunicipality - use.repo.GetTrafficByMunicipality(ctx, %s, %s, %v)", state, county, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTrafficByODN(odn string, date *commonModel.TranficRangeDate) ([]*commonModel.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByODN(ctx, odn, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx, %s, %v)", odn, date),
			err,
		)
		return nil, err
	}

	traffics := []*commonModel.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &commonModel.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
			BytesIn:   e.BytesIn,
			BytesOut:  e.BytesOut,
		})
	}

	return traffics, err
}

func (use TrafficUsecase) GetTotalTrafficByStateByMonth(month string) ([]*commonModel.TrafficState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	traffics := new([]*commonModel.TrafficState)
	err := use.cache.FindOne(ctx, "trafficState", traffics)
	if len(*traffics) > 0 {
		return *traffics, nil
	} else if err != nil && err != redis.Nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState(%s) - use.cache.FindOne(ctx,  \"trafficState\", %v)", month, traffics),
			err,
		)
		return nil, err
	}

	states, err := use.repoinfo.GetLocationStates(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState(%s) - use.repoinfo.GetLocationStates(ctx)"),
			err,
		)
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, state := range states {
		wg.Add(1)
		go use.upsertStateDevicesInCache(ctx, &wg, &mu, *state, month, traffics)
	}

	wg.Wait()
	err = use.cache.InsertOne(ctx, "trafficState", 24*time.Hour, *traffics)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState(%s) - use.cache.InsertOne(ctx, \"trafficState\", 24*time.Hour, %v)", month, traffics),
			err,
		)
		return nil, err
	}

	return *traffics, err
}

func (use TrafficUsecase) GetTopTrafficByStateByMonth(month string, topN int) ([]*commonModel.TrafficState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	traffics := new([]*commonModel.TrafficState)
	err := use.cache.FindOne(ctx, "trafficState", traffics)
	if len(*traffics) > 0 {
		sort.Slice(*traffics, func(i, j int) bool {
			return (*traffics)[i].In+(*traffics)[i].Out > (*traffics)[j].In+(*traffics)[j].Out
		})
		return (*traffics)[:topN], nil
	} else if err != nil && err != redis.Nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState(%s) - use.cache.FindOne(ctx,  \"trafficState\", %v)", month, traffics),
			err,
		)
		return nil, err
	}

	states, err := use.repoinfo.GetLocationStates(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTopTrafficByStateByMonth(%s, %d) - use.repoinfo.GetLocationStates(ctx)", month, topN),
			err,
		)
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, state := range states {
		wg.Add(1)
		go use.upsertStateDevicesInCache(ctx, &wg, &mu, *state, month, traffics)
	}

	wg.Wait()
	err = use.cache.InsertOne(ctx, "trafficState", 24*time.Hour, *traffics)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState(%s) - use.cache.InsertOne(ctx, \"trafficState\",24*time.Hour, %v)", month, traffics),
			err,
		)
		return nil, err
	}

	sort.Slice(*traffics, func(i, j int) bool {
		return (*traffics)[i].In+(*traffics)[i].Out > (*traffics)[j].In+(*traffics)[j].Out
	})

	return (*traffics)[:topN], nil
}

func (use TrafficUsecase) GetTotalTrafficOdnByMonth(month string) ([]*commonModel.TrafficODN, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	traffics := new([]*commonModel.TrafficODN)
	err := use.cache.FindOne(ctx, "trafficODN", traffics)
	if len(*traffics) > 0 {
		return *traffics, nil
	} else if err != nil && err != redis.Nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND(%s) - use.cache.FindOne(ctx,  \"trafficODN\", %v)", month, traffics),
			err,
		)
		return nil, err
	}

	odns, err := use.repoinfo.GetAllODN(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND(%s) - use.repoinfo.GetAllODN(ctx)", month),
			err,
		)
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, odn := range odns {
		wg.Add(1)
		go use.upsertODNDevicesInCache(ctx, &wg, &mu, *odn, month, traffics)
	}

	wg.Wait()
	err = use.cache.InsertOne(ctx, "trafficODN", 24*time.Hour, traffics)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND(%s) - use.cache.InsertOne(ctx, \"trafficODN\", 24*time.Hour, %v)", month, traffics),
			err,
		)
		return nil, err
	}

	return *traffics, err
}

func (use TrafficUsecase) upsertStateDevicesInCache(ctx context.Context, wg *sync.WaitGroup, mu *sync.Mutex, state, month string, traffic *[]*commonModel.TrafficState) {
	defer wg.Done()
	devices := use.getDevicesByState(ctx, state)
	if devices == nil {
		return
	}

	res, err := use.repo.GetTotalTrafficByState(ctx, devices, month)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState - use.repo.GetTotalTrafficByState(ctx, %v, %s)", devices, month),
			err,
		)
		return
	}

	mu.Lock()
	res.State = state
	*traffic = append(*traffic, res)
	mu.Unlock()
}

func (use TrafficUsecase) upsertODNDevicesInCache(ctx context.Context, wg *sync.WaitGroup, mu *sync.Mutex, odn, month string, traffic *[]*commonModel.TrafficODN) {
	defer wg.Done()
	devices := use.getDeviceByODN(ctx, odn)
	if devices == nil {
		return
	}

	res, err := use.repo.GetTotalTrafficByODN(ctx, devices, month)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState - use.repo.GetTotalTrafficByState(ctx, %v, %s)", devices, month),
			err,
		)
		return
	}

	mu.Lock()
	res.ODN = odn
	*traffic = append(*traffic, res)
	mu.Unlock()
}

func (use TrafficUsecase) getDeviceByODN(ctx context.Context, odn string) []uint {
	retrieved := new(model.DevicesID)
	err := use.cache.FindOne(ctx, odn, retrieved)
	if err == redis.Nil {
		devicesEntity, err := use.repoinfo.GetDevicesByOND(ctx, odn)
		if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				fmt.Sprintf("(trafficUsecase).GetTrafficByODN - use.repoinfo.getDeviceByODN(ctx, %s)", odn),
				err,
			)
			return nil
		}
		if devicesEntity == nil {
			return nil
		}

		devices := make([]uint, len(devicesEntity))
		for i, ptr := range devices {
			devices[i] = ptr
		}

		odnDevices := model.DevicesID{
			Devices: devices,
		}

		if err = use.cache.InsertOne(ctx, odn, 24*time.Hour, odnDevices); err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				fmt.Sprintf("(trafficUsecase).GetTrafficByODN - use.cache.InsertOne(ctx, %s, 24*time.Hour, %v)", odn, odnDevices),
				err,
			)
			return nil
		}
		return devices

	} else if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState -use.cache.FindOne(ctx, %s, %v)", odn, retrieved),
			err,
		)
		return nil
	}

	return retrieved.Devices
}

func (use TrafficUsecase) getDevicesByState(ctx context.Context, state string) []uint {
	retrieved := new(model.DevicesID)
	err := use.cache.FindOne(ctx, state, retrieved)
	if err == redis.Nil {
		devicesEntity, err := use.repoinfo.GetDeviceByState(ctx, state)
		if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				fmt.Sprintf("(trafficUsecase).GetTrafficByODN - use.repoinfo.GetDeviceByState(ctx, %s)", state),
				err,
			)
			return nil
		}
		if devicesEntity == nil {
			return nil
		}

		devices := make([]uint, len(devicesEntity))
		for i, ptr := range devices {
			devices[i] = ptr
		}

		locationsDevices := model.DevicesID{
			Devices: devices,
		}

		if err = use.cache.InsertOne(ctx, state, 24*time.Hour, locationsDevices); err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_TRAFFIC,
				constants.CATEGORY_DATABASE,
				fmt.Sprintf("(trafficUsecase).GetTrafficByODN - use.cache.InsertOne(ctx, %s, 24*time.Hour, %v).Err()", state, locationsDevices),
				err,
			)
			return nil
		}
		return devices

	} else if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTotalTrafficByState -use.cache.FindOne(ctx, %s, %v)", state, retrieved),
			err,
		)
		return nil
	}

	return retrieved.Devices
}
