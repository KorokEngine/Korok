package ease

import "math"

func InExpo(t float64) float64 {
	if t == 0 {
		return 0
	} else {
		return math.Pow(2, 10*(t-1))
	}
}

func OutExpo(t float64) float64 {
	if t == 1 {
		return 1
	} else {
		return 1 - math.Pow(2, -10*t)
	}
}

func InOutExpo(t float64) float64 {
	if t == 0 {
		return 0
	} else if t == 1 {
		return 1
	} else {
		if t < 0.5 {
			return 0.5 * math.Pow(2, (20*t)-10)
		} else {
			return 1 - 0.5*math.Pow(2, (-20*t)+10)
		}
	}
}

