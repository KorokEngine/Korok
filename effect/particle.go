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

// Particle Component
type ParticleComp struct {
	engi.Entity
	sim Simulator
	zOrder int16

	tex gfx.Tex2D
}

func (pc *ParticleComp) SetSimulator(sim Simulator) {
	pc.sim = sim
}

func (pc *ParticleComp) Simulator() Simulator {
	return pc.sim
}

func (pc *ParticleComp) SetTexture(tex gfx.Tex2D) {
	pc.tex = tex
}

func (pc *ParticleComp) Texture() gfx.Tex2D {
	return pc.tex
}

func (pc *ParticleComp) Play() {
	pc.sim.Play()
}

func (pc *ParticleComp) Stop() {
	pc.sim.Stop()
}

func (pc *ParticleComp) SetZOrder(z int16) {
	pc.zOrder = z
}

func (pc *ParticleComp) Z() int16 {
	return pc.zOrder
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
	id int

	et *ParticleSystemTable
	xt *gfx.TransformTable

	bigBuffer struct{
		vertex uint16
		index  uint16
	}
	BufferContext
}

// 此处初始化所有的依赖
func (f *ParticleRenderFeature) Register(rs *gfx.RenderSystem) {
	// init render
	for _, r := range rs.RenderList {
		switch mr := r.(type) {
		case *gfx.MeshRender:
			f.mr = mr; break
		}
	}
	// init table
	for _, t := range rs.TableList {
		switch table := t.(type){
		case *ParticleSystemTable:
			f.et = table
		case *gfx.TransformTable:
			f.xt = table
		}
	}
	// add new feature
	f.id = rs.Accept(f)
}
var mat = f32.Ident4()

func (f *ParticleRenderFeature) Extract(v *gfx.View) {

}

func (f *ParticleRenderFeature) Draw(nodes gfx.RenderNodes) {

}

// 此处执行渲染
// BatchRender 需要的是一组排过序的渲染对象！！！
func (f *ParticleRenderFeature) draw(filter []engi.Entity) {
	xt, mt, n := f.xt, f.et, f.et.index
	mr := f.mr

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
	f.AllocBuffer(requireVertexSize, requireIndexSize)

	vertex := f.BufferContext.vertex
	renderObjs := make([]renderObject, n)
	vertexOffset := int(0)

	for i := 0; i < n; i++ {
		comp := &mt.comps[i]
		entity := comp.Entity
		xform  := xt.Comp(entity)

		ro := &renderObjs[i]
		ro.Transform = xform
		ro.depth = int32(comp.zOrder)
		live, _ := comp.sim.Size()
		vn, in := live * 4, live * 6

		// write vertex
		buf := vertex[vertexOffset:vertexOffset+vn]
		comp.sim.Visualize(buf, comp.tex)

		ro.Mesh = gfx.Mesh{
			IndexId:     f.BufferContext.indexId,
			VertexId:    f.BufferContext.vertexId,
			FirstVertex: uint16(vertexOffset),
			NumVertex:   uint16(vn),
			FirstIndex:  0,
			NumIndex:    uint16(in),
		}
		ro.Mesh.SetTexture(comp.tex.Tex())
		vertexOffset += vn
	}

	updateSize := uint32(requireVertexSize * 20)

	dbg.Move(400, 300)
	dbg.DrawStrScaled(fmt.Sprintf("lives: %d", vertexOffset>>2), .6)

	f.vb.Update(0, updateSize, unsafe.Pointer(&vertex[0]), false)

	for i := range renderObjs {
		ro := &renderObjs[i]
		p := ro.Transform.Position()
		mat.Set(0, 3, p[0])
		mat.Set(1, 3, p[1])
		mr.Draw(&ro.Mesh, &mat, ro.depth)
	}
}

type renderObject struct {
	depth int32
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
		if ctx.vertexId != 0 {
			bk.R.Free(ctx.vertexId)
		}
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
		ctx.indexId, ctx.indexSize = gfx.Context.SharedIndexBuffer()
	}
}

func (ctx *BufferContext) Release() {

}
