package bk

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"unsafe"
)
const (
	UNIFORM_BUFFER_SIZE = 16 << 10
)

type UniformType uint8

/// 支持哪些类型呢?
/// https://github.com/bkaradzic/bgfx/issues/653
/// float 在 GPU 中表现为 vec4
const (
	UniformStart UniformType = iota
	UniformMat4 		// mat4 array
	UniformMat3 		// mat3 array
	UniformVec4 		// vec4 array
	UniformVec1 		// float array
	UniformInt1 		// int array

	UniformSampler 		// sampler

	UniformEnd
)

type Uniform struct {
	Name string // Uniform Name

	Slot  uint8 // slot in shader
	Size  uint8
	Type  UniformType
	Count uint8
}

func (um *Uniform) create(program uint32, name string, xType UniformType, num uint32) (slot int32) {
	slot = gl.GetUniformLocation(program, gl.Str(name))
	um.Slot = uint8(slot)
	um.Name = name
	um.Type = xType
	um.Count = uint8(num)
	um.Size = g_uniform_type2size[xType]
	return
}

type UniformBuffer struct {
	buffer [UNIFORM_BUFFER_SIZE]uint8

	size uint32
	pos  uint32
}

func NewUniformBuffer() *UniformBuffer {
	return &UniformBuffer{}
}

func (ub *UniformBuffer) GetPos() uint32 {
	return ub.pos
}

func (ub *UniformBuffer) Reset() {
	ub.pos = 0
}

func (ub *UniformBuffer) Seek(pos uint32) {
	ub.pos = pos
}

func (ub *UniformBuffer) IsEmpty() bool {
	return ub.pos == 0
}

func (ub *UniformBuffer) WriteUInt32(value uint32) {
	u32 := (*uint32)(unsafe.Pointer(&ub.buffer[ub.pos]))
	*u32 = value
	ub.pos += 4
}

func (ub *UniformBuffer) Copy(ptr unsafe.Pointer, size uint32) {
	data := (*[1024]uint8)(ptr)[:size]
	copy(ub.buffer[ub.pos:], data)
	ub.pos += size
}

func (ub *UniformBuffer) ReadUInt32() uint32 {
	u32 := (*uint32)(unsafe.Pointer(&ub.buffer[ub.pos]))
	ub.pos += 4
	return *u32
}

func (ub *UniformBuffer) ReadPointer(size uint32) unsafe.Pointer {
	ptr := &ub.buffer[ub.pos]
	ub.pos += size
	return unsafe.Pointer(ptr)
}

/// Uniform ENCODE FORMAT
// 	 0x FF FF FF FF
//	    ^  ^  ^  ^
//	    |  |  |  |
// type-+  |  |  +--- Num
// loc  ---+  +------ size

func Uniform_decode(code uint32, uType, loc, size, num *uint8) {
	*uType = uint8((code >> 24) & 0xFF)
	*loc = uint8((code >> 16) & 0xFF)
	*size = uint8((code >> 8) & 0xFF)
	*num = uint8((code >> 0) & 0xFF)
}

func Uniform_encode(uType UniformType, loc, size, num uint8) uint32 {
	return uint32(uType)<<24 |
		uint32(loc)<<16 |
		uint32(size)<<8 |
		uint32(num)<<0
}

var g_uniform_type2size = [UniformEnd]uint8{
	0,				// ignore
	16 * 4,			// mat4
	9  * 4,			// mat3
	4  * 4,			// vec4
	4,				// vec1(float32)
	4,				// int32
	4,				// sampler
}
