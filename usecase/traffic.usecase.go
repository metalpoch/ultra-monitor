package usecase

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/prometheus"
	"github.com/metalpoch/ultra-monitor/internal/trend"
	"github.com/metalpoch/ultra-monitor/repository"
	"github.com/redis/go-redis/v9"
)

type TrafficUsecase struct {
	repo       repository.TrafficRepository
	cache      *cache.Redis
	prometheus prometheus.Prometheus
}

// isToday checks if a date is today
type dateRange struct {
	start time.Time
	end   time.Time
}

// generateCacheKey creates a cache key based on query parameters and date
func (use *TrafficUsecase) generateCacheKey(prefix, queryParam string, date time.Time) string {
	// Format date as YYYY-MM-DD for daily caching
	dateStr := date.Format("2006-01-02")
	if queryParam != "" {
		return fmt.Sprintf("%s:%s:%s", prefix, queryParam, dateStr)
	}
	return fmt.Sprintf("%s:%s", prefix, dateStr)
}

// getDateRanges splits a date range into individual days
func (use *TrafficUsecase) getDateRanges(initDate, finalDate time.Time) []dateRange {
	var ranges []dateRange
	current := time.Date(initDate.Year(), initDate.Month(), initDate.Day(), 0, 0, 0, 0, initDate.Location())
	final := time.Date(finalDate.Year(), finalDate.Month(), finalDate.Day(), 0, 0, 0, 0, finalDate.Location())

	for current.Before(final) || current.Equal(final) {
		startOfDay := current
		endOfDay := current.Add(24*time.Hour - time.Second)
		ranges = append(ranges, dateRange{start: startOfDay, end: endOfDay})
		current = current.Add(24 * time.Hour)
	}

	return ranges
}

// isToday checks if a date is today
func (use *TrafficUsecase) isToday(date time.Time) bool {
	today := time.Now()
	return date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day()
}

// getCachedTraffic retrieves cached traffic data for a date range
func (use *TrafficUsecase) getCachedTraffic(ctx context.Context, prefix, queryParam string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	var allTraffic []dto.Traffic

	dateRanges := use.getDateRanges(initDate, finalDate)
	for _, dateRange := range dateRanges {
		// Skip caching for today
		if use.isToday(dateRange.start) {
			continue
		}

		cacheKey := use.generateCacheKey(prefix, queryParam, dateRange.start)
		var dailyTraffic []dto.Traffic
		if err := use.cache.FindOne(ctx, cacheKey, &dailyTraffic); err == nil {
			// Cache hit, add all intervals to results
			allTraffic = append(allTraffic, dailyTraffic...)
		} else if err != redis.Nil {
			return nil, err
		}
		// Cache miss - we'll handle this in the main query logic
	}

	return allTraffic, nil
}

// cacheTrafficData stores all traffic data intervals in cache by day
func (use *TrafficUsecase) cacheTrafficData(ctx context.Context, prefix, queryParam string, trafficData []dto.Traffic) error {
	// Group traffic data by day, preserving all time intervals
	trafficByDay := make(map[string][]dto.Traffic)
	for _, traffic := range trafficData {
		dateKey := traffic.Time.Format("2006-01-02")
		trafficByDay[dateKey] = append(trafficByDay[dateKey], traffic)
	}

	// Cache each day's data with all intervals
	for dateStr, dailyTraffic := range trafficByDay {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		// Skip caching for today
		if use.isToday(date) {
			continue
		}

		cacheKey := use.generateCacheKey(prefix, queryParam, date)
		// Cache for 30 days (historical data)
		if err := use.cache.InsertOne(ctx, cacheKey, 30*24*time.Hour, dailyTraffic); err != nil {
			return err
		}
	}

	return nil
}

func NewTrafficUsecase(db *sqlx.DB, cache *cache.Redis, prometheus *prometheus.Prometheus) *TrafficUsecase {
	return &TrafficUsecase{repository.NewTrafficRepository(db), cache, *prometheus}
}

func (use *TrafficUsecase) DeviceLocation() ([]dto.DeviceLocation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	devices, err := use.prometheus.DeviceLocations(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.DeviceLocation
	for _, d := range devices {
		res = append(res, (dto.DeviceLocation)(d))
	}

	sort.SliceStable(res, func(i, j int) bool {
		if res[i].Region != res[j].Region {
			return res[i].Region < res[j].Region
		}

		if res[i].State != res[j].State {
			return res[i].State < res[j].State
		}
		return res[i].SysName < res[j].SysName
	})

	return res, nil
}

func (use *TrafficUsecase) InfoInstance(ip string) ([]dto.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	devices, err := use.prometheus.InstanceScan(ctx, ip)
	if err != nil {
		return nil, err
	}

	var res []dto.InfoDevice
	for _, d := range devices {
		res = append(res, (dto.InfoDevice)(d))
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].IfName < res[j].IfName
	})

	return res, nil
}

func (use *TrafficUsecase) UpdateSummaryTraffic(initDate, finalDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	trafficData, err := use.prometheus.TrafficByInstanceStateRegion(ctx, initDate, finalDate)
	if err != nil {
		return err
	}

	maxTrafficByIP := make(map[string]prometheus.TrafficByInstance)

	for _, record := range trafficData {
		totalTraffic := record.BpsIn + record.BpsOut

		// Check if we have a record for this IP already
		if existing, exists := maxTrafficByIP[record.IP]; exists {
			existingTotal := existing.BpsIn + existing.BpsOut
			if totalTraffic > existingTotal {
				maxTrafficByIP[record.IP] = record
			}
		} else {
			maxTrafficByIP[record.IP] = record
		}
	}

	var result []entity.SumaryTraffic
	for _, record := range maxTrafficByIP {
		result = append(result, entity.SumaryTraffic{
			Time:     record.Time,
			IP:       record.IP,
			State:    record.State,
			Region:   record.Region,
			Sysname:  record.SysName,
			BpsIn:    record.BpsIn,
			BpsOut:   record.BpsOut,
			BytesIn:  record.BytesIn,
			BytesOut: record.BytesOut,
		})
	}

	if err := use.repo.SaveSummaryTraffic(ctx, result); err != nil {
		log.Printf("UpdateSummaryTraffic: Error saving to repository: %v", err)
		return err
	}

	return nil
}

func (use *TrafficUsecase) GetTrafficByCriteria(criteria, value string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx := context.Background()

	// Try to get cached data first
	cachedTraffic, err := use.getCachedTraffic(ctx, "traffic:criteria:"+criteria, value, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	// Determine which dates need to be queried from Prometheus
	dateRanges := use.getDateRanges(initDate, finalDate)
	var missingRanges []dateRange

	for _, dateRange := range dateRanges {
		// Skip today (never cache today)
		if use.isToday(dateRange.start) {
			missingRanges = append(missingRanges, dateRange)
			continue
		}

		// Check if we have this date in cache
		cacheKey := use.generateCacheKey("traffic:criteria:"+criteria, value, dateRange.start)
		var dailyTraffic []dto.Traffic
		if err := use.cache.FindOne(ctx, cacheKey, &dailyTraffic); err != nil && err != redis.Nil {
			return nil, err
		} else if err == redis.Nil {
			// Cache miss, need to query this date
			missingRanges = append(missingRanges, dateRange)
		}
	}

	var prometheusTraffic []dto.Traffic
	if len(missingRanges) > 0 {
		// Query Prometheus for missing dates
		for _, dateRange := range missingRanges {
			traffic, err := use.prometheus.TrafficByCriteria(ctx, criteria, value, dateRange.start, dateRange.end)
			if err != nil {
				return nil, err
			}

			var result []dto.Traffic
			for _, t := range traffic {
				result = append(result, (dto.Traffic)(*t))
			}
			prometheusTraffic = append(prometheusTraffic, result...)

			// Cache the data (excluding today)
			if !use.isToday(dateRange.start) {
				if err := use.cacheTrafficData(ctx, "traffic:criteria:"+criteria, value, result); err != nil {
					log.Printf("Warning: failed to cache traffic data: %v", err)
				}
			}
		}
	}

	// Combine cached and fresh data
	allTraffic := append(cachedTraffic, prometheusTraffic...)

	// Sort by time
	sort.Slice(allTraffic, func(i, j int) bool {
		return allTraffic[i].Time.Before(allTraffic[j].Time)
	})

	return allTraffic, nil
}

func (use *TrafficUsecase) Regions(initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	results := make(dto.TrafficByLabel)

	keyCache := fmt.Sprintf("regions-%d-%d", initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &results); err == nil {
		return results, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffics, err := use.prometheus.TrafficGroupedByField(context.Background(), "", "", "region", initDate, finalDate)
	if err != nil {
		return nil, err
	}

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, results)

	return results, nil
}

func (use *TrafficUsecase) StatesByRegion(region string, initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	results := make(dto.TrafficByLabel)

	keyCache := fmt.Sprintf("statesByRegion-%s-%d-%d", region, initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &results); err == nil {
		return results, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffics, err := use.prometheus.TrafficGroupedByField(context.Background(), "region", region, "state", initDate, finalDate)
	if err != nil {
		return nil, err
	}

	if len(traffics) < 1 {
		return results, nil
	}

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, results)

	return results, nil
}

func (use *TrafficUsecase) SysnameByState(state string, initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	results := make(dto.TrafficByLabel)

	keyCache := fmt.Sprintf("oltsByState-%s-%d-%d", state, initDate.Unix(), finalDate.Unix())
	if err := use.cache.FindOne(context.Background(), keyCache, &results); err == nil {
		return results, nil
	} else if err != redis.Nil {
		return nil, err
	}

	traffics, err := use.prometheus.SysnameByState(context.Background(), state, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	if len(traffics) < 1 {
		return results, nil
	}

	for state, traffic := range traffics {
		var trafficState []dto.Traffic
		for _, t := range traffic {
			trafficState = append(trafficState, (dto.Traffic)(*t))
		}

		results[state] = trafficState
	}

	use.cache.InsertOne(context.Background(), keyCache, 8*time.Hour, results)

	return results, nil
}

func (use *TrafficUsecase) RegionStats(ip string, initDate, finalDate time.Time) ([]dto.StateStats, error) {
	stats, err := use.prometheus.StatesStatsByRegion(context.Background(), ip, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.StateStats
	for _, s := range stats {
		result = append(result, (dto.StateStats)(s))
	}

	return result, nil
}

func (use *TrafficUsecase) GetNationalTrend(prediction dto.TrendPrediction) (*dto.TrendResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTraffic(ctx, prediction.InitDate, prediction.FinalDate)
	if err != nil {
		return nil, err
	}

	return use.generateTrendResponse(trafficData, prediction, "national")
}

func (use *TrafficUsecase) GetRegionalTrend(region string, prediction dto.TrendPrediction) (*dto.TrendResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTrafficByRegion(ctx, region, prediction.InitDate, prediction.FinalDate)
	if err != nil {
		return nil, err
	}

	return use.generateTrendResponse(trafficData, prediction, "regional")
}

func (use *TrafficUsecase) GetStateTrend(state string, prediction dto.TrendPrediction) (*dto.TrendResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTrafficByState(ctx, state, prediction.InitDate, prediction.FinalDate)
	if err != nil {
		return nil, err
	}

	return use.generateTrendResponse(trafficData, prediction, "state")
}

func (use *TrafficUsecase) GetOLTTrend(ip string, prediction dto.TrendPrediction) (*dto.TrendResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTrafficByIP(ctx, ip, prediction.InitDate, prediction.FinalDate)
	if err != nil {
		return nil, err
	}

	return use.generateTrendResponse(trafficData, prediction, "olt")
}

func (use *TrafficUsecase) generateTrendResponse(trafficData []entity.TrafficSummary, prediction dto.TrendPrediction, trendType string) (*dto.TrendResponse, error) {
	if len(trafficData) < 2 {
		return nil, fmt.Errorf("insufficient data for trend analysis: need at least 2 data points, got %d", len(trafficData))
	}

	// Extract total traffic values (bps_in + bps_out) for trend analysis
	trafficValues := make([]float64, len(trafficData))
	for i, data := range trafficData {
		trafficValues[i] = data.TotalBpsIn + data.TotalBpsOut
	}

	// Create trend analyzer
	trendAnalyzer, err := trend.NewTrend(trafficValues)
	if err != nil {
		return nil, err
	}

	// Get trend metrics
	slope, intercept, rSquared, err := trendAnalyzer.GetTrendMetrics()
	if err != nil {
		return nil, err
	}

	// Generate predictions
	var predictions []dto.TrendDataPoint
	if prediction.Confidence > 0 {
		predictedValues, lowerBounds, upperBounds, err := trendAnalyzer.PredictionWithConfidence(prediction.FutureDays, prediction.Confidence)
		if err != nil {
			return nil, err
		}

		// Generate future dates starting from the last data point
		lastDate := trafficData[len(trafficData)-1].Time
		for i := 0; i < prediction.FutureDays; i++ {
			futureDate := lastDate.Add(time.Duration(i+1) * 24 * time.Hour)
			predictions = append(predictions, dto.TrendDataPoint{
				Date:         futureDate,
				PredictedBps: predictedValues[i],
				LowerBound:   lowerBounds[i],
				UpperBound:   upperBounds[i],
			})
		}
	} else {
		predictedValues, err := trendAnalyzer.Prediction(prediction.FutureDays)
		if err != nil {
			return nil, err
		}

		// Generate future dates starting from the last data point
		lastDate := trafficData[len(trafficData)-1].Time
		for i := 0; i < prediction.FutureDays; i++ {
			futureDate := lastDate.Add(time.Duration(i+1) * 24 * time.Hour)
			predictions = append(predictions, dto.TrendDataPoint{
				Date:         futureDate,
				PredictedBps: predictedValues[i],
			})
		}
	}

	// Determine trend direction
	isIncreasing, _ := trendAnalyzer.IsIncreasing()
	isDecreasing, _ := trendAnalyzer.IsDecreasing()

	response := &dto.TrendResponse{
		Predictions: predictions,
		Metrics: dto.TrendMetrics{
			Slope:        slope,
			Intercept:    intercept,
			RSquared:     rSquared,
			IsIncreasing: isIncreasing,
			IsDecreasing: isDecreasing,
		},
		TrendType: trendType,
	}

	return response, nil
}

func (use *TrafficUsecase) StateStats(ip string, initDate, finalDate time.Time) ([]dto.OltStats, error) {
	stats, err := use.prometheus.OltStatsByState(context.Background(), ip, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.OltStats
	for _, s := range stats {
		result = append(result, (dto.OltStats)(s))
	}

	return result, nil
}

func (use *TrafficUsecase) GponStats(ip string, initDate, finalDate time.Time) ([]dto.GponStats, error) {
	stats, err := use.prometheus.GponStatsByInstance(context.Background(), ip, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.GponStats
	for _, s := range stats {
		result = append(result, (dto.GponStats)(s))
	}

	return result, nil
}

func (use *TrafficUsecase) GroupIP(ips []string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx := context.Background()

	// Create a unique identifier for the IPs group
	ipsKey := strings.Join(ips, ",")

	// Try to get cached data first
	cachedTraffic, err := use.getCachedTraffic(ctx, "traffic:groupip", ipsKey, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	// Determine which dates need to be queried from Prometheus
	dateRanges := use.getDateRanges(initDate, finalDate)
	var missingRanges []dateRange

	for _, dateRange := range dateRanges {
		// Skip today (never cache today)
		if use.isToday(dateRange.start) {
			missingRanges = append(missingRanges, dateRange)
			continue
		}

		// Check if we have this date in cache
		cacheKey := use.generateCacheKey("traffic:groupip", ipsKey, dateRange.start)
		var dailyTraffic []dto.Traffic
		if err := use.cache.FindOne(ctx, cacheKey, &dailyTraffic); err != nil && err != redis.Nil {
			return nil, err
		} else if err == redis.Nil {
			// Cache miss, need to query this date
			missingRanges = append(missingRanges, dateRange)
		}
	}

	var prometheusTraffic []dto.Traffic
	if len(missingRanges) > 0 {
		// Query Prometheus for missing dates
		for _, dateRange := range missingRanges {
			traffic, err := use.prometheus.TrafficGroupInstance(ctx, ips, dateRange.start, dateRange.end)
			if err != nil {
				return nil, err
			}

			var result []dto.Traffic
			for _, t := range traffic {
				result = append(result, (dto.Traffic)(*t))
			}
			prometheusTraffic = append(prometheusTraffic, result...)

			// Cache the data (excluding today)
			if !use.isToday(dateRange.start) {
				if err := use.cacheTrafficData(ctx, "traffic:groupip", ipsKey, result); err != nil {
					log.Printf("Warning: failed to cache traffic data: %v", err)
				}
			}
		}
	}

	// Combine cached and fresh data
	allTraffic := append(cachedTraffic, prometheusTraffic...)

	// Sort by time
	sort.Slice(allTraffic, func(i, j int) bool {
		return allTraffic[i].Time.Before(allTraffic[j].Time)
	})

	return allTraffic, nil
}

func (use *TrafficUsecase) ByMunicipality(state, municipality string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetSnmpIndexByMunicipality(ctx, state, municipality)
	if err != nil {
		return nil, err
	}

	instancesMap := make(map[string]string)
	for _, r := range res {
		_, ok := instancesMap[r.IP]
		if ok {
			instancesMap[r.IP] += "|"
		}
		instancesMap[r.IP] += r.Idx
	}

	// Try to get cached data first
	cacheKey := fmt.Sprintf("municipality:%s:%s", state, municipality)
	cachedTraffic, err := use.getCachedTraffic(ctx, "traffic:municipality", cacheKey, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	// Determine which dates need to be queried from Prometheus
	dateRanges := use.getDateRanges(initDate, finalDate)
	var missingRanges []dateRange

	for _, dateRange := range dateRanges {
		// Skip today (never cache today)
		if use.isToday(dateRange.start) {
			missingRanges = append(missingRanges, dateRange)
			continue
		}

		// Check if we have this date in cache
		dailyCacheKey := use.generateCacheKey("traffic:municipality", cacheKey, dateRange.start)
		var dailyTraffic []dto.Traffic
		if err := use.cache.FindOne(ctx, dailyCacheKey, &dailyTraffic); err != nil && err != redis.Nil {
			return nil, err
		} else if err == redis.Nil {
			// Cache miss, need to query this date
			missingRanges = append(missingRanges, dateRange)
		}
	}

	var prometheusTraffic []dto.Traffic
	if len(missingRanges) > 0 {
		// Query Prometheus for missing dates and preserve all time intervals
		for _, dateRange := range missingRanges {
			var dailyTraffic []dto.Traffic
			for ip, indexes := range instancesMap {
				traffic, err := use.prometheus.TrafficInstanceByIndex(ctx, ip, indexes, dateRange.start, dateRange.end)
				if err != nil {
					return nil, err
				}

				// Convert and preserve all individual time intervals
				for _, t := range traffic {
					dailyTraffic = append(dailyTraffic, (dto.Traffic)(*t))
				}
			}

			// Add to prometheus traffic results
			prometheusTraffic = append(prometheusTraffic, dailyTraffic...)

			// Cache the data (excluding today) - preserving all intervals
			if !use.isToday(dateRange.start) && len(dailyTraffic) > 0 {
				if err := use.cacheTrafficData(ctx, "traffic:municipality", cacheKey, dailyTraffic); err != nil {
					log.Printf("Warning: failed to cache traffic data: %v", err)
				}
			}
		}
	}

	// Combine cached and fresh data
	allTraffic := append(cachedTraffic, prometheusTraffic...)

	// Sort by time
	sort.Slice(allTraffic, func(i, j int) bool {
		return allTraffic[i].Time.Before(allTraffic[j].Time)
	})

	return allTraffic, nil
}

func (use *TrafficUsecase) ByCounty(state, municipality, county string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetSnmpIndexByCounty(ctx, state, municipality, county)
	if err != nil {
		return nil, err
	}

	instancesMap := make(map[string]string)
	for _, r := range res {
		_, ok := instancesMap[r.IP]
		if ok {
			instancesMap[r.IP] += "|"
		}
		instancesMap[r.IP] += r.Idx
	}
	accum := make(map[time.Time]prometheus.Traffic)
	for ip, indexes := range instancesMap {
		traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, indexes, initDate, finalDate)
		if err != nil {
			return nil, err
		}

		for _, t := range traffic {
			key := t.Time.Truncate(15 * time.Minute)

			if data, ok := accum[key]; ok {
				data.BpsIn += t.BpsIn
				data.BpsOut += t.BpsOut
				data.BytesIn += t.BytesIn
				data.BytesOut += t.BytesOut
				data.Bandwidth += t.Bandwidth
				accum[key] = data
			} else {
				cloned := *t
				cloned.Time = key
				accum[key] = cloned
			}
		}
	}
	var result []dto.Traffic
	for _, val := range accum {
		result = append(result, (dto.Traffic)(val))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil

}

func (use *TrafficUsecase) ByODN(state, municipality, odn string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetSnmpIndexByODN(ctx, state, municipality, odn)
	if err != nil {
		return nil, err
	}

	instancesMap := make(map[string]string)
	for _, r := range res {
		_, ok := instancesMap[r.IP]
		if ok {
			instancesMap[r.IP] += "|"
		}
		instancesMap[r.IP] += r.Idx
	}
	accum := make(map[time.Time]prometheus.Traffic)
	for ip, indexes := range instancesMap {
		traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, indexes, initDate, finalDate)
		if err != nil {
			return nil, err
		}

		for _, t := range traffic {
			key := t.Time.Truncate(15 * time.Minute)

			if data, ok := accum[key]; ok {
				data.BpsIn += t.BpsIn
				data.BpsOut += t.BpsOut
				data.BytesIn += t.BytesIn
				data.BytesOut += t.BytesOut
				data.Bandwidth += t.Bandwidth
				accum[key] = data
			} else {
				cloned := *t
				cloned.Time = key
				accum[key] = cloned
			}
		}
	}
	var result []dto.Traffic
	for _, val := range accum {
		result = append(result, (dto.Traffic)(val))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (use *TrafficUsecase) ByIdx(ip, idx string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	traffic, err := use.prometheus.TrafficInstanceByIndex(context.Background(), ip, idx, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, t := range traffic {
		result = append(result, (dto.Traffic)(*t))
	}

	return result, nil
}

func (use *TrafficUsecase) GetNationalTraffic(initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTraffic(ctx, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, data := range trafficData {
		result = append(result, dto.Traffic{
			Time:     data.Time,
			BpsIn:    data.TotalBpsIn,
			BpsOut:   data.TotalBpsOut,
			BytesIn:  data.TotalBytesIn,
			BytesOut: data.TotalBytesOut,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (use *TrafficUsecase) GetRegionalTraffic(region string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTrafficByRegion(ctx, region, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, data := range trafficData {
		result = append(result, dto.Traffic{
			Time:     data.Time,
			BpsIn:    data.TotalBpsIn,
			BpsOut:   data.TotalBpsOut,
			BytesIn:  data.TotalBytesIn,
			BytesOut: data.TotalBytesOut,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (use *TrafficUsecase) GetStateTraffic(state string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTrafficByState(ctx, state, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, data := range trafficData {
		result = append(result, dto.Traffic{
			Time:     data.Time,
			BpsIn:    data.TotalBpsIn,
			BpsOut:   data.TotalBpsOut,
			BytesIn:  data.TotalBytesIn,
			BytesOut: data.TotalBytesOut,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (use *TrafficUsecase) GetOLTByIPTraffic(ip string, initDate, finalDate time.Time) ([]dto.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTotalTrafficByIP(ctx, ip, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	var result []dto.Traffic
	for _, data := range trafficData {
		result = append(result, dto.Traffic{
			Time:     data.Time,
			BpsIn:    data.TotalBpsIn,
			BpsOut:   data.TotalBpsOut,
			BytesIn:  data.TotalBytesIn,
			BytesOut: data.TotalBytesOut,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

func (use *TrafficUsecase) GetTrafficByRegions(initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTrafficGroupedByRegion(ctx, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	result := make(dto.TrafficByLabel)
	for region, trafficList := range trafficData {
		var trafficDTO []dto.Traffic
		for _, traffic := range trafficList {
			trafficDTO = append(trafficDTO, dto.Traffic{
				Time:     traffic.Time,
				BpsIn:    traffic.TotalBpsIn,
				BpsOut:   traffic.TotalBpsOut,
				BytesIn:  traffic.TotalBytesIn,
				BytesOut: traffic.TotalBytesOut,
			})
		}
		result[region] = trafficDTO
	}

	return result, nil
}

func (use *TrafficUsecase) GetTrafficByStates(region string, initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTrafficGroupedByState(ctx, region, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	result := make(dto.TrafficByLabel)
	for state, trafficList := range trafficData {
		var trafficDTO []dto.Traffic
		for _, traffic := range trafficList {
			trafficDTO = append(trafficDTO, dto.Traffic{
				Time:     traffic.Time,
				BpsIn:    traffic.TotalBpsIn,
				BpsOut:   traffic.TotalBpsOut,
				BytesIn:  traffic.TotalBytesIn,
				BytesOut: traffic.TotalBytesOut,
			})
		}
		result[state] = trafficDTO
	}

	return result, nil
}

func (use *TrafficUsecase) GetTrafficByIPs(state string, initDate, finalDate time.Time) (dto.TrafficByLabel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trafficData, err := use.repo.GetTrafficGroupedByIP(ctx, state, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	result := make(dto.TrafficByLabel)
	for ip, trafficList := range trafficData {
		var trafficDTO []dto.Traffic
		for _, traffic := range trafficList {
			trafficDTO = append(trafficDTO, dto.Traffic{
				Time:     traffic.Time,
				BpsIn:    traffic.TotalBpsIn,
				BpsOut:   traffic.TotalBpsOut,
				BytesIn:  traffic.TotalBytesIn,
				BytesOut: traffic.TotalBytesOut,
			})
		}
		result[ip] = trafficDTO
	}

	return result, nil
}

func (use *TrafficUsecase) GetLocationHierarchy(initDate, finalDate time.Time) (*dto.LocationHierarchy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hierarchy, err := use.repo.GetLocationHierarchy(ctx, initDate, finalDate)
	if err != nil {
		return nil, err
	}

	return hierarchy, nil
}
