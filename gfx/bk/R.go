package bk

import "log"

/// 在 2D 引擎中，GPU 资源的使用时很有限的，
/// 大部分图元会在 Batch 环节得到优化，最终生成有限的渲染指令.

const InvalidId uint16 = 0xFFFF
const UINT16_MAX uint16 = 0xFFFF

type Memory struct {
	Data interface {}
	Size uint32
}

// ID FORMAT
// 0x F   FFF
//    ^    ^
//    |    +---- id value
//    +--------- id type

const (
	ID_TYPE_INDEX    uint16 = iota
	ID_TYPE_VERTEX
	ID_TYPE_TEXTURE
	ID_TYPE_LAYOUT
	ID_TYPE_UNIFORM
	ID_TYPE_SHADER
)

const (
	MAX_INDEX   = 2 << 10
	MAX_VERTEX  = 2 << 10
	MAX_TEXTURE = 1 << 10
	MAX_LAYOUT  = 32
	MAX_UNIFORM = 64
	MAX_SHADER  = 32
)

// predefined vertex format
const (
	VERTEX_POS2_TEX_COLOR uint16 = (ID_TYPE_LAYOUT << 12) & 0
	VERTEX_POS2_TEX       uint16 = (ID_TYPE_LAYOUT << 12) & 1
	VERTEX_POS2_COLOR     uint16 = (ID_TYPE_LAYOUT << 12) & 2

	//VERTEX_POS3_TEX_COLOR uint16 = (ID_TYPE_LAYOUT << 12) & 3
	//VERTEX_POS3_TEX       uint16 = (ID_TYPE_LAYOUT << 12) & 4
	//VERTEX_POS3_COLOR     uint16 = (ID_TYPE_LAYOUT << 12) & 5
)

type ResManager struct {
	indexBuffers 	[MAX_INDEX]IndexBuffer
	vertexBuffers 	[MAX_VERTEX]VertexBuffer
	textures 		[MAX_TEXTURE]Texture2D

	vertexLayouts  	[MAX_LAYOUT]VertexLayout
	uniforms 		[MAX_UNIFORM]Uniform

	shaders 		[MAX_SHADER]Shader

	ibIndex uint16
	vbIndex uint16
	ttIndex uint16
	vlIndex uint16
	umIndex uint16
	shIndex uint16
}

func NewResManager() *ResManager{
	return &ResManager{}
}

func (rm *ResManager) Init() {
	rm.setupPredefine()
}

func (rm *ResManager) Destroy() {

}

/// Create Method
func (rm *ResManager) AllocIndexBuffer(mem Memory) (id uint16, ib *IndexBuffer) {
	id, ib = rm.ibIndex, &rm.indexBuffers[rm.ibIndex]
	rm.ibIndex ++
	id = id & (ID_TYPE_INDEX << 12)
	ib.create(mem.Size, mem.Data, 0)
	return
}

func (rm *ResManager) AllocVertexBuffer(mem Memory, lytId uint16) (id uint16, vb *VertexBuffer) {
	id, vb = rm.vbIndex, &rm.vertexBuffers[rm.vbIndex]
	rm.vbIndex ++
	id = id & (ID_TYPE_VERTEX << 12)
	vb.create(mem.Size, mem.Data, lytId, 0)
	return
}

func (rm *ResManager) AllocVertexLayout(vLayout *VertexLayout) (id uint16, vl *VertexLayout) {
	id, vl = rm.vlIndex, &rm.vertexLayouts[rm.vlIndex]
	rm.vlIndex ++
	id = id & (ID_TYPE_LAYOUT << 12)
	// copy data
	*vl = *vLayout
	return
}

func (rm *ResManager) AllocUniform(name string, xType UniformType, num uint32) (id uint16, um *Uniform) {
	id, um = rm.umIndex, &rm.uniforms[rm.umIndex]
	rm.umIndex ++
	id = id & (ID_TYPE_UNIFORM << 12)

	um.Name = name
	um.Type = xType
	um.Count = num
	um.Size = num * 0 // TODO num * sizeOf(UniformType)
	return
}

func (rm *ResManager) AllocShader(vsh, fsh string) (id uint16, sh *Shader) {
	id, sh = rm.shIndex, &rm.shaders[rm.shIndex]
	rm.shIndex ++
	id = id & (ID_TYPE_SHADER << 12)
	if program, err := Compile(vsh, fsh); err == nil {
		sh.Program = program
	} else {
		log.Println("Failed to alloc shader..")
	}
	return
}

/// Destroy Method
func (rm *ResManager) Free(id uint16) {
	t := (id >> 12) & 0x000F
	v := id & 0x0FFF

	switch t {
	case ID_TYPE_INDEX:
		rm.indexBuffers[v].destroy()
		rm.ibIndex --
	case ID_TYPE_VERTEX:
		rm.vertexBuffers[v].destroy()
		rm.vbIndex --
	case ID_TYPE_TEXTURE:
		rm.textures[v].Destroy()
		rm.ttIndex --
	case ID_TYPE_LAYOUT:
		rm.vlIndex --
	case ID_TYPE_UNIFORM:
		rm.umIndex --
	case ID_TYPE_SHADER:
		rm.shaders[v].Destroy()
		rm.shIndex --
	}
}

/// Retrieve Method
func (rm *ResManager) IndexBuffer(id uint16) (ok bool, ib *IndexBuffer) {
	t, v := id >> 12, id & 0x0FFF
	if t != ID_TYPE_INDEX || v >= MAX_INDEX {
		return false, nil
	}
	return true, &rm.indexBuffers[v]
}

func (rm *ResManager) VertexBuffer(id uint16) (ok bool, vb *VertexBuffer) {
	t, v := id >> 12, id & 0x0FFF
	if t != ID_TYPE_VERTEX || v >= MAX_VERTEX {
		return false, nil
	}
	return true, &rm.vertexBuffers[v]
}

func (rm *ResManager) VertexLayout(id uint16) (ok bool, vb *VertexLayout) {
	t, v := id >> 12, id & 0x0FFF
	if t != ID_TYPE_LAYOUT || v >= MAX_LAYOUT {
		return false, nil
	}
	return true, &rm.vertexLayouts[v]
}

func (rm *ResManager) Uniform(id uint16) (ok bool, um *Uniform) {
	t, v := id >> 12, id & 0x0FFF
	if t != ID_TYPE_UNIFORM || v >= MAX_UNIFORM {
		return false, nil
	}
	return true, &rm.uniforms[v]
}

func (rm *ResManager) Shader(id uint16) (ok bool, sh *Shader) {
	t, v := id >> 12, id & 0x0FFF
	if t != ID_TYPE_SHADER || v >= MAX_SHADER {
		return false, nil
	}
	return true, &rm.shaders[v]
}

func (rm *ResManager) setupPredefine() {
	// predefined vertex layout
	rm.vlIndex = 6
	p2_tex_color := &rm.vertexLayouts[VERTEX_POS2_TEX_COLOR & 0x0F]
	p2_tex_color.Begin().
				Add(2, ATTRIB_TYPE_FLOAT, false, false).
				Add(2, ATTRIB_TYPE_FLOAT, false, false).
				Add(4, ATTRIB_TYPE_UINT8, true, false).End()

	p2_tex       := &rm.vertexLayouts[VERTEX_POS2_TEX & 0x0F]
	p2_tex.Begin().
				Add(2, ATTRIB_TYPE_FLOAT, false, false).
				Add(2, ATTRIB_TYPE_FLOAT, false, false).End()

	p2_color     := &rm.vertexLayouts[VERTEX_POS2_COLOR & 0x0F]
	p2_color.Begin().
				Add(2, ATTRIB_TYPE_FLOAT, false, false).
				Add(4, ATTRIB_TYPE_UINT8, true, false).End()
}

////// MAX SIZE
var MAX = struct {

}{}


////// STATE MASK AND VALUE DEFINES

var ST = struct {
	DEPTH_WRITE 	 uint64
	DEPTH_TEST_MASK  uint64
	DEPTH_TEST_SHIFT uint64

	ALPHA_WRITE uint64
	RGB_WRITE 	uint64

	BLEND_MASK 	uint64
	BLEND_SHIFT uint64

	BLEND_EQUATION_MASK  uint64
	BLEND_EQUATION_SHIFT uint64

	PT_MASK	 uint64
	PT_SHIFT uint64

}{
	DEPTH_WRITE: 0,
	DEPTH_TEST_MASK: 0,
	DEPTH_TEST_SHIFT: 0,

	ALPHA_WRITE:0,
	RGB_WRITE:0,

	BLEND_MASK:0,
	BLEND_SHIFT:0,

	BLEND_EQUATION_MASK: 0,
	BLEND_EQUATION_SHIFT: 0,

	PT_MASK	:0,
	PT_SHIFT:0,
}


var R_CmpFunc = []uint32 {

}

var R_BlendFactor = []uint32 {

}

var R_BlendEquation = []uint32 {

}

var R_PrimInfo = []uint32 {

}