package snmp

import (
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

type OntData struct {
	Despt            string
	SerialNumber     string
	LineProfName     string
	ControlRanging   int32
	ControlMacCount  int8
	ControlRunStatus int8
	BytesIn          uint64
	BytesOut         uint64
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

func (s snmp) toUint64(value interface{}) (uint64, bool) {
	if !gosnmp.ToBigInt(value).IsUint64() {
		return 0, false
	}
	return gosnmp.ToBigInt(value).Uint64(), true
}

func (s snmp) toInt64(value interface{}) (int64, bool) {
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
				return nil
			}

			idx, err := strconv.Atoi(index)
			if err != nil {
				return err
			}

			pon := data[int64(idx)]
			if err := oidHandler.handler(&pon, pdu); err != nil {
				log.Printf("Error on proccess OID %s: %v", pdu.Name, err)
				return nil
			}
			data[int64(idx)] = pon
			return nil
		})
	}

	return data, err
}

func (s snmp) OntQuery(ponIdx int64) (map[int64]OntData, error) {
	err := s.client.Connect()
	if err != nil {
		return nil, fmt.Errorf("conexión fallida: %v", err)
	}
	defer s.client.Conn.Close()

	oidHandlers := s.ontOidHandlers(ponIdx)
	data := make(map[int64]OntData)

	for _, oidHandler := range oidHandlers {
		err = s.client.BulkWalk(oidHandler.oid, func(pdu gosnmp.SnmpPDU) error {
			index := s.extractOntIdx(pdu.Name)
			if index == "" {
				return nil
			}

			idx, err := strconv.Atoi(index)
			if err != nil {
				return err
			}

			ont := data[int64(idx)]
			if err := oidHandler.handler(&ont, pdu); err != nil {
				return err
			}
			data[int64(idx)] = ont
			return nil
		})
	}

	return data, err
}
