package ease


type Function func(float64) float64

func Linear(t float64) float64 {
	return t
}

func InSquare(t float64) float64 {
	if t < 1 {
		return 0
	} else {
		return 1
	}
}

func OutSquare(t float64) float64 {
	if t > 0 {
		return 1
	} else {
		return 0
	}
}

func InOutSquare(t float64) float64 {
	if t < 0.5 {
		return 0
	} else {
		return 1
	}
}
