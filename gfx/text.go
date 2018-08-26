package gfx

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx/font"
	"korok.io/korok/math/f32"
)

// 文字应该采用 BatchRender 绘制
// 如果使用 BatchRender 那么此处生成模型即可

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
	font font.Font
	zOrder
	batchId

	size float32
	color uint32

	width float32
	height float32
	gravity struct{
		x, y float32
	}
	visible bool

	// TextModel
	text  string
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

func (tc *TextComp) SetGravity(x, y float32) {
	tc.gravity.x = x
	tc.gravity.y = y
}

func (tc *TextComp) SetVisible(v bool) {
	tc.visible = v
}

func (tc *TextComp) SetFontSize(sz float32) {
	tc.size = sz
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
	chars := make([]TextQuad, len(tc.text))
	tc.vertex = chars
	_, gh := tc.font.Bounds()

	var (
		size = struct { w, h float32}{}
		scale = tc.size/gh
		xOffset float32
		yOffset float32
	)

	if tc.size == 0 {
		scale = 1.0
	} else {
		gh = tc.size
	}
	yOffset = gh

	for i, r := range tc.text {
		if glyph, ok := tc.font.Glyph(r); ok {
			advance := float32(glyph.Advance) * scale
			vw := glyph.Width * scale
			vh := glyph.Height * scale
			ox := glyph.XOffset * scale
			oy := glyph.YOffset * scale

			u1, v1, u2, v2 := tc.font.Frame(r)
			char := &chars[i]

			char.xOffset = xOffset + ox
			char.yOffset = yOffset - (oy+vh)
			char.w, char.h = vw, vh
			char.region.X1, char.region.Y1 = u1, v1
			char.region.X2, char.region.Y2 = u2, v2

			// left to right shift only!
			xOffset += advance
			yOffset += 0
		}
	}
	size.w = xOffset + chars[len(chars)-1].w
	tc.width = size.w
	tc.height = gh
}

// should have default font!!
func (tc *TextComp) SetFont(fnt font.Font) {
	tc.font = fnt
	if fnt != nil && tc.batchId.value != 0 {
		tex, _ := fnt.Tex2D()
		tc.batchId.value = tex
	}
}

// TextTable
type TextTable struct {
	comps []TextComp
	_map   map[uint32]int
	index, cap int

}

func NewTextTable(cap int) *TextTable {
	return &TextTable{cap: cap,_map: make(map[uint32]int)}
}

func (tt *TextTable) NewComp(entity engi.Entity) (tc *TextComp) {
	if size := len(tt.comps); tt.index >= size {
		tt.comps = textResize(tt.comps, size + STEP)
	}
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		return &tt.comps[v]
	}

	tc = &tt.comps[tt.index]
	tc.Entity = entity
	tc.color = 0xFFFFFFFF
	tc.gravity.x = .5
	tc.gravity.y = .5
	tc.visible = true
	tt._map[ei] = tt.index;
	tt.index ++
	return
}

func (tt *TextTable) Alive(entity engi.Entity) bool {
	if v, ok := tt._map[entity.Index()]; ok {
		return tt.comps[v].Entity == 0
	}
	return false
}

func (tt *TextTable) Comp(entity engi.Entity) (tc *TextComp) {
	if v, ok := tt._map[entity.Index()]; ok {
		tc = &tt.comps[v]
	}
	return
}

func (tt *TextTable) Delete(entity engi.Entity) (tc *TextComp) {
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		if tail := tt.index -1; v != tail && tail > 0 {
			tt.comps[v] = tt.comps[tail]
			// remap index
			tComp := &tt.comps[tail]
			ei := tComp.Entity.Index()
			tt._map[ei] = v
			tComp.Entity = 0
		} else {
			tt.comps[tail].Entity = 0
		}

		tt.index -= 1
		delete(tt._map, ei)
	}
	return
}

// Destroy Table
func (tt *TextTable) Destroy() {
	tt.comps = make([]TextComp, 0)
	tt._map = make(map[uint32]int)
	tt.index = 0
}

func (tt *TextTable) Size() (size, cap int) {
	return tt.index, tt.cap
}

func textResize(slice []TextComp, size int) []TextComp {
	newSlice := make([]TextComp, size)
	copy(newSlice, slice)
	return newSlice
}


type TextRenderFeature struct {
	Stack *StackAllocator
	R *BatchRender
	id int

	tt *TextTable
	xt *TransformTable
}

// 此处初始化所有的依赖
func (f *TextRenderFeature) Register(rs *RenderSystem) {
	// init render
	for _, r := range rs.RenderList {
		switch br := r.(type) {
		case *BatchRender:
			f.R = br; break
		}
	}

	// init table
	for _, t := range rs.TableList {
		switch table := t.(type){
		case *TextTable:
			f.tt = table
		case *TransformTable:
			f.xt = table
		}
	}
	f.id = rs.Accept(f)
}

func (f *TextRenderFeature) Extract(v *View) {
	var (
		camera = v.Camera
		xt     = f.xt
		fi     = uint32(f.id)<<16
	)

	for i, spr := range f.tt.comps[:f.tt.index] {
		xf := xt.Comp(spr.Entity)
		sz := f32.Vec2{spr.width, spr.height}
		g  := f32.Vec2{spr.gravity.x, spr.gravity.y}

		if spr.visible && camera.InView(xf, sz, g) {
			sid := PackSortId(spr.zOrder.value, spr.batchId.value)
			val := fi + uint32(i)
			v.RenderNodes = append(v.RenderNodes, SortObject{sid, val})
		}
	}
}

func (f *TextRenderFeature) Draw(nodes RenderNodes) {
	var (
		tt, xt = f.tt, f.xt
		sortId  = uint32(0xFFFFFFFF)
		begin = false
		render = f.R
	)

	// batch draw!
	var textBatchObject = textBatchObject{}
	for _, b := range nodes {
		ii := b.Value & 0xFFFF
		if sid := b.SortId & 0xFFFF; sortId != sid {
			if begin {
				render.End()
			}
			sortId = sid
			begin = true
			tex2d, _ := tt.comps[ii].font.Tex2D()
			depth, _ := UnpackSortId(b.SortId)
			render.Begin(tex2d, depth)
		}
		textBatchObject.TextComp = &tt.comps[ii]
		textBatchObject.Transform = xt.Comp(tt.comps[ii].Entity)
		render.Draw(textBatchObject)
	}
	if begin {
		render.End()
	}
	render.Flush()
}

func (f *TextRenderFeature) Flush() {

}

type textBatchObject struct {
	*TextComp
	*Transform
}

// batch system winding order
//
//		3----------2
//		| . 	   |
//      |   .	   |
//		|     .    |
// 		|		.  |
// 		0----------1
// 1 * 1 quad for each char
// order: 3 0 1 3 1 2
func (tbo textBatchObject) Fill(buf []PosTexColorVertex) {
	srt :=  &tbo.Transform.world
	p := tbo.Transform.world.Position
	t := tbo.TextComp

	// Center of model
	ox := t.width * t.gravity.x
	oy := t.height * t.gravity.y

	// Transform matrix
	m := f32.Mat3{}; m.Initialize(p[0], p[1], srt.Rotation, srt.Scale[0], srt.Scale[1], ox, oy, 0,0)


	for i, char := range tbo.vertex {
		vi := i * 4

		// index (0, 0) <x,y,u,v>
		v := &buf[vi+0]
		v.X, v.Y = m.Transform(char.xOffset, char.yOffset)
		v.U, v.V = char.region.X1, char.region.Y2
		v.RGBA = t.color

		// index (1,0) <x,y,u,v>
		v = &buf[vi+1]
		v.X, v.Y = m.Transform(char.xOffset + char.w, char.yOffset)
		v.U, v.V = char.region.X2, char.region.Y2
		v.RGBA = t.color

		// index(1,1) <x,y,u,v>
		v = &buf[vi+2]
		v.X, v.Y = m.Transform(char.xOffset + char.w, char.yOffset + char.h)
		v.U, v.V = char.region.X2, char.region.Y1
		v.RGBA = t.color

		// index(0, 1) <x,y,u,v>
		v = &buf[vi+3]
		v.X, v.Y = m.Transform(char.xOffset, char.yOffset + char.h)
		v.U, v.V = char.region.X1, char.region.Y1
		v.RGBA = t.color
	}
}

func (tbo textBatchObject) Size() int {
	return 4 * len(tbo.vertex)
}


