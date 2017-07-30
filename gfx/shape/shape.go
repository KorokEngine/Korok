package shape

import "github.com/go-gl/mathgl/mgl32"

type shape struct {
	// 填充颜色
	fill mgl32.Vec3

	// 边界颜色
	border mgl32.Vec3
}

type Shape interface {
	Fill() mgl32.Vec3
	Border() mgl32.Vec3

	Vertex() []mgl32.Vec2
}
