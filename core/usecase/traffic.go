package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/core/repository"
	"gorm.io/gorm"
)

type trafficUsecase struct {
	repo     repository.TrafficRepository
	telegram tracking.Telegram
}

func NewTrafficUsecase(db *gorm.DB, telegram tracking.Telegram) *trafficUsecase {
	return &trafficUsecase{repository.NewTrafficRepository(db), telegram}
}

func (use trafficUsecase) GetTrafficByInterface(id uint, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	traffics := []*model.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &model.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}

func (use trafficUsecase) GetTrafficByDevice(id uint, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	traffics := []*model.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &model.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}

func (use trafficUsecase) GetTrafficByFat(id uint, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	traffics := []*model.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &model.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}

func (use trafficUsecase) GetTrafficByLocationID(id uint, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	traffics := []*model.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &model.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}

func (use trafficUsecase) GetTrafficByState(state string, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	traffics := []*model.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &model.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}

func (use trafficUsecase) GetTrafficByCounty(state, county string, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	traffics := []*model.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &model.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}

func (use trafficUsecase) GetTrafficByMunicipality(state, county, municipality string, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	traffics := []*model.TrafficResponse{}
	for _, e := range res {
		traffics = append(traffics, &model.TrafficResponse{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}
