package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
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
