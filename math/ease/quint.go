package ease

func InQuint(t float64) float64 {
	return t * t * t * t * t
}

func OutQuint(t float64) float64 {
	t -= 1
	return t*t*t*t*t + 1
}

func InOutQuint(t float64) float64 {
	t *= 2
	if t < 1 {
		return 0.5 * t * t * t * t * t
	} else {
		t -= 2
		return 0.5 * (t*t*t*t*t + 2)
	}
}
