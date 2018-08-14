package f32

import (
	"golang.org/x/image/math/f32"
	"math"
)

type Vec2 f32.Vec2
type Vec3 f32.Vec3
type Vec4 f32.Vec4


func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1[0] + v2[0], v1[1] + v2[1]}
}

func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v1[0] - v2[0], v1[1] - v2[1]}
}

func (v1 Vec2) Mul(c float32) Vec2 {
	return Vec2{v1[0] * c, v1[1] * c}
}

func (v1 Vec2) Len() float32 {
	return float32(math.Sqrt(float64(v1[0]*v1[0]+v1[1]*v1[1])))
}

func (v1 Vec2) IsZero() bool {
	return v1[0]==0 && v1[1]==0
}

func (v1 Vec2) Norm() Vec2 {
	d := float32(math.Sqrt(float64(v1[0]*v1[0]+v1[1]*v1[1])))
	return Vec2{v1[0]/d, v1[1]/d}
}

func (v1 Vec2) Dot(v2 Vec2) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func (v1 Vec2) Cross(v2 Vec2) float32 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}

func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

func (v1 Vec3) Mul(c float32) Vec3 {
	return Vec3{v1[0] * c, v1[1] * c, v1[2] * c}
}

