package prometheus

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/metalpoch/ultra-monitor/internal/utils"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Prometheus interface {
	DeviceScan(ctx context.Context) ([]InfoDevice, error)
	InstanceScan(ctx context.Context, ip string) ([]InfoDevice, error)
	DeviceLocations(ctx context.Context) ([]DeviceLocation, error)
	TrafficTotal(ctx context.Context, initDate, finalDate time.Time) ([]*Traffic, error)
	TrafficByRegion(ctx context.Context, region string, initDate, finalDate time.Time) ([]*Traffic, error)
	TrafficByState(ctx context.Context, state string, initDate, finalDate time.Time) ([]*Traffic, error)
	TrafficGroupInstance(ctx context.Context, instances []string, initDate, finalDate time.Time) ([]*Traffic, error)
	TrafficInstanceByIndex(ctx context.Context, instance, index string, initDate, finalDate time.Time) ([]*Traffic, error)
}

type prometheus struct {
	client v1.API
}

func NewPrometheusClient(host string) *prometheus {
	client, err := api.NewClient(api.Config{Address: host})
	if err != nil {
		log.Fatal(err)
	}
	return &prometheus{client: v1.NewAPI(client)}
}

func (p *prometheus) DeviceScan(ctx context.Context) ([]InfoDevice, error) {
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

func (p *prometheus) InstanceScan(ctx context.Context, ip string) ([]InfoDevice, error) {
	ifNameVec, err := p.queryVector(ctx, fmt.Sprintf("ifName{instance='%s'}", ip), time.Now())
	if err != nil {
		return nil, err
	}

	var devices []InfoDevice
	for _, s := range ifNameVec {
		devices = append(devices, InfoDevice{
			Region:  s.Labels["region"],
			State:   s.Labels["state"],
			IP:      s.Labels["instance"],
			IfName:  s.Labels["ifName"],
			IfIndex: utils.ParseInt64(s.Labels["ifIndex"]),
		})
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices found in Prometheus")
	}

	return devices, nil
}

func (p *prometheus) DeviceLocations(ctx context.Context) ([]DeviceLocation, error) {
	locationVec, err := p.queryVector(ctx, "count(sysName) by (region, state, instance, sysName)", time.Now())
	if err != nil {
		return nil, err
	}

	var devices []DeviceLocation
	for _, s := range locationVec {
		devices = append(devices, DeviceLocation{
			Region:  s.Labels["region"],
			State:   s.Labels["state"],
			IP:      s.Labels["instance"],
			SysName: s.Labels["sysName"],
		})
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices found in Prometheus")
	}

	return devices, nil
}

func (p *prometheus) TrafficByRegion(ctx context.Context, region string, initDate, finalDate time.Time) ([]*Traffic, error) {
	queryBW := fmt.Sprintf("sum(ifSpeed{region='%s'})", region)
	queryBpsIn := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{region='%s'}[10m]) * 8)", region)
	queryBpsOut := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticSendBytes_count{region='%s'}[10m]) * 8)", region)
	queryBytesIn := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticReceivedBytes_count{region='%s'}[10m]))", region)
	queryBytesOut := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticSendBytes_count{region='%s'}[10m]))", region)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  5 * time.Minute,
	}

	mbpsBwResult, _, _ := p.client.QueryRange(ctx, queryBW, r)
	bpsInResult, _, _ := p.client.QueryRange(ctx, queryBpsIn, r)
	bpsOutResult, _, _ := p.client.QueryRange(ctx, queryBpsOut, r)
	bytesInResult, _, _ := p.client.QueryRange(ctx, queryBytesIn, r)
	bytesOutResult, _, _ := p.client.QueryRange(ctx, queryBytesOut, r)

	mbpsBwMatrix, _ := mbpsBwResult.(model.Matrix)
	bpsInMatrix, _ := bpsInResult.(model.Matrix)
	bpsOutMatrix, _ := bpsOutResult.(model.Matrix)
	bytesInMatrix, _ := bytesInResult.(model.Matrix)
	bytesOutMatrix, _ := bytesOutResult.(model.Matrix)

	trafficMap := make(map[int64]*Traffic)

	processMatrix := func(matrix model.Matrix, updateFunc func(*Traffic, float64)) {
		for _, serie := range matrix {
			for _, point := range serie.Values {
				key := int64(point.Timestamp) / 1000
				if _, ok := trafficMap[key]; !ok {
					trafficMap[key] = &Traffic{Time: time.Unix(key, 0)}
				}
				updateFunc(trafficMap[key], float64(point.Value))
			}
		}
	}

	processMatrix(mbpsBwMatrix, func(t *Traffic, val float64) { t.Bandwidth = val })
	processMatrix(bpsInMatrix, func(t *Traffic, val float64) { t.BpsIn = val })
	processMatrix(bpsOutMatrix, func(t *Traffic, val float64) { t.BpsOut = val })
	processMatrix(bytesInMatrix, func(t *Traffic, val float64) { t.BytesIn = val })
	processMatrix(bytesOutMatrix, func(t *Traffic, val float64) { t.BytesOut = val })

	result := make([]*Traffic, 0, len(trafficMap))
	for _, traffic := range trafficMap {
		result = append(result, traffic)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (p *prometheus) TrafficTotal(ctx context.Context, initDate, finalDate time.Time) ([]*Traffic, error) {
	queryBW := "sum(ifSpeed)"
	queryBpsIn := "sum(rate(hwGponOltEthernetStatisticReceivedBytes_count[10m]) * 8)"
	queryBpsOut := "sum(rate(hwGponOltEthernetStatisticSendBytes_count[10m]) * 8)"
	queryBytesIn := "sum(increase(hwGponOltEthernetStatisticReceivedBytes_count[10m]))"
	queryBytesOut := "sum(increase(hwGponOltEthernetStatisticSendBytes_count[10m]))"

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  5 * time.Minute,
	}

	mbpsBwResult, _, _ := p.client.QueryRange(ctx, queryBW, r)
	bpsInResult, _, _ := p.client.QueryRange(ctx, queryBpsIn, r)
	bpsOutResult, _, _ := p.client.QueryRange(ctx, queryBpsOut, r)
	bytesInResult, _, _ := p.client.QueryRange(ctx, queryBytesIn, r)
	bytesOutResult, _, _ := p.client.QueryRange(ctx, queryBytesOut, r)

	mbpsBwMatrix, _ := mbpsBwResult.(model.Matrix)
	bpsInMatrix, _ := bpsInResult.(model.Matrix)
	bpsOutMatrix, _ := bpsOutResult.(model.Matrix)
	bytesInMatrix, _ := bytesInResult.(model.Matrix)
	bytesOutMatrix, _ := bytesOutResult.(model.Matrix)

	trafficMap := make(map[int64]*Traffic)

	processMatrix := func(matrix model.Matrix, updateFunc func(*Traffic, float64)) {
		for _, serie := range matrix {
			for _, point := range serie.Values {
				key := int64(point.Timestamp) / 1000
				if _, ok := trafficMap[key]; !ok {
					trafficMap[key] = &Traffic{Time: time.Unix(key, 0)}
				}
				updateFunc(trafficMap[key], float64(point.Value))
			}
		}
	}

	processMatrix(mbpsBwMatrix, func(t *Traffic, val float64) { t.Bandwidth = val })
	processMatrix(bpsInMatrix, func(t *Traffic, val float64) { t.BpsIn = val })
	processMatrix(bpsOutMatrix, func(t *Traffic, val float64) { t.BpsOut = val })
	processMatrix(bytesInMatrix, func(t *Traffic, val float64) { t.BytesIn = val })
	processMatrix(bytesOutMatrix, func(t *Traffic, val float64) { t.BytesOut = val })

	result := make([]*Traffic, 0, len(trafficMap))
	for _, traffic := range trafficMap {
		result = append(result, traffic)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (p *prometheus) TrafficByState(ctx context.Context, state string, initDate, finalDate time.Time) ([]*Traffic, error) {
	queryBW := fmt.Sprintf("sum(ifSpeed{state='%s'})", state)
	queryBpsIn := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{state='%s'}[10m]) * 8)", state)
	queryBpsOut := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticSendBytes_count{state='%s'}[10m]) * 8)", state)
	queryBytesIn := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticReceivedBytes_count{state='%s'}[10m]))", state)
	queryBytesOut := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticSendBytes_count{state='%s'}[10m]))", state)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  5 * time.Minute,
	}

	mbpsBwResult, _, _ := p.client.QueryRange(ctx, queryBW, r)
	bpsInResult, _, _ := p.client.QueryRange(ctx, queryBpsIn, r)
	bpsOutResult, _, _ := p.client.QueryRange(ctx, queryBpsOut, r)
	bytesInResult, _, _ := p.client.QueryRange(ctx, queryBytesIn, r)
	bytesOutResult, _, _ := p.client.QueryRange(ctx, queryBytesOut, r)

	mbpsBwMatrix, _ := mbpsBwResult.(model.Matrix)
	bpsInMatrix, _ := bpsInResult.(model.Matrix)
	bpsOutMatrix, _ := bpsOutResult.(model.Matrix)
	bytesInMatrix, _ := bytesInResult.(model.Matrix)
	bytesOutMatrix, _ := bytesOutResult.(model.Matrix)

	trafficMap := make(map[int64]*Traffic)

	processMatrix := func(matrix model.Matrix, updateFunc func(*Traffic, float64)) {
		for _, serie := range matrix {
			for _, point := range serie.Values {
				key := int64(point.Timestamp) / 1000
				if _, ok := trafficMap[key]; !ok {
					trafficMap[key] = &Traffic{Time: time.Unix(key, 0)}
				}
				updateFunc(trafficMap[key], float64(point.Value))
			}
		}
	}

	processMatrix(mbpsBwMatrix, func(t *Traffic, val float64) { t.Bandwidth = val })
	processMatrix(bpsInMatrix, func(t *Traffic, val float64) { t.BpsIn = val })
	processMatrix(bpsOutMatrix, func(t *Traffic, val float64) { t.BpsOut = val })
	processMatrix(bytesInMatrix, func(t *Traffic, val float64) { t.BytesIn = val })
	processMatrix(bytesOutMatrix, func(t *Traffic, val float64) { t.BytesOut = val })

	result := make([]*Traffic, 0, len(trafficMap))
	for _, traffic := range trafficMap {
		result = append(result, traffic)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (p *prometheus) TrafficGroupInstance(ctx context.Context, instances []string, initDate, finalDate time.Time) ([]*Traffic, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances provided")
	}

	instancesStr := strings.Join(instances, "|")
	queryBW := fmt.Sprintf("sum(ifSpeed{instance=~'%s'})", instancesStr)
	queryBpsIn := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance=~'%s'}[10m]) * 8)", instancesStr)
	queryBpsOut := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[10m]) * 8)", instancesStr)
	queryBytesIn := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticReceivedBytes_count{instance=~'%s'}[10m]))", instancesStr)
	queryBytesOut := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[10m]))", instancesStr)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  5 * time.Minute,
	}

	mbpsBwResult, _, _ := p.client.QueryRange(ctx, queryBW, r)
	bpsInResult, _, _ := p.client.QueryRange(ctx, queryBpsIn, r)
	bpsOutResult, _, _ := p.client.QueryRange(ctx, queryBpsOut, r)
	bytesInResult, _, _ := p.client.QueryRange(ctx, queryBytesIn, r)
	bytesOutResult, _, _ := p.client.QueryRange(ctx, queryBytesOut, r)

	mbpsBwMatrix, _ := mbpsBwResult.(model.Matrix)
	bpsInMatrix, _ := bpsInResult.(model.Matrix)
	bpsOutMatrix, _ := bpsOutResult.(model.Matrix)
	bytesInMatrix, _ := bytesInResult.(model.Matrix)
	bytesOutMatrix, _ := bytesOutResult.(model.Matrix)

	trafficMap := make(map[int64]*Traffic)

	processMatrix := func(matrix model.Matrix, updateFunc func(*Traffic, float64)) {
		for _, serie := range matrix {
			for _, point := range serie.Values {
				key := int64(point.Timestamp) / 1000
				if _, ok := trafficMap[key]; !ok {
					trafficMap[key] = &Traffic{Time: time.Unix(key, 0)}
				}
				updateFunc(trafficMap[key], float64(point.Value))
			}
		}
	}

	processMatrix(mbpsBwMatrix, func(t *Traffic, val float64) { t.Bandwidth = val })
	processMatrix(bpsInMatrix, func(t *Traffic, val float64) { t.BpsIn = val })
	processMatrix(bpsOutMatrix, func(t *Traffic, val float64) { t.BpsOut = val })
	processMatrix(bytesInMatrix, func(t *Traffic, val float64) { t.BytesIn = val })
	processMatrix(bytesOutMatrix, func(t *Traffic, val float64) { t.BytesOut = val })

	result := make([]*Traffic, 0, len(trafficMap))
	for _, traffic := range trafficMap {
		result = append(result, traffic)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (p *prometheus) TrafficInstanceByIndex(ctx context.Context, instance, index string, initDate, finalDate time.Time) ([]*Traffic, error) {
	queryBW := fmt.Sprintf("sum(ifSpeed{instance='%s', ifIndex='%s'})", instance, index)
	queryBpsIn := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance='%s', ponPortIndex='%s'}[10m]) * 8)", instance, index)
	queryBpsOut := fmt.Sprintf("sum(rate(hwGponOltEthernetStatisticSendBytes_count{instance='%s', ponPortIndex='%s'}[10m]) * 8)", instance, index)
	queryBytesIn := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticReceivedBytes_count{instance='%s', ponPortIndex='%s'}[10m]))", instance, index)
	queryBytesOut := fmt.Sprintf("sum(increase(hwGponOltEthernetStatisticSendBytes_count{instance='%s', ponPortIndex='%s'}[10m]))", instance, index)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  5 * time.Minute,
	}

	mbpsBwResult, _, _ := p.client.QueryRange(ctx, queryBW, r)
	bpsInResult, _, _ := p.client.QueryRange(ctx, queryBpsIn, r)
	bpsOutResult, _, _ := p.client.QueryRange(ctx, queryBpsOut, r)
	bytesInResult, _, _ := p.client.QueryRange(ctx, queryBytesIn, r)
	bytesOutResult, _, _ := p.client.QueryRange(ctx, queryBytesOut, r)

	mbpsBwMatrix, _ := mbpsBwResult.(model.Matrix)
	bpsInMatrix, _ := bpsInResult.(model.Matrix)
	bpsOutMatrix, _ := bpsOutResult.(model.Matrix)
	bytesInMatrix, _ := bytesInResult.(model.Matrix)
	bytesOutMatrix, _ := bytesOutResult.(model.Matrix)

	trafficMap := make(map[int64]*Traffic)

	processMatrix := func(matrix model.Matrix, updateFunc func(*Traffic, float64)) {
		for _, serie := range matrix {
			for _, point := range serie.Values {
				key := int64(point.Timestamp) / 1000
				if _, ok := trafficMap[key]; !ok {
					trafficMap[key] = &Traffic{Time: time.Unix(key, 0)}
				}
				updateFunc(trafficMap[key], float64(point.Value))
			}
		}
	}

	processMatrix(mbpsBwMatrix, func(t *Traffic, val float64) { t.Bandwidth = val })
	processMatrix(bpsInMatrix, func(t *Traffic, val float64) { t.BpsIn = val })
	processMatrix(bpsOutMatrix, func(t *Traffic, val float64) { t.BpsOut = val })
	processMatrix(bytesInMatrix, func(t *Traffic, val float64) { t.BytesIn = val })
	processMatrix(bytesOutMatrix, func(t *Traffic, val float64) { t.BytesOut = val })

	result := make([]*Traffic, 0, len(trafficMap))
	for _, traffic := range trafficMap {
		result = append(result, traffic)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (p *prometheus) queryVector(ctx context.Context, query string, ts time.Time) ([]dataProm, error) {
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
