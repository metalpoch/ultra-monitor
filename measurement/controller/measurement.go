package controller

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/measurement/constants"
	"github.com/metalpoch/olt-blueprint/measurement/model"
	"github.com/metalpoch/olt-blueprint/measurement/pkg/snmp"
	"github.com/metalpoch/olt-blueprint/measurement/usecase"
	"github.com/metalpoch/olt-blueprint/measurement/utils"
	"gorm.io/gorm"
)

type trafficController struct {
	Measurement usecase.MeasurementUsecase
	Traffic     usecase.TrafficUsecase
}

func newTrafficController(db *gorm.DB, telegram tracking.Telegram) *trafficController {
	return &trafficController{
		Measurement: *usecase.NewMeasurementUsecase(db, telegram),
		Traffic:     *usecase.NewTrafficUsecase(db, telegram),
	}
}

func deviceUpdater(db *gorm.DB, telegram tracking.Telegram, device *model.DeviceWithOID) (bool, error) {
	var isAlive bool
	checkDevice := &model.Device{
		ID:          device.ID,
		SysName:     device.SysName,
		SysLocation: device.SysLocation,
	}

	info, err := snmp.GetInfo(device.IP, device.Community)
	if err != nil {
		log.Printf("deviceUpdaterSNMPError: on device %s with the community %s - %v\n", device.IP, device.Community, err)
	} else {
		isAlive = true
		checkDevice.SysName = info.SysName
		checkDevice.SysLocation = info.SysLocation
	}

	checkDevice.IsAlive = isAlive
	checkDevice.LastCheck = time.Now()
	return isAlive, newDeviceController(db, telegram).Usecase.Check(checkDevice)
}

func measurements(db *gorm.DB, telegram tracking.Telegram, device *model.DeviceWithOID) error {
	measurementUsecase := newTrafficController(db, telegram).Measurement
	trafficUsecase := newTrafficController(db, telegram).Traffic
	var (
		err error
		wg  sync.WaitGroup
	)

	result := model.MapSnmp{
		"bw":      model.Snmp{},
		"in":      model.Snmp{},
		"out":     model.Snmp{},
		"ifname":  model.Snmp{},
		"ifdescr": model.Snmp{},
		"ifalias": model.Snmp{},
	}

	oidMap := map[string]string{
		"bw":      constants.IF_HIGH_SPEED_OID,
		"in":      device.OidIn,
		"out":     device.OidOut,
		"ifalias": constants.IF_ALIAS_OID,
		"ifdescr": constants.IF_DESCR_OID,
		"ifname":  constants.IF_NAME_OID,
	}
	date := time.Now()
	for name, oid := range oidMap {
		wg.Add(1)
		go func(oid string) {
			defer wg.Done()
			res, errSnmp := snmp.Walk(device.IP, device.Community, oid)
			if err != nil {
				err = errSnmp
			} else {
				result[name] = res
			}
		}(oid)
	}
	wg.Wait()

	if err != nil {
		return err
	}

	interfaces, measurements := utils.SnmpElements(device.ID, date, result)
	for idx := range interfaces {
		var isFirstMeasurement bool
		i := interfaces[idx]
		m := measurements[idx]
		err := newInterfaceController(db, telegram).Usecase.Upsert(i)
		if err != nil {
			log.Printf("interfaceUpdaterError: on update the interface %s of deviceID %d:%v\n", i.IfName, device.ID, err)
			continue
		}

		id := i.ID
		m.InterfaceID = id

		old_m, err := measurementUsecase.Get(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isFirstMeasurement = true
			} else {
				log.Printf("measurementGetError: on get the measurement of %d:%v\n", id, err)
			}
		}

		err = measurementUsecase.Upsert(m)
		if err != nil {
			log.Printf("measurementUpdaterError: on update the measurement %d:%v\n", id, err)
		}

		if isFirstMeasurement { // There is no prior "measurement" to obtain the traffic
			continue
		}

		diffTime := uint(m.Date.Sub(old_m.Date).Seconds())

		if err := trafficUsecase.Add(&model.Traffic{
			InterfaceID: id,
			Date:        m.Date,
			Bandwidth:   m.Bandwidth,
			In:          utils.BytesToBbps(old_m.In, m.In, diffTime),
			Out:         utils.BytesToBbps(old_m.Out, m.Out, diffTime),
		}); err != nil {
			log.Println("error guardando el trafico:", err.Error())
		}
	}
	return nil
}

func Scan(db *gorm.DB, telegram tracking.Telegram, devices []*model.DeviceWithOID) {

	for _, d := range devices {
		go func(d *model.DeviceWithOID) {
			if ok, err := deviceUpdater(db, telegram, d); err != nil {
				log.Printf("deviceUpdaterError: on update the device %d:%v\n", d.ID, err)
				return
			} else if !ok {
				return
			}

			err := measurements(db, telegram, d)
			if err != nil {
				log.Printf("deviceUpdaterError: on update the measurement of %d:%v\n", d.ID, err)
				return
			}
		}(d)

	}
}
