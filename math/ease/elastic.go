package ease

import "math"

func InElastic(t float64) float64 {
	return InElasticFunction(0.5)(t)
}

func OutElastic(t float64) float64 {
	return OutElasticFunction(0.5)(t)
}

func InOutElastic(t float64) float64 {
	return InOutElasticFunction(0.5)(t)
}

func InElasticFunction(period float64) Function {
	p := period
	return func(t float64) float64 {
		t -= 1
		return -1 * (math.Pow(2, 10*t) * math.Sin((t-p/4)*(2*math.Pi)/p))
	}
}

func OutElasticFunction(period float64) Function {
	p := period
	return func(t float64) float64 {
		return math.Pow(2, -10*t)*math.Sin((t-p/4)*(2*math.Pi/p)) + 1
	}
}

func InOutElasticFunction(period float64) Function {
	p := period
	return func(t float64) float64 {
		t *= 2
		if t < 1 {
			t -= 1
			return -0.5 * (math.Pow(2, 10*t) * math.Sin((t-p/4)*2*math.Pi/p))
		} else {
			t -= 1
			return math.Pow(2, -10*t)*math.Sin((t-p/4)*2*math.Pi/p)*0.5 + 1
		}
	}
}
