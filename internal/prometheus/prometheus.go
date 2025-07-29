package prometheus

import (
	"context"
	"fmt"
	"log"
	"strings"
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

func (p *Prometheus) PrometheusTrafficRegion(ctx context.Context, initDate, finalDate time.Time) (map[string][]*Traffic, error) {
	queryBW := "sum(ifSpeed) by (region)"
	queryBpsIn := "sum(rate(hwGponOltEthernetStatisticReceivedBytes_count[10m]) * 8) by (region)"
	queryBpsOut := "sum(rate(hwGponOltEthernetStatisticSendBytes_count[10m]) * 8) by (region)"
	queryBytesIn := "sum(increase(hwGponOltEthernetStatisticReceivedBytes_count[10m])) by (region)"
	queryBytesOut := "sum(increase(hwGponOltEthernetStatisticSendBytes_count[10m])) by (region)"

	result := make(map[string][]*Traffic)

	for t := initDate; t.Before(finalDate); t = t.Add(5 * time.Minute) {
		mbpsBwVec, _ := p.queryVector(ctx, queryBW, t)
		bpsInVec, _ := p.queryVector(ctx, queryBpsIn, t)
		bpsOutVec, _ := p.queryVector(ctx, queryBpsOut, t)
		bytesInVec, _ := p.queryVector(ctx, queryBytesIn, t)
		bytesOutVec, _ := p.queryVector(ctx, queryBytesOut, t)

		tempData := make(map[string]*Traffic)

		for _, s := range bpsInVec {
			key := s.Labels["region"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsIn = s.Value
		}
		for _, s := range bpsOutVec {
			key := s.Labels["region"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsOut = s.Value
		}
		for _, s := range mbpsBwVec {
			key := s.Labels["region"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].Bandwidth = s.Value
		}
		for _, s := range bytesInVec {
			key := s.Labels["region"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesIn = s.Value
		}
		for _, s := range bytesOutVec {
			key := s.Labels["region"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesOut = s.Value
		}

		for region, traffic := range tempData {
			result[region] = append(result[region], traffic)
		}
	}

	return result, nil
}

func (p *Prometheus) PrometheusTrafficState(ctx context.Context, initDate, finalDate time.Time) (map[string][]*Traffic, error) {
	queryBW := "sum(ifSpeed) by (state)"
	queryBpsIn := "sum(rate(hwGponOltEthernetStatisticReceivedBytes_count[10m]) * 8) by (state)"
	queryBpsOut := "sum(rate(hwGponOltEthernetStatisticSendBytes_count[10m]) * 8) by (state)"
	queryBytesIn := "sum(increase(hwGponOltEthernetStatisticReceivedBytes_count[10m])) by (state)"
	queryBytesOut := "sum(increase(hwGponOltEthernetStatisticSendBytes_count[10m])) by (state)"

	result := make(map[string][]*Traffic)

	for t := initDate; t.Before(finalDate); t = t.Add(5 * time.Minute) {
		mbpsBwVec, _ := p.queryVector(ctx, queryBW, t)
		bpsInVec, _ := p.queryVector(ctx, queryBpsIn, t)
		bpsOutVec, _ := p.queryVector(ctx, queryBpsOut, t)
		bytesInVec, _ := p.queryVector(ctx, queryBytesIn, t)
		bytesOutVec, _ := p.queryVector(ctx, queryBytesOut, t)

		tempData := make(map[string]*Traffic)

		for _, s := range bpsInVec {
			key := s.Labels["state"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsIn = s.Value
		}
		for _, s := range bpsOutVec {
			key := s.Labels["state"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsOut = s.Value
		}
		for _, s := range mbpsBwVec {
			key := s.Labels["state"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].Bandwidth = s.Value
		}
		for _, s := range bytesInVec {
			key := s.Labels["state"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesIn = s.Value
		}
		for _, s := range bytesOutVec {
			key := s.Labels["state"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesOut = s.Value
		}

		for region, traffic := range tempData {
			result[region] = append(result[region], traffic)
		}
	}

	return result, nil
}

func (p *Prometheus) PrometheusTrafficGroupInstance(ctx context.Context, instances []string, initDate, finalDate time.Time) (map[string][]*Traffic, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances provided")
	}

	instancesStr := strings.Join(instances, "|")

	queryBW := fmt.Sprintf("sum(ifSpeed{instance=~'%s'}) by (instance)", instancesStr)
	queryBpsIn := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance=~'%s'}[10m]) * 8) by (instance)", instancesStr)
	queryBpsOut := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[10m]) * 8) by (instance)", instancesStr)
	queryBytesIn := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticReceivedBytes_count{instance=~'%s'}[10m])) by (instance)", instancesStr)
	queryBytesOut := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[10m])) by (instance)", instancesStr)

	result := make(map[string][]*Traffic)

	for t := initDate; t.Before(finalDate); t = t.Add(5 * time.Minute) {
		mbpsBwVec, _ := p.queryVector(ctx, queryBW, t)
		bpsInVec, _ := p.queryVector(ctx, queryBpsIn, t)
		bpsOutVec, _ := p.queryVector(ctx, queryBpsOut, t)
		bytesInVec, _ := p.queryVector(ctx, queryBytesIn, t)
		bytesOutVec, _ := p.queryVector(ctx, queryBytesOut, t)

		tempData := make(map[string]*Traffic)

		for _, s := range bpsInVec {
			key := s.Labels["instance"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsIn = s.Value
		}
		for _, s := range bpsOutVec {
			key := s.Labels["instance"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsOut = s.Value
		}
		for _, s := range mbpsBwVec {
			key := s.Labels["instance"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].Bandwidth = s.Value
		}
		for _, s := range bytesInVec {
			key := s.Labels["instance"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesIn = s.Value
		}
		for _, s := range bytesOutVec {
			key := s.Labels["instance"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesOut = s.Value
		}

		for region, traffic := range tempData {
			result[region] = append(result[region], traffic)
		}
	}

	return result, nil
}

func (p *Prometheus) PrometheusTrafficInstance(ctx context.Context, instance string, initDate, finalDate time.Time) (map[string][]*Traffic, error) {
	queryBW := fmt.Sprintf("sum(ifSpeed{instance='%s'}) by (ifIndex)", instance)
	queryIfName := fmt.Sprintf("ifName{instance='%s'}", instance) // ifName{ifIndex="2097152", ifName="GPON X/Y/Z", instance="10.125.120.231", job="olt_distrito-capital", region="Capital", state="Distrito Capital"}
	queryBpsIn := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance='%s'}[10m]) * 8) by (ponPortIndex)", instance)
	queryBpsOut := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[10m]) * 8) by (ponPortIndex)", instance)
	queryBytesIn := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticReceivedBytes_count{instance=~'%s'}[10m])) by (ponPortIndex)", instance)
	queryBytesOut := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[10m])) by (ponPortIndex)", instance)

	result := make(map[string][]*Traffic)

	for t := initDate; t.Before(finalDate); t = t.Add(5 * time.Minute) {
		mbpsBwVec, _ := p.queryVector(ctx, queryBW, t)
		ifNameVec, _ := p.queryVector(ctx, queryIfName, t)
		bpsInVec, _ := p.queryVector(ctx, queryBpsIn, t)
		bpsOutVec, _ := p.queryVector(ctx, queryBpsOut, t)
		bytesInVec, _ := p.queryVector(ctx, queryBytesIn, t)
		bytesOutVec, _ := p.queryVector(ctx, queryBytesOut, t)

		tempData := make(map[string]*Traffic)

		for _, s := range ifNameVec {
			key := s.Labels["ifIndex"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].Description = s.Labels["ifName"]
		}

		for _, s := range mbpsBwVec {
			key := s.Labels["ifIndex"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].Bandwidth = s.Value
		}

		for _, s := range bpsInVec {
			key := s.Labels["ponPortIndex"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsIn = s.Value
		}
		for _, s := range bpsOutVec {
			key := s.Labels["ponPortIndex"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BpsOut = s.Value
		}

		for _, s := range bytesInVec {
			key := s.Labels["ponPortIndex"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesIn = s.Value
		}
		for _, s := range bytesOutVec {
			key := s.Labels["ponPortIndex"]
			if _, ok := tempData[key]; !ok {
				tempData[key] = &Traffic{Time: t}
			}
			tempData[key].BytesOut = s.Value
		}

		for region, traffic := range tempData {
			result[region] = append(result[region], traffic)
		}
	}

	return result, nil
}

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
