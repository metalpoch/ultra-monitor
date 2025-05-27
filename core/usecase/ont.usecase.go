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

type OntUsecase struct {
	repo  repository.OntRepository
	cache cache.Redis
}

func NewOntUsecase(db *gorm.DB, cache cache.Redis) *OntUsecase {
	return &OntUsecase{repository.NewOntRepository(db), cache}
}

func (use OntUsecase) OntStatus(dates *commonModel.TrafficRangeDate) ([]model.OntStatusCounts, error) {
	var status []model.OntStatusCounts
	key := fmt.Sprintf("ontStatus:%d:%d", dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := use.repo.GetOntStatus(context.Background(), dates)
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

func (use OntUsecase) OntStatusByState(state string, dates *commonModel.TrafficRangeDate) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByState:%s:%d:%d", state, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := use.repo.GetOntStatusByState(context.Background(), state, dates)
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

func (use OntUsecase) OntStatusByODN(state, odn string, dates *commonModel.TrafficRangeDate) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByODN:%s:%s:%d:%d", state, odn, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := use.repo.GetOntStatusByODN(context.Background(), state, odn, dates)
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

func (use OntUsecase) TrafficOnt(interfaceID, idx string, dates *commonModel.TrafficRangeDate) ([]model.OntTraffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var traffic []model.OntTraffic
	res, err := use.repo.GetTrafficOnt(ctx, interfaceID, idx, dates)
	if err != nil {
		return nil, err
	}
	for _, s := range res {
		traffic = append(traffic, (model.OntTraffic)(s))
	}

	return traffic, err
}
