package prometheus

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/metalpoch/ultra-monitor/internal/utils"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Prometheus struct {
	client v1.API
}

func NewPrometheusClient(host string) *Prometheus {
	client, err := api.NewClient(api.Config{Address: host})
	if err != nil {
		log.Fatal(err)
	}
	return &Prometheus{client: v1.NewAPI(client)}
}

func (p *Prometheus) PrometheusDeviceScan(ctx context.Context) ([]InfoDevice, error) {
	ifNameVec, err := p.queryVector(ctx, "ifName", time.Now())
	if err != nil {
		return nil, err
	}

	var devices []InfoDevice
	for _, s := range ifNameVec {
		oltIP := s.Labels["instance"]
		ifIndex := s.Labels["ifIndex"]
		oltRegion := s.Labels["region"]
		oltState := s.Labels["state"]
		if oltIP == "" || ifIndex == "" {
			continue
		}
		devices = append(devices, InfoDevice{
			Region:  oltRegion,
			State:   oltState,
			IP:      oltIP,
			IfName:  s.Labels["ifName"],
			IfIndex: utils.ParseInt64(ifIndex),
		})
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices found in Prometheus")
	}

	return devices, nil
}

func (p *Prometheus) QueryPonTraffic(ctx context.Context, date time.Time) ([]TrafficResult, error) {
	var allResults []TrafficResult

	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	for t := start; t.Before(end); t = t.Add(5 * time.Minute) {
		queryIn := `rate(hwGponOltEthernetStatisticReceivedBytes_count[10m])`
		queryOut := `rate(hwGponOltEthernetStatisticSendBytes_count[10m])`
		queryIfName := `ifName`
		queryIfDescr := `ifDescr`
		queryIfAlias := `ifAlias`
		queryIfSpeed := `ifSpeed`
		querySysLocation := `sysLocation`
		querySysName := `sysName`

		bpsInVec, err := p.queryVector(ctx, queryIn, t)
		if err != nil {
			return nil, err
		}
		bpsOutVec, err := p.queryVector(ctx, queryOut, t)
		if err != nil {
			return nil, err
		}
		ifNameVec, _ := p.queryVector(ctx, queryIfName, t)
		ifDescrVec, _ := p.queryVector(ctx, queryIfDescr, t)
		ifAliasVec, _ := p.queryVector(ctx, queryIfAlias, t)
		ifSpeedVec, _ := p.queryVector(ctx, queryIfSpeed, t)
		sysLocVec, _ := p.queryVector(ctx, querySysLocation, t)
		sysNameVec, _ := p.queryVector(ctx, querySysName, t)

		results := make(map[string]*TrafficResult)
		for _, s := range bpsInVec {
			oltIP := s.Labels["instance"]
			ifIndex := s.Labels["ponPortIndex"]
			oltRegion := s.Labels["region"]
			oltState := s.Labels["state"]
			key := oltIP + ":" + ifIndex
			results[key] = &TrafficResult{
				OltIP:     oltIP,
				OltRegion: oltRegion,
				OltState:  oltState,
				IfIndex:   utils.ParseInt64(ifIndex),
				BpsIn:     s.Value * 8,
				Time:      t,
			}
		}
		for _, s := range bpsOutVec {
			oltIP := s.Labels["instance"]
			ifIndex := s.Labels["ponPortIndex"]
			oltRegion := s.Labels["region"]
			oltState := s.Labels["state"]
			key := oltIP + ":" + ifIndex
			if r, ok := results[key]; ok {
				r.BpsOut = s.Value * 8
			} else {
				results[key] = &TrafficResult{
					OltIP:     oltIP,
					OltRegion: oltRegion,
					OltState:  oltState,
					IfIndex:   utils.ParseInt64(ifIndex),
					BpsOut:    s.Value * 8,
					Time:      t,
				}
			}
		}
		for _, s := range ifNameVec {
			oltIP := s.Labels["instance"]
			ifIndex := s.Labels["ifIndex"]
			key := oltIP + ":" + ifIndex
			if r, ok := results[key]; ok {
				r.IfName = s.Labels["ifName"]
			}
		}
		for _, s := range ifDescrVec {
			oltIP := s.Labels["instance"]
			ifIndex := s.Labels["ifIndex"]
			key := oltIP + ":" + ifIndex
			if r, ok := results[key]; ok {
				r.IfDescr = s.Labels["ifDescr"]
			}
		}
		for _, s := range ifAliasVec {
			oltIP := s.Labels["instance"]
			ifIndex := s.Labels["ifIndex"]
			key := oltIP + ":" + ifIndex
			if r, ok := results[key]; ok {
				r.IfAlias = s.Labels["ifAlias"]
			}
		}
		for _, s := range ifSpeedVec {
			oltIP := s.Labels["instance"]
			ifIndex := s.Labels["ifIndex"]
			key := oltIP + ":" + ifIndex
			if r, ok := results[key]; ok {
				r.IfSpeed = s.Value
			}
		}

		for _, s := range sysLocVec {
			oltIP := s.Labels["instance"]
			sysLocation := s.Labels["sysLocation"]
			for _, r := range results {
				if r.OltIP == oltIP {
					r.SysLocation = sysLocation
				}
			}
		}
		for _, s := range sysNameVec {
			oltIP := s.Labels["instance"]
			sysName := s.Labels["sysName"]
			for _, r := range results {
				if r.OltIP == oltIP {
					r.SysName = sysName
				}
			}
		}

		for _, r := range results {
			r.BytesIn = (r.BpsIn / 8) * 300 // 5 minutos = 300 segundos
			r.BytesOut = (r.BpsOut / 8) * 300
			allResults = append(allResults, *r)
		}
	}

	return allResults, nil
}

// queryVector ejecuta una consulta instantánea y devuelve los resultados como []Result
func (p *Prometheus) queryVector(ctx context.Context, query string, ts time.Time) ([]dataProm, error) {
	val, warn, err := p.client.Query(ctx, query, ts)
	if err != nil {
		return nil, err
	}
	if len(warn) > 0 {
		log.Printf("Prometheus warnings: %v", warn)
	}
	vector, ok := val.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", val)
	}
	var vectors []dataProm
	for _, sample := range vector {
		labels := make(map[string]string)
		for k, v := range sample.Metric {
			labels[string(k)] = string(v)
		}
		vectors = append(vectors, dataProm{
			Labels: labels,
			Value:  float64(sample.Value),
			Time:   sample.Timestamp.Time(),
		})
	}
	return vectors, nil
}
