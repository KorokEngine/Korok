package gfx

import (
	"korok/engi"
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
}

func (sc *SpriteComp) SetZOrder(z int16) {
	sc.zOrder = z
}

func (sc *SpriteComp) SetBatchId(b int16) {
	sc.batchId = b
}

type SpriteTable struct {
	_comps [1024]SpriteComp
	_index uint32
	_map   [1024]uint32
}

func (st *SpriteTable) NewComp(entity engi.Entity, tex *SubTex) (sc *SpriteComp) {
	sc = &st._comps[st._index]
	sc.SubTex = tex
	sc.Entity = entity
	st._map[entity] = st._index
	st._index ++
	return
}

func (st *SpriteTable) Comp(id uint32) *SpriteComp {
	return &st._comps[st._map[id]]
}

/////
type SpriteRenderFeature struct {
	Stack *StackAllocator

	R *BatchRender
	st *SpriteTable
	xt *TransformTable
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
	xt, st, n := srf.xt, srf.st, srf.st._index
	bList := make([]spriteBatchObject, n)

	// get batch list
	for i := uint32(0); i < n; i++ {
		sprite := &st._comps[i]
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

	var batchId int16 = 0xFFFF
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
		render.flushBuffer()
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
	buf[0].RGBA = 0x00ffffff

	buf[1].X, buf[1].Y = p[0] + w, p[1]
	buf[1].U, buf[1].V = r.X2, r.Y1
	buf[1].RGBA = 0x00ffffff

	buf[2].X, buf[2].Y = p[0] + w, p[1] + h
	buf[2].U, buf[2].V = r.X2, r.Y2
	buf[2].RGBA = 0x00ffffff

	buf[3].X, buf[3].Y = p[0], p[1] + h
	buf[3].U, buf[3].V = r.X1, r.Y2
	buf[3].RGBA = 0x00ffffff
}

func (sbo spriteBatchObject) Size() int {
	return 4
}


