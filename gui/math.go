package gui

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const PI float32 = 3.14

func InvLength(v mgl32.Vec2, fail float32) float32 {
	return 1/float32(math.Sqrt(float64(v[0] * v[0] + v[1] * v[1])))
}


