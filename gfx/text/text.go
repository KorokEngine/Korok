package text

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok/gfx"
)

type TextQuad struct {
	// local shit
	xOffset, yOffset float32

	// size
	w, h float32

	// texture
	region gfx.Region
}

type TextData struct {
	Chars []TextQuad

	// texture-id
	tex uint16

	// rgba
	color uint32
}

func (td *TextData) Fill(buf []gfx.PosTexColorVertex, p mgl32.Vec2) {
	for i, char := range td.Chars {
		vi := i * 4

		v := &buf[vi+0]
		v.X = p[0] + char.xOffset
		v.Y = p[1] + char.yOffset
		v.U = char.region.X1
		v.V = char.region.Y1

		// index (1,0) <x,y,u,v>
		v = &buf[vi+1]
		v.X = p[0] + char.xOffset + char.w
		v.Y = p[1] + char.yOffset
		v.U = char.region.X2
		v.V = char.region.Y2

		// index(1,1) <x,y,u,v>
		v = &buf[vi+2]
		v.X = p[0] + char.xOffset + char.w
		v.Y = p[1] + char.yOffset + char.h
		v.U = char.region.X2
		v.V = char.region.Y1

		// index(0, 1) <x,y,u,v>
		v = &buf[vi+3]
		v.X = p[0] + char.xOffset
		v.Y = p[1] + char.yOffset + char.h
		v.U = char.region.X1
		v.V = char.region.Y1
	}
}

