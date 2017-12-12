package gfx

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/gfx/font"

	"sort"
	"log"
)

// 文字应该采用 BatchRender 绘制
// 如果使用 BatchRender 那么此处生成模型即可

type FontSystem interface {
	Glyph(rune rune) *font.Glyph
	Tex() (uint16, *bk.Texture2D)
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

	batchId int16
	zOrder  int16

	// TextModel
	vertex []TextQuad
	runeCount int32
}

func (tc *TextComp) SetBatchId(bid int16) {
	tc.batchId = bid
}

func (tc *TextComp) SetZOrder(z int16) {
	tc.zOrder = z
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
	id, tex := tc.font.Tex()

	if id == bk.InvalidId {
		log.Println("failt to get font texture!!")
	}

	log.Println("fill data:", len(tc.text))

	for i, r := range tc.text {
		if glyph := tc.font.Glyph(r); glyph != nil {
			advance := float32(glyph.Advance)
			vw := glyph.Width
			vh := glyph.Height

			log.Println("fill glyph:", glyph)

			min := font.Point{float32(glyph.X) / tex.Width, float32(glyph.Y) / tex.Height}
			max := font.Point{float32(glyph.X+glyph.Width)/ tex.Width, float32(glyph.Y+glyph.Height) / tex.Height}

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

// should have default font!!
func (tc *TextComp) SetFont(fs FontSystem) {
	if fs != nil {
		log.Println("set font success!!")
	}
	tc.font = fs
}

// TextTable
type TextTable struct {
	_comps [10]TextComp
	_index uint32
	_map   map[int]uint32

}

func NewTextTable() *TextTable {
	return &TextTable{_map:make(map[int]uint32)}
}

func (tt *TextTable) NewComp(entity engi.Entity) (tc *TextComp) {
	tc = &tt._comps[tt._index];
	tc.Entity = entity
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


type TextRenderFeature struct {
	Stack *StackAllocator
	R *BatchRender

	tt *TextTable
	xt *TransformTable
}

// 此处初始化所有的依赖
func (trf *TextRenderFeature) Register(rs *RenderSystem) {
	// init render
	for _, r := range rs.RenderList {
		switch br := r.(type) {
		case *BatchRender:
			trf.R = br; break
		}
	}

	// init table
	for _, t := range rs.TableList {
		switch table := t.(type){
		case *TextTable:
			trf.tt = table
		case *TransformTable:
			trf.xt = table
		}
	}
	rs.Accept(trf)
}


// 此处执行渲染
// BatchRender 需要的是一组排过序的渲染对象！！！
func (trf *TextRenderFeature) Draw(filter []engi.Entity) {
	xt, tt, n := trf.xt, trf.tt, trf.tt._index
	bList := make([]textBatchObject, n)

	// get batch list
	for i := uint32(0); i < n; i++ {
		text := &tt._comps[i]
		entity := text.Entity

		xform  := xt.Comp(entity)
		bList[i] = textBatchObject{text.batchId, text, xform}
	}



	// sort
	sort.Slice(bList, func(i, j int) bool {
		return bList[i].batchId < bList[j].batchId
	})

	var batchId int16 = 0x0FFF
	var begin = false
	var render = trf.R

	// batch draw!
	for _, b := range bList{
		bid := b.batchId

		if batchId != bid {
			if begin {
				render.End()
			}
			batchId = bid
			begin = true

			id, _ := b.TextComp.font.Tex()
			render.Begin(id)
		}

		render.Draw(b)
	}

	if begin {
		render.End()
		render.flushBuffer()
	}

	render.Flush()
}

type textBatchObject struct {
	batchId int16
	*TextComp
	*Transform
}

func (tbo textBatchObject) Fill(buf []PosTexColorVertex) {
	p := tbo.Transform.Position

	for i, char := range tbo.vertex {
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

func (tbo textBatchObject) Size() int {
	return 4 * len(tbo.vertex)
}


