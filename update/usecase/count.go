package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type countUsecase struct {
	repo repository.CountRepository
}

func NewCountUsecase(repo repository.CountRepository) *countUsecase {
	return &countUsecase{repo}
}

func (use countUsecase) Find(olt, interfaceName string) (*model.Count, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.Find(ctx, olt, interfaceName)
	if err != mongo.ErrNoDocuments && err != nil {
		return nil, err
	}

	count := &model.Count{
		OLT:       res.OLT,
		Interface: res.Interface,
		Date:      res.Date,
		BytesIn:   res.BytesIn,
		BytesOut:  res.BytesOut,
		Bandwidth: res.Bandwidth,
	}

	return count, nil
}

func (use countUsecase) Create(count model.Count) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return use.repo.Create(ctx, entity.Count{
		OLT:       count.OLT,
		Interface: count.Interface,
		Date:      count.Date,
		BytesIn:   count.BytesIn,
		BytesOut:  count.BytesOut,
		Bandwidth: count.Bandwidth,
	})
}
