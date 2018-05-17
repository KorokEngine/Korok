package gfx

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/engi"
)

// Sprite is Tex2D
type Sprite Tex2D

// SpriteComp & SpriteTable
// Usually, sprite can be rendered with a BatchRenderer
type SpriteComp struct {
	engi.Entity
	Sprite
	zOrder
	batchId

	color uint32
	flipX uint16
	flipY uint16

	width float32
	height float32
	gravity struct{
		x, y float32
	}
}

func (sc *SpriteComp) SetSprite(spt Sprite) {
	sc.Sprite = spt
	sc.batchId.value = spt.Tex()

	// optional size
	if sc.width == 0 || sc.height == 0 {
		size := spt.Size()
		sc.width = size.Width
		sc.height = size.Height
	}
}

func (sc *SpriteComp) SetSize(w, h float32) {
	sc.width, sc.height = w, h
}

func (sc *SpriteComp) Size() (w, h float32) {
	w, h = sc.width, sc.height
	return
}

func (sc *SpriteComp) SetGravity(x, y float32) {
	sc.gravity.x = x
	sc.gravity.y = y
}

func (sc *SpriteComp) Color() uint32 {
	return sc.color
}

func (sc *SpriteComp) SetColor(color uint32) {
	sc.color = color
}

func (sc *SpriteComp) Flip(flipX, flipY bool) {
	if flipX {
		sc.flipX = 1
	} else {
		sc.flipX = 0
	}
	if flipY {
		sc.flipY = 1
	} else {
		sc.flipY = 0
	}
}

type SpriteTable struct {
	comps []SpriteComp
	_map   map[uint32]int
	index, cap int
}

func NewSpriteTable(cap int) *SpriteTable {
	return &SpriteTable{
		cap:cap,
		_map:make(map[uint32]int),
	}
}

func (st *SpriteTable) NewComp(entity engi.Entity) (sc *SpriteComp) {
	if size := len(st.comps); st.index >= size {
		st.comps = spriteResize(st.comps, size + STEP)
	}
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		sc = &st.comps[v]
		return
	}
	sc = &st.comps[st.index]
	sc.Entity = entity
	sc.gravity.x, sc.gravity.y = .5, .5
	sc.color = 0xFFFFFFFF
	st._map[ei] = st.index
	st.index ++
	return
}

// New SpriteComp with parameter
func (st *SpriteTable) NewCompX(entity engi.Entity, spt Tex2D) (sc *SpriteComp) {
	sc = st.NewComp(entity)
	sc.SetSprite(spt)
	return
}

func (st *SpriteTable) Alive(entity engi.Entity) bool {
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		return st.comps[v].Entity != 0
	}
	return false
}

func (st *SpriteTable) Comp(entity engi.Entity) (sc *SpriteComp) {
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		sc = &st.comps[v]
	}
	return
}

func (st *SpriteTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		if tail := st.index -1; v != tail && tail > 0 {
			st.comps[v] = st.comps[tail]
			// remap index
			tComp := &st.comps[tail]
			ei := tComp.Entity.Index()
			st._map[ei] = v
			tComp.Entity = 0
		} else {
			st.comps[tail].Entity = 0
		}

		st.index -= 1
		delete(st._map, ei)
	}
}

func (st *SpriteTable) Size() (size, cap int) {
	return st.index, st.cap
}

func (st *SpriteTable) Destroy() {
	st.comps = make([]SpriteComp, 0)
	st._map = make(map[uint32]int)
	st.index = 0
}

func spriteResize(slice []SpriteComp, size int) []SpriteComp {
	newSlice := make([]SpriteComp, size)
	copy(newSlice, slice)
	return newSlice
}

/////
type SpriteRenderFeature struct {
	Stack *StackAllocator
	id int

	R *BatchRender
	st *SpriteTable
	xt *TransformTable
}

func (f *SpriteRenderFeature) SetRender(render *BatchRender) {
	f.R = render
}

func (f *SpriteRenderFeature) SetTable(st *SpriteTable, xt *TransformTable) {
	f.st, f.xt = st, xt
}

// 此处初始化所有的依赖
func (f *SpriteRenderFeature) Register(rs *RenderSystem) {
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
		case *SpriteTable:
			f.st = table
		case *TransformTable:
			f.xt = table
		}
	}
	// add new feature, use the index as id
	f.id = rs.Accept(f)
}

func (f *SpriteRenderFeature) Extract(v *View) {
	var (
		camera = v.Camera
		xt     = f.xt
		fi = uint32(f.id) << 16
	)
	for i, spr := range f.st.comps[:f.st.index] {
		xf := xt.Comp(spr.Entity)
		sz := f32.Vec2{spr.width, spr.height}
		g  := f32.Vec2{spr.gravity.x, spr.gravity.y}

		if camera.InView(xf, sz , g) {
			sid := packSortId(spr.zOrder.value, spr.batchId.value)
			val := fi + uint32(i)
			v.RenderNodes = append(v.RenderNodes, sortObject{sid, val})
		}
	}
}

func (f *SpriteRenderFeature) Draw(nodes RenderNodes) {
	var (
		st, xt = f.st, f.xt
		sortId = uint32(0xFFFFFFFF)
		begin = false
		render = f.R
	)

	// batch draw!
	var spriteBatchObject = spriteBatchObject{}
	for _, b := range nodes {
		ii := b.value & 0xFFFF
		if sid := b.sortId; sortId != sid {
			if begin {
				render.End()
			}
			sortId = sid
			begin = true
			tex2d := st.comps[ii].Sprite.Tex()
			depth, _ := unpackSortId(b.sortId)
			render.Begin(tex2d, depth)
		}
		spriteBatchObject.SpriteComp = &st.comps[ii]
		spriteBatchObject.Transform = xt.Comp(st.comps[ii].Entity)
		render.Draw(spriteBatchObject)
	}
	if begin {
		render.End()
	}
	render.Flush()
}

type spriteBatchObject struct {
	*SpriteComp
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
//
// Transform Method:
//
// |
// |
// |
func (sbo spriteBatchObject) Fill(buf []PosTexColorVertex) {
	var (
		srt = sbo.Transform.world
		p = srt.Position
		c = sbo.SpriteComp
		w = sbo.width
		h = sbo.height
	)
	// Texture
	rg := c.Sprite.Region()
	if rg.Rotated {
		if c.flipX == 1 {
			rg.Y1, rg.Y2 = rg.Y2, rg.Y1
		}
		if c.flipY == 1 {
			rg.X1, rg.X2 = rg.X2, rg.X1
		}
		buf[1].U, buf[1].V = rg.X1, rg.Y2
		buf[2].U, buf[2].V = rg.X2, rg.Y2
		buf[3].U, buf[3].V = rg.X2, rg.Y1
		buf[0].U, buf[0].V = rg.X1, rg.Y1
	} else {
		if c.flipY == 1 {
			rg.Y1, rg.Y2 = rg.Y2, rg.Y1
		}
		if c.flipX == 1 {
			rg.X1, rg.X2 = rg.X2, rg.X1
		}
		buf[0].U, buf[0].V = rg.X1, rg.Y2
		buf[1].U, buf[1].V = rg.X2, rg.Y2
		buf[2].U, buf[2].V = rg.X2, rg.Y1
		buf[3].U, buf[3].V = rg.X1, rg.Y1
	}

	// Color
	buf[0].RGBA = c.color
	buf[1].RGBA = c.color
	buf[2].RGBA = c.color
	buf[3].RGBA = c.color

	// Center of model
	ox := w * c.gravity.x
	oy := h * c.gravity.y

	// Transform matrix
	m := f32.Mat3{}; m.Initialize(p[0], p[1], srt.Rotation, srt.Scale[0], srt.Scale[1], ox, oy, 0,0)

	// Let's go!
	buf[0].X, buf[0].Y = m.Transform(0, 0)
	buf[1].X, buf[1].Y = m.Transform(w, 0)
	buf[2].X, buf[2].Y = m.Transform(w, h)
	buf[3].X, buf[3].Y = m.Transform(0, h)

}

func (sbo spriteBatchObject) Size() int {
	return 4
}


