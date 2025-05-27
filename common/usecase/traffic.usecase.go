package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/common/repository"
	"gorm.io/gorm"
)

type TrafficUsecase struct {
	repo     repository.TrafficRepository
	telegram tracking.SmartModule
}

func NewTrafficUsecase(db *gorm.DB, telegram tracking.SmartModule) *TrafficUsecase {
	return &TrafficUsecase{repository.NewTrafficRepository(db), telegram}
}

func (use TrafficUsecase) Add(traffic *model.Traffic) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := use.repo.AddOlt(ctx, (*entity.TrafficOlt)(traffic))
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).Add - use.repo.Add(ctx, %v)", *(*entity.TrafficOlt)(traffic)),
			err,
		)
	}
	return err
}

func (use TrafficUsecase) AddOnt(traffic *model.TrafficOnt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data := entity.TrafficOnt{
		InterfaceID:      traffic.InterfaceID,
		Idx:              traffic.Idx,
		Despt:            traffic.Despt,
		SerialNumber:     traffic.SerialNumber,
		LineProfName:     traffic.LineProfName,
		ControlMacCount:  traffic.ControlMacCount,
		OltDistance:      traffic.OltDistance,
		ControlRunStatus: traffic.ControlRunStatus,
		InBps:            traffic.InBps,
		OutBps:           traffic.OutBps,
		BytesIn:          traffic.BytesIn,
		BytesOut:         traffic.BytesOut,
		Date:             traffic.Date,
	}
	err := use.repo.AddOnt(ctx, &data)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).AddOnt - use.repo.Add(ctx, %v)", data),
			err,
		)
	}
	return err
}
