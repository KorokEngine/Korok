package ease


func InQuad(t float64) float64 {
	return t * t
}

func OutQuad(t float64) float64 {
	return -t * (t - 2)
}

func InOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	} else {
		t = 2*t - 1
		return -0.5 * (t*(t-2) - 1)
	}
}
