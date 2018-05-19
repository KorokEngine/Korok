package gfx

import (
	"korok.io/korok/gfx/bk"
	"unsafe"
)

type CompRef struct {
	Type int32
	*Transform
	*SpriteComp
}

type AABB struct {
	x, y float32
	width, height float32
}

func OverlapAB(a, b *AABB) bool {
	if a.x < b.x+b.width && a.x+a.width>b.x && a.y < b.y+b.height && a.y+a.height > b.y {
		return true
	}
	return false
}

type zOrder struct {
	value int16
}

func (zo *zOrder) SetZOrder(z int16) {
	zo.value = z
}

func (zo *zOrder) Z() int16 {
	return zo.value
}

type batchId struct {
	value uint16
}

func (b *batchId) SetBatchId(id uint16) {
	b.value = id
}

func (b *batchId) BatchId() uint16 {
	return b.value
}

func PackSortId(z int16, b uint16) (sid uint32) {
	 sid = uint32(int32(z) + 0xFFFF>>1)
	 sid = (sid << 16) + uint32(b)
	 return
}

func UnpackSortId(sortId uint32) (z int16, b uint16) {
	b = uint16(sortId & 0xFFFF)
	z = int16(int32(sortId>>16)-0xFFFF>>1)
	return
}

// format <x,y,u,v rgba>
var P4C4 = []bk.VertexComp{
	{4, bk.AttrFloat, 0, 0},
	{4, bk.AttrUInt8, 16, 1},
}

// vertex struct
type PosTexColorVertex struct {
	X, Y, U, V float32
	RGBA       uint32
}

//
var PosTexColorVertexSize = unsafe.Sizeof(PosTexColorVertex{})
var UInt16Size = unsafe.Sizeof(uint16(0))





