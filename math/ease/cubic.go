package ease


func InCubic(t float64) float64 {
	return t * t * t
}

func OutCubic(t float64) float64 {
	t -= 1
	return t*t*t + 1
}

func InOutCubic(t float64) float64 {
	t *= 2
	if t < 1 {
		return 0.5 * t * t * t
	} else {
		t -= 2
		return 0.5 * (t*t*t + 2)
	}
}
