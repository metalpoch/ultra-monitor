package snmp

import (
	"encoding/hex"
	"fmt"

	"github.com/gosnmp/gosnmp"
)

type oltHandler []struct {
	oid     string
	handler func(*OltData, gosnmp.SnmpPDU) error
}

type ponHandler []struct {
	oid     string
	handler func(*PonData, gosnmp.SnmpPDU) error
}

type ontSerialAndDesptHandler []struct {
	oid     string
	handler func(*OntSerialsAndDespts, gosnmp.SnmpPDU) error
}

func (s snmp) oltOidHandlers() oltHandler {
	return oltHandler{
		{
			SysName,
			func(olt *OltData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					olt.SysName = string(value)
					return nil
				}
				return fmt.Errorf("invalid type for SysName OID: %v", pdu.Value)
			},
		},
		{
			SysLocation,
			func(olt *OltData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					olt.SysLocation = string(value)
					return nil
				}
				return fmt.Errorf("invalid type for SysLocation OID: %v", pdu.Value)
			},
		},
	}
}

func (s snmp) ponOidHandlers() ponHandler {
	return ponHandler{
		{
			IfName,
			func(pon *PonData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					pon.IfName = string(value)
					return nil
				}
				return fmt.Errorf("invalid type for IfName OID: %v", pdu.Value)
			},
		},

		{
			IfDescr,
			func(pon *PonData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					pon.IfDescr = string(value)
					return nil
				}
				return fmt.Errorf("invalid type for IfDescr OID: %v", pdu.Value)
			},
		},
		{
			IfAlias,
			func(pon *PonData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					pon.IfAlias = string(value)
					return nil
				}
				return fmt.Errorf("invalid type for IfAlias OID: %v", pdu.Value)
			},
		},
		{
			IfSpeed,
			func(pon *PonData, pdu gosnmp.SnmpPDU) error {
				value, ok := s.toInt64(pdu.Value)
				if !ok {
					return fmt.Errorf("invalid type for IfSpeed OID: %v", pdu.Value)
				}
				pon.Bandwidth = value
				return nil
			},
		},
		{
			HwGponOltEthernetStatisticReceivedBytes,
			func(pon *PonData, pdu gosnmp.SnmpPDU) error {
				value, ok := s.toUint64(pdu.Value)
				if !ok {
					return fmt.Errorf("invalid type for HwGponOltEthernetStatisticReceivedBytes OID: %v", pdu.Value)
				}
				pon.CounterBytesIn = value
				return nil
			},
		},
		{
			HwGponOltEthernetStatisticSendBytes,
			func(pon *PonData, pdu gosnmp.SnmpPDU) error {
				value, ok := s.toUint64(pdu.Value)
				if !ok {
					return fmt.Errorf("invalid type for HwGponOltEthernetStatisticSendBytes OID: %v", pdu.Value)
				}
				pon.CounterBytesOut = value
				return nil
			},
		},
	}
}

func (s snmp) ontOidSerialDesptHandlers() ontSerialAndDesptHandler {
	return ontSerialAndDesptHandler{
		{
			HwGponDeviceOntSerialNumber,
			func(ont *OntSerialsAndDespts, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					ont.SerialNumber = hex.EncodeToString(value)
					return nil
				}
				return fmt.Errorf("invalid type for Depst OID: %v", pdu.Value)
			},
		},
		{
			HwGponDeviceOntDespt,
			func(ont *OntSerialsAndDespts, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					ont.Despt = string(value)
					return nil
				}
				return fmt.Errorf("invalid type for Serial OID: %v", pdu.Value)
			},
		},
	}
}
