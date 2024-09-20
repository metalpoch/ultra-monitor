package snmp

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gosnmp/gosnmp"
	"github.com/metalpoch/olt-blueprint/update/constants"
	"github.com/metalpoch/olt-blueprint/update/model"
)

func GetSysname(ip, community string) (string, error) {
	gosnmp.Default.Target = ip
	gosnmp.Default.Community = community

	err := gosnmp.Default.Connect()
	if err != nil {
		return "", err
	}
	defer gosnmp.Default.Conn.Close()

	result, err := gosnmp.Default.Get([]string{constants.SYSNAME_OID})
	if err != nil {
		return "", err
	}

	return string(result.Variables[0].Value.([]byte)), nil
}

func Walk(ip, community, oid string) (model.Snmp, error) {
	response := model.Snmp{}
	gosnmp.Default.Target = ip
	gosnmp.Default.Community = community

	err := gosnmp.Default.Connect()
	if err != nil {
		return nil, err
	}
	defer gosnmp.Default.Conn.Close()

	if err = gosnmp.Default.BulkWalk(oid, func(pdu gosnmp.SnmpPDU) error {
		parts := strings.Split(pdu.Name, ".")
		strID := parts[len(parts)-1]

		id, err := strconv.Atoi(strID)
		if err != nil {
			return err
		}

		switch pdu.Type {

		case gosnmp.OctetString:
			response[uint(id)] = string(pdu.Value.([]byte))
		case gosnmp.Counter32:
			response[uint(id)] = pdu.Value.(uint)
		case gosnmp.Counter64:
			response[uint(id)] = pdu.Value.(uint)
		case gosnmp.Gauge32:
			response[uint(id)] = pdu.Value.(uint)
		case gosnmp.Uinteger32:
			response[uint(id)] = pdu.Value.(uint)
		default:
			return errors.New("invalid snmp response type")
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return response, err
}
