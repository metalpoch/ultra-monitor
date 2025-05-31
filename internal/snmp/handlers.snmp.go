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

type ontHandler []struct {
	oid     string
	handler func(*OntData, gosnmp.SnmpPDU) error
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

func (s snmp) ontOidHandlers(ponIdx string) ontHandler {
	return ontHandler{
		{
			fmt.Sprintf("%s.%s", HwGponDeviceOntDespt, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					ont.Despt = string(value)
				}
				return fmt.Errorf("invalid type for HwGponDeviceOntDespt OID: %v", pdu.Value)
			},
		},
		{
			fmt.Sprintf("%s.%s", HwGponDeviceOntSerialNumber, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					ont.SerialNumber = hex.EncodeToString(value)
					return nil
				}
				return fmt.Errorf("invalid type for HwGponDeviceOntSerialNumber OID: %v", pdu.Value)
			},
		},
		{
			fmt.Sprintf("%s.%s", HwGponDeviceOntLineProfName, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := pdu.Value.([]byte); ok {
					ont.LineProfName = string(value)
					return nil
				}
				return fmt.Errorf("invalid type for HwGponDeviceOntLineProfName OID: %v", pdu.Value)
			},
		},
		{
			fmt.Sprintf("%s.%s", HwGponDeviceOntControlRanging, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := s.toInt64(pdu.Value); ok {
					ont.ControlRanging = int32(value)
					return nil
				}
				return fmt.Errorf("invalid type for HwGponDeviceOntControlRanging OID: %v", pdu.Value)
			},
		},
		{
			fmt.Sprintf("%s.%s", HwGponDeviceOntControlMacCount, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := s.toInt64(pdu.Value); ok {
					ont.ControlMacCount = int8(value)
					return nil
				}
				return fmt.Errorf("invalid type for HwGponDeviceOntControlMacCount OID: %v", pdu.Value)
			},
		},
		{
			fmt.Sprintf("%s.%s", HwGponDeviceOntControlRunStatus, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := s.toInt64(pdu.Value); ok {
					ont.ControlRunStatus = int8(value)
					return nil
				}
				return fmt.Errorf("invalid type for HwGponDeviceOntControlRunStatus OID: %v", pdu.Value)
			},
		},
		{
			fmt.Sprintf("%s.%s", HwGponOntStatisticUpBytes, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := s.toUint64(pdu.Value); ok {
					ont.BytesOut = value
					return nil
				}
				return fmt.Errorf("invalid type for HwGponOntStatisticUpBytes OID: %v", pdu.Value)
			},
		},
		{
			fmt.Sprintf("%s.%s", HwGponOntStatisticDownBytes, ponIdx),
			func(ont *OntData, pdu gosnmp.SnmpPDU) error {
				if value, ok := s.toUint64(pdu.Value); ok {
					ont.BytesIn = value
					return nil
				}
				return fmt.Errorf("invalid type for HwGponOntStatisticDownBytes OID: %v", pdu.Value)
			},
		},
	}
}
