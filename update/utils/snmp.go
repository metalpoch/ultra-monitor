package utils

import (
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
)

func MeasurementToEntitys(date time.Time, measurements model.MapSnmp) []*entity.Measurement {
	var rows []*entity.Measurement
	fields := [3]string{"bw", "in", "out"}

	for ifIndex := range measurements["ifname"] {
		row := new(entity.Measurement)
		row.InterfaceID = uint(ifIndex)
		row.Date = date

		for i := 0; i < 3; i++ {
			f := fields[i]
			switch f {
			case "bw":
				row.Bandwidth = uint(measurements[f][ifIndex].(int))
			case "out":
				row.Out = uint(measurements[f][ifIndex].(int))
			case "in":
				row.In = uint(measurements[f][ifIndex].(int))
			}
		}
		fmt.Println(row)
		rows = append(rows, row)
	}
	return rows
}

// func InterfacesToEntitys(deviceID uint, date time.Time, measurements model.MapSnmp) []*model.Interface {
// 	var rows []*model.Interface
// 	fields := [3]string{"ifname", "ifdescr", "ifalias"}
// 	for ifIndex := range measurements["ifname"] {
// 		m := new(model.Measurement)
// 		i := new(model.Interface)

// 		i.DeviceID = deviceID
// 		i.IfIndex = uint(ifIndex)
// 		i.UpdatedAt = date

// 		for idx := 0; idx < 3; idx++ {
// 			f := fields[idx]
// 			switch f {
// 			case "ifname":
// 				i.IfName = measurements[f][ifIndex].(string)
// 			case "ifdescr":
// 				i.IfDescr = measurements[f][ifIndex].(string)
// 			case "ifalias":
// 				i.IfAlias = measurements[f][ifIndex].(string)

// 			}
// 		}
// 		rows = append(rows, i)
// 	}
// 	return rows
// }

func isString(field string) bool {
	if field == "ifname" || field == "ifdescr" || field == "ifalias" {
		return true
	}
	return false
}

func SnmpElements(deviceID uint, date time.Time, snmp model.MapSnmp) ([]*model.Interface, []*model.Measurement) {
	var interfaces []*model.Interface
	var measurements []*model.Measurement

	fields := [6]string{"ifname", "ifdescr", "ifalias", "bw", "in", "out"}
	for ifIndex := range snmp["ifname"] {
		i := new(model.Interface)
		i.DeviceID = deviceID
		i.IfIndex = uint(ifIndex)
		i.UpdatedAt = date

		m := new(model.Measurement)
		m.Date = date

		for idx := 0; idx < 6; idx++ {
			f := fields[idx]
			if isString(f) {
				switch f {
				case "ifname":
					i.IfName = snmp[f][ifIndex].(string)
				case "ifdescr":
					i.IfDescr = snmp[f][ifIndex].(string)
				case "ifalias":
					i.IfAlias = snmp[f][ifIndex].(string)
				}

			} else {
				switch f {
				case "bw":
					m.Bandwidth = uint(snmp[f][ifIndex].(int))
				case "out":
					m.Out = uint(snmp[f][ifIndex].(int))
				case "in":
					m.In = uint(snmp[f][ifIndex].(int))
				}
			}
		}
		measurements = append(measurements, m)
		interfaces = append(interfaces, i)
	}
	return interfaces, measurements
}
