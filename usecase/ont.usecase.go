package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/internal/cache"
	"github.com/metalpoch/olt-blueprint/internal/dto"
	"github.com/metalpoch/olt-blueprint/model"
	"github.com/metalpoch/olt-blueprint/repository"
	"github.com/redis/go-redis/v9"
)

type OntUsecase struct {
	repo  repository.OntRepository
	cache *cache.Redis
}

func NewOntUsecase(db *sqlx.DB, cache *cache.Redis) *OntUsecase {
	return &OntUsecase{repository.NewOntRepository(db), cache}
}

func (uc *OntUsecase) AllOntStatus(dates dto.RangeDate) ([]model.OntStatusCounts, error) {
	var status []model.OntStatusCounts
	key := fmt.Sprintf("ontStatus:%d:%d", dates.InitDate.Unix(), dates.EndDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.AllOntStatus(context.Background(), dates.InitDate, dates.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.OntStatusCounts)(s))
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

func (uc *OntUsecase) OntStatusByOdn(state, odn string, dates dto.RangeDate) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByODN:%s:%s:%d:%d", state, odn, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.GetOntStatusByODN(context.Background(), state, odn, dates.InitDate, dates.EndDate)
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
