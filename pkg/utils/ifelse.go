package utils

func IIf(expr bool, whenTrue, whenFalse string) string {
	if expr {
		return whenTrue
	}
	return whenFalse
}
