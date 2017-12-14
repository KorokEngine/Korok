package gfx

import (
	"korok.io/korok/gfx/bk"
	//"log"
)

// USED to model sprite

// Anchor type
type Anchor uint8
const(
	ANCHOR_CENTER Anchor = 0x00
	ANCHOR_LEFT          = 0x01
	ANCHOR_RIGHT         = 0x02
	ANCHOR_UP  			 = 0x04
	ANCHOR_DOWN          = 0x08
)

type Region struct {
	X1, Y1 float32
	X2, Y2 float32
}

type SubTex struct {
	TexId   uint16
	padding uint16
	Width   uint16
	Height  uint16

	Region
}

// Sprite 只需要记录 Width/Height/U/V/Tex 即可
// Anchor
type QuadSprite struct {
	// Size
	width, height uint16

	// texture region
	region Region

	// texture id
	tex uint16
	anchor  Anchor

	// next free
	next int32
}

func (*QuadSprite) Type() int32 {
	return 1
}

func (q *QuadSprite) SetSize(w, h uint16) {
	q.width = w
	q.height = h
}

func (q *QuadSprite) SetTexture(id uint16, tex *bk.SubTex) {
	q.tex = id
	q.region.X1 = tex.Min[0]/tex.Width
	q.region.Y1 = tex.Min[1]/tex.Height
	q.region.X2 = tex.Max[0]/tex.Width
	q.region.Y2 = tex.Max[1]/tex.Height
}

func (q *QuadSprite) SetAnchor(anchor Anchor) {
	q.anchor = anchor
}

// a buffer to keep quad-Data together
type QuadBuffer struct {
	// quad Data
	buf []QuadSprite

	// freelist and used count
	freeList int32
	count int32
}

func NewQuadBuffer(cap int32) *QuadBuffer{
	qb := new(QuadBuffer)
	qb.buf = make([]QuadSprite, cap)
	qb.count = 0

	// init freelist
	for i := int32(0); i < cap - 1; i++ {
		qb.buf[i].next = i + 1
	}
	qb.buf[cap - 1].next = -1
	qb.freeList = 0

	return qb
}

func (qb *QuadBuffer) Alloc() *QuadSprite {
	cap := int32(len(qb.buf))
	if qb.count >= cap {
		qb.Resize(cap + 100)
	}

	free := qb.freeList
	qb.freeList = qb.buf[free].next

	return &qb.buf[free]
}

func (qb *QuadBuffer) Free(id int32) {
	qb.buf[id].next = qb.freeList
	qb.freeList = id
	qb.count -= 1
}

func (qb *QuadBuffer) Resize(cap int32) {
	// TODO
}








