package dto

import "time"

type TrendPrediction struct {
	FutureDays    int       `query:"futureDays" validate:"required,min=1,max=30"`
	Confidence    float64   `query:"confidence" validate:"omitempty,min=0.5,max=0.99"`
	InitDate      time.Time `query:"initDate" validate:"required"`
	FinalDate     time.Time `query:"finalDate" validate:"required"`
}

type TrendResponse struct {
	Predictions []TrendDataPoint `json:"predictions"`
	Metrics     TrendMetrics     `json:"metrics"`
	TrendType   string           `json:"trend_type"`
}

type TrendDataPoint struct {
	Date          time.Time `json:"date"`
	PredictedBps  float64   `json:"predicted_bps"`
	LowerBound    float64   `json:"lower_bound,omitempty"`
	UpperBound    float64   `json:"upper_bound,omitempty"`
}

type TrendMetrics struct {
	Slope      float64 `json:"slope"`
	Intercept  float64 `json:"intercept"`
	RSquared   float64 `json:"r_squared"`
	IsIncreasing bool   `json:"is_increasing"`
	IsDecreasing bool   `json:"is_decreasing"`
}

