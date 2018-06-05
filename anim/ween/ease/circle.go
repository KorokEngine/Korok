package ease

import "math"

func InCirc(t float64) float64 {
	return -1 * (math.Sqrt(1-t*t) - 1)
}

func OutCirc(t float64) float64 {
	t -= 1
	return math.Sqrt(1 - (t * t))
}

func InOutCirc(t float64) float64 {
	t *= 2
	if t < 1 {
		return -0.5 * (math.Sqrt(1-t*t) - 1)
	} else {
		t = t - 2
		return 0.5 * (math.Sqrt(1-t*t) + 1)
	}
}
