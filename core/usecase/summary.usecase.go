package usecase

import (
	"context"
	"fmt"
	"time"

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

func (use SummaryUsecase) UserStatus(query *model.UserStatusQuery) ([]model.UserStatusCounts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var status []model.UserStatusCounts

	key := fmt.Sprintf("userStatus:%d:%d", query.InitDate.Unix(), query.EndDate.Unix())
	err := use.cache.FindOne(ctx, key, &status)
	if err == redis.Nil {
		res, err := use.repo.UserStatus(ctx, query.InitDate, query.EndDate)
		if err != nil {
			return nil, err
		}
		for _, s := range res {
			status = append(status, (model.UserStatusCounts)(s))
		}
		err = use.cache.InsertOne(ctx, key, 24*time.Hour, status)
		if err != nil {
			return nil, err
		}
		return status, nil

	} else if err != nil {
		return nil, err
	}

	return status, err
}
