package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/snmp"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type MeasurementController struct {
	usecase *usecase.MeasurementUsecase
	olt     *usecase.OltUsecase
	pon     *usecase.PonUsecase
}

func NewMeasurementController(db *sqlx.DB, cache *cache.Redis) *MeasurementController {
	return &MeasurementController{
		usecase.NewMeasurementUsecase(db),
		usecase.NewOltUsecase(db),
		usecase.NewPonUsecase(db, cache),
	}
}

func (ctrl MeasurementController) calculateDelta(prev, curr uint64) uint64 {
	if curr >= prev {
		return curr - prev
	}
	return (math.MaxUint64-prev)*curr + 1
}

func (ctrl MeasurementController) bytesToBytesPerSec(prev, curr, diffDate uint64) float64 {
	return float64(ctrl.calculateDelta(prev, curr)) / float64(diffDate)
}

func (ctrl MeasurementController) bytesToBbps(prev, curr, bandwidth, diffDate uint64) float64 {
	maxPossible := float64(bandwidth) + (float64(bandwidth) / 10) // +10% de tolerancia
	bps := 8 * ctrl.bytesToBytesPerSec(prev, curr, diffDate)
	if bps > maxPossible {
		return float64(bandwidth)
	}
	return bps
}

func (ctrl MeasurementController) GetIPs() ([]string, error) {
	return ctrl.olt.GetAllIP()
}

func (ctrl MeasurementController) GetPonsBySysname(sysname string) ([]model.Pon, error) {
	return ctrl.pon.GetAllByDevice(sysname)
}

func (ctrl MeasurementController) PonScan(olt model.Olt) {
	client := snmp.NewSnmp(snmp.Config{
		IP:        olt.IP,
		Community: olt.Community,
		Timeout:   time.Duration(10 * time.Second),
		Retries:   5,
	})

	ponData, err := client.PonQuery()
	if err != nil {
		log.Printf("error al ejecutar client.PonQuery en %s: %v", olt.SysName, err)
	}

	for idx, data := range ponData {
		ponID, err := ctrl.usecase.UpsertPon(model.Pon{
			IfIndex: idx,
			IfName:  data.IfName,
			IfDescr: data.IfDescr,
			IfAlias: data.IfAlias,
			OltIP:   olt.IP,
		})
		if err != nil {
			log.Printf("error al ejecutar ctrl.usecase.UpsertPon en %s en el ifIndex %d: %v", olt.SysName, idx, err)
		}

		var isFirstMeasurement bool
		oldData, err := ctrl.usecase.GetTemportalMeasurementPon(ponID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				isFirstMeasurement = true
			} else {
				log.Printf("error al ejecutar ctrl.usecase.GetTemportalMeasurementPon(%d) en %s: %v", idx, olt.SysName, err)
			}
		}

		now := time.Now()
		err = ctrl.usecase.UpsertTemportalMeasurementPon(model.MeasurementPon{
			PonID:     ponID,
			Date:      now,
			Bandwidth: uint64(data.Bandwidth),
			In:        data.CounterBytesIn,
			Out:       data.CounterBytesOut,
		})
		if err != nil {
			log.Printf("error al ctrl.usecase.UpsertTemportalMeasurementPon(...) en %s-%d: %v", olt.SysName, ponID, err)
		}

		if isFirstMeasurement {
			continue
		}

		diffTime := uint64(now.Sub(oldData.Date).Seconds())
		err = ctrl.usecase.InsertTrafficPon(model.TrafficPon{
			PonID:     ponID,
			Bandwidth: float64(data.Bandwidth),
			BpsIn:     ctrl.bytesToBbps(oldData.In, data.CounterBytesIn, uint64(data.Bandwidth), diffTime),
			BpsOut:    ctrl.bytesToBbps(oldData.Out, data.CounterBytesOut, uint64(data.Bandwidth), diffTime),
			BytesIn:   ctrl.bytesToBytesPerSec(oldData.Out, data.CounterBytesOut, diffTime),
			BytesOut:  ctrl.bytesToBytesPerSec(oldData.Out, data.CounterBytesOut, diffTime),
			Date:      now,
		})
		if err != nil {
			log.Println("error saving the traffic:", err.Error())
		}
	}
}

func (ctrl MeasurementController) OntScan(olt model.Olt, ponID int32, idx int64) ([]model.MeasurementOnt, error) {
	client := snmp.NewSnmp(snmp.Config{
		IP:        olt.IP,
		Community: olt.Community,
		Timeout:   time.Duration(5 * time.Second),
		Retries:   3,
	})

	ontData, err := client.OntQuery(idx)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar client.OntQuery(%d) en %s: %v", idx, olt.SysName, err)
	}
	var records []model.MeasurementOnt
	for ontIdx, data := range ontData {
		records = append(records, model.MeasurementOnt{
			PonID:            ponID,
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

func (ctrl MeasurementController) ProcessOntBatchData(measurements []string) error {
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

	return ctrl.usecase.InsertManyOnt(records)
}
