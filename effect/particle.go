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
	visible int16

	tex gfx.Tex2D
	size f32.Vec2
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

func (pc *ParticleComp) Visible() bool {
	if pc.visible == 0 {
		return false
	}
	return true
}

func (pc *ParticleComp) SetVisible(v bool) {
	if v {
		pc.visible = 1
	} else {
		pc.visible = 0
	}
}

// The width and height of the particle system. We'll use it to
// make visibility test. The default value is {w:64, h:64}
func (pc *ParticleComp) SetSize(w, h float32) {
	pc.size[0], pc.size[1] = w, h
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
	ec.visible = 1
	ec.size = f32.Vec2{64, 64}
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
	*gfx.MeshRender
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
			f.MeshRender = mr; break
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

func (f *ParticleRenderFeature) Extract(v *gfx.View) {
	var (
		camera = v.Camera
		xt     = f.xt
		fi = uint32(f.id) << 16
	)
	for i, pc := range f.et.comps[:f.et.index] {
		if xf := xt.Comp(pc.Entity); pc.visible != 0 && camera.InView(xf, pc.size, f32.Vec2{.5, .5}) {
			sid := gfx.PackSortId(pc.zOrder, 0)
			v.RenderNodes = append(v.RenderNodes, gfx.SortObject{sid,fi+uint32(i)})
		}
	}
}

func (f *ParticleRenderFeature) Draw(nodes gfx.RenderNodes) {
	var (
		requireVertexSize int
		requireIndexSize int
	)
	for _, node := range nodes {
		_, cap := f.et.comps[node.Value&0xFFFF].sim.Size()
		requireVertexSize += cap * 4
		if cap > requireIndexSize {
			requireIndexSize = cap
		}
	}
	f.AllocBuffer(requireVertexSize, requireIndexSize)

	// setup mesh & matrix
	mesh := &gfx.Mesh{
		IndexId:f.BufferContext.indexId,
		VertexId:f.BufferContext.vertexId,
	}
	mat4 := &f32.Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}

	// fill vertex buffer
	var (
		offset int
	)
	for _, node := range nodes {
		z, _ := gfx.UnpackSortId(node.SortId)
		ps := f.et.comps[node.Value&0xFFFF]
		xf := f.xt.Comp(ps.Entity)

		live, _ := ps.sim.Size()
		vsz, isz := live*4, live*6
		buff := f.vertex[offset:offset+vsz]; offset += vsz
		ps.sim.Visualize(buff, ps.tex)

		mesh.FirstVertex = uint16(offset)
		mesh.NumVertex = uint16(vsz)
		mesh.FirstIndex = 0
		mesh.NumIndex = uint16(isz)
		mesh.SetTexture(ps.tex.Tex())

		p := xf.Position()
		mat4.Set(0, 3, p[0])
		mat4.Set(1, 3, p[1])
		f.MeshRender.Draw(mesh, mat4, int32(z))
	}

	// update buffer
	updateSize := uint32(requireVertexSize * 20)
	f.vb.Update(0, updateSize, unsafe.Pointer(&f.vertex[0]), false)

	dbg.Move(400, 300)
	dbg.DrawStrScaled(fmt.Sprintf("lives: %d", offset>>2), .6)
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
		id, sz, vb := gfx.Context.TempVertexBuffer(vertexSize, 20)
		ctx.vertexId = id
		ctx.vertexSize = sz
		ctx.vb = vb
		ctx.vertex = make([]gfx.PosTexColorVertex, sz)
	}

	if indexSize > ctx.indexSize {
		ctx.indexId, ctx.indexSize = gfx.Context.SharedIndexBuffer()
	}
}

func (ctx *BufferContext) Release() {

}
