package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	core "github.com/metalpoch/olt-blueprint/core/model"
	"github.com/metalpoch/olt-blueprint/core/repository"
	"github.com/metalpoch/olt-blueprint/core/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type trafficUsecase struct {
	repo     repository.TrafficRepository
	repoinfo repository.InfoRepository
	telegram tracking.SmartModule
	redis    *redis.Client
}

func NewTrafficUsecase(db *gorm.DB, telegram tracking.SmartModule, redis *redis.Client) *trafficUsecase {
	return &trafficUsecase{repository.NewTrafficRepository(db), repository.NewInfoRepository(db), telegram, redis}
}

func (use trafficUsecase) GetTrafficByInterface(id uint, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
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

func (use trafficUsecase) GetTrafficByODN(odn string, date *model.TranficRangeDate) ([]*model.TrafficResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := use.repo.GetTrafficByODN(ctx, odn, date)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx, %s, %v)", odn, date),
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

func (use trafficUsecase) GetTotalTrafficByState(month string) ([]*model.TrafficState, error) {
	var traffics []*model.TrafficState
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	//Pregunta por todos los estados almacenados en la base de datos
	states, err := use.repoinfo.GetLocationStates(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
			err,
		)
		return nil, err
	}
	//Recorre la lista de estados
	for _, state := range states {
		//Verifica si existe el estado en redis
		exist := utils.VerifyExistence(use.redis, *state, ctx)
		if exist {
			//Si existe, obtiene el json con  los ids de todos los devices para ese estado
			retrievedData, err := use.redis.Get(ctx, *state).Result()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			//Parsea el json obtenido a un objeto LocationsDevice
			var retrievedLocations core.RedisDevice
			retrievedLocations, err = utils.DeserializeModel(retrievedData)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			//Pregunta por el total de trafico por device para ese estado
			res, err := use.repo.GetTotalTrafficByState(ctx, retrievedLocations.Devices, month)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			res.State = *state
			traffics = append(traffics, res)
		} else {
			//Si no existe, pregunta a la base de datos por los ids de todos los devices para ese estado
			var devices []uint
			devicesEntity, err := use.repoinfo.GetDeviceByState(ctx, *state)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			if devicesEntity == nil {
				continue
			}
			for _, deviceEntity := range devicesEntity {
				devices = append(devices, deviceEntity.ID)
			}
			//Crea un objeto LocationsDevice con los ids de todos los devices para ese estado
			locationsDevices := core.RedisDevice{
				Devices: devices,
			}
			//Parsea el objeto LocationsDevice a json
			jsonData, err := utils.SerializeModel(locationsDevices)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			//Guarda el json en redis
			err = use.redis.Set(ctx, *state, jsonData, time.Hour*24).Err()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			//Repite el mismo proceso para extraer los ids de todos los devices para ese estado
			retrievedData, err := use.redis.Get(ctx, *state).Result()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			var retrievedLocations core.RedisDevice
			retrievedLocations, err = utils.DeserializeModel(retrievedData)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			res, err := use.repo.GetTotalTrafficByState(ctx, retrievedLocations.Devices, month)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			res.State = *state
			traffics = append(traffics, res)
		}
	}

	return traffics, err
}

func (use trafficUsecase) GetTotalTrafficByState_N(month string, n int) ([]*model.TrafficState, error) {
	var traffics []*model.TrafficState
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	//Pregunta por todos los estados almacenados en la base de datos
	states, err := use.repoinfo.GetLocationStates(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
			err,
		)
		return nil, err
	}
	//Recorre la lista de estados
	for _, state := range states {
		//Verifica si existe el estado en redis
		exist := utils.VerifyExistence(use.redis, *state, ctx)
		if exist {
			//Si existe, obtiene el json con  los ids de todos los devices para ese estado
			retrievedData, err := use.redis.Get(ctx, *state).Result()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			//Parsea el json obtenido a un objeto LocationsDevice
			var retrievedLocations core.RedisDevice
			retrievedLocations, err = utils.DeserializeModel(retrievedData)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			//Pregunta por el total de trafico por device para ese estado
			res, err := use.repo.GetTotalTrafficByState(ctx, retrievedLocations.Devices, month)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			res.State = *state
			traffics = append(traffics, res)
		} else {
			//Si no existe, pregunta a la base de datos por los ids de todos los devices para ese estado
			var devices []uint
			devicesEntity, err := use.repoinfo.GetDeviceByState(ctx, *state)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			if devicesEntity == nil {
				continue
			}
			for _, deviceEntity := range devicesEntity {
				devices = append(devices, deviceEntity.ID)
			}
			//Crea un objeto LocationsDevice con los ids de todos los devices para ese estado
			locationsDevices := core.RedisDevice{
				Devices: devices,
			}
			//Parsea el objeto LocationsDevice a json
			jsonData, err := utils.SerializeModel(locationsDevices)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			//Guarda el json en redis
			err = use.redis.Set(ctx, *state, jsonData, time.Hour*24).Err()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			//Repite el mismo proceso para extraer los ids de todos los devices para ese estado
			retrievedData, err := use.redis.Get(ctx, *state).Result()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			var retrievedLocations core.RedisDevice
			retrievedLocations, err = utils.DeserializeModel(retrievedData)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTotalTrafficByState - use.repoinfo.GetLocationStates(ctx)",
					err,
				)
				return nil, err
			}
			res, err := use.repo.GetTotalTrafficByState(ctx, retrievedLocations.Devices, month)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					"(trafficUsecase).GetTrafficByODN - use.repo.GetTrafficByODN(ctx)",
					err,
				)
				return nil, err
			}
			res.State = *state
			traffics = append(traffics, res)
		}
	}

	topN := utils.SortTrafficStatesByOut(traffics, n)
	return topN, err
}

func (use trafficUsecase) GetTotalTrafficByOND(month string) ([]*model.TrafficODN, error) {
	var traffics []*model.TrafficODN
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	//Pregunta por todos los ODN almacenados en la base de datos
	odns, err := use.repoinfo.GetAllODN(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(trafficUsecase).GetTotalTrafficByOND - use.repoinfo.GetAllODN(ctx)",
			err,
		)
		return nil, err
	}
	//Recorre la lista de ODN
	for _, odn := range odns {
		//Verifica si existe el ODN en redis
		exist := utils.VerifyExistence(use.redis, *odn, ctx)
		if exist {
			//Si existe, obtiene el json con  los ids de todos los devices para ese ODN
			retrievedData, err := use.redis.Get(ctx, *odn).Result()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - use.redis.Get(ctx, %s)", *odn),
					err,
				)
				return nil, err
			}
			//Parsea el json obtenido a un objeto RedisDevice
			var retrievedODN core.RedisDevice
			retrievedODN, err = utils.DeserializeModel(retrievedData)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_UTILS,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - utils.DeserializeModelODN(%s)", retrievedData),
					err,
				)
				return nil, err
			}
			listUnique := utils.DeleteDuplicate(retrievedODN.Devices)

			//Pregunta por el total de trafico por device para ese ODN
			res, err := use.repo.GetTotalTrafficByODN(ctx, listUnique, month)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - use.repo.GetTotalTrafficByODN(ctx, %v, %s)", retrievedODN.Devices, month),
					err,
				)
				return nil, err
			}
			res.ODN = *odn
			traffics = append(traffics, res)
		} else {
			//Si no existe, pregunta a la base de datos por los ids de todos los devices para ese ODN
			var devices []uint
			idsDevices, err := use.repoinfo.GetDevicesByOND(ctx, *odn)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - use.repo.GetDevicesByOND(ctx, %s)", *odn),
					err,
				)
				return nil, err
			}
			if idsDevices == nil {
				continue
			}
			for _, idDevice := range idsDevices {
				devices = append(devices, *idDevice)
			}
			//Crea un objeto RedisDevice con los ids de todos los devices para ese ODN
			ODNsDevices := core.RedisDevice{
				Devices: devices,
			}
			//Parsea el objeto RedisDevice a json
			jsonData, err := utils.SerializeModel(ODNsDevices)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_UTILS,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - utils.SerializeModel(%v)", ODNsDevices),
					err,
				)
				return nil, err
			}
			//Guarda el json en redis
			err = use.redis.Set(ctx, *odn, jsonData, time.Hour*24).Err()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - use.redis.Set(ctx, %s,%s,%s)", *odn, jsonData, time.Hour*24),
					err,
				)
				return nil, err
			}
			//Repite el mismo proceso para extraer los ids de todos los devices para ese ODN
			retrievedData, err := use.redis.Get(ctx, *odn).Result()
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - use.redis.Get(ctx, %s)", *odn),
					err,
				)
				return nil, err
			}
			//Parsea el json obtenido a un objeto RedisDevice
			var retrievedODN core.RedisDevice
			retrievedODN, err = utils.DeserializeModel(retrievedData)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_UTILS,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - utils.DeserializeModelODN(%s)", retrievedData),
					err,
				)
				return nil, err
			}
			listUnique := utils.DeleteDuplicate(retrievedODN.Devices)
			//Pregunta por el total de trafico por device para ese ODN
			res, err := use.repo.GetTotalTrafficByODN(ctx, listUnique, month)
			if err != nil {
				go use.telegram.SendMessage(
					constants.MODULE_TRAFFIC,
					constants.CATEGORY_DATABASE,
					fmt.Sprintf("(trafficUsecase).GetTotalTrafficByOND - use.repo.GetTotalTrafficByODN(ctx, %v, %s)", retrievedODN.Devices, month),
					err,
				)
				return nil, err
			}
			res.ODN = *odn
			traffics = append(traffics, res)
		}
	}
	return traffics, err
}
