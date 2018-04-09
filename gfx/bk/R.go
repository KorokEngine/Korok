package bk

import (
	"korok.io/korok/hid/gl"

	"log"
	"unsafe"
	"image"
)

/// 在 2D 引擎中，GPU 资源的使用时很有限的，
/// 大部分图元会在 Batch 环节得到优化，最终生成有限的渲染指令.
/// Id从1开始，0作为Invalid-Id比0xFFFF有很多天然的好处

const InvalidId uint16 = 0x0000
const UINT16_MAX uint16 = 0xFFFF

const (
	ID_MASK       uint16 = 0x0FFF
	ID_TYPE_SHIFT        = 12
)

type Memory struct {
	Data unsafe.Pointer
	Size uint32
}

// ID FORMAT
// 0x F   FFF
//    ^    ^
//    |    +---- id value
//    +--------- id type

const (
	ID_TYPE_INDEX uint16 = iota
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
	MAX_UNIFORM = 32 * 8
	MAX_SHADER  = 32
)

type FreeList struct {
	slots []uint16
}

func (fl *FreeList) Pop() (slot uint16, ok bool) {
	if size := len(fl.slots); size > 0 {
		slot = fl.slots[size-1]
		ok = true
		fl.slots = fl.slots[:size-1]
	}
	return
}

func (fl *FreeList) Push(slot uint16) {
	fl.slots = append(fl.slots, slot)
}

type ResManager struct {
	indexBuffers  [MAX_INDEX]IndexBuffer
	vertexBuffers [MAX_VERTEX]VertexBuffer
	textures      [MAX_TEXTURE]Texture2D

	uniforms [MAX_UNIFORM]Uniform
	shaders  [MAX_SHADER]Shader

	ibIndex uint16
	vbIndex uint16
	ttIndex uint16
	vlIndex uint16
	umIndex uint16
	shIndex uint16

	// free list
	ibFrees FreeList
	vbFrees FreeList
	ttFrees FreeList
	umFrees FreeList
	shFrees FreeList
}

func NewResManager() *ResManager {
	return &ResManager{}
}

// skip first index - 0
func (rm *ResManager) Init() {
	if rm.ibIndex == 0 {
		rm.ibIndex ++
		rm.vbIndex ++
		rm.ttIndex ++
		rm.vlIndex ++
		rm.umIndex ++
		rm.shIndex ++
	}
}

func (rm *ResManager) Destroy() {

}

// AllocIndexBuffer alloc a new Index-Buffer, Return the resource handler.
func (rm *ResManager) AllocIndexBuffer(mem Memory) (id uint16, ib *IndexBuffer) {
	if index, ok := rm.ibFrees.Pop(); ok {
		id = index
		ib = &rm.indexBuffers[index]
	} else {
		id, ib = rm.ibIndex, &rm.indexBuffers[rm.ibIndex]
		rm.ibIndex++
	}
	id = id | (ID_TYPE_INDEX << ID_TYPE_SHIFT)

	if err := ib.Create(mem.Size, mem.Data, 0); err != nil {
		log.Println("fail to alloc index-buffer, ", err)
	} else {
		if g_debug&DEBUG_R != 0 {
			log.Printf("alloc index-buffer: (%d, %d)", id&ID_MASK, ib.Id)
		}
	}
	return
}

// AllocVertexBuffer alloc a new Vertex-Buffer, Return the resource handler.
func (rm *ResManager) AllocVertexBuffer(mem Memory, stride uint16) (id uint16, vb *VertexBuffer) {
	if index, ok := rm.vbFrees.Pop(); ok {
		id = index
		vb = &rm.vertexBuffers[index]
	} else {
		id, vb = rm.vbIndex, &rm.vertexBuffers[rm.vbIndex]
		rm.vbIndex++
	}
	id = id | (ID_TYPE_VERTEX << ID_TYPE_SHIFT)
	if err := vb.Create(mem.Size, mem.Data, stride, 0); err != nil {
		log.Println("fail to alloc vertex-buffer, ", err)
	} else {
		if  g_debug&DEBUG_R != 0 {
			log.Printf("alloc vertex-buffer: (%d, %d)", id&ID_MASK, vb.Id)
		}
	}
	return
}

// AllocUniform get the uniform slot in a shader program, Return the resource handler.
func (rm *ResManager) AllocUniform(shId uint16, name string, xType UniformType, num uint32) (id uint16, um *Uniform) {
	if index, ok := rm.umFrees.Pop(); ok {
		id = index
		um = &rm.uniforms[index]
	} else {
		id, um = rm.umIndex, &rm.uniforms[rm.umIndex]
		rm.umIndex++
	}
	id = id | (ID_TYPE_UNIFORM << ID_TYPE_SHIFT)
	if ok, sh := rm.Shader(shId); ok {
		if um.create(sh.Program, name, xType, num) < 0 {
			log.Printf("fail to alloc uniform - %s, make sure shader %d in use", name, shId&ID_MASK)
		} else {
			if (g_debug & DEBUG_R) != 0 {
				log.Printf("alloc uniform: %s(%d, %d)", name, id&ID_MASK, um.Slot)
			}
		}
	}
	return
}

// AllocTexture upload image to GPU, Return the resource handler.
func (rm *ResManager) AllocTexture(img image.Image) (id uint16, tex *Texture2D) {
	if index, ok := rm.ttFrees.Pop(); ok {
		id = index
		tex = &rm.textures[index]
	} else {
		id, tex = rm.ttIndex, &rm.textures[rm.ttIndex]
		rm.ttIndex ++
	}
	id = id | (ID_TYPE_TEXTURE << ID_TYPE_SHIFT)
	if err := tex.Create(img); err != nil {
		log.Printf("fail to alloc texture, %s", err)
	} else {
		if (g_debug & DEBUG_R) != 0 {
			log.Printf("alloc texture id: (%d, %d)", id&ID_MASK, tex.Id)
		}
	}
	return
}

// AllocShader compile and link the Shader source code, Return the resource handler.
func (rm *ResManager) AllocShader(vsh, fsh string) (id uint16, sh *Shader) {
	if index, ok := rm.shFrees.Pop(); ok {
		id = index
		sh = &rm.shaders[index]
	} else {
		id, sh = rm.shIndex, &rm.shaders[rm.shIndex]
		rm.shIndex++
	}
	id = id | (ID_TYPE_SHADER << ID_TYPE_SHIFT)

	if err := sh.Create(vsh, fsh); err != nil {
		log.Println("fail to alloc shader, ", err)
	} else {
		if (g_debug & DEBUG_R) != 0 {
			log.Printf("alloc shader id:(%d, %d) ", id&ID_MASK, sh.Program)
		}
	}
	return
}

// Free free all resource managed by R. Including index-buffer, vertex-buffer,
// texture, uniform and shader program.
func (rm *ResManager) Free(id uint16) {
	t := (id >> ID_TYPE_SHIFT) & 0x000F
	v := id & ID_MASK

	switch t {
	case ID_TYPE_INDEX:
		rm.indexBuffers[v].Destroy()
		rm.ibFrees.Push(v)
	case ID_TYPE_VERTEX:
		rm.vertexBuffers[v].Destroy()
		rm.vbFrees.Push(v)
	case ID_TYPE_TEXTURE:
		rm.textures[v].Destroy()
		rm.ttFrees.Push(v)
	case ID_TYPE_LAYOUT:
		// todo
	case ID_TYPE_UNIFORM:
		rm.umFrees.Push(v)
	case ID_TYPE_SHADER:
		rm.shaders[v].Destroy()
		rm.shFrees.Push(v)
	}
}

// IndexBuffer returns the low-level IndexBuffer struct.
func (rm *ResManager) IndexBuffer(id uint16) (ok bool, ib *IndexBuffer) {
	t, v := id>>ID_TYPE_SHIFT, id&ID_MASK
	if t != ID_TYPE_INDEX || v >= MAX_INDEX {
		return false, nil
	}
	return true, &rm.indexBuffers[v]
}

// VertexBuffer returns the low-level VertexBuffer struct.
func (rm *ResManager) VertexBuffer(id uint16) (ok bool, vb *VertexBuffer) {
	t, v := id>>ID_TYPE_SHIFT, id&ID_MASK
	if t != ID_TYPE_VERTEX || v >= MAX_VERTEX {
		return false, nil
	}
	return true, &rm.vertexBuffers[v]
}

// Texture returns the low-level Texture struct.
func (rm *ResManager) Texture(id uint16) (ok bool, tex *Texture2D) {
	t, v := id >>ID_TYPE_SHIFT, id&ID_MASK
	if t != ID_TYPE_TEXTURE || v >= MAX_TEXTURE {
		return false, nil
	}
	return true, &rm.textures[v]
}

// Uniform returns the low-level Uniform struct.
func (rm *ResManager) Uniform(id uint16) (ok bool, um *Uniform) {
	t, v := id>>ID_TYPE_SHIFT, id&ID_MASK
	if t != ID_TYPE_UNIFORM || v >= MAX_UNIFORM {
		return false, nil
	}
	return true, &rm.uniforms[v]
}

// Shader returns the low-level Shader struct.
func (rm *ResManager) Shader(id uint16) (ok bool, sh *Shader) {
	t, v := id>>ID_TYPE_SHIFT, id&ID_MASK
	if t != ID_TYPE_SHADER || v >= MAX_SHADER {
		if (g_debug & DEBUG_R) != 0 {
			log.Printf("Invalid shader id:(%d, %d, %d)", id, t, v)
		}
		return false, nil
	}
	return true, &rm.shaders[v]
}

////// MAX SIZE
var MAX = struct {
}{}

////// STATE MASK AND VALUE DEFINES

var ST = struct {
	RGB_WRITE   uint64
	ALPHA_WRITE uint64
	DEPTH_WRITE uint64

	DEPTH_TEST_MASK  uint64
	DEPTH_TEST_SHIFT uint64

	BLEND_MASK  uint64
	BLEND_SHIFT uint64

	PT_MASK  uint64
	PT_SHIFT uint64
}{
	RGB_WRITE:   0x0000000000000001,
	ALPHA_WRITE: 0x0000000000000002,
	DEPTH_WRITE: 0x0000000000000004,

	DEPTH_TEST_MASK:  0x00000000000000F0,
	DEPTH_TEST_SHIFT: 4,

	BLEND_MASK:  0x0000000000000F00,
	BLEND_SHIFT: 8,

	PT_MASK:  0x000000000000F000,
	PT_SHIFT: 12,
}

// zero means no depth-test
var ST_DEPTH = struct {
	LESS     uint64
	LEQUAL   uint64
	EQUAL    uint64
	GEQUAL   uint64
	GREATER  uint64
	NOTEQUAL uint64
	NEVER    uint64
	ALWAYS   uint64
}{
	LESS:     0x0000000000000010,
	LEQUAL:   0x0000000000000020,
	EQUAL:    0x0000000000000030,
	GEQUAL:   0x0000000000000040,
	GREATER:  0x0000000000000050,
	NOTEQUAL: 0x0000000000000060,
	NEVER:    0x0000000000000070,
	ALWAYS:   0x0000000000000080,
}

var g_CmpFunc = []uint32{
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

// zero means no blend
var ST_BLEND = struct {
	DEFAULT                 uint64
	ISABLE                  uint64
	ALPHA_PREMULTIPLIED     uint64
	ALPHA_NON_PREMULTIPLIED uint64
	ADDITIVE                uint64
}{
	ISABLE:                  0x0000000000000100,
	ALPHA_PREMULTIPLIED:     0x0000000000000200,
	ALPHA_NON_PREMULTIPLIED: 0x0000000000000300,
	ADDITIVE:                0x0000000000000400,
}

var g_Blend = []struct {
	Src, Dst uint32
}{
	{0, 0},
	{gl.ONE, gl.ZERO},
	{gl.ONE, gl.ONE_MINUS_SRC_ALPHA},
	{gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA},
	{gl.SRC_ALPHA, gl.ONE},
}

var ST_PT = struct {
	TRIANGLES      uint64
	TRIANGLE_STRIP uint64
	LINES          uint64
	LINE_STRIP     uint64
	POINTS         uint64
}{
	TRIANGLES:      0x0000000000000000,
	TRIANGLE_STRIP: 0x0000000000001000,
	LINES:          0x0000000000002000,
	LINE_STRIP:     0x0000000000003000,
	POINTS:         0x0000000000004000,
}

var g_PrimInfo = []uint32{
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
