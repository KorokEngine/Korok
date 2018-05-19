package gfx

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/engi"
	"korok.io/korok/math"
)

type CameraMode uint8
const (
	Perspective CameraMode = iota
	Orthographic
)

//
type Camera struct {
	Eye f32.Vec3
	//
	bound struct{
		left, top, right, bottom float32
	}
	pos struct{
		x, y float32
	}
	view struct{
		w, h float32
	}
	follow engi.Entity
}

func (c *Camera) Flow(entity engi.Entity) {
	c.follow = entity
}

func (c *Camera) MoveTo(x, y float32) {
	c.pos.x, c.pos.y = x, y
	c.clamp()
}

func (c *Camera) MoveBy(dx, dy float32) {
	c.pos.x += dx
	c.pos.y += dy
	c.clamp()
}

func (c *Camera) SetBound(left, top, right, bottom float32) {
	c.bound.left = left
	c.bound.right = right
	c.bound.top = top
	c.bound.bottom = bottom
	c.clamp()
}

func (c *Camera) SetViewPort(w, h float32) {
	c.view.w = w
	c.view.h = h
	c.clamp()
}

func (c *Camera) clamp() {
	// x
	if left := c.pos.x - c.view.w/2; left < c.bound.left {
		c.pos.x += c.bound.left - left
	} else if right := c.pos.x + c.view.w/2; right > c.bound.right {
		c.pos.x += c.bound.right - right
	}

	// y
	if bottom := c.pos.y - c.view.h/2; bottom < c.bound.bottom {
		c.pos.y += c.bound.bottom - bottom
	} else if top := c.pos.y + c.view.h/2; top > c.bound.top {
		c.pos.y += c.bound.top - top
	}
}

func (c *Camera) InView(xf *Transform, size, gravity f32.Vec2) bool {
	if xf.world.Rotation == 0 { // happy path
		p := xf.world.Position
		size[0], size[1] = size[0]*xf.world.Scale[0], size[1]*xf.world.Scale[1]
		a := AABB{p[0]-size[0]*gravity[0], p[1]-size[1]*gravity[1], size[0], size[1]}
		b := AABB{c.pos.x-c.view.w/2, c.pos.y-c.view.h/2, c.view.w, c.view.h}
		return OverlapAB(&a, &b)
	} else {
		srt := xf.world
		m := mat3{}; m.Initialize(srt.Position[0], srt.Position[1], srt.Rotation, srt.Scale[0], srt.Scale[1])
		// center and extent
		cx, cy := -size[0]*gravity[0] + size[0]/2, -size[1]*gravity[1] + size[1]/2
		ex, ey := size[0]/2, size[1]/2

		// transform center
		cx, cy = m.TransformCoord(cx, cy)

		// transform extent
		for i, v := range m {
			if v < 0 {
				m[i] = -v
			}
		}
		ex, ey = m.TransformNormal(ex, ey)
		a := AABB{cx-ex, cy-ey, ex*2, ey*2}
		b := AABB{c.pos.x-c.view.w/2, c.pos.y-c.view.h/2, c.view.w, c.view.h}
		return OverlapAB(&a, &b)
	}
}

type mat3 [9]float32 // fast culling matrix, (0, 0) as the center of the local model

func (m *mat3) Initialize(x, y, angle, sx, sy float32) {
	c, s := math.Cos(angle), math.Sin(angle)

	m[0] = c * sx
	m[1] = s * sx
	m[3] =  - s * sy
	m[4] =  + c * sy
	m[6] = x
	m[7] = y

	m[2], m[5] = 0, 0
	m[8] = 1.0
}

func (m mat3) TransformCoord(x, y float32) (x1, y1 float32) {
	x1 = m[0]*x + m[3]*y + m[6]
	y1 = m[1]*x + m[4]*y + m[7]
	return
}

func (m mat3) TransformNormal(x, y float32) (x1, y1 float32) {
	x1 = m[0]*x + m[3]*y
	y1 = m[1]*x + m[4]*y
	return
}