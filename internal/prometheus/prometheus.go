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

	//Dashboard
	TrafficTotalByField(ctx context.Context, fieldQuery, value string, initDate, finalDate time.Time) ([]*Traffic, error)
	TrafficGroupedByField(ctx context.Context, fieldQuery, value, groupBy string, initDate, finalDate time.Time) (map[string][]*Traffic, error)
	SysnameByState(ctx context.Context, state string, initDate, finalDate time.Time) (map[string][]*Traffic, error)

	// Details
	TrafficGroupInstance(ctx context.Context, instances []string, initDate, finalDate time.Time) ([]*Traffic, error)
	TrafficInstanceByIndex(ctx context.Context, instance, index string, initDate, finalDate time.Time) ([]*Traffic, error)

	// stats
	StatesStatsByRegion(ctx context.Context, region string, initDate, finalDate time.Time) ([]RegionStats, error)
	OltStatsByState(ctx context.Context, state string, initDate, finalDate time.Time) ([]OltStats, error)
	GponStatsByInstance(ctx context.Context, instance string, initDate, finalDate time.Time) ([]GponStats, error)

	// Traffic by instance, state, region
	TrafficByInstanceStateRegion(ctx context.Context, initDate, finalDate time.Time) ([]TrafficByInstance, error)
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
	ifNameVec, err := p.queryVector(ctx, "ifOperStatus * on(ifIndex, instance) group_left(ifName) ifName", time.Now())
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
			Region:       oltRegion,
			State:        oltState,
			IP:           oltIP,
			IfName:       s.Labels["ifName"],
			IfIndex:      utils.ParseInt64(ifIndex),
			IfOperStatus: int8(s.Value),
		})
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices found in Prometheus")
	}

	return devices, nil
}

func (p *prometheus) InstanceScan(ctx context.Context, ip string) ([]InfoDevice, error) {
	queryVec, err := p.queryVector(ctx,
		fmt.Sprintf("ifOperStatus * on(ifIndex, instance) group_left(ifName) ifName{instance='%s'}", ip),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	var devices []InfoDevice
	for _, s := range queryVec {
		devices = append(devices, InfoDevice{
			Region:       s.Labels["region"],
			State:        s.Labels["state"],
			IP:           s.Labels["instance"],
			IfName:       s.Labels["ifName"],
			IfIndex:      utils.ParseInt64(s.Labels["ifIndex"]),
			IfOperStatus: int8(s.Value),
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

func (p *prometheus) TrafficTotalByField(ctx context.Context, fieldQuery, value string, initDate, finalDate time.Time) ([]*Traffic, error) {
	var query string
	if fieldQuery != "" {
		query = fmt.Sprintf("%s='%s'", fieldQuery, value)
	}

	queryBW := fmt.Sprintf("sum(ifSpeed{%s})", query)
	queryBpsIn := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticReceivedBytes_count{%s}[1h])[3h:1h]) * 8)", query)
	queryBpsOut := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticSendBytes_count{%s}[1h])[3h:1h]) * 8)", query)
	queryBytesIn := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticReceivedBytes_count{%s}[1h])[3h:1h]))", query)
	queryBytesOut := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticSendBytes_count{%s}[1h])[3h:1h]))", query)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  time.Hour,
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

func (p *prometheus) TrafficGroupedByField(ctx context.Context, fieldQuery, value, groupBy string, initDate, finalDate time.Time) (map[string][]*Traffic, error) {
	var query string
	if fieldQuery != "" {
		query = fmt.Sprintf("%s='%s'", fieldQuery, value)
	}

	queryBW := fmt.Sprintf("sum(ifSpeed{%s}) by (%s)", query, groupBy)
	queryBpsIn := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticReceivedBytes_count{%s}[1h])[3h:1h]) * 8) by (%s)", query, groupBy)
	queryBpsOut := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticSendBytes_count{%s}[1h])[3h:1h]) * 8) by (%s)", query, groupBy)
	queryBytesIn := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticReceivedBytes_count{%s}[1h])[3h:1h])) by (%s)", query, groupBy)
	queryBytesOut := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticSendBytes_count{%s}[1h])[3h:1h])) by (%s)", query, groupBy)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  time.Hour,
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

	trafficByStateAndTime := make(map[string]map[int64]*Traffic)

	processMatrix := func(matrix model.Matrix, field string) {
		for _, serie := range matrix {
			fieldName := string(serie.Metric[model.LabelName(groupBy)])
			if _, ok := trafficByStateAndTime[fieldName]; !ok {
				trafficByStateAndTime[fieldName] = make(map[int64]*Traffic)
			}
			for _, point := range serie.Values {
				key := int64(point.Timestamp) / 1000
				if _, ok := trafficByStateAndTime[fieldName][key]; !ok {
					trafficByStateAndTime[fieldName][key] = &Traffic{
						Time: time.Unix(key, 0),
					}
				}
				traffic := trafficByStateAndTime[fieldName][key]
				switch field {
				case "Bandwidth":
					traffic.Bandwidth = float64(point.Value)
				case "BpsIn":
					traffic.BpsIn = float64(point.Value)
				case "BpsOut":
					traffic.BpsOut = float64(point.Value)
				case "BytesIn":
					traffic.BytesIn = float64(point.Value)
				case "BytesOut":
					traffic.BytesOut = float64(point.Value)
				}
			}
		}
	}

	processMatrix(mbpsBwMatrix, "Bandwidth")
	processMatrix(bpsInMatrix, "BpsIn")
	processMatrix(bpsOutMatrix, "BpsOut")
	processMatrix(bytesInMatrix, "BytesIn")
	processMatrix(bytesOutMatrix, "BytesOut")

	result := make(map[string][]*Traffic)
	for fieldName, trafficMap := range trafficByStateAndTime {
		slice := make([]*Traffic, 0, len(trafficMap))
		for _, traffic := range trafficMap {
			slice = append(slice, traffic)
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Time.Before(slice[j].Time)
		})
		result[fieldName] = slice
	}

	return result, nil
}

func (p *prometheus) SysnameByState(ctx context.Context, state string, initDate, finalDate time.Time) (map[string][]*Traffic, error) {
	queryBW := fmt.Sprintf("sum(ifSpeed{state='%s'}) by (instance) * on(instance) group_left(sysName) sysName", state)
	queryBpsIn := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticReceivedBytes_count{state='%s'}[1h])[3h:1h]) * 8) by (instance) * on(instance) group_left(sysName) sysName", state)
	queryBpsOut := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticSendBytes_count{state='%s'}[1h])[3h:1h]) * 8) by (instance) * on(instance) group_left(sysName) sysName", state)
	queryBytesIn := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticReceivedBytes_count{state='%s'}[1h])[3h:1h])) by (instance) * on(instance) group_left(sysName) sysName", state)
	queryBytesOut := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticSendBytes_count{state='%s'}[1h])[3h:1h])) by (instance) * on(instance) group_left(sysName) sysName", state)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  time.Hour,
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

	trafficByStateAndTime := make(map[string]map[int64]*Traffic)

	processMatrix := func(matrix model.Matrix, field string) {
		for _, serie := range matrix {
			fieldName := string(serie.Metric["sysName"])
			if _, ok := trafficByStateAndTime[fieldName]; !ok {
				trafficByStateAndTime[fieldName] = make(map[int64]*Traffic)
			}
			for _, point := range serie.Values {
				key := int64(point.Timestamp) / 1000
				if _, ok := trafficByStateAndTime[fieldName][key]; !ok {
					trafficByStateAndTime[fieldName][key] = &Traffic{
						Time: time.Unix(key, 0),
					}
				}
				traffic := trafficByStateAndTime[fieldName][key]
				switch field {
				case "Bandwidth":
					traffic.Bandwidth = float64(point.Value)
				case "BpsIn":
					traffic.BpsIn = float64(point.Value)
				case "BpsOut":
					traffic.BpsOut = float64(point.Value)
				case "BytesIn":
					traffic.BytesIn = float64(point.Value)
				case "BytesOut":
					traffic.BytesOut = float64(point.Value)
				}
			}
		}
	}

	processMatrix(mbpsBwMatrix, "Bandwidth")
	processMatrix(bpsInMatrix, "BpsIn")
	processMatrix(bpsOutMatrix, "BpsOut")
	processMatrix(bytesInMatrix, "BytesIn")
	processMatrix(bytesOutMatrix, "BytesOut")

	result := make(map[string][]*Traffic)
	for fieldName, trafficMap := range trafficByStateAndTime {
		slice := make([]*Traffic, 0, len(trafficMap))
		for _, traffic := range trafficMap {
			slice = append(slice, traffic)
		}
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Time.Before(slice[j].Time)
		})
		result[fieldName] = slice
	}

	return result, nil
}

func (p *prometheus) TrafficGroupInstance(ctx context.Context, instances []string, initDate, finalDate time.Time) ([]*Traffic, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("no instances provided")
	}

	instancesStr := strings.Join(instances, "|")
	queryBW := fmt.Sprintf("sum(ifSpeed{instance=~'%s'})", instancesStr)
	queryBpsIn := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance=~'%s'}[15m])[30m:15m]) * 8)", instancesStr)
	queryBpsOut := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[15m])[30m:15m]) * 8)", instancesStr)
	queryBytesIn := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticReceivedBytes_count{instance=~'%s'}[15m])[30m:15m]))", instancesStr)
	queryBytesOut := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticSendBytes_count{instance=~'%s'}[15m])[30m:15m]))", instancesStr)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  15 * time.Minute,
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

func (p *prometheus) TrafficInstanceByIndex(ctx context.Context, instance, indexes string, initDate, finalDate time.Time) ([]*Traffic, error) {
	queryBW := fmt.Sprintf("sum(ifSpeed{instance='%s', ifIndex=~'%s'})", instance, indexes)
	queryBpsIn := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance='%s', ponPortIndex=~'%s'}[15m])[30m:15m]) * 8)", instance, indexes)
	queryBpsOut := fmt.Sprintf("sum(avg_over_time(rate(hwGponOltEthernetStatisticSendBytes_count{instance='%s', ponPortIndex=~'%s'}[15m])[30m:15m]) * 8)", instance, indexes)
	queryBytesIn := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticReceivedBytes_count{instance='%s', ponPortIndex=~'%s'}[15m])[30m:15m]))", instance, indexes)
	queryBytesOut := fmt.Sprintf("sum(avg_over_time(increase(hwGponOltEthernetStatisticSendBytes_count{instance='%s', ponPortIndex=~'%s'}[15m])[30m:15m]))", instance, indexes)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  15 * time.Minute,
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

func (p *prometheus) StatesStatsByRegion(ctx context.Context, region string, initDate, finalDate time.Time) ([]RegionStats, error) {
	queryIn := fmt.Sprintf(QUERY_STATES_BY_REGION_IN, region)
	queryOut := fmt.Sprintf(QUERY_STATES_BY_REGION_OUT, region)
	querySpeed := fmt.Sprintf(QUERY_STATES_BY_REGION_IFSPEED, region)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  15 * time.Minute,
	}

	inResult, inWarn, inErr := p.client.QueryRange(ctx, queryIn, r)
	outResult, outWarn, outErr := p.client.QueryRange(ctx, queryOut, r)
	speedResult, speedWarn, speedErr := p.client.QueryRange(ctx, querySpeed, r)

	if inErr != nil || outErr != nil || speedErr != nil {
		return nil, fmt.Errorf("error in queries: in=%v, out=%v, speed=%v", inErr, outErr, speedErr)
	}

	if len(inWarn) > 0 {
		log.Printf("Prometheus IN warnings for region %s: %v", region, inWarn)
	}
	if len(outWarn) > 0 {
		log.Printf("Prometheus OUT warnings for region %s: %v", region, outWarn)
	}
	if len(speedWarn) > 0 {
		log.Printf("Prometheus SPEED warnings for region %s: %v", region, speedWarn)
	}

	inMatrix, okIn := inResult.(model.Matrix)
	outMatrix, okOut := outResult.(model.Matrix)
	speedMatrix, okSpeed := speedResult.(model.Matrix)

	if !okIn || !okOut || !okSpeed {
		return nil, fmt.Errorf("unexpected result types: in=%T, out=%T, speed=%T", inResult, outResult, speedResult)
	}

	type trafficStats struct {
		IfSpeed   float64
		MaxInBps  float64
		AvgInBps  float64
		MaxOutBps float64
		AvgOutBps float64
		UsageIn   float64
		UsageOut  float64
	}

	statsMap := make(map[string]*trafficStats)

	processStats := func(matrix model.Matrix, isIn bool) {
		for _, serie := range matrix {
			state := string(serie.Metric["state"])
			if state == "" {
				continue
			}
			stat, exists := statsMap[state]
			if !exists {
				stat = &trafficStats{}
				statsMap[state] = stat
			}
			var sum float64
			for _, v := range serie.Values {
				val := float64(v.Value)
				sum += val
				if isIn && val > stat.MaxInBps {
					stat.MaxInBps = val
				}
				if !isIn && val > stat.MaxOutBps {
					stat.MaxOutBps = val
				}
			}
			avg := sum / float64(len(serie.Values))
			if isIn {
				stat.AvgInBps = avg
			} else {
				stat.AvgOutBps = avg
			}
		}
	}

	processStats(inMatrix, true)
	processStats(outMatrix, false)

	for _, serie := range speedMatrix {
		state := string(serie.Metric["state"])
		stat, exists := statsMap[state]
		if !exists {
			continue
		}
		if len(serie.Values) > 0 {
			stat.IfSpeed = float64(serie.Values[0].Value)
		}
	}

	for _, stat := range statsMap {
		if stat.IfSpeed > 0 {
			stat.UsageIn = (stat.MaxInBps / stat.IfSpeed) * 100
			stat.UsageOut = (stat.MaxOutBps / stat.IfSpeed) * 100
		}
	}

	stats := make([]RegionStats, 0, len(statsMap))
	for state, s := range statsMap {
		stats = append(stats, RegionStats{
			State:     state,
			IfSpeed:   s.IfSpeed,
			MaxInBps:  s.MaxInBps,
			AvgInBps:  s.AvgInBps,
			MaxOutBps: s.MaxOutBps,
			AvgOutBps: s.AvgOutBps,
			UsageIn:   s.UsageIn,
			UsageOut:  s.UsageOut,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].State < stats[j].State
	})

	return stats, nil
}

func (p *prometheus) OltStatsByState(ctx context.Context, state string, initDate, finalDate time.Time) ([]OltStats, error) {
	queryIn := fmt.Sprintf(QUERY_OLT_BY_STATES_IN, state)
	queryOut := fmt.Sprintf(QUERY_OLT_BY_STATES_OUT, state)
	querySpeed := fmt.Sprintf(QUERY_OLT_BY_STATES_IFSPEED, state)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  15 * time.Minute,
	}

	inResult, inWarn, inErr := p.client.QueryRange(ctx, queryIn, r)
	outResult, outWarn, outErr := p.client.QueryRange(ctx, queryOut, r)
	speedResult, speedWarn, speedErr := p.client.QueryRange(ctx, querySpeed, r)

	if inErr != nil || outErr != nil || speedErr != nil {
		return nil, fmt.Errorf("error in queries: in=%v, out=%v, speed=%v", inErr, outErr, speedErr)
	}

	logWarnings := func(label string, warns v1.Warnings) {
		if len(warns) > 0 {
			log.Printf("Prometheus %s warnings for state %s: %v", label, state, warns)
		}
	}
	logWarnings("IN", inWarn)
	logWarnings("OUT", outWarn)
	logWarnings("SPEED", speedWarn)

	inMatrix, okIn := inResult.(model.Matrix)
	outMatrix, okOut := outResult.(model.Matrix)
	speedMatrix, okSpeed := speedResult.(model.Matrix)

	if !okIn || !okOut || !okSpeed {
		return nil, fmt.Errorf("unexpected result types: in=%T, out=%T, speed=%T", inResult, outResult, speedResult)
	}

	type trafficStats struct {
		SysName   string
		IfSpeed   float64
		MaxInBps  float64
		AvgInBps  float64
		MaxOutBps float64
		AvgOutBps float64
		UsageIn   float64
		UsageOut  float64
	}

	statsMap := make(map[string]*trafficStats)

	processStats := func(matrix model.Matrix, isIn bool) {
		for _, serie := range matrix {
			instance := string(serie.Metric["instance"])
			sysName := string(serie.Metric["sysName"])
			if instance == "" {
				continue
			}
			stat, exists := statsMap[instance]
			if !exists {
				stat = &trafficStats{SysName: sysName}
				statsMap[instance] = stat
			}
			var sum float64
			for _, v := range serie.Values {
				val := float64(v.Value)
				sum += val
				if isIn && val > stat.MaxInBps {
					stat.MaxInBps = val
				}
				if !isIn && val > stat.MaxOutBps {
					stat.MaxOutBps = val
				}
			}
			avg := sum / float64(len(serie.Values))
			if isIn {
				stat.AvgInBps = avg
			} else {
				stat.AvgOutBps = avg
			}
		}
	}

	processStats(inMatrix, true)
	processStats(outMatrix, false)

	for _, serie := range speedMatrix {
		instance := string(serie.Metric["instance"])
		stat, exists := statsMap[instance]
		if !exists {
			continue
		}
		if len(serie.Values) > 0 {
			stat.IfSpeed = float64(serie.Values[0].Value)
		}
	}

	for _, stat := range statsMap {
		if stat.IfSpeed > 0 {
			stat.UsageIn = (stat.MaxInBps / stat.IfSpeed) * 100
			stat.UsageOut = (stat.MaxOutBps / stat.IfSpeed) * 100
		}
	}

	stats := make([]OltStats, 0, len(statsMap))
	for instance, s := range statsMap {
		stats = append(stats, OltStats{
			Instance:  instance,
			SysName:   s.SysName,
			IfSpeed:   s.IfSpeed,
			MaxInBps:  s.MaxInBps,
			AvgInBps:  s.AvgInBps,
			MaxOutBps: s.MaxOutBps,
			AvgOutBps: s.AvgOutBps,
			UsageIn:   s.UsageIn,
			UsageOut:  s.UsageOut,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Instance < stats[j].Instance
	})

	return stats, nil
}

func (p *prometheus) GponStatsByInstance(ctx context.Context, instance string, initDate, finalDate time.Time) ([]GponStats, error) {
	queryIn := fmt.Sprintf(QUERY_GPON_STATS_IN, instance, instance)
	queryOut := fmt.Sprintf(QUERY_GPON_STATS_OUT, instance, instance)
	querySpeed := fmt.Sprintf(QUERY_GPON_STATS_IFSPEED, instance)

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  15 * time.Minute,
	}
	inResult, inWarn, inErr := p.client.QueryRange(ctx, queryIn, r)
	outResult, outWarn, outErr := p.client.QueryRange(ctx, queryOut, r)
	speedResult, speedWarn, speedErr := p.client.QueryRange(ctx, querySpeed, r)

	if inErr != nil || outErr != nil || speedErr != nil {
		return nil, fmt.Errorf("error in queries: in=%v, out=%v, speed=%v", inErr, outErr, speedErr)
	}

	if len(inWarn) > 0 {
		log.Printf("Prometheus IN warnings for %s: %v", instance, inWarn)
	}
	if len(outWarn) > 0 {
		log.Printf("Prometheus OUT warnings for %s: %v", instance, outWarn)
	}
	if len(speedWarn) > 0 {
		log.Printf("Prometheus SPEED warnings for %s: %v", instance, speedWarn)
	}

	inMatrix, okIn := inResult.(model.Matrix)
	outMatrix, okOut := outResult.(model.Matrix)
	speedMatrix, okSpeed := speedResult.(model.Matrix)

	if !okOut || !okIn || !okSpeed {
		return nil, fmt.Errorf("unexpected result for %s type: in=%T, out=%T, speed=%T", instance, inResult, outResult, speedResult)
	}

	statsMap := make(map[string]*trafficStats)
	p.processStats(inMatrix, true, instance, statsMap) // isIn = True
	p.processStats(outMatrix, false, instance, statsMap)

	for _, serie := range speedMatrix {
		port := string(serie.Metric["ifIndex"])
		stat, exist := statsMap[port]
		if !exist {
			continue
		}
		if len(serie.Values) > 0 {
			stat.IfSpeed = float64(serie.Values[0].Value)
		}
	}

	// Calcular porcentaje de uso
	for _, stat := range statsMap {
		if stat.IfSpeed > 0 {
			stat.UsageIn = (stat.MaxInBps / stat.IfSpeed) * 100
			stat.UsageOut = (stat.MaxOutBps / stat.IfSpeed) * 100
		}
	}

	stats := make([]GponStats, 0, len(statsMap))
	for _, s := range statsMap {
		stats = append(stats, GponStats{
			Port:      s.Port,
			IfName:    s.IfName,
			IfSpeed:   s.IfSpeed,
			AvgInBps:  s.AvgInBps,
			MaxInBps:  s.MaxInBps,
			UsageIn:   s.UsageIn,
			UsageOut:  s.UsageOut,
			AvgOutBps: s.AvgOutBps,
			MaxOutBps: s.MaxOutBps,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Port < stats[j].Port
	})

	return stats, nil
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

func (p *prometheus) processStats(matrix model.Matrix, isIn bool, instance string, statsMap map[string]*trafficStats) {
	for _, serie := range matrix {
		port := string(serie.Metric["ponPortIndex"])
		ifName := string(serie.Metric["ifName"])
		stat, exists := statsMap[port]
		if !exists {
			stat = &trafficStats{
				Port:     port,
				IfName:   ifName,
				Instance: instance,
			}
			statsMap[port] = stat
		}
		var max, sum float64
		for _, point := range serie.Values {
			val := float64(point.Value)
			sum += val
			if val > max {
				max = val
			}
		}
		if isIn {
			stat.MaxInBps = max
			stat.AvgInBps = sum / float64(len(serie.Values))
		} else {
			stat.MaxOutBps = max
			stat.AvgOutBps = sum / float64(len(serie.Values))
		}
		stat.Samples = len(serie.Values)
	}
}

func (p *prometheus) TrafficByInstanceStateRegion(ctx context.Context, initDate, finalDate time.Time) ([]TrafficByInstance, error) {
	// Query for BpsOut (send bytes)
	queryBpsOut := "sum(avg_over_time(rate(hwGponOltEthernetStatisticSendBytes_count{}[1h])[3h:1h]) * 8) by (instance, state, region)"
	// Query for BpsIn (received bytes)
	queryBpsIn := "sum(avg_over_time(rate(hwGponOltEthernetStatisticReceivedBytes_count{}[1h])[3h:1h]) * 8) by (instance, state, region)"
	// Query for Bandwidth (ifSpeed)
	queryBandwidth := "sum(ifSpeed{}) by (instance, state, region)"
	// Query for BytesIn (received bytes total)
	queryBytesIn := "sum(avg_over_time(increase(hwGponOltEthernetStatisticReceivedBytes_count{}[1h])[3h:1h])) by (instance, state, region)"
	// Query for BytesOut (send bytes total)
	queryBytesOut := "sum(avg_over_time(increase(hwGponOltEthernetStatisticSendBytes_count{}[1h])[3h:1h])) by (instance, state, region)"

	r := v1.Range{
		Start: initDate,
		End:   finalDate,
		Step:  time.Hour,
	}

	// Execute all queries
	bpsOutResult, _, err := p.client.QueryRange(ctx, queryBpsOut, r)
	if err != nil {
		return nil, err
	}

	bpsInResult, _, err := p.client.QueryRange(ctx, queryBpsIn, r)
	if err != nil {
		return nil, err
	}

	bandwidthResult, _, err := p.client.QueryRange(ctx, queryBandwidth, r)
	if err != nil {
		return nil, err
	}

	bytesInResult, _, err := p.client.QueryRange(ctx, queryBytesIn, r)
	if err != nil {
		return nil, err
	}

	bytesOutResult, _, err := p.client.QueryRange(ctx, queryBytesOut, r)
	if err != nil {
		return nil, err
	}

	bpsOutMatrix, _ := bpsOutResult.(model.Matrix)
	bpsInMatrix, _ := bpsInResult.(model.Matrix)
	bandwidthMatrix, _ := bandwidthResult.(model.Matrix)
	bytesInMatrix, _ := bytesInResult.(model.Matrix)
	bytesOutMatrix, _ := bytesOutResult.(model.Matrix)

	// Create a map to aggregate data by instance, state, region and time
	trafficMap := make(map[string]map[int64]*TrafficByInstance)

	// Process BpsOut data
	for _, serie := range bpsOutMatrix {
		ip := string(serie.Metric["instance"])
		state := string(serie.Metric["state"])
		region := string(serie.Metric["region"])
		key := fmt.Sprintf("%s-%s-%s", ip, state, region)

		if _, ok := trafficMap[key]; !ok {
			trafficMap[key] = make(map[int64]*TrafficByInstance)
		}

		for _, point := range serie.Values {
			timeKey := int64(point.Timestamp) / 1000
			if _, ok := trafficMap[key][timeKey]; !ok {
				trafficMap[key][timeKey] = &TrafficByInstance{
					IP:     ip,
					State:  state,
					Region: region,
					Time:   time.Unix(timeKey, 0),
				}
			}
			trafficMap[key][timeKey].BpsOut = float64(point.Value)
		}
	}

	// Process BpsIn data
	for _, serie := range bpsInMatrix {
		ip := string(serie.Metric["instance"])
		state := string(serie.Metric["state"])
		region := string(serie.Metric["region"])
		key := fmt.Sprintf("%s-%s-%s", ip, state, region)

		if _, ok := trafficMap[key]; !ok {
			trafficMap[key] = make(map[int64]*TrafficByInstance)
		}

		for _, point := range serie.Values {
			timeKey := int64(point.Timestamp) / 1000
			if _, ok := trafficMap[key][timeKey]; !ok {
				trafficMap[key][timeKey] = &TrafficByInstance{
					IP:     ip,
					State:  state,
					Region: region,
					Time:   time.Unix(timeKey, 0),
				}
			}
			trafficMap[key][timeKey].BpsIn = float64(point.Value)
		}
	}

	// Process Bandwidth data
	for _, serie := range bandwidthMatrix {
		ip := string(serie.Metric["instance"])
		state := string(serie.Metric["state"])
		region := string(serie.Metric["region"])
		key := fmt.Sprintf("%s-%s-%s", ip, state, region)

		if _, ok := trafficMap[key]; !ok {
			trafficMap[key] = make(map[int64]*TrafficByInstance)
		}

		for _, point := range serie.Values {
			timeKey := int64(point.Timestamp) / 1000
			if _, ok := trafficMap[key][timeKey]; !ok {
				trafficMap[key][timeKey] = &TrafficByInstance{
					IP:     ip,
					State:  state,
					Region: region,
					Time:   time.Unix(timeKey, 0),
				}
			}
			trafficMap[key][timeKey].Bandwidth = float64(point.Value)
		}
	}

	// Process BytesIn data
	for _, serie := range bytesInMatrix {
		ip := string(serie.Metric["instance"])
		state := string(serie.Metric["state"])
		region := string(serie.Metric["region"])
		key := fmt.Sprintf("%s-%s-%s", ip, state, region)

		if _, ok := trafficMap[key]; !ok {
			trafficMap[key] = make(map[int64]*TrafficByInstance)
		}

		for _, point := range serie.Values {
			timeKey := int64(point.Timestamp) / 1000
			if _, ok := trafficMap[key][timeKey]; !ok {
				trafficMap[key][timeKey] = &TrafficByInstance{
					IP:     ip,
					State:  state,
					Region: region,
					Time:   time.Unix(timeKey, 0),
				}
			}
			trafficMap[key][timeKey].BytesIn = float64(point.Value)
		}
	}

	// Process BytesOut data
	for _, serie := range bytesOutMatrix {
		ip := string(serie.Metric["instance"])
		state := string(serie.Metric["state"])
		region := string(serie.Metric["region"])
		key := fmt.Sprintf("%s-%s-%s", ip, state, region)

		if _, ok := trafficMap[key]; !ok {
			trafficMap[key] = make(map[int64]*TrafficByInstance)
		}

		for _, point := range serie.Values {
			timeKey := int64(point.Timestamp) / 1000
			if _, ok := trafficMap[key][timeKey]; !ok {
				trafficMap[key][timeKey] = &TrafficByInstance{
					IP:     ip,
					State:  state,
					Region: region,
					Time:   time.Unix(timeKey, 0),
				}
			}
			trafficMap[key][timeKey].BytesOut = float64(point.Value)
		}
	}

	// Convert the map to a slice
	var trafficData []TrafficByInstance
	for _, timeMap := range trafficMap {
		for _, traffic := range timeMap {
			trafficData = append(trafficData, *traffic)
		}
	}

	return trafficData, nil
}
