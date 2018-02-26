package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok.io/korok/engi"
)

type CameraMode uint8
const (
	Perspective CameraMode = iota
	Orthographic
)

//
type Camera struct {
	Eye mgl32.Vec3
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

