package util

import "strconv"

func ParseFloat64(v string) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return f
}
