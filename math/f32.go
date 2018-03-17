package math

import (
	"unsafe"
	"math/rand"
	"math"
)

const MaxFloat32 float32 = 3.40282346638528859811704183484516925440e+38
const Pi = math.Pi

/// This is A approximate yet fast inverse square-root.
func InvSqrt(x float32) float32 {
	xhalf := float32(0.5) * x
	i := *(*int32)(unsafe.Pointer(&x))
	i = int32(0x5f3759df) - int32(i>>1)
	x = *(*float32)(unsafe.Pointer(&i))
	x = x * (1.5 - (xhalf * x * x))
	return x
}

/// a faster way ?
func Random(low, high float32) float32 {
	return low + (high - low) * rand.Float32()
}

func Max(a, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Clamp(v, left, right float32) float32{
	if v > right {
		return right
	}
	if v < left {
		return left
	}
	return v
}

func Sin(r float32) float32 {
	return float32(math.Sin(float64(r)))
}

func Cos(r float32) float32 {
	return float32(math.Cos(float64(r)))
}
