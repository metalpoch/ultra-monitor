package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/internal/cache"
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

func (uc OntUsecase) AllOntStatus(initDate, endDate time.Time) ([]model.OntStatusCounts, error) {
	var status []model.OntStatusCounts
	key := fmt.Sprintf("ontStatus:%d:%d", initDate.Unix(), endDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.AllOntStatus(context.Background(), initDate, endDate)
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

func (uc OntUsecase) OntStatusByState(state string, initDate, endDate time.Time) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByState:%s:%d:%d", state, initDate.Unix(), endDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.GetOntStatusByState(context.Background(), state, initDate, endDate)
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

func (uc OntUsecase) OntStatusByODN(state, odn string, initDate, endDate time.Time) ([]model.OntStatusCountsByState, error) {
	var status []model.OntStatusCountsByState
	key := fmt.Sprintf("ontStatusByODN:%s:%s:%d:%d", state, odn, initDate.Unix(), endDate.Unix())
	err := uc.cache.FindOne(context.Background(), key, &status)
	if err == redis.Nil {
		res, err := uc.repo.GetOntStatusByODN(context.Background(), state, odn, initDate, endDate)
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

func (uc OntUsecase) TrafficOnt(ponID uint64, idx string, initDate, endDate time.Time) ([]model.TrafficOnt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var traffic []model.TrafficOnt
	res, err := uc.repo.TrafficOnt(ctx, ponID, idx, initDate, endDate)
	if err != nil {
		return nil, err
	}
	for _, s := range res {
		traffic = append(traffic, (model.TrafficOnt)(s))
	}

	return traffic, err
}
