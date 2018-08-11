package dbg

import (
	"korok.io/korok/gfx/bk"
	"korok.io/korok/math/f32"

	"unsafe"
	"log"
	"fmt"
	"korok.io/korok/math"
)

type DebugEnum uint32

const (
	FPS DebugEnum = 1 << iota
	Stats
	Draw

	ALL = FPS|Stats|Draw
	None = DebugEnum(0)
)


var screen struct{
	w, h float32
}
// max value of z-order
const zOrder = int32(0xFFFF>>1)

// dbg - draw debug info
// provide self-contained, im-gui api
// can be used to show debug info, fps..
// dependents: bk-api

var p2t2c4 = []bk.VertexComp{
	{4, bk.AttrFloat, 0, 0},
	{4, bk.AttrUInt8, 16, 1},
}

type PosTexColorVertex struct {
	X, Y, U, V float32
	RGBA       uint32
}

func Init(w, h int) {
	if gRender == nil {
		gRender = NewDebugRender(vsh, fsh)
		gBuffer = &gRender.Buffer
		hud = &HudLog{}
	}
	screen.w = float32(w)
	screen.h = float32(h)

	log.Println("dbg init w,h", w, h)
}

func SetDebug(enum DebugEnum) {
	DEBUG = enum
}

func SetCamera(x,y, w,h float32) {
	gRender.SetViewPort(x, y, w, h)
}

func Destroy() {
	if gBuffer != nil {
		gBuffer.Destroy(); gBuffer = nil
	}
	if gRender != nil {
		gRender.Destroy(); gRender = nil
	}
}

func Color(argb uint32) {
	gBuffer.color = argb
}

// draw a rect
func DrawRect(x, y, w, h float32) {
	if (DEBUG & Draw) != 0 {
		gBuffer.Rect(x, y, w, h)
	}
}

func DrawBorder(x, y, w, h, thickness float32) {
	if (DEBUG & Draw) != 0 {
		gBuffer.Border(x, y, w, h, thickness)
	}
}

// draw a circle
func DrawCircle(x,y float32, r float32) {
	if (DEBUG & Draw) != 0 {
		gBuffer.Circle(x, y, r)
	}
}

func DrawLine(from, to f32.Vec2) {
	if (DEBUG & Draw) != 0 {
		gBuffer.Line(from, to)
	}
}

// draw string
func DrawStr(x,y float32, str string, args ...interface{}) {
	if (DEBUG & Draw) != 0 {
		gBuffer.String(x, y, fmt.Sprintf(str, args...), 1)
	}
}

func DrawStrScaled(x, y float32, scale float32, str string, args ...interface{}) {
	if (DEBUG & Draw) != 0 {
		gBuffer.String(x, y, fmt.Sprintf(str, args...), scale)
	}
}

func NextFrame() {
	// draw hud
	hud.draw()
	hud.reset()

	// flush
	gBuffer.Update()
	gRender.Draw()
	gBuffer.Reset()
}


type DebugRender struct {
	stateFlags uint64
	rgba       uint32
	view struct{
		x, y, w, h float32
	}

	// shader program
	program uint16

	// uniform handle
	umhProjection uint16 // Projection
	umhSampler0   uint16 // Sampler0

	// buffer
	Buffer TextShapeBuffer
}

func NewDebugRender(vsh, fsh string) *DebugRender {
	dr := new(DebugRender)
	// blend func
	dr.stateFlags |= bk.ST_BLEND.ALPHA_NON_PREMULTIPLIED

	// setup shader
	if id, sh := bk.R.AllocShader(vsh, fsh); id != bk.InvalidId {
		dr.program = id
		sh.Use()

		// setup attribute
		sh.AddAttributeBinding("xyuv\x00", 0, p2t2c4[0])
		sh.AddAttributeBinding("rgba\x00", 0, p2t2c4[1])

		s0 := int32(0)
		// setup uniform
		if pid, _ := bk.R.AllocUniform(id, "projection\x00", bk.UniformMat4, 1); pid != bk.InvalidId {
			dr.umhProjection = pid
		}
		if sid,_ := bk.R.AllocUniform(id, "tex\x00", bk.UniformSampler, 1); sid != bk.InvalidId {
			dr.umhSampler0 = sid
			bk.SetUniform(sid, unsafe.Pointer(&s0))
		}

		// submit render state
		//bk.Touch(0)
		bk.Submit(0, id, zOrder)
	}
	// setup buffer, we can draw 512 rect at most!!
	dr.Buffer.init(2048*4)
	return dr
}

func (dr *DebugRender) Destroy() {
	bk.R.Free(dr.program)
	bk.R.Free(dr.umhProjection)
	bk.R.Free(dr.umhSampler0)
}

func (dr *DebugRender) SetViewPort(x,y, w,h float32) {
	dr.view.x = x
	dr.view.y = y
	dr.view.w = w
	dr.view.h = h

	var (
		left   = x - dr.view.w/2
		right  = x + dr.view.w/2
		bottom = y - dr.view.h/2
		top    = y + dr.view.h/2
	)

	p := f32.Ortho2D(left, right, bottom, top)

	// setup uniform
	bk.SetUniform(dr.umhProjection, unsafe.Pointer(&p[0]))
	bk.Submit(0, dr.program, zOrder)
}

//func (dr *DebugRender) SetViewPort(x, y, w, h float32) {
//	dr.view.w, dr.view.h = w, h
//	p := f32.Ortho2D(0, w, 0, h)
//	bk.SetUniform(dr.umhProjection, unsafe.Pointer(&p[0]))
//	bk.Submit(0, dr.program, zOrder)
//}

func (dr *DebugRender) Draw() {
	bk.SetState(dr.stateFlags, dr.rgba)
	bk.SetTexture(0, dr.umhSampler0, uint16(dr.Buffer.fontTexId), 0)

	b := &dr.Buffer
	// set vertex
	bk.SetVertexBuffer(0, b.vertexId, 0, b.pos)
	bk.SetIndexBuffer(dr.Buffer.indexId, 0, b.pos * 6 >> 2)
	// submit
	bk.Submit(0, dr.program, zOrder)
}

// Rect:
//   3 ---- 2
//   | `    |
//   |   `  |
//   0------1
// Order:
// 3, 0, 1, 3, 1, 2
type TextShapeBuffer struct {
	// real data
	vertex []PosTexColorVertex
	index  []uint16

	// gpu res
	indexId, vertexId uint16
	ib *bk.IndexBuffer
	vb *bk.VertexBuffer
	fontTexId uint16

	// current painter color
	color uint32

	// current buffer position
	pos uint32
}

func (buff *TextShapeBuffer) init(maxVertex uint32) {
	iboSize := maxVertex * 6 / 4
	buff.index = make([]uint16, iboSize)
	iFormat := [6]uint16 {3, 0, 1, 3, 1, 2}
	for i := uint32(0); i < iboSize; i += 6 {
		copy(buff.index[i:], iFormat[:])
		iFormat[0] += 4
		iFormat[1] += 4
		iFormat[2] += 4
		iFormat[3] += 4
		iFormat[4] += 4
		iFormat[5] += 4
	}
	if id, ib := bk.R.AllocIndexBuffer(bk.Memory{unsafe.Pointer(&buff.index[0]), iboSize}); id != bk.InvalidId {
		buff.indexId = id
		buff.ib = ib
	}

	buff.vertex = make([]PosTexColorVertex, maxVertex)
	vboSize := maxVertex * 20
	if id, vb := bk.R.AllocVertexBuffer(bk.Memory{nil, vboSize}, 20); id != bk.InvalidId {
		buff.vertexId = id
		buff.vb = vb
	}

	// texture
	img, fmt, err := LoadFontImage()
	if err != nil {
		log.Println("fail to load font image.. fmt:", fmt)
	}
	if id, _ := bk.R.AllocTexture(img); id != bk.InvalidId {
		buff.fontTexId = id
	}
	buff.color = 0xFF000000
}

func (buff *TextShapeBuffer) String(x, y float32, chars string, scale float32) {
	w, h := font_width * scale, font_height * scale

	for i, N := 0, len(chars); i < N; i++ {
		b := buff.vertex[buff.pos: buff.pos+4]
		buff.pos += 4

		// vv := chars[0]
		var left, right, bottom, top float32 = GlyphRegion(chars[i])
		bottom, top = top, bottom

		b[0].X, b[0].Y = x, y
		b[0].U, b[0].V = left, bottom
		b[0].RGBA = buff.color

		b[1].X, b[1].Y = x + w, y
		b[1].U, b[1].V = right, bottom
		b[1].RGBA = buff.color

		b[2].X, b[2].Y = x + w, y + h
		b[2].U, b[2].V = right, top
		b[2].RGBA = buff.color

		b[3].X, b[3].Y = x, y + h
		b[3].U, b[3].V = left, top
		b[3].RGBA = buff.color

		// advance x,y
		x += w
	}
}


//
//  3-------2
//  |       |
//  |       |
//  0-------1
func (buff *TextShapeBuffer) Rect(x,y, w, h float32) {
	b := buff.vertex[buff.pos: buff.pos+4]
	buff.pos += 4

	b[0].X, b[0].Y = x, y
	b[0].U, b[0].V = 2, 0
	b[0].RGBA = buff.color

	b[1].X, b[1].Y = x + w, y
	b[1].U, b[1].V = 2, 0
	b[1].RGBA = buff.color

	b[2].X, b[2].Y = x + w, y + h
	b[2].U, b[2].V = 2, 0
	b[2].RGBA = buff.color

	b[3].X, b[3].Y = x, y + h
	b[3].U, b[3].V = 2, 0
	b[3].RGBA = buff.color
}

func (buff *TextShapeBuffer) Line(from, to f32.Vec2) {
	b := buff.vertex[buff.pos: buff.pos+4]
	buff.pos += 4

	diff := to.Sub(from)
	invLength := math.InvLength(diff[0], diff[1], 1.0)
	diff = diff.Mul(invLength)
	thickness := float32(1)

	dx := diff[1] * (thickness * 0.5)
	dy := diff[0] * (thickness * 0.5)

	b[0].X, b[0].Y = from[0]+dx, from[1]-dy
	b[0].U, b[0].V = 2, 0
	b[0].RGBA = buff.color

	b[1].X, b[1].Y = to[0]+dx, to[1]-dy
	b[1].U, b[1].V = 2, 0
	b[1].RGBA = buff.color

	b[2].X, b[2].Y = to[0]-dx, to[1]+dy
	b[2].U, b[2].V = 2, 0
	b[2].RGBA = buff.color

	b[3].X, b[3].Y = from[0]-dx, from[1]+dy
	b[3].U, b[3].V = 2, 0
	b[3].RGBA = buff.color
}

func (buff *TextShapeBuffer) Border(x, y, w, h, thick float32) {
	buff.Rect(x,y,w,thick)
	buff.Rect(x,y+h-thick,w,thick)
	buff.Rect(x, y, thick, h)
	buff.Rect(x+w-thick,y,thick,h)
}

func (buff *TextShapeBuffer) Circle(x, y float32, radius float32) {
	var (
		segments = 12
		path = [24]f32.Vec2{}
		angle = float32(3.14*2)
	)

	switch {
	case radius < 4:
		segments = 4
	case radius < 100:
		segments = int(radius/100 * 16) + 8
	default:
		segments = 24
	}

	for i := 0; i < segments; i++ {
		a := float32(i)/float32(segments) * angle
		x1 := x + math.Cos(a) * radius
		y1 := y + math.Sin(a) * radius
		path[i] = f32.Vec2{x1, y1}
	}
	for i := 0; i < segments; i++ {
		j := i+1
		if j == segments {
			j = 0
		}
		p1, p2 := path[i],path[j]
		buff.Line(p1, p2)
	}
}

func (buff *TextShapeBuffer) Update() {
	buff.vb.Update(0, buff.pos * 20, unsafe.Pointer(&buff.vertex[0]), false)
}

func (buff *TextShapeBuffer) Reset() {
	buff.pos = 0
}

func (buff *TextShapeBuffer) Destroy() {
	buff.vertex = nil
	buff.index = nil
	bk.R.Free(buff.vertexId)
	bk.R.Free(buff.indexId)
	bk.R.Free(buff.fontTexId)
}

//// static filed
var gRender *DebugRender
var gBuffer *TextShapeBuffer
var hud *HudLog
var DEBUG DebugEnum = FPS|Draw
