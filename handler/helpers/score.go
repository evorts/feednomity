package helpers

import "github.com/evorts/feednomity/domain/assessments"

func CalculateScore(factors *assessments.Factor) float64 {
	var score float64
	if factors.Weight > 0 && factors.Rating > 0 {
		score += (float64(factors.Weight) / 100) * float64(factors.Rating)
	}
	if len(factors.Items) < 1 {
		return score
	}
	for _, v := range factors.Items {
		score += CalculateScore(v)
	}
	return score
}
