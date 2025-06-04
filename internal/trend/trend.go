package trend

type trend struct {
	data []float64
}

func NewTrend(data []float64) *trend {
	return &trend{data}
}

func (t *trend) linearRegression() (slope, intercept float64) {
	n := float64(len(t.data))
	var sumIndex, sumIndex2, sumValue, sumIndexValue float64

	for i, value := range t.data {
		index := float64(i + 1)
		sumIndex += index
		sumIndex2 += index * index
		sumValue += value
		sumIndexValue += index * value
	}

	denominator := n*sumIndex2 - sumIndex*sumIndex
	if denominator == 0 {
		return 0, 0
	}

	slope = (n*sumIndexValue - sumIndex*sumValue) / denominator
	intercept = (sumValue - slope*sumIndex) / n

	return
}

func (t *trend) Prediction(futureDays int) []float64 {
	slope, intercept := t.linearRegression()

	predictions := make([]float64, futureDays)
	for i := 0; i < futureDays; i++ {

		index := float64(len(t.data) + i + 1)
		predictions[i] = slope*index + intercept
	}

	return predictions
}
