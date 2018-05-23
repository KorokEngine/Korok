package gfx

import (
	"korok.io/korok/gfx/bk"
	"unsafe"
	"image/color"
)

type Color color.RGBA

var (
	While  = Color{0xFF, 0xFF, 0xFF, 0xFF}
	Blank  = Color{0,0,0,0xFF}
	LTGray = Color{0xCC,0xCC,0xCC, 0xFF}
	Gray   = Color{0x88, 0x88, 0x88, 0xFF}
	DKGray = Color{0x44, 0x44, 0x44, 0xFF}

	Red   = Color{0xFF,0,0, 0xFF}
	Green = Color{0,0xFF,0, 0xFF}
	Blue  = Color{0,0,0xFF, 0xFF}

	Cyan    = Color{0, 0xFF, 0xFF, 0xFF}
	Magenta = Color{0xFF, 00, 0xFF, 0xFF}
	Yellow  = Color{0xFF, 0xFF, 0x00, 0xFF}

	Transparent = Color{}
	Opaque      = Color{0xFF, 0xFF, 0xFF, 0xFF}
)

func (c Color) U32() uint32 {
	return *(*uint32)(unsafe.Pointer(&c))
}

func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

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





