package gfx

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/engi"

	"unsafe"
)

/// MeshComp and MeshTable

type Mesh struct {
	// vertex data <x,y,u,v>
	vertex []PosTexColorVertex
	index  []uint16

	// res handle
	textureId uint16
	padding   uint16

	IndexId   uint16
	VertexId  uint16

	FirstVertex uint16
	NumVertex   uint16

	FirstIndex uint16
	NumIndex   uint16
}

type MeshComp struct {
	engi.Entity
	Mesh
}

func (*Mesh) Type() int32{
	return 0
}

func (m *Mesh) Setup() {
	mem_v := bk.Memory{unsafe.Pointer(&m.vertex[0]), uint32(len(m.vertex)) * 20 }
	if id, _:= bk.R.AllocVertexBuffer(mem_v, 20); id != bk.InvalidId {
		m.VertexId = id
	}

	mem_i := bk.Memory{unsafe.Pointer(&m.index[0]), uint32(len(m.index)) * 2}
	if id, _:= bk.R.AllocIndexBuffer(mem_i); id != bk.InvalidId {
		m.IndexId = id
	}

	m.FirstVertex = 0
	m.NumVertex = uint16(len(m.vertex))
	m.FirstIndex = 0
	m.NumIndex = uint16(len(m.index))
}

func (m*Mesh) SetTexture(id uint16) {
	m.textureId = id
}

func (m*Mesh) SetVertex(v []PosTexColorVertex) {
	m.vertex = v
}

func (m*Mesh) SetIndex(v []uint16) {
	m.index = v
}


func (m *Mesh) Update() {
	if ok, ib := bk.R.IndexBuffer(m.IndexId); ok {
		ib.Update(0, uint32(len(m.index)) * uint32(UInt16Size), unsafe.Pointer(&m.index[0]), false)
	}

	if ok, vb := bk.R.VertexBuffer(m.VertexId); ok {
		vb.Update(0, uint32(len(m.vertex)) * uint32(PosTexColorVertexSize), unsafe.Pointer(&m.vertex[0]), false)
	}
}

func (m*Mesh) Delete() {
	if ok, ib := bk.R.IndexBuffer(m.IndexId); ok {
		ib.Destroy()
	}
	if ok, vb := bk.R.VertexBuffer(m.VertexId); ok {
		vb.Destroy()
	}
	if ok, tex := bk.R.Texture(m.textureId); ok {
		tex.Destroy()
	}
}

// Configure VAO/VBO TODO
// 如果每个Sprite都创建一个VBO还是挺浪费的，
// 但是如果不创建新的VBO，那么怎么处理纹理坐标呢？
// 2D 场景中会出现大量模型和纹理相同的物体，仅仅
// 位置不同，比如满屏的子弹
// 或许可以通过工厂来构建mesh，这样自动把重复的mesh丢弃
// mesh 数量最终 <= 精灵的数量
var vertices = []float32{
	// Pos      // Tex
	0.0, 1.0, 0.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 0.0,

	0.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
}


type MeshTable struct {
	comps []MeshComp
	_map  map[uint32]int
	index, cap int
}

func NewMeshTable(cap int) *MeshTable {
	return &MeshTable{cap:cap, _map:make(map[uint32]int)}
}

func (mt *MeshTable) NewComp(entity engi.Entity) (mc *MeshComp) {
	if size := len(mt.comps); mt.index >= size {
		mt.comps = meshResize(mt.comps, size + STEP)
	}
	ei := entity.Index()
	if v, ok := mt._map[ei]; ok {
		return &mt.comps[v]
	}

	mc = &mt.comps[mt.index]
	mc.Entity = entity
	mt._map[ei] = mt.index
	mt.index ++
	return
}

func (mt *MeshTable) Alive(entity engi.Entity) bool {
	if v, ok := mt._map[entity.Index()]; ok {
		return mt.comps[v].Entity == 0
	}
	return false
}

func (mt *MeshTable) Comp(entity engi.Entity) (mc *MeshComp) {
	if v, ok := mt._map[entity.Index()]; ok {
		mc = &mt.comps[v]
	}
	return
}

func (mt *MeshTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := mt._map[ei]; ok {
		if tail := mt.index -1; v != tail && tail > 0 {
			mt.comps[v] = mt.comps[tail]
			// remap index
			tComp := &mt.comps[tail]
			ei := tComp.Entity.Index()
			mt._map[ei] = v
			tComp.Entity = 0
		} else {
			mt.comps[tail].Entity = 0
		}

		mt.index -= 1
		delete(mt._map, ei)
	}
}

func (mt *MeshTable) Destroy() {
	mt.comps = make([]MeshComp, 0)
	mt._map = make(map[uint32]int)
	mt.index = 0
}

func (mt *MeshTable) Size() (size, cap int) {
	return mt.index, mt.cap
}

func meshResize(slice []MeshComp, size int) []MeshComp {
	newSlice := make([]MeshComp, size)
	copy(newSlice, slice)
	return newSlice
}

/////
type MeshRenderFeature struct {
	Stack *StackAllocator

	R *MeshRender
	mt *MeshTable
	xt *TransformTable
}

// 此处初始化所有的依赖
func (srf *MeshRenderFeature) Register(rs *RenderSystem) {
	// init render
	for _, r := range rs.RenderList {
		switch br := r.(type) {
		case *MeshRender:
			srf.R = br; break
		}
	}
	// init table
	for _, t := range rs.TableList {
		switch table := t.(type){
		case *MeshTable:
			srf.mt = table
		case *TransformTable:
			srf.xt = table
		}
	}
	// add new feature
	rs.Accept(srf)
}

var mat4 = f32.Translate3D(300, 100, 0)
// 此处执行渲染
// BatchRender 需要的是一组排过序的渲染对象！！！
func (srf *MeshRenderFeature) Draw(filter []engi.Entity) {
	xt, mt, n := srf.xt, srf.mt, srf.mt.index
	mr := srf.R

	for i := 0; i < n; i++ {
		mesh := &mt.comps[i]
		entity := mesh.Entity
		xform  := xt.Comp(entity)

		// TODO transform!!
		mat4.Set(0, 3, xform.world.Position[0])
		mat4.Set(1, 3, xform.world.Position[1])

		mr.Draw(&mesh.Mesh, &mat4)
	}
}




