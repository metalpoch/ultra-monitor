package trend

import (
	"errors"
	"math"
)

// Trend represents a linear regression trend analyzer
type Trend struct {
	data []float64
}

// NewTrend creates a new Trend instance with the provided data
func NewTrend(data []float64) (*Trend, error) {
	if len(data) < 2 {
		return nil, errors.New("at least 2 data points are required for trend analysis")
	}
	return &Trend{data: data}, nil
}

// LinearRegression performs linear regression on the data
func (t *Trend) LinearRegression() (slope, intercept, rSquared float64, err error) {
	n := float64(len(t.data))
	if n < 2 {
		return 0, 0, 0, errors.New("insufficient data for linear regression")
	}

	var sumIndex, sumIndex2, sumValue, sumIndexValue, sumValue2 float64

	for i, value := range t.data {
		index := float64(i + 1)
		sumIndex += index
		sumIndex2 += index * index
		sumValue += value
		sumIndexValue += index * value
		sumValue2 += value * value
	}

	denominator := n*sumIndex2 - sumIndex*sumIndex
	if denominator == 0 {
		return 0, 0, 0, errors.New("cannot compute linear regression for this data")
	}

	slope = (n*sumIndexValue - sumIndex*sumValue) / denominator
	intercept = (sumValue - slope*sumIndex) / n

	// Calculate R-squared (coefficient of determination)
	var ssTotal, ssResidual float64
	meanValue := sumValue / n

	for i, value := range t.data {
		index := float64(i + 1)
		predicted := slope*index + intercept
		ssResidual += (value - predicted) * (value - predicted)
		ssTotal += (value - meanValue) * (value - meanValue)
	}

	if ssTotal == 0 {
		rSquared = 1.0 // Perfect fit if all values are the same
	} else {
		rSquared = 1 - (ssResidual / ssTotal)
	}

	return slope, intercept, rSquared, nil
}

// Prediction generates future predictions based on the linear regression
func (t *Trend) Prediction(futureDays int) ([]float64, error) {
	if futureDays <= 0 {
		return nil, errors.New("futureDays must be positive")
	}

	slope, intercept, _, err := t.LinearRegression()
	if err != nil {
		return nil, err
	}

	predictions := make([]float64, futureDays)
	for i := 0; i < futureDays; i++ {
		index := float64(len(t.data) + i + 1)
		predictions[i] = slope*index + intercept
		// Ensure predictions are non-negative (traffic can't be negative)
		if predictions[i] < 0 {
			predictions[i] = 0
		}
	}

	return predictions, nil
}

// PredictionWithConfidence generates predictions with confidence intervals
func (t *Trend) PredictionWithConfidence(futureDays int, confidenceLevel float64) ([]float64, []float64, []float64, error) {
	if futureDays <= 0 {
		return nil, nil, nil, errors.New("futureDays must be positive")
	}
	if confidenceLevel <= 0 || confidenceLevel >= 1 {
		return nil, nil, nil, errors.New("confidenceLevel must be between 0 and 1")
	}

	slope, intercept, _, err := t.LinearRegression()
	if err != nil {
		return nil, nil, nil, err
	}

	n := float64(len(t.data))
	var ssResidual float64
	for i, value := range t.data {
		index := float64(i + 1)
		predicted := slope*index + intercept
		ssResidual += (value - predicted) * (value - predicted)
	}

	stdError := math.Sqrt(ssResidual / (n - 2))

	// Calculate t-value for confidence interval (simplified)
	tValue := 1.96 // For 95% confidence, can be made more sophisticated

	predictions := make([]float64, futureDays)
	lowerBounds := make([]float64, futureDays)
	upperBounds := make([]float64, futureDays)

	for i := 0; i < futureDays; i++ {
		index := float64(len(t.data) + i + 1)
		prediction := slope*index + intercept
		if prediction < 0 {
			prediction = 0
		}
		predictions[i] = prediction

		// Simplified confidence interval calculation
		margin := tValue * stdError * math.Sqrt(1+1/n)
		lowerBounds[i] = math.Max(0, prediction-margin)
		upperBounds[i] = prediction + margin
	}

	return predictions, lowerBounds, upperBounds, nil
}

// GetTrendMetrics returns trend analysis metrics
func (t *Trend) GetTrendMetrics() (slope, intercept, rSquared float64, err error) {
	return t.LinearRegression()
}

// IsIncreasing returns true if the trend is increasing
func (t *Trend) IsIncreasing() (bool, error) {
	slope, _, _, err := t.LinearRegression()
	if err != nil {
		return false, err
	}
	return slope > 0, nil
}

// IsDecreasing returns true if the trend is decreasing
func (t *Trend) IsDecreasing() (bool, error) {
	slope, _, _, err := t.LinearRegression()
	if err != nil {
		return false, err
	}
	return slope < 0, nil
}
