package controller

import (
	"log"
	"sync"

	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/pkg/snmp"
	"github.com/metalpoch/olt-blueprint/update/usecase"
	"gorm.io/gorm"
)

type trafficController struct {
	Measurement usecase.MeasurementUsecase
}

func newTrafficController(db *gorm.DB) *trafficController {
	return &trafficController{
		Measurement: *usecase.NewMeasurementUsecase(db),
	}
}

func deviceUpdater(db *gorm.DB, device *model.DeviceWithOID) (bool, error) {
	var isAlive bool
	checkDevice := &model.CheckDevice{
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

	return isAlive, newDeviceController(db).Usecase.Check(checkDevice)
}

func measurements(ctrl *trafficController, device *model.DeviceWithOID) error {
	var (
		err error
		wg  sync.WaitGroup
	)

	result := map[string]model.Snmp{
		"bw":      model.Snmp{},
		"in":      model.Snmp{},
		"out":     model.Snmp{},
		"ifname":  model.Snmp{},
		"ifdescr": model.Snmp{},
		"ifalias": model.Snmp{},
	}

	oidMap := map[string]string{
		"bw":      device.OidBw,
		"in":      device.OidIn,
		"out":     device.OidOut,
		"ifalias": constants.IF_ALIAS_OID,
		"ifdescr": constants.IF_DESCR_OID,
		"ifname":  constants.IF_NAME_OID,
	}

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

	// TODO: save data in db
	//...
	return nil
}

func Scan(db *gorm.DB, devices []*model.DeviceWithOID) {
	ctrl := newTrafficController(db)
	for _, d := range devices {
		go func(d *model.DeviceWithOID) {
			if ok, err := deviceUpdater(db, d); err != nil {
				log.Printf("deviceUpdaterError: on update the device %d:%v\n", d.ID, err)
				return
			} else if !ok {
				return
			}

			err := measurements(ctrl, d)
			if err != nil {
				log.Printf("deviceUpdaterError: on update the measurement of %d:%v\n", d.ID, err)
				return
			}
			// controll.Interface.Upsert()
		}(d)

	}
}
