package controller

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/pkg/snmp"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"github.com/metalpoch/olt-blueprint/update/usecase"
	"github.com/metalpoch/olt-blueprint/update/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type handlers struct {
	oltElement     usecase.OltElementUsecase
	countUsecase   usecase.CountUsecase
	trafficUsecase usecase.TrafficUsecase
}

func newHandler(client *mongo.Client) *handlers {
	return &handlers{
		oltElement:     usecase.OltElementUsecase(usecase.NewOltElementUsecase(repository.NewOltElementRepository(client))),
		countUsecase:   usecase.CountUsecase(usecase.NewCountUsecase(repository.NewCountRepository(client))),
		trafficUsecase: usecase.TrafficUsecase(usecase.NewTrafficUsecase(repository.NewTrafficRepository(client))),
	}
}

func TrafficUpdate(devices []*model.Device, client *mongo.Client) {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	hdlr := newHandler(client)

	length := len(devices)
	if length == 0 {
		log.Fatalln("no device data to scan")
	}

	var wg sync.WaitGroup
	wg.Add(length)

	for _, device := range devices {
		go func(d *model.Device) {
			unix_time := time.Now().Unix()
			measurements := snmp.Measurements(d)

			for idx, ifname := range measurements.IfName {

				if !strings.HasPrefix(ifname, "GPON") {
					log.Printf("error the ifname %s no is GPON\n", ifname)
					continue
				}

				elem := utils.GetGPONInterface(d.Sysname, ifname)
				if _, err := hdlr.oltElement.Create(elem); !mongo.IsDuplicateKeyError(err) && err != nil {
					log.Printf("error when trying to save %s %s: %s\n", d.Sysname, ifname, err.Error())
				}

				prevCount, err := hdlr.countUsecase.Find(d.Sysname, ifname)
				if err != mongo.ErrNoDocuments && err != nil {
					log.Printf("error when trying to get previous measurement of the element %s %s: %s\n", d.Sysname, ifname, err.Error())
					continue
				}

				_, err = hdlr.countUsecase.Create(model.Count{
					OLT:       d.Sysname,
					Interface: ifname,
					Date:      unix_time,
					BytesIn:   measurements.ByteIn[idx],
					BytesOut:  measurements.ByteOut[idx],
					Bandwidth: prevCount.Bandwidth,
				})

				if err != nil {
					log.Printf("error when trying to save the measurement of %s %s: %s\n", d.Sysname, ifname, err.Error())
					continue
				}

				// if prevCount dont exist, skip after save the measurement.
				if prevCount.Date == 0 {
					continue
				}

				countDiff := model.CountDiff{
					OLT:           d.Sysname,
					Interface:     ifname,
					PrevDate:      prevCount.Date,
					PrevBytesIn:   prevCount.BytesIn,
					PrevBytesOut:  prevCount.BytesOut,
					CurrDate:      unix_time,
					CurrBytesIn:   measurements.ByteIn[idx],
					CurrBytesOut:  measurements.ByteOut[idx],
					CurrBandwidth: measurements.Bandwidth[idx],
				}

				if _, err := hdlr.trafficUsecase.Create(countDiff); err != nil {
					log.Printf("error when trying to save the traffic of %s on day %d: %s\n", ifname, countDiff.CurrDate, err.Error())
				}
			}

		}(device)

	}
	wg.Wait()
}
