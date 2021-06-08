package utils

func IIf(expr bool, whenTrue, whenFalse string) string {
	if expr {
		return whenTrue
	}
	return whenFalse
}

func IIfInt64(expr bool, whenTrue, whenFalse int64) int64 {
	if expr {
		return whenTrue
	}
	return whenFalse
}

func IIfF64(expr bool, whenTrue, whenFalse float64) float64 {
	if expr {
		return whenTrue
	}
	return whenFalse
}
