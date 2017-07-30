package shape

import "github.com/go-gl/mathgl/mgl32"

type Circle struct {
	shape
	radius float32
}

func (c *Circle) Fill() mgl32.Vec3 {
	return c.fill
}

func (c *Circle) Border() mgl32.Vec3 {
	return c.border
}

// TODO!!
func (c *Circle) Vertex() []mgl32.Vec2 {
	return nil
}

