package gui

import (
	"korok.io/korok/math/f32"
	"math"
)

const PI float32 = 3.14

func InvLength(v f32.Vec2, fail float32) float32 {
	return 1/float32(math.Sqrt(float64(v[0] * v[0] + v[1] * v[1])))
}


