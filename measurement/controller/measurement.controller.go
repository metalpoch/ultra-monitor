package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	commonUsecase "github.com/metalpoch/olt-blueprint/common/usecase"
	"github.com/metalpoch/olt-blueprint/measurement/internal/snmp"
	"github.com/metalpoch/olt-blueprint/measurement/model"
	"github.com/metalpoch/olt-blueprint/measurement/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type (
	measurementController struct {
		Measurement usecase.MeasurementUsecase
		Traffic     commonUsecase.TrafficUsecase
		Device      commonUsecase.DeviceUsecase
		Interface   commonUsecase.InterfaceUsecase
		rdb         *redis.Client
	}
)

func NewMeasurementController(db *gorm.DB, telegram tracking.SmartModule) *measurementController {
	return &measurementController{
		Measurement: *usecase.NewMeasurementUsecase(db, telegram),
		Traffic:     *commonUsecase.NewTrafficUsecase(db, telegram),
		Device:      *commonUsecase.NewDeviceUsecase(db, telegram),
		Interface:   *commonUsecase.NewInterfaceUsecase(db, telegram),
	}
}

func (m measurementController) calculateDelta(prev, curr uint64) uint64 {
	if curr >= prev {
		return curr - prev
	}
	return (math.MaxUint64-prev)*curr + 1
}

func (m measurementController) bytesToBbps(prev, curr, bandwidth, diffDate uint64) uint64 {
	maxPossible := bandwidth + (bandwidth / 10) // +10% de tolerancia
	delta := m.calculateDelta(prev, curr)
	bps := (8 * delta) / diffDate

	if bps > maxPossible {
		return bandwidth
	}

	return bps
}

func (m measurementController) OltScan(device *commonModel.DeviceWithOID) {
	client := snmp.NewSnmp(snmp.Config{
		IP:        device.IP,
		Community: device.Community,
		Timeout:   time.Duration(10 * time.Second),
		Retries:   5,
	})

	ponData, err := client.PonQuery()
	if err != nil {
		log.Printf("error al ejecutar client.PonQuery en %s: %v", device.SysName, err)
	}

	for idx, data := range ponData {
		interfaceID, err := m.Interface.Upsert(&commonModel.Interface{
			IfIndex:   idx,
			IfName:    data.IfName,
			IfDescr:   data.IfDescr,
			IfAlias:   data.IfAlias,
			DeviceID:  device.ID,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			log.Printf("error al ejecutar m.Interface.Upsert en %s en el ifIndex %d: %v", device.SysName, idx, err)
		}

		var isFirstMeasurement bool
		oldData, err := m.Measurement.GetOlt(interfaceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isFirstMeasurement = true
			}
			log.Printf("error al ejecutar m.Measurement.Get(%d) en %s: %v", idx, device.SysName, err)
		}

		now := time.Now()
		err = m.Measurement.UpsertOlt(&model.MeasurementOlt{
			InterfaceID: interfaceID,
			Date:        now,
			Bandwidth:   uint64(data.Bandwidth),
			In:          data.CounterBytesIn,
			Out:         data.CounterBytesOut,
		})
		if err != nil {
			log.Printf("error al ejecutar m.Measurement.Upsert en %s-%v: %v", device.SysName, &model.MeasurementOlt{
				InterfaceID: interfaceID,
				Date:        now,
				Bandwidth:   uint64(data.Bandwidth),
				In:          data.CounterBytesIn,
				Out:         data.CounterBytesOut,
			}, err)
		}

		if isFirstMeasurement {
			continue
		}

		diffTime := uint64(now.Sub(oldData.Date).Seconds())
		err = m.Traffic.Add(&commonModel.Traffic{
			InterfaceID: interfaceID,
			Date:        now,
			Bandwidth:   uint64(data.Bandwidth),
			In:          m.bytesToBbps(oldData.In, data.CounterBytesIn, uint64(data.Bandwidth), diffTime),
			Out:         m.bytesToBbps(oldData.Out, data.CounterBytesOut, uint64(data.Bandwidth), diffTime),
			BytesIn:     m.calculateDelta(oldData.In, data.CounterBytesIn),
			BytesOut:    m.calculateDelta(oldData.Out, data.CounterBytesOut),
		})
		if err != nil {
			log.Println("error saving the traffic:", err.Error())
		}
	}
}

func (m measurementController) OntScan(device *commonModel.DeviceWithOID, interfaceID, idx uint64) ([]model.MeasurementOnt, error) {
	client := snmp.NewSnmp(snmp.Config{
		IP:        device.IP,
		Community: device.Community,
		Timeout:   time.Duration(5 * time.Second),
		Retries:   3,
	})

	ontData, err := client.OntQuery(idx)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar client.OntQuery(%d) en %s: %v", idx, device.SysName, err)
	}
	var records []model.MeasurementOnt
	for ontIdx, data := range ontData {
		records = append(records, model.MeasurementOnt{
			InterfaceID:      interfaceID,
			Idx:              ontIdx,
			Despt:            data.Despt,
			SerialNumber:     data.SerialNumber,
			LineProfName:     data.LineProfName,
			OltDistance:      data.ControlRanging,
			ControlMacCount:  data.ControlMacCount,
			ControlRunStatus: data.ControlRunStatus,
			BytesIn:          data.BytesIn,
			BytesOut:         data.BytesOut,
			Date:             time.Now(),
		})
	}

	return records, err
}

func (m measurementController) ProcessOntBatchData(measurements []string) error {
	var records []model.MeasurementOnt

	for _, item := range measurements {
		var record model.MeasurementOnt
		if err := json.Unmarshal([]byte(item), &record); err != nil {
			log.Println("error unmarshalling SNMP data:", err)
			continue
		}
		records = append(records, record)
	}

	if len(records) == 0 {
		return nil
	}
	return m.Measurement.InsertManyOnt(records)
}

func (m measurementController) calculateDeltaOnt(value, prev, curr *uint64) {
	if prev == nil || curr == nil {
		return
	}

	delta := m.calculateDelta(*prev, *curr)
	value = &delta
}

func (m measurementController) bytesToBbpsOnt(value, prev, curr *uint64, diffDate uint64) {
	if prev == nil || curr == nil {
		return
	}

	delta := m.calculateDelta(*prev, *curr)
	res := (8 * delta) / diffDate

	value = &res
}
