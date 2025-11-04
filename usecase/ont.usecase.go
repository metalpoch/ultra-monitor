package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/snmp"
	"github.com/metalpoch/ultra-monitor/internal/utils"
	"github.com/metalpoch/ultra-monitor/repository"
)

type OntUsecase struct {
	repo  repository.OntRepository
	cache *cache.Redis
}

func NewOntUsecase(db *sqlx.DB, cache *cache.Redis) *OntUsecase {
	return &OntUsecase{repository.NewOntRepository(db), cache}
}

// UpdateTrafficForONT updates traffic for a single ONT
func (use *OntUsecase) UpdateTrafficForONT(ctx context.Context, ontID int32, community string) error {
	ont, err := use.repo.GetByID(ctx, ontID)
	if err != nil {
		return fmt.Errorf("error getting ONT: %v", err)
	}

	if !ont.Enabled {
		return fmt.Errorf("ONT %d is disabled", ontID)
	}

	currentTime := time.Now()
	ontTraffic, err := use.getONTData(ctx, ont, currentTime, community)
	if err != nil {
		return fmt.Errorf("error getting traffic data: %v", err)
	}

	if ontTraffic != nil {
		if err := use.repo.CreateTraffic(ctx, *ontTraffic); err != nil {
			return fmt.Errorf("error saving traffic data: %v", err)
		}
		log.Printf("Successfully saved traffic data for ONT %d", ontID)
	}

	return nil
}

// UpdateTrafficForAllONTs processes all ONTs in parallel (for backward compatibility)
func (use *OntUsecase) UpdateTrafficForAllONTs(ctx context.Context, community string) error {
	onts, err := use.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting ONTs: %v", err)
	}

	log.Printf("Processing traffic for %d ONTs", len(onts))

	var trafficData []entity.OntTraffic
	currentTime := time.Now()

	// Use goroutines for parallel processing
	type result struct {
		traffic *entity.OntTraffic
		err     error
	}

	results := make(chan result, len(onts))
	var wg sync.WaitGroup

	// Process each ONT in parallel
	for _, ont := range onts {
		if !ont.Enabled {
			continue
		}

		wg.Add(1)
		go func(ont entity.Ont) {
			defer wg.Done()
			ontTraffic, err := use.getONTData(ctx, ont, currentTime, community)
			results <- result{ontTraffic, err}
		}(ont)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Collect results
	for res := range results {
		if res.err != nil {
			log.Printf("Error processing ONT: %v", res.err)
			continue
		}
		if res.traffic != nil {
			trafficData = append(trafficData, *res.traffic)
		}
	}

	if len(trafficData) > 0 {
		if err := use.repo.CreateTrafficBatch(ctx, trafficData); err != nil {
			return fmt.Errorf("error saving traffic data: %v", err)
		}
		log.Printf("Successfully saved traffic data for %d ONTs", len(trafficData))
	}

	return nil
}

func (use *OntUsecase) getONTData(ctx context.Context, ont entity.Ont, currentTime time.Time, community string) (*entity.OntTraffic, error) {
	// Parse ont_idx to extract pon_idx and ont_idx
	ponIdx, ontIdx, err := utils.ParseOntIDX(ont.OntIDX)
	if err != nil {
		return nil, fmt.Errorf("error parsing ont_idx '%s': %v", ont.OntIDX, err)
	}

	// Create SNMP client
	snmpClient := snmp.NewSnmp(snmp.Config{
		IP:        ont.IP,
		Community: community,
		Timeout:   30 * time.Second,
		Retries:   3,
	})

	// Get ONT data via SNMP
	ontData, err := snmpClient.OntQuery(ponIdx, ontIdx)
	if err != nil {
		return nil, fmt.Errorf("SNMP query failed: %v", err)
	}

	// Update ONT status and last_check
	status := ontData.ControlRunStatus == 1
	if err := use.updateONTStatus(ctx, ont.ID, status, currentTime); err != nil {
		log.Printf("Warning: failed to update ONT %d status: %v", ont.ID, err)
	}

	// Calculate traffic using previous data from cache
	bpsIn, bpsOut, bytesIn, bytesOut, shouldSave, err := use.calculateTraffic(ctx, ont.ID, ontData.BytesIn, ontData.BytesOut, currentTime)
	if err != nil {
		return nil, fmt.Errorf("error calculating traffic: %v", err)
	}

	// If this is the first measurement, don't save to database
	if !shouldSave {
		return nil, nil
	}

	ontTraffic := &entity.OntTraffic{
		OntID:       ont.ID,
		Time:        currentTime,
		BpsIn:       bpsIn,
		BpsOut:      bpsOut,
		BytesIn:     bytesIn,
		BytesOut:    bytesOut,
		Temperature: ontData.Temperature,
		Rx:          ontData.Rx,
		Tx:          ontData.Tx,
	}

	return ontTraffic, nil
}

func (use *OntUsecase) calculateTraffic(ctx context.Context, ontID int32, currentBytesIn, currentBytesOut uint64, currentTime time.Time) (float64, float64, float64, float64, bool, error) {
	cacheKey := fmt.Sprintf("ont_traffic:%d", ontID)

	var cachedData dto.OntTrafficCache
	err := use.cache.FindOne(ctx, cacheKey, &cachedData)

	// If no cached data, store current data and return zeros for traffic
	if err != nil {
		newCacheData := dto.OntTrafficCache{
			BytesIn:   currentBytesIn,
			BytesOut:  currentBytesOut,
			LastCheck: currentTime,
		}

		// Cache for 1 hour
		if err := use.cache.InsertOne(ctx, cacheKey, time.Hour, newCacheData); err != nil {
			log.Printf("Warning: failed to cache traffic data for ONT %d: %v", ontID, err)
		}

		// Return false to indicate this is the first measurement (should not be saved to DB)
		return 0, 0, 0, 0, false, nil
	}

	// Calculate time difference in seconds
	timeDiff := currentTime.Sub(cachedData.LastCheck).Seconds()
	if timeDiff <= 0 {
		return 0, 0, 0, 0, false, nil
	}

	// Calculate bytes difference
	bytesInDiff := float64(currentBytesIn) - float64(cachedData.BytesIn)
	bytesOutDiff := float64(currentBytesOut) - float64(cachedData.BytesOut)

	// Calculate bps (bits per second)
	bpsIn := (bytesInDiff * 8) / timeDiff
	bpsOut := (bytesOutDiff * 8) / timeDiff

	// Update cache with current data
	newCacheData := dto.OntTrafficCache{
		BytesIn:   currentBytesIn,
		BytesOut:  currentBytesOut,
		LastCheck: currentTime,
	}

	if err := use.cache.InsertOne(ctx, cacheKey, time.Hour, newCacheData); err != nil {
		log.Printf("Warning: failed to update cache for ONT %d: %v", ontID, err)
	}

	// Return true to indicate this is a valid measurement (should be saved to DB)
	return bpsIn, bpsOut, bytesInDiff, bytesOutDiff, true, nil
}

func (use *OntUsecase) updateONTStatus(ctx context.Context, ontID int32, status bool, lastCheck time.Time) error {
	return use.repo.UpdateStatus(ctx, ontID, status, lastCheck)
}

func (use *OntUsecase) SnmpSerialDespt(device dto.AllSerialDesptByPon) ([]dto.OntSerialsAndDespts, error) {
	snmpClient := snmp.NewSnmp(snmp.Config{
		IP:        device.IP,
		Community: device.Community,
		Timeout:   10 * time.Second,
		Retries:   3,
	})

	ontData, err := snmpClient.OntSerialsAndDespts(device.PonIdx)
	if err != nil {
		return nil, fmt.Errorf("SNMP query failed: %v", err)
	}

	var res []dto.OntSerialsAndDespts
	for idx, data := range ontData {
		res = append(res, dto.OntSerialsAndDespts{
			PonIdx:       device.PonIdx,
			OntIdx:       uint8(idx),
			Despt:        data.Despt,
			SerialNumber: data.SerialNumber,
		})
	}

	return res, nil
}

func (use *OntUsecase) SaveFromSNMP(device dto.CreateOntRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	snmpClient := snmp.NewSnmp(snmp.Config{
		IP:        device.IP,
		Community: device.Community,
		Timeout:   10 * time.Second,
		Retries:   3,
	})

	ontData, err := snmpClient.OntQuery(device.PonIdx, device.OntIdx)
	if err != nil {
		return fmt.Errorf("SNMP query failed: %v", err)
	}

	return use.repo.Create(ctx, entity.Ont{
		IP:          device.IP,
		OntIDX:      fmt.Sprintf("%d.%d", device.PonIdx, device.OntIdx),
		Serial:      ontData.SerialNumber,
		Despt:       ontData.Despt,
		LineProf:    ontData.LineProfName,
		Description: device.Description,
		OltDistance: ontData.ControlRanging,
		Status:      ontData.ControlRunStatus == 1,
		LastCheck:   time.Now(),
	})
}

func (use *OntUsecase) GetByID(id int32) (dto.OntResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ont, err := use.repo.GetByID(ctx, id)
	if err != nil {
		return dto.OntResponse{}, err
	}

	return dto.OntResponse{
		ID:          ont.ID,
		IP:          ont.IP,
		OntIDX:      ont.OntIDX,
		Serial:      ont.Serial,
		Despt:       ont.Despt,
		LineProf:    ont.LineProf,
		Description: ont.Description,
		Enabled:     ont.Enabled,
		Status:      ont.Status,
		LastCheck:   ont.LastCheck.Format(time.RFC3339),
		CreatedAt:   ont.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (use *OntUsecase) GetAll() ([]dto.OntResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	onts, err := use.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.OntResponse
	for _, ont := range onts {
		responses = append(responses, dto.OntResponse{
			ID:          ont.ID,
			IP:          ont.IP,
			OntIDX:      ont.OntIDX,
			Serial:      ont.Serial,
			Despt:       ont.Despt,
			LineProf:    ont.LineProf,
			Description: ont.Description,
			Enabled:     ont.Enabled,
			Status:      ont.Status,
			OltDistance: ont.OltDistance,
			LastCheck:   ont.LastCheck.Format(time.RFC3339),
			CreatedAt:   ont.CreatedAt.Format(time.RFC3339),
		})
	}

	return responses, nil
}

func (use *OntUsecase) Delete(id int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return use.repo.Delete(ctx, id)
}

func (use *OntUsecase) Enable(id int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return use.repo.Enable(ctx, id)
}

func (use *OntUsecase) Disable(id int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return use.repo.Disable(ctx, id)
}
