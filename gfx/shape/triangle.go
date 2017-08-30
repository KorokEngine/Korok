package shape

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok/gfx"
)

/**
	图形的特点：
	1. 可以没有纹理坐标，有颜色
	2. 有形状, 填充色和边界
	TODO：有待思考实现方式
 */

type Triangle struct {
	id uint32
	shape
	mesh gfx.Mesh

	// 三角形顶点
	vertex [3]mgl32.Vec2
}

func (tri *Triangle) Fill() mgl32.Vec3 {
	return tri.fill
}

func (tri *Triangle) Border() mgl32.Vec3 {
	return tri.border
}

func (tri *Triangle) Vertex() []mgl32.Vec2 {
	return tri.vertex[:]
}

//
func (tri *Triangle) SetVertex(v [3]mgl32.Vec2) {
	tri.vertex = v
}

func (tri *Triangle) Setup() {
	// tri.mesh.SetVertex()
}

func NewTriangle(v [3]mgl32.Vec2) *Triangle {
	tri := new(Triangle)
	tri.SetVertex(v)
	tri.Setup()
	return tri
}
