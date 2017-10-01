package bk

import (
	"unsafe"
	"github.com/go-gl/gl/v3.2-core/gl"
)

type UniformType uint8

/// 支持哪些类型呢?
const (
	UniformStart UniformType = iota
	UniformMat4

	UniformFloat1 	// float
	UniformFloat2   // vec2
	UniformFloat3   // vec3
	UniformFloat4   // vec4
	UniformFloatN   // float[]

	UniformInt1 	// int
	UniformInt2		// int_vec2
	UniformInt3 	// int_vec3
	UniformInt4 	// int_vec4
	UniformIntN  	// int[]

	UniformEnd
)

type Uniform struct {
	Name string 	// Uniform Name

	Slot uint8 		// slot in shader
	Size uint8
	Type UniformType
	Count uint8
}

func (um *Uniform) create(program uint32, name string, xType UniformType, num uint32) {
	um.Slot = uint8(gl.GetAttribLocation(program, gl.Str(name)))
	um.Name = name
	um.Type = xType
	um.Count = uint8(num)
	um.Size = g_uniform_type2size[xType]
}

type UniformBuffer struct {
	buffer []uint8

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

func Uniform_decode(code uint32, uType, loc, size, num *uint8 ) {
	*uType = uint8((code >> 24) & 0xFF)
	*loc   = uint8((code >> 16) & 0xFF)
	*size  = uint8((code >>  8) & 0xFF)
	*num   = uint8((code >>  0) & 0xFF)
}

func Uniform_encode(uType UniformType, loc , size , num uint8) uint32 {
	return uint32(uType) << 24 |
		   uint32(loc)   << 16 |
		   uint32(size)  <<  8 |
		   uint32(num)   <<  0
}

var g_uniform_type2size = [UniformEnd]uint8 {
	0,				// ignore
	16, 			// mat4
	4, 4, 4, 4, 4, 	// float32
	4, 4, 4, 4, 4,  // int32
}