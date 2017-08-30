package gfx

import "github.com/go-gl/mathgl/mgl32"

// USED to model sprite
// <x,y, r,g,b,a, u,v> compress Data TODO

type QuadVertex struct {
	XY [2]float32
	UV [2]float32
	Color int32
}

type Quad struct {
	buf [4]QuadVertex
	tex uint32

	next int32
}

func (*Quad) Type() int32 {
	return 1
}

func (q *Quad) SetXY(xy mgl32.Vec2) {
	for i := 0; i < 4; i++ {
		q.buf[i].XY[0] += xy[0]
		q.buf[i].XY[1] += xy[1]
	}
}

func (q *Quad) SetTexture(tex *SubTex) {
	q.tex = tex.Id

}

// a buffer to keep quad-Data together
type QuadBuffer struct {
	// quad Data
	buf []Quad

	// freelist and used count
	freeList int32
	count int32
}

func NewQuadBuffer(cap int32) *QuadBuffer{
	qb := new(QuadBuffer)
	qb.buf = make([]Quad, cap)
	qb.count = 0

	// init freelist
	for i := int32(0); i < cap - 1; i++ {
		qb.buf[i].next = i + 1
	}
	qb.buf[cap - 1].next = -1
	qb.freeList = 0

	return qb
}

func (qb *QuadBuffer) Alloc() *Quad {
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








