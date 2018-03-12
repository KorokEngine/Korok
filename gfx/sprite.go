package gfx

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx/dbg"

	"sort"
	"fmt"
)

/// SpriteComp & SpriteTable
/// Usually, sprite can be rendered with a BatchRenderer

type SpriteComp struct {
	engi.Entity
	*SubTex
	anim uint16

	scale float32
	color uint32

	width float32
	height float32

	zOrder  int16
	batchId int16
}

func (sc *SpriteComp) SetTexture(tex *SubTex) {
	sc.SubTex = tex
	if tex != nil {
		sc.batchId = int16(tex.TexId)
		sc.width = float32(tex.Width)
		sc.height = float32(tex.Width)
	}
}

func (sc *SpriteComp) SetZOrder(z int16) {
	sc.zOrder = z
}

func (sc *SpriteComp) SetBatchId(b int16) {
	sc.batchId = b
}

func (sc *SpriteComp) SetSize(w, h float32) {
	sc.width, sc.height = w, h
}

func (sc *SpriteComp) Size() (w, h float32) {
	w, h = sc.width, sc.height
	return
}

func (sc *SpriteComp) Color() uint32 {
	return sc.color
}

func (sc *SpriteComp) SetColor(color uint32) {
	sc.color = color
}

func (sc *SpriteComp) Scale() float32 {
	return sc.scale
}

func (sc *SpriteComp) SetScale(sk float32)  {
	sc.scale = sk
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
	st._map[ei] = st.index
	st.index ++
	return
}

// New SpriteComp with parameter
func (st *SpriteTable) NewCompX(entity engi.Entity, tex *SubTex) (sc *SpriteComp) {
	sc = st.NewComp(entity)
	sc.SetTexture(tex)
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

	R *BatchRender
	st *SpriteTable
	xt *TransformTable
}

func (srf *SpriteRenderFeature) SetRender(render *BatchRender) {
	srf.R = render
}

func (srf *SpriteRenderFeature) SetTable(st *SpriteTable, xt *TransformTable) {
	srf.st, srf.xt = st, xt
}

// 此处初始化所有的依赖
func (srf *SpriteRenderFeature) Register(rs *RenderSystem) {
	// init render
	for _, r := range rs.RenderList {
		switch br := r.(type) {
		case *BatchRender:
			srf.R = br; break
		}
	}
	// init table
	for _, t := range rs.TableList {
		switch table := t.(type){
		case *SpriteTable:
			srf.st = table
		case *TransformTable:
			srf.xt = table
		}
	}
	// add new feature
	rs.Accept(srf)
}


// 此处执行渲染
// BatchRender 需要的是一组排过序的渲染对象！！！
func (srf *SpriteRenderFeature) Draw(filter []engi.Entity) {
	xt, st, n := srf.xt, srf.st, srf.st.index
	bList := make([]spriteBatchObject, n)

	// get batch list
	for i := 0; i < n; i++ {
		sprite := &st.comps[i]
		entity := sprite.Entity
		xform  := xt.Comp(entity)

		// sortId =  order << 16 + batch
		sortId := uint32(int32(sprite.zOrder) + 0xFFFF>>1)
		sortId = sortId << 16
		sortId += uint32(sprite.batchId)

		bList[i] = spriteBatchObject{
			sortId,
			sprite.batchId,
			sprite,
			xform,
		}
	}

	// sort
	sort.Slice(bList, func(i, j int) bool {
		return bList[i].sortId < bList[j].sortId
	})

	var batchId int16 = 0x0FFF
	var begin = false
	var render = srf.R

	// batch draw!
	for _, b := range bList{
		bid := b.batchId

		if batchId != bid {
			if begin {
				render.End()
			}
			batchId = bid
			begin = true

			render.Begin(b.SpriteComp.SubTex.TexId)
		}

		render.Draw(b)
	}

	if begin {
		render.End()
	}

	num := render.Flush()

	dbg.Move(10, 300)
	dbg.DrawStrScaled(fmt.Sprintf("Batch num: %d", num), .6)
}

// TODO uint32 = (z-order << 16 + batch-id)
type spriteBatchObject struct {
	sortId uint32
	batchId int16
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
func (sbo spriteBatchObject) Fill(buf []PosTexColorVertex) {
	p := sbo.Transform.world.Position
	r := sbo.SpriteComp.Region
	w := sbo.width
	h := sbo.height

	buf[0].X, buf[0].Y = p[0], p[1]
	buf[0].U, buf[0].V = r.X1, r.Y2
	buf[0].RGBA = 0xffffffff

	buf[1].X, buf[1].Y = p[0] + w, p[1]
	buf[1].U, buf[1].V = r.X2, r.Y2
	buf[1].RGBA = 0xffffffff

	buf[2].X, buf[2].Y = p[0] + w, p[1] + h
	buf[2].U, buf[2].V = r.X2, r.Y1
	buf[2].RGBA = 0xffffffff

	buf[3].X, buf[3].Y = p[0], p[1] + h
	buf[3].U, buf[3].V = r.X1, r.Y1
	buf[3].RGBA = 0x00ffffff
}

func (sbo spriteBatchObject) Size() int {
	return 4
}


