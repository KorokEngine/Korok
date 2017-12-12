package dbg

import (
	"korok.io/korok/gfx/bk"

	"github.com/go-gl/mathgl/mgl32"

	"unsafe"
	"log"
)

type DrawType uint8
const (
	TEXT DrawType = iota
	RECT
	CIRCLE
)

// dbg - draw debug info
// provide self-contained, im-gui api
// can be used to show debug info, fps..
// dependents: bk-api

var P4C4 = []bk.VertexComp{
	{4, bk.ATTR_TYPE_FLOAT, 0, 0},
	{4, bk.ATTR_TYPE_UINT8, 16, 1},
}

type PosTexColorVertex struct {
	X, Y, U, V float32
	RGBA       uint32
}

func Init() {
	if g_render == nil {
		g_render = NewDebugRender(vsh, fsh)
		g_buffer = &g_render.Buffer
	}
}

func Destroy() {
	if g_buffer != nil {
		g_buffer.Destroy()
	}
}

func FPS(fps int32) {

}

func Move(x, y float32) {
	g_buffer.x, g_buffer.y = x, y
}

func Color(argb uint32) {
	g_buffer.color = argb
}

// draw a rect
func DrawRect(x, y, w, h float32) {
	g_buffer.Rect(x, y, w, h)
}

// draw a circle
func DrawCircle(x,y float32, r float32) {

}

// draw string
func DrawStr(str string) {
	g_buffer.String(str, 1)
}

func DrawStrScaled(str string, scale float32) {
	g_buffer.String(str, scale)
}

func NextFrame() {
	g_buffer.Update()
	g_render.Draw()
	g_buffer.Reset()
}


type DebugRender struct {
	stateFlags uint64
	rgba       uint32

	// shader program
	program uint16

	// uniform handle
	umh_P  uint16 // Projection
	umh_S0 uint16 // Sampler0

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
		sh.AddAttributeBinding("xyuv\x00", 0, P4C4[0])
		sh.AddAttributeBinding("rgba\x00", 0, P4C4[1])

		p := mgl32.Ortho2D(0, 480, 0, 320)
		s0 := int32(0)

		// setup uniform
		if pid, _ := bk.R.AllocUniform(id, "proj\x00", bk.UniformMat4, 1); pid != bk.InvalidId {
			dr.umh_P = pid
			bk.SetUniform(pid, unsafe.Pointer(&p[0]))
		}
		if sid,_ := bk.R.AllocUniform(id, "tex\x00", bk.UniformSampler, 1); sid != bk.InvalidId {
			dr.umh_S0 = sid
			bk.SetUniform(sid, unsafe.Pointer(&s0))
		}

		// submit render state
		bk.Touch(0)
	}
	// setup buffer, we can draw 512 rect at most!!
	dr.Buffer.init(512)
	return dr
}

func (dr *DebugRender) Draw() {
	bk.SetState(dr.stateFlags, dr.rgba)
	bk.SetTexture(0, dr.umh_S0, uint16(dr.Buffer.fontTexId), 0)

	b := &dr.Buffer
	// set vertex
	bk.SetVertexBuffer(0, b.vertexId, 0, b.pos)
	bk.SetIndexBuffer(dr.Buffer.indexId, 0, b.pos * 6 >> 2)
	// submit
	bk.Submit(0, dr.program, 0)
}

type TextShapeBuffer struct {
	// real data
	vertex []PosTexColorVertex
	index  []uint16

	// gpu res
	indexId, vertexId uint16
	ib *bk.IndexBuffer
	vb *bk.VertexBuffer
	fontTexId uint16

	// current cursor position and painter color
	x, y float32
	color uint32

	// current buffer position
	pos uint32
}

func (buff *TextShapeBuffer) init(maxVertex uint32) {
	iboSize := maxVertex * 6 / 4
	buff.index = make([]uint16, iboSize)
	iFormat := [6]uint16 {3, 1, 2, 3, 2, 0}
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
}

//
//  3-------0
//  |       |
//  |       |
//  1-------2
func (buff *TextShapeBuffer) String(chars string, scale float32) {
	x, y := float32(0), float32(0)
	w, h := font_width * scale, font_height * scale

	for i, N := 0, len(chars); i < N; i++ {
		b := buff.vertex[buff.pos: buff.pos+4]
		buff.pos += 4

		// vv := chars[0]
		var left, right, bottom, top float32 = GlyphRegion(chars[i])
		bottom, top = top, bottom

		b[0].X, b[0].Y = buff.x + x + w, buff.y + y + h
		b[0].U, b[0].V = right, top
		b[0].RGBA = buff.color

		b[1].X, b[1].Y = buff.x + x + 0, buff.y + y + 0
		b[1].U, b[1].V = left, bottom
		b[1].RGBA = buff.color

		b[2].X, b[2].Y = buff.x + x + w, buff.y + y + 0
		b[2].U, b[2].V = right, bottom
		b[2].RGBA = buff.color

		b[3].X, b[3].Y = buff.x + x + 0, buff.y + y + h
		b[3].U, b[3].V = left, top
		b[3].RGBA = buff.color

		// advance x,y
		x += w
	}
}


//
//  3-------0
//  |       |
//  |       |
//  1-------2
func (buff *TextShapeBuffer) Rect(x,y, w, h float32) {
	b := buff.vertex[buff.pos: buff.pos+4]
	buff.pos += 4

	b[0].X, b[0].Y = buff.x + x + w, buff.y + y + h
	b[0].U, b[0].V = 2, 0
	b[0].RGBA = buff.color

	b[1].X, b[1].Y = buff.x + x + 0, buff.y + y + 0
	b[1].U, b[1].V = 2, 0
	b[1].RGBA = buff.color

	b[2].X, b[2].Y = buff.x + x + w, buff.y + y + 0
	b[2].U, b[2].V = 2, 0
	b[2].RGBA = buff.color

	b[3].X, b[3].Y = buff.x + x + 0, buff.y + y + h
	b[3].U, b[3].V = 2, 0
	b[3].RGBA = buff.color
}

func (buff *TextShapeBuffer) Update() {
	buff.vb.Update(0, buff.pos * 20, unsafe.Pointer(&buff.vertex[0]), false)
}

func (buff *TextShapeBuffer) Reset() {
	buff.pos = 0
}

func (buff *TextShapeBuffer) Destroy() {

}

//// static filed
var g_render *DebugRender
var g_buffer *TextShapeBuffer


//// frag & vertex shader
var vsh = `
#version 330

uniform mat4 proj;

in vec4 xyuv;
in vec4 rgba;

out vec4 outColor;
out vec2 fragTexCoord;

void main() {
    outColor = rgba;
	fragTexCoord = xyuv.zw;
    gl_Position = proj * vec4(xyuv.xy, 1, 1);
}
` + "\x00"

var fsh = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec4 outColor;
out vec4 outputColor;
void main() {
	if (fragTexCoord.x == 2) {
		outputColor = outColor;
	} else {
	    outputColor = outColor * texture(tex, fragTexCoord);
	}
}
` + "\x00"