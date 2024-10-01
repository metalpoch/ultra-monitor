package snmp

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/metalpoch/olt-blueprint/measurement/constants"
	"github.com/metalpoch/olt-blueprint/measurement/model"
)

func GetInfo(ip, community string) (model.SnmpInfo, error) {
	query := gosnmp.GoSNMP{
		Target:             ip,
		Community:          community,
		Port:               161,
		Transport:          "udp",
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            0,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
	err := query.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer query.Conn.Close()

	result, err := query.Get([]string{constants.SYSNAME_OID, constants.SYSLOCATION_OID})
	if err != nil {
		return model.SnmpInfo{}, err
	}

	info := model.SnmpInfo{
		SysName:     string(result.Variables[0].Value.([]byte)),
		SysLocation: string(result.Variables[1].Value.([]byte)),
	}

	return info, nil
}

func Walk(ip, community, oid string) (model.Snmp, error) {
	result := model.Snmp{}
	query := gosnmp.GoSNMP{
		Target:             ip,
		Community:          community,
		Port:               161,
		Transport:          "udp",
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            0,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}
	err := query.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer query.Conn.Close()

	if err = query.BulkWalk(oid, func(pdu gosnmp.SnmpPDU) error {
		parts := strings.Split(pdu.Name, ".")
		strID := parts[len(parts)-1]

		id, err := strconv.Atoi(strID)
		if err != nil {
			return err
		}

		switch pdu.Type {

		case gosnmp.OctetString:
			result[id] = string(pdu.Value.([]byte))

		case gosnmp.Counter32:
			result[id] = int(pdu.Value.(uint))

		case gosnmp.Counter64:
			result[id] = int(pdu.Value.(uint64))

		case gosnmp.Gauge32:
			result[id] = int(pdu.Value.(uint))

		case gosnmp.Uinteger32:
			result[id] = int(pdu.Value.(uint))
		default:
			return errors.New("invalid snmp response type")
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return result, err
}
