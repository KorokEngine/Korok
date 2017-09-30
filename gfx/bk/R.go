package bk

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

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
	sh.create(vsh, fsh)
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
	RGB_WRITE 	uint64
	ALPHA_WRITE uint64
	DEPTH_WRITE uint64

	DEPTH_TEST_MASK  uint64
	DEPTH_TEST_SHIFT uint64

	BLEND_MASK 	uint64
	BLEND_SHIFT uint64

	PT_MASK	 uint64
	PT_SHIFT uint64

}{
	RGB_WRITE:		  0x0000000000000001,
	ALPHA_WRITE:      0x0000000000000002,
	DEPTH_WRITE:      0x0000000000000004,

	DEPTH_TEST_MASK:  0x00000000000000F0,
	DEPTH_TEST_SHIFT: 4,

	BLEND_MASK:		  0x0000000000000F00,
	BLEND_SHIFT: 	  8,

	PT_MASK	:		  0x000000000000F000,
	PT_SHIFT:		  12,
}

var ST_DEPTH = struct {
	LESS 	 uint64
	LEQUAL   uint64
	EQUAL    uint64
	GEQUAL   uint64
	GREATER  uint64
	NOTEQUAL uint64
	NEVER    uint64
	ALWAYS   uint64
}{
	LESS: 	  0x0000000000000010,
	LEQUAL:   0x0000000000000020,
	EQUAL: 	  0x0000000000000030,
	GEQUAL:   0x0000000000000040,
	GREATER:  0x0000000000000050,
	NOTEQUAL: 0x0000000000000060,
	NEVER:    0x0000000000000070,
	ALWAYS:   0x0000000000000080,
}

var g_CmpFunc = []uint32 {
	0, // ignored
	gl.LESS,
	gl.LEQUAL,
	gl.EQUAL,
	gl.GEQUAL,
	gl.GREATER,
	gl.NOTEQUAL,
	gl.NEVER,
	gl.ALWAYS,
}

var ST_BLEND = struct {
	DEFAULT 				uint64
	ISABLE 				 	uint64
	ALPHA_PREMULTIPLIED 	uint64
	ALPHA_NON_PREMULTIPLIED uint64
	ADDITIVE 				uint64
}{
	ISABLE:					0x0000000000000100,
	ALPHA_PREMULTIPLIED:	0x0000000000000200,
	ALPHA_NON_PREMULTIPLIED:0x0000000000000300,
	ADDITIVE: 				0x0000000000000400,
}

var g_Blend = []struct {
	Src, Dst uint32
} {
	{gl.ONE, 		gl.ZERO},
	{gl.ONE, 		gl.ONE_MINUS_SRC_ALPHA},
	{gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA},
	{gl.SRC_ALPHA, gl.ONE},
}

var ST_PT = struct {
	TRIANGLES 	   uint64
	TRIANGLE_STRIP uint64
	LINES 		   uint64
	LINE_STRIP     uint64
	POINTS         uint64
}{
	TRIANGLES: 		0x0000000000000000,
	TRIANGLE_STRIP: 0x0000000000001000,
	LINES:          0x0000000000002000,
	LINE_STRIP:     0x0000000000003000,
	POINTS:         0x0000000000004000,
}

var g_PrimInfo = []uint32 {
	gl.TRIANGLES,
	gl.TRIANGLE_STRIP,
	gl.LINES,
	gl.LINE_STRIP,
	gl.POINTS,
}

/// STATE ENCODE FORMAT - bgfx
// 64bit:
//
//                                0-4 func ----------+            rgb-dst
//                                               +   |  a-dsr        |  rgb-src +--------- 1 - 8 depth-function
//                  independent|alpha_cover --+  |   |    |    a-src |    |     |     +--- depth_write | alpha_write | rgb_write
//                                            |  |   |    |     |    |    |     |     |
// 000000000 000 0000 0-000 0000-0000-00 00 + 00 00-0000 0000-0000-0000-0000  0000-0 000
//            |    |     |      |        |
//            |    |     |      |        +---- CULL_CCW | CW
//            |    |     |      +------------- Alpha Ref Value (255)
//            |    |     +-------------------- Primitive Type
//            |    +-------------------------- Point Size (16)
//            +------------------------------- Conservative Raster | LineAA | MSAA
