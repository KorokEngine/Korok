package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/engi"
	"korok.io/korok/gfx"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/gfx/dbg"

	"unsafe"
	"fmt"
)

// particle component
type ParticleComp struct {
	engi.Entity
	sim Simulator

	tex *gfx.SubTex
	color uint32
	scale float32
}

func (ec *ParticleComp) SetSimulator(sim Simulator) {
	ec.sim = sim
}

func (ec *ParticleComp) SetTexture(tex *gfx.SubTex) {
	ec.tex = tex
}

func (ec *ParticleComp) SetColor(color uint32) {
	ec.color = color
}

func (ec *ParticleComp) SetScale(scale float32) {
	ec.scale = scale
}

func (ec *ParticleComp) Play() {

}

func (ec *ParticleComp) Pause() {

}

func (ec *ParticleComp) Resume() {

}

func (ec *ParticleComp) Stop() {

}

// component manager
type ParticleSystemTable struct {
	comps []ParticleComp
	_map   map[uint32]int
	index, cap int
}

func NewParticleSystemTable(cap int) *ParticleSystemTable {
	return &ParticleSystemTable{
		_map:make(map[uint32]int),
		cap:cap,
	}
}

func (et *ParticleSystemTable) NewComp(entity engi.Entity) (ec *ParticleComp) {
	if size := len(et.comps); et.index >= size {
		et.comps = effectCompResize(et.comps, size + 64)
	}
	ei := entity.Index()
	if v, ok := et._map[ei]; ok {
		return &et.comps[v]
	}
	ec = &et.comps[et.index]
	ec.Entity = entity
	et._map[ei] = et.index
	et.index ++
	return
}

func (et *ParticleSystemTable) Alive(entity engi.Entity) bool {
	if v, ok := et._map[entity.Index()]; ok {
		return et.comps[v].Entity != 0
	}
	return false
}

func (et *ParticleSystemTable) Comp(entity engi.Entity) (ec *ParticleComp) {
	if v, ok := et._map[entity.Index()]; ok {
		ec = &et.comps[v]
	}
	return
}

func (et *ParticleSystemTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := et._map[ei]; ok {
		if tail := et.index -1; v != tail && tail > 0 {
			et.comps[v] = et.comps[tail]
			// remap index
			tComp := &et.comps[tail]
			ei := tComp.Entity.Index()
			et._map[ei] = v
			tComp.Entity = 0
		} else {
			et.comps[tail].Entity = 0
		}

		et.index -= 1
		delete(et._map, ei)
	}
}

func (et *ParticleSystemTable) Size() (size, cap int) {
	return et.index, et.cap
}

func effectCompResize(slice []ParticleComp, size int) []ParticleComp {
	newSlice := make([]ParticleComp, size)
	copy(newSlice, slice)
	return newSlice
}

type ParticleRenderFeature struct {
	mr *gfx.MeshRender

	et *ParticleSystemTable
	xt *gfx.TransformTable

	bigBuffer struct{
		vertex uint16
		index  uint16
	}
	BufferContext
}

// 此处初始化所有的依赖
func (prf *ParticleRenderFeature) Register(rs *gfx.RenderSystem) {
	// init render
	for _, r := range rs.RenderList {
		switch mr := r.(type) {
		case *gfx.MeshRender:
			prf.mr = mr; break
		}
	}
	// init table
	for _, t := range rs.TableList {
		switch table := t.(type){
		case *ParticleSystemTable:
			prf.et = table
		case *gfx.TransformTable:
			prf.xt = table
		}
	}
	// add new feature
	rs.Accept(prf)
}
var mat = f32.Ident4()

// 此处执行渲染
// BatchRender 需要的是一组排过序的渲染对象！！！
func (prf *ParticleRenderFeature) Draw(filter []engi.Entity) {
	xt, mt, n := prf.xt, prf.et, prf.et.index
	mr := prf.mr

	if n == 0 {
		return
	}

	var requireVertexSize int
	var requireIndexSize int
	for i := 0; i < n; i++ {
		_, cap := mt.comps[i].sim.Size()
		requireVertexSize += cap * 4
		if cap > requireIndexSize {
			requireIndexSize = cap
		}
	}
	prf.AllocBuffer(requireVertexSize, requireIndexSize)

	vertex := prf.BufferContext.vertex
	renderObjs := make([]renderObject, n)
	vertexOffset := int(0)

	for i := 0; i < n; i++ {
		comp := &mt.comps[i]
		entity := comp.Entity
		xform  := xt.Comp(entity)

		ro := &renderObjs[i]
		ro.Transform = xform
		live, _ := comp.sim.Size()
		vn, in := live * 4, live * 6

		// write vertex
		buf := vertex[vertexOffset:vertexOffset+vn]
		comp.sim.Visualize(buf)

		ro.Mesh = gfx.Mesh{
			TextureId:comp.tex.TexId,
			IndexId:prf.BufferContext.indexId,
			VertexId:prf.BufferContext.vertexId,
			FirstVertex:uint16(vertexOffset),
			NumVertex:uint16(vn),
			FirstIndex:0,
			NumIndex:uint16(in),
		}
		vertexOffset += vn
	}

	updateSize := uint32(requireVertexSize * 20)

	dbg.Move(400, 300)
	dbg.DrawStrScaled(fmt.Sprintf("lives: %d", vertexOffset>>2), .6)

	prf.vb.Update(0, updateSize, unsafe.Pointer(&vertex[0]), false)

	for i := range renderObjs {
		ro := &renderObjs[i]
		p := ro.Transform.Position()
		mat.Set(0, 3, p[0])
		mat.Set(1, 3, p[1])
		mr.Draw(&ro.Mesh, &mat)
	}
}

type renderObject struct {
	gfx.Mesh
	*gfx.Transform
}

// 目前所有的粒子都会使用一个VBO进行渲染 TODO
// 这么做同时可渲染的例子数量会受限于VBO的大小，需要一些经验数据支持
type BufferContext struct {
	vertexId uint16
	indexId uint16

	vertexSize int
	indexSize int

	// 目前我们使用 MeshRender 来渲染粒子
	// 所以必须支持如下的数据结构
	vertex []gfx.PosTexColorVertex

	vb *bk.VertexBuffer
}

func (ctx *BufferContext) AllocBuffer(vertexSize, indexSize int) {
	if vertexSize > ctx.vertexSize {
		bk.R.Free(ctx.vertexId)
		{
			vertexSize--
			vertexSize |= vertexSize >> 1
			vertexSize |= vertexSize >> 2
			vertexSize |= vertexSize >> 3
			vertexSize |= vertexSize >> 8
			vertexSize |= vertexSize >> 16
			vertexSize++
		}

		ctx.vertexSize = vertexSize
		ctx.vertex = make([]gfx.PosTexColorVertex, vertexSize)
		if id, vb := bk.R.AllocVertexBuffer(bk.Memory{nil,uint32(vertexSize) * 20}, 20); id != bk.InvalidId {
			ctx.vertexId = id
			ctx.vb = vb
		}
	}

	if indexSize > ctx.indexSize {
		bk.R.Free(ctx.indexId)
		ctx.indexId, ctx.indexSize = gfx.Context.SharedIndexBuffer()
	}
}

func (ctx *BufferContext) Release() {

}
