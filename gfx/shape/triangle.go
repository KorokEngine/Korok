package shape

import "github.com/go-gl/mathgl/mgl32"

type Triangle struct {
	shape

	// 三角形顶点
	vertex [3]mgl32.Vec2
}

func (tri*Triangle)Fill() mgl32.Vec3 {
	return tri.fill
}

func (tri*Triangle)Border() mgl32.Vec3 {
	return tri.border
}

func (tri*Triangle)Vertex() []mgl32.Vec2{
	return tri.vertex[:]
}




