package math

import (
	"unsafe"
	"math/rand"
	"math"
	"korok.io/korok/math/f32"
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

func InvLength(x, y , fail float32) float32 {
	return 1/float32(math.Sqrt(float64(x*x + y*y)))
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

func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

// Radian converts degree to radian.
func Radian(d float32) float32 {
	return d * Pi / 180
}

// Degree converts radian to degree.
func Degree(r float32) float32 {
	return r * 180 / Pi
}


func AngleTo(v1, v2 f32.Vec2) (dot float32) {
	l1 := InvLength(v1[0], v1[1], 1)
	l2 := InvLength(v2[0], v2[1], 1)

	x1, y1 := v1[0]*l1, v1[1]*l1
	x2, y2 := v2[0]*l2, v2[1]*l2
	dot = x1 * x2 + y1 * y2
	return
}

func Direction(v1, v2 f32.Vec2) float32 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}

func Angle(v f32.Vec2) float32 {
	return float32(math.Atan2(float64(v[1]), float64(v[0])))
}

func Vector(a float32) f32.Vec2 {
	return f32.Vec2{Cos(a), Sin(a)}
}