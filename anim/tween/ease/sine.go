package ease

import "math"

func InSine(t float64) float64 {
	return -1*math.Cos(t*math.Pi/2) + 1
}

func OutSine(t float64) float64 {
	return math.Sin(t * math.Pi / 2)
}

func InOutSine(t float64) float64 {
	return -0.5 * (math.Cos(math.Pi*t) - 1)
}

