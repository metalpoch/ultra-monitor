package usecase

import (
	"context"
	"fmt"
	"time"

	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/core/model"
	"github.com/metalpoch/olt-blueprint/core/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OltUsecase struct {
	repo  repository.OltRepository
	cache cache.Redis
}

func NewOltUsecase(db *gorm.DB, cache cache.Redis) *OltUsecase {
	return &OltUsecase{repository.NewOltRepository(db), cache}
}

func (use OltUsecase) Olts(page, limit uint8) ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.Olts(ctx, page, limit)
	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}

func (use OltUsecase) OltsByState(state string, page, limit uint8) ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.OltsByState(ctx, state, page, limit)
	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}

func (use OltUsecase) OltsByCounty(state, county string, page, limit uint8) ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.OltsByCounty(ctx, state, county, page, limit)
	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}

func (use OltUsecase) OltsByMunicipality(state, county, municipality string, page, limit uint8) ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.OltsByMunicipality(ctx, state, county, municipality, page, limit)
	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}

func (use OltUsecase) Traffic(dates *commonModel.TrafficRangeDate) ([]model.TrafficOlt, error) {
	var traffic []model.TrafficOlt
	key := fmt.Sprintf("oltTraffic:%d:%d", dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.Traffic(context.Background(), dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, e := range res {
			traffic = append(traffic, (model.TrafficOlt)(e))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic)
		if err != nil {
			return nil, err
		}
		return traffic, nil

	} else if err != nil {
		return nil, err
	}

	return traffic, err
}

func (use OltUsecase) TrafficByState(state string, dates *commonModel.TrafficRangeDate) ([]model.TrafficOlt, error) {
	var traffic []model.TrafficOlt
	key := fmt.Sprintf("oltTrafficByState:%s:%d:%d", state, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.TrafficByState(context.Background(), state, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			traffic = append(traffic, (model.TrafficOlt)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic)
		if err != nil {
			return nil, err
		}
		return traffic, nil

	} else if err != nil {
		return nil, err
	}

	return traffic, err
}

func (use OltUsecase) TrafficByCounty(state, county string, dates *commonModel.TrafficRangeDate) ([]model.TrafficOlt, error) {
	var traffic []model.TrafficOlt
	key := fmt.Sprintf("oltTrafficByCounty:%s:%s:%d:%d", state, county, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.TrafficByCounty(context.Background(), state, county, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			traffic = append(traffic, (model.TrafficOlt)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic)
		if err != nil {
			return nil, err
		}
		return traffic, nil

	} else if err != nil {
		return nil, err
	}

	return traffic, err
}

func (use OltUsecase) TrafficByMunicipality(state, county, municipality string, dates *commonModel.TrafficRangeDate) ([]model.TrafficOlt, error) {
	var traffic []model.TrafficOlt
	key := fmt.Sprintf("oltTrafficByMunicipality:%s:%s:%s:%d:%d", state, county, municipality, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.TrafficByMunicipality(context.Background(), state, county, municipality, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			traffic = append(traffic, (model.TrafficOlt)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic)
		if err != nil {
			return nil, err
		}
		return traffic, nil

	} else if err != nil {
		return nil, err
	}

	return traffic, err
}

func (use OltUsecase) TrafficByODN(state, odn string, dates *commonModel.TrafficRangeDate) ([]model.TrafficOlt, error) {
	var traffic []model.TrafficOlt
	key := fmt.Sprintf("oltTrafficByODN:%s:%s:%d:%d", state, odn, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.TrafficByODN(context.Background(), state, odn, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			traffic = append(traffic, (model.TrafficOlt)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic)
		if err != nil {
			return nil, err
		}
		return traffic, nil

	} else if err != nil {
		return nil, err
	}

	return traffic, err
}

func (use OltUsecase) TrafficByOLT(sysname string, dates *commonModel.TrafficRangeDate) ([]model.TrafficOlt, error) {
	var traffic []model.TrafficOlt
	key := fmt.Sprintf("oltTrafficByOLT:%s::%d:%d", sysname, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.TrafficByOLT(context.Background(), sysname, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			traffic = append(traffic, (model.TrafficOlt)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic)
		if err != nil {
			return nil, err
		}
		return traffic, nil

	} else if err != nil {
		return nil, err
	}

	return traffic, err
}

func (use OltUsecase) TrafficByPON(sysname string, shell, card, port uint8, dates *commonModel.TrafficRangeDate) ([]model.TrafficOlt, error) {
	var traffic []model.TrafficOlt
	pon := fmt.Sprintf("GPON %d/%d/%d", shell, card, shell)
	key := fmt.Sprintf("oltTrafficByPON:%s:%s:%d:%d", sysname, pon, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.TrafficByPON(context.Background(), sysname, pon, dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			traffic = append(traffic, (model.TrafficOlt)(s))
		}
		err = use.cache.InsertOne(context.Background(), key, 12*time.Hour, traffic)
		if err != nil {
			return nil, err
		}
		return traffic, nil

	} else if err != nil {
		return nil, err
	}

	return traffic, err
}
