package gfx

import (
	"korok.io/korok/engi"

	"sort"
)

/// SpriteComp & SpriteTable
/// Usually, sprite can be rendered with a BatchRenderer

type SpriteComp struct {
	engi.Entity
	*SubTex

	Scale float32
	Color uint32
	Width float32
	Height float32

	zOrder  int16
	batchId int16
}

func (sc *SpriteComp) SetTexture(tex *SubTex) {
	sc.SubTex = tex
	if tex != nil {
		sc.batchId = int16(tex.TexId)
		sc.Width = float32(tex.Width)
		sc.Height = float32(tex.Width)
	}
}

func (sc *SpriteComp) SetZOrder(z int16) {
	sc.zOrder = z
}

func (sc *SpriteComp) SetBatchId(b int16) {
	sc.batchId = b
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

func (st *SpriteTable) NewComp(entity engi.Entity, tex *SubTex) (sc *SpriteComp) {
	if size := len(st.comps); st.index >= size {
		st.comps = spriteResize(st.comps, size + STEP)
	}
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		sc = &st.comps[v]
		return
	}
	sc = &st.comps[st.index]
	sc.SubTex = tex
	sc.Entity = entity
	if tex != nil {
		sc.batchId = int16(tex.TexId)
		sc.Width = float32(tex.Width)
		sc.Height = float32(tex.Width)
	}
	st._map[ei] = st.index
	st.index ++
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
		bList[i] = spriteBatchObject{
			sprite.batchId,
			sprite,
			xform,
		}
	}

	// sort
	sort.Slice(bList, func(i, j int) bool {
		return bList[i].batchId < bList[j].batchId
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

	render.Flush()
}

type spriteBatchObject struct {
	batchId int16
	*SpriteComp
	*Transform
}

func (sbo spriteBatchObject) Fill(buf []PosTexColorVertex) {
	p := sbo.Transform.Position
	r := sbo.SpriteComp.Region
	w := sbo.Width
	h := sbo.Height

	buf[0].X, buf[0].Y = p[0], p[1]
	buf[0].U, buf[0].V = r.X1, r.Y1
	buf[0].RGBA = 0xffffffff

	buf[1].X, buf[1].Y = p[0] + w, p[1]
	buf[1].U, buf[1].V = r.X2, r.Y1
	buf[1].RGBA = 0xffffffff

	buf[2].X, buf[2].Y = p[0] + w, p[1] + h
	buf[2].U, buf[2].V = r.X2, r.Y2
	buf[2].RGBA = 0xffffffff

	buf[3].X, buf[3].Y = p[0], p[1] + h
	buf[3].U, buf[3].V = r.X1, r.Y2
	buf[3].RGBA = 0x00ffffff
}

func (sbo spriteBatchObject) Size() int {
	return 4
}


