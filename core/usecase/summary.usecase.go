package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/cache"
	"github.com/metalpoch/olt-blueprint/core/model"
	"github.com/metalpoch/olt-blueprint/core/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SummaryUsecase struct {
	repo  repository.SummaryRepository
	cache cache.Redis
}

func NewSummaryUsecase(db *gorm.DB, cache cache.Redis) *SummaryUsecase {
	return &SummaryUsecase{repository.NewSummaryRepository(db), cache}
}

func (use SummaryUsecase) UserStatus(query *commonModel.TranficRangeDate) ([]model.UserStatusCounts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var status []model.UserStatusCounts

	key := fmt.Sprintf("userStatus:%d:%d", query.InitDate.Unix(), query.EndDate.Unix())
	err := use.cache.FindOne(ctx, key, &status)
	if err == redis.Nil {
		res, err := use.repo.UserStatus(ctx, query)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.UserStatusCounts)(s))
		}
		err = use.cache.InsertOne(ctx, key, 12*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil

	} else if err != nil {
		return nil, err
	}

	return status, err
}

func (use SummaryUsecase) UserStatusByState(query *model.UserStatusByState) ([]model.UserStatusCounts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var status []model.UserStatusCounts

	key := fmt.Sprintf("userStatus:%s:%d:%d", strings.Join(query.States, "-"), query.InitDate.Unix(), query.EndDate.Unix())
	err := use.cache.FindOne(ctx, key, &status)
	if err == redis.Nil {
		res, err := use.repo.UserStatusByState(ctx, query)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.UserStatusCounts)(s))
		}
		err = use.cache.InsertOne(ctx, key, 12*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil

	} else if err != nil {
		return nil, err
	}

	return status, err
}

func (use SummaryUsecase) Traffic(dates *commonModel.TranficRangeDate) ([]model.TrafficResponse, error) {
	var traffic []model.TrafficResponse
	key := fmt.Sprintf("trafficState:all:%d:%d", dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.Traffic(context.Background(), dates)
		if err != nil {
			return nil, err
		}
		for _, trf := range res {
			traffic = append(traffic, (model.TrafficResponse)(trf))
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

func (use SummaryUsecase) TrafficByState(state string, dates *commonModel.TranficRangeDate) ([]model.TrafficResponse, error) {
	/*ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)*/
	/*defer cancel()*/

	var traffic []model.TrafficResponse
	key := fmt.Sprintf("trafficState:%s:%d:%d", state, dates.InitDate.Unix(), dates.EndDate.Unix())
	err := use.cache.FindOne(context.Background(), key, &traffic)
	if err == redis.Nil {
		res, err := use.repo.TrafficByState(context.Background(), state, dates)
		if err != nil {
			return nil, err
		}
		for _, trf := range res {
			traffic = append(traffic, (model.TrafficResponse)(trf))
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
