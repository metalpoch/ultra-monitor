package utils

/*
 *import (
 *        "fmt"
 *        "time"
 *
 *        "github.com/metalpoch/olt-blueprint/common/entity"
 *        commonModel "github.com/metalpoch/olt-blueprint/common/model"
 *        "github.com/metalpoch/olt-blueprint/measurement/internal/snmp"
 *        "github.com/metalpoch/olt-blueprint/measurement/model"
 *)
 *
 *func MeasurementToEntitys(date time.Time, measurements snmp.MapSnmp) []*entity.Measurement {
 *        var rows []*entity.Measurement
 *        fields := [3]string{"bw", "in", "out"}
 *
 *        for ifIndex := range measurements["ifname"] {
 *                row := new(entity.Measurement)
 *                row.InterfaceID = uint(ifIndex)
 *                row.Date = date
 *
 *                for i := 0; i < 3; i++ {
 *                        f := fields[i]
 *                        switch f {
 *                        case "bw":
 *                                row.Bandwidth = uint64(measurements[f][ifIndex].(int))
 *                        case "out":
 *                                row.Out = uint64(measurements[f][ifIndex].(int))
 *                        case "in":
 *                                row.In = uint64(measurements[f][ifIndex].(int))
 *                        }
 *                }
 *                fmt.Println(row)
 *                rows = append(rows, row)
 *        }
 *        return rows
 *}
 *
 *func isString(field string) bool {
 *        if field == "ifname" || field == "ifdescr" || field == "ifalias" {
 *                return true
 *        }
 *        return false
 *}
 *
 *func SnmpElements(deviceID uint, date time.Time, snmp snmp.MapSnmp) ([]*commonModel.Interface, []*model.MeasurementOlt) {
 *        var interfaces []*commonModel.Interface
 *        var measurements []*model.MeasurementOlt
 *
 *        fields := [6]string{"ifname", "ifdescr", "ifalias", "bw", "in", "out"}
 *        for ifIndex := range snmp["ifname"] {
 *                i := new(commonModel.Interface)
 *                i.DeviceID = deviceID
 *                i.IfIndex = uint(ifIndex)
 *                i.UpdatedAt = date
 *
 *                m := new(model.MeasurementOlt)
 *                m.Date = date
 *
 *                for idx := 0; idx < 6; idx++ {
 *                        f := fields[idx]
 *                        if snmp[f][ifIndex] == nil {
 *                                continue
 *                        }
 *                        if isString(f) {
 *                                switch f {
 *                                case "ifname":
 *                                        i.IfName = snmp[f][ifIndex].(string)
 *                                case "ifdescr":
 *                                        i.IfDescr = snmp[f][ifIndex].(string)
 *                                case "ifalias":
 *                                        i.IfAlias = snmp[f][ifIndex].(string)
 *                                }
 *                        } else {
 *                                switch f {
 *                                case "bw":
 *                                        m.Bandwidth = uint64(snmp[f][ifIndex].(int))
 *                                case "out":
 *                                        m.Out = uint64(snmp[f][ifIndex].(int))
 *                                case "in":
 *                                        m.In = uint64(snmp[f][ifIndex].(int))
 *                                }
 *                        }
 *                }
 *                measurements = append(measurements, m)
 *                interfaces = append(interfaces, i)
 *        }
 *        return interfaces, measurements
 *}*/
