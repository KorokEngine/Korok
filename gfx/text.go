package gfx

import (
	"korok/engi"
	"github.com/go-gl/mathgl/mgl32"
	"korok/gfx/font"
)

// 文字应该采用 BatchRender 绘制
// 如果使用 BatchRender 那么此处生成模型即可

type FontSystem interface {
	Glyph(rune rune) *font.Glyph
}

type TextQuad struct {
	// local shit
	xOffset, yOffset float32

	// size
	w, h float32

	// texture
	region Region
}

// TextSprite
type TextComp struct {
	engi.Entity
	font FontSystem

	scale float32
	color uint32
	text  string

	// TextModel
	vertex []TextQuad
	runeCount int32
}

func (tc *TextComp) SetColor(color uint32) {
	tc.color = color
}

func (tc *TextComp) SetText(text string) {
	tc.text = text
	// init ebo, vbo
	tc.runeCount  = int32(len(text))

	// fill data
	tc.fillData()
}

// generate text-vertex with the string
//
//		+----------+
//		| . 	   |
//      |   .	   |
//		|     .    |
// 		|		.  |
// 		+----------+
// 1 * 1 quad for each char
func (tc *TextComp) fillData() {
	var xOffset float32
	var yOffset float32

	chars := make([]TextQuad, len(tc.text))
	tc.vertex = chars

	for i, r := range tc.text {
		if glyph := tc.font.Glyph(r); glyph != nil {
			advance := float32(glyph.Advance)
			vw := glyph.Width
			vh := glyph.Height

			min, max := glyph.GetTexturePosition(t.Font)
			char := &chars[i]

			char.xOffset = xOffset
			char.yOffset = yOffset
			char.w, char.h = float32(vw), float32(vh)
			char.region.X1, char.region.Y1 = min.X, min.Y
			char.region.X2, char.region.Y2 = max.X, max.Y

			// left to right shit
			xOffset += advance
			yOffset += 0
		}
	}
}

func (tc *TextComp) Fill(buf []PosTexColorVertex, p mgl32.Vec2) {
	for i, char := range tc.vertex {
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

// should have default font!!
func (tc *TextComp) SetFont(fs FontSystem) {
	tc.font = fs
}

// TextTable
type TextTable struct {
	_comps []TextComp
	_index uint32
	_map   map[int]uint32

}

func (tt *TextTable) NewComp(entity engi.Entity) (tc *TextComp) {
	tc = &tt._comps[tt._index];
	tt._map[int(entity)] = tt._index;
	tt._index ++
	return
}

func (tt *TextTable) Comp(entity engi.Entity) (tc *TextComp) {
	if v, ok := tt._map[int(entity)]; ok {
		tc = &tt._comps[v]
	}
	return
}

func (tt *TextTable) Delete(entity engi.Entity) (tc *TextComp){
	if v, ok := tt._map[int(entity)]; ok {
		tc = &tt._comps[v]
		delete(tt._map, int(entity))
		// todo swap erase

	}
	return
}

// Destroy Table
func (tt *TextTable) Destroy() {

}




