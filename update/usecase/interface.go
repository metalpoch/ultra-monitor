package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"gorm.io/gorm"
)

type interfaceUsecase struct {
	repo repository.InterfaceRepository
}

func NewInterfaceUsecase(db *gorm.DB) *interfaceUsecase {
	return &interfaceUsecase{repository.NewInterfaceRepository(db)}
}

func (use interfaceUsecase) Upsert(element *model.Interface) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return use.repo.Upsert(ctx, (*entity.Interface)(element))
}

func (use interfaceUsecase) GetAllByDevice(id uint) ([]*model.Interface, error) {
	var interfaces []*model.Interface
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetAllByDevice(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		interfaces = append(interfaces, (*model.Interface)(e))
	}

	return interfaces, nil
}

func (use interfaceUsecase) GetAll() ([]*model.Interface, error) {
	var interfaces []*model.Interface
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		interfaces = append(interfaces, (*model.Interface)(e))
	}

	return interfaces, nil
}

// func (use interfaceUsecase) Scan(device model.DeviceWithOID) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	var (
// 		names      model.Snmp
// 		descrs     model.Snmp
// 		alias      model.Snmp
// 		bw         model.Snmp
// 		in         model.Snmp
// 		out        model.Snmp
// 		wg         sync.WaitGroup
// 		date       time.Time = time.Now()
// 		isTemporal bool      = true
// 		err        error
// 	)

// 	wg.Add(6)
// 	go func(oid string) {
// 		defer wg.Done()
// 		if res, snmpErr := snmp.Walk(device.IP, device.Community, oid); snmpErr != nil {
// 			err = snmpErr
// 		} else {
// 			names = res
// 		}
// 	}(constants.IF_NAME_OID)

// 	go func(oid string) {
// 		defer wg.Done()
// 		if res, snmpErr := snmp.Walk(device.IP, device.Community, oid); snmpErr != nil {
// 			err = snmpErr
// 		} else {
// 			descrs = res
// 		}
// 	}(constants.IF_DESCR_OID)

// 	go func(oid string) {
// 		defer wg.Done()
// 		if res, snmpErr := snmp.Walk(device.IP, device.Community, oid); snmpErr != nil {
// 			err = snmpErr
// 		} else {
// 			alias = res
// 		}
// 	}(constants.IF_ALIAS_OID)

// 	go func(oid string) {
// 		defer wg.Done()
// 		if res, snmpErr := snmp.Walk(device.IP, device.Community, oid); snmpErr != nil {
// 			err = snmpErr
// 		} else {
// 			bw = res
// 		}
// 	}(device.OidBw)

// 	go func(oid string) {
// 		defer wg.Done()
// 		if res, snmpErr := snmp.Walk(device.IP, device.Community, oid); snmpErr != nil {
// 			err = snmpErr
// 		} else {
// 			in = res
// 		}
// 	}(device.OidIn)

// 	go func(oid string) {
// 		defer wg.Done()
// 		if res, snmpErr := snmp.Walk(device.IP, device.Community, oid); snmpErr != nil {
// 			err = snmpErr
// 		} else {
// 			out = res
// 		}
// 	}(device.OidOut)
// 	wg.Wait()

// 	if err != nil {
// 		log.Println("error on", device.IP, err.Error())
// 		return
// 	}

// 	measurements := []entity.Measurement{}
// 	for id, ifname := range names {
// 		result, err := use.repo.UpsertInterface(ctx, &entity.Interface{
// 			ID:        int32(id),
// 			IfName:    ifname.(string),
// 			IfDescr:   descrs[id].(string),
// 			IfAlias:   alias[id].(string),
// 			CreatedAt: date,
// 			UpdatedAt: date,
// 		})

// 		fmt.Println(result, err)

// 		measurements = append(measurements, entity.Measurement{
// 			Date:        date,
// 			Bandwidth:   int64(bw[id].(int)),
// 			In:          int64(in[id].(int)),
// 			Out:         int64(out[id].(int)),
// 			InterfaceID: int32(id),
// 		})
// 	}

// 	if err := use.repo.SaveMeasurement(ctx, isTemporal, device.IP, measurements); err != nil {
// 		log.Println("error on try save the element", device.IP, err.Error())
// 	}
// }
