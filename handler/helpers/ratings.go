package helpers

import "github.com/evorts/feednomity/pkg/utils"

func GetRating(e *utils.Eval, labels []string, threshold [][]string, value float64) string {
	if e == nil {
		return ""
	}
	if len(labels) != len(threshold) {
		return ""
	}
	loopThreshold:
	for idx, th := range threshold {
		for _, expr := range th {
			if !e.SetExpression(expr).Evaluate(value) {
				continue loopThreshold
			}
		}
		return labels[idx]
	}
	return ""
}
