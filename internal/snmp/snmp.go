package snmp

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
)

const (
	// Iftable & Systable
	SysName     string = ".1.3.6.1.2.1.1.5.0"
	SysLocation string = ".1.3.6.1.2.1.1.6.0"
	IfName      string = ".1.3.6.1.2.1.31.1.1.1.1"
	IfSpeed     string = ".1.3.6.1.2.1.2.2.1.5"
	IfDescr     string = ".1.3.6.1.2.1.2.2.1.2"
	IfAlias     string = ".1.3.6.1.2.1.31.1.1.1.18"

	// OLT pon Traffic
	HwGponOltEthernetStatisticReceivedBytes string = ".1.3.6.1.4.1.2011.6.128.1.1.4.21.1.15"
	HwGponOltEthernetStatisticSendBytes     string = ".1.3.6.1.4.1.2011.6.128.1.1.4.21.1.30"

	// ONT queries
	HwGponDeviceOntDespt            string = ".1.3.6.1.4.1.2011.6.128.1.1.2.43.1.9"
	HwGponDeviceOntSerialNumber     string = ".1.3.6.1.4.1.2011.6.128.1.1.2.43.1.3"
	HwGponDeviceOntLineProfName     string = ".1.3.6.1.4.1.2011.6.128.1.1.2.43.1.7"
	HwGponDeviceOntControlRanging   string = ".1.3.6.1.4.1.2011.6.128.1.1.2.46.1.20"
	HwGponDeviceOntControlMacCount  string = ".1.3.6.1.4.1.2011.6.128.1.1.2.46.1.21"
	HwGponDeviceOntControlRunStatus string = ".1.3.6.1.4.1.2011.6.128.1.1.2.46.1.15"
	HwGponOntStatisticUpBytes       string = ".1.3.6.1.4.1.2011.6.128.1.1.4.23.1.3"
	HwGponOntStatisticDownBytes     string = ".1.3.6.1.4.1.2011.6.128.1.1.4.23.1.4"
	hwGponOntOpticalDdmTemperature  string = "1.3.6.1.4.1.2011.6.128.1.1.2.51.1.1"
	HwGponOntOpticalDdmTxPower      string = "1.3.6.1.4.1.2011.6.128.1.1.2.51.1.3"
	HwGponOntOpticalDdmRxPower      string = "1.3.6.1.4.1.2011.6.128.1.1.2.51.1.4"
)

type snmp struct {
	client *gosnmp.GoSNMP
}

type Config struct {
	IP        string
	Community string
	Timeout   time.Duration
	Retries   int
}

type OltData struct {
	SysName     string
	SysLocation string
}

type PonData struct {
	IfName          string
	IfDescr         string
	IfAlias         string
	CounterBytesIn  uint64
	CounterBytesOut uint64
	Bandwidth       int64
}

type OntSerialsAndDespts struct {
	Despt        string
	SerialNumber string
}

type OntData struct {
	Despt            string
	SerialNumber     string
	LineProfName     string
	ControlRanging   int32
	ControlMacCount  int8
	ControlRunStatus int8
	BytesIn          uint64
	BytesOut         uint64
	Temperature      int32
	Tx               int32
	Rx               int32
}

func NewSnmp(config Config) *snmp {
	return &snmp{
		client: &gosnmp.GoSNMP{
			Target:             config.IP,
			Port:               161,
			Community:          config.Community,
			Version:            gosnmp.Version2c,
			Timeout:            config.Timeout,
			Retries:            config.Retries,
			ExponentialTimeout: false,
		},
	}
}

func (s snmp) extractOntIdx(fullOID string) string {
	lastDot := strings.LastIndex(fullOID, ".")
	if lastDot == -1 {
		return ""
	}

	return fullOID[lastDot+1:]
}

func (s snmp) toUint64(value any) (uint64, bool) {
	if !gosnmp.ToBigInt(value).IsUint64() {
		return 0, false
	}
	return gosnmp.ToBigInt(value).Uint64(), true
}

func (s snmp) toInt64(value any) (int64, bool) {
	if !gosnmp.ToBigInt(value).IsInt64() {
		return 0, false
	}
	return gosnmp.ToBigInt(value).Int64(), true
}

func (s snmp) OltSysQuery() (*OltData, error) {
	err := s.client.Connect()
	if err != nil {
		return nil, fmt.Errorf("conexión fallida: %v", err)
	}
	defer s.client.Conn.Close()

	oidHandlers := s.oltOidHandlers()
	errChan := make(chan error, len(oidHandlers))
	data := new(OltData)

	var mutex sync.Mutex
	var wg sync.WaitGroup
	for _, oidHandler := range oidHandlers {
		wg.Add(1)

		go func(oid string, handler func(*OltData, gosnmp.SnmpPDU) error) {
			defer wg.Done()
			err := s.client.BulkWalk(oid, func(pdu gosnmp.SnmpPDU) error {
				index := s.extractOntIdx(pdu.Name)
				if index == "" {
					return nil
				}

				mutex.Lock()
				defer mutex.Unlock()

				if err := oidHandler.handler(data, pdu); err != nil {
					log.Printf("Error on proccess OID %s: %v", pdu.Name, err)
					return nil
				}
				return nil
			})
			if err != nil {
				errChan <- fmt.Errorf("SNMP Walk error for OID %s: %v", oid, err)
			}
		}(oidHandler.oid, oidHandler.handler)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s snmp) PonQuery() (map[int64]PonData, error) {
	err := s.client.Connect()
	if err != nil {
		return nil, fmt.Errorf("conexión fallida: %v", err)
	}
	defer s.client.Conn.Close()
	data := make(map[int64]PonData)

	for _, oidHandler := range s.ponOidHandlers() {
		err = s.client.BulkWalk(oidHandler.oid, func(pdu gosnmp.SnmpPDU) error {
			index := s.extractOntIdx(pdu.Name)
			if index == "" {
				log.Printf("Error on proccess OID %s: index empty", pdu.Name)
				return nil
			}

			idx, err := strconv.Atoi(index)
			if err != nil {
				log.Printf("Error on proccess OID %s in the index %s: %v", pdu.Name, index, err)
				return err
			}

			pon := data[int64(idx)]
			if err := oidHandler.handler(&pon, pdu); err != nil {
				log.Printf("Error on proccess OID %s in the index %s: %v", pdu.Name, index, err)
				return nil
			}
			data[int64(idx)] = pon
			return nil
		})
	}

	return data, err
}

func (s snmp) OntSerialsAndDespts(ponIdx int64) (map[int64]OntSerialsAndDespts, error) {
	err := s.client.Connect()
	if err != nil {
		return nil, fmt.Errorf("conexión fallida: %v", err)
	}
	defer s.client.Conn.Close()

	data := make(map[int64]OntSerialsAndDespts)

	for _, oidHandler := range s.ontOidSerialDesptHandlers() {
		err = s.client.BulkWalk(fmt.Sprintf("%s.%d", oidHandler.oid, ponIdx), func(pdu gosnmp.SnmpPDU) error {
			index := s.extractOntIdx(pdu.Name)
			if index == "" {
				log.Printf("Error on proccess OID %s: index empty", pdu.Name)
				return nil
			}

			idx, err := strconv.Atoi(index)
			if err != nil {
				log.Printf("Error on proccess OID %s in the index %s: %v", pdu.Name, index, err)
				return err
			}

			ont := data[int64(idx)]
			if err := oidHandler.handler(&ont, pdu); err != nil {
				log.Printf("Error on proccess OID %s in the index %s: %v", pdu.Name, index, err)
				return nil
			}
			data[int64(idx)] = ont
			return nil
		})
	}

	return data, err
}

func (s snmp) OntQuery(ponIdx int64, ontIdx uint8) (OntData, error) {
	err := s.client.Connect()
	if err != nil {
		return OntData{}, fmt.Errorf("conexión fallida: %v", err)
	}
	defer s.client.Conn.Close()

	var ont OntData

	oids := []string{
		fmt.Sprintf("%s.%d.%d", HwGponDeviceOntDespt, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponDeviceOntSerialNumber, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponDeviceOntLineProfName, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponDeviceOntControlRanging, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponDeviceOntControlMacCount, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponDeviceOntControlRunStatus, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponOntStatisticUpBytes, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponOntStatisticDownBytes, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", hwGponOntOpticalDdmTemperature, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponOntOpticalDdmTxPower, ponIdx, ontIdx),
		fmt.Sprintf("%s.%d.%d", HwGponOntOpticalDdmRxPower, ponIdx, ontIdx),
	}

	result, err := s.client.Get(oids)
	if err != nil {
		return OntData{}, fmt.Errorf("error en SNMP Get: %v", err)
	}

	for i, pdu := range result.Variables {
		if pdu.Type == gosnmp.NoSuchInstance || pdu.Type == gosnmp.NoSuchObject {
			continue
		}

		switch i {
		case 0: // HwGponDeviceOntDespt
			if value, ok := pdu.Value.([]byte); ok {
				ont.Despt = string(value)
			}
		case 1: // HwGponDeviceOntSerialNumber
			if value, ok := pdu.Value.([]byte); ok {
				ont.SerialNumber = hex.EncodeToString(value)
			}
		case 2: // HwGponDeviceOntLineProfName
			if value, ok := pdu.Value.([]byte); ok {
				ont.LineProfName = string(value)
			}
		case 3: // HwGponDeviceOntControlRanging
			if value, ok := s.toInt64(pdu.Value); ok {
				ont.ControlRanging = int32(value)
			}
		case 4: // HwGponDeviceOntControlMacCount
			if value, ok := s.toInt64(pdu.Value); ok {
				ont.ControlMacCount = int8(value)
			}
		case 5: // HwGponDeviceOntControlRunStatus
			if value, ok := s.toInt64(pdu.Value); ok {
				ont.ControlRunStatus = int8(value)
			}
		case 6: // HwGponOntStatisticUpBytes
			if value, ok := s.toUint64(pdu.Value); ok {
				ont.BytesOut = value
			}
		case 7: // HwGponOntStatisticDownBytes
			if value, ok := s.toUint64(pdu.Value); ok {
				ont.BytesIn = value
			}
		case 8: // hwGponOntOpticalDdmTemperature
			if value, ok := s.toInt64(pdu.Value); ok {
				ont.Temperature = int32(value)
			}
		case 9: // HwGponOntOpticalDdmTxPower
			if value, ok := s.toInt64(pdu.Value); ok {
				ont.Tx = int32(value)
			}
		case 10: // HwGponOntOpticalDdmRxPower
			if value, ok := s.toInt64(pdu.Value); ok {
				ont.Rx = int32(value)
			}
		}
	}

	return ont, nil
}
