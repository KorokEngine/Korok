package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"fmt"
)

//
type Mesh struct {
	// render handle
	vao uint32
	vbo, ebo Buffer

	// vertex <x,y,u,v>
	vertex []float32
	index  []int32

	//
	tex uint32
}

func (*Mesh) Type() int32{
	return 0
}

func (m *Mesh) Setup() {
	m.vbo = NewArrayBuffer(Format_POS_COLOR_UV)
	m.vbo.Update(gl.Ptr(m.vertex), len(m.vertex) * 4)

	//// ebo optional
	//if m.index != nil {
	//	m.ebo = NewIndexBuffer()
	//	m.ebo.Update(gl.Ptr(m.index), len(m.index) * 4)
	//}

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	return
}

func (m*Mesh) SRT(pos mgl32.Vec2, rot float32, scale mgl32.Vec2) {
	for i := 0; i < 4; i++ {
		m.vertex[i*4 + 0] += pos[0]
		m.vertex[i*4 + 1] += pos[1]
	}
}

func (m*Mesh) VertexBuffer() Buffer {
	return m.vbo
}

func (m*Mesh) IndexBuffer() Buffer {
	return m.ebo
}

func (m*Mesh) VAO() uint32 {
	return m.vao
}

func (m*Mesh) SetVertex(v []float32) {
	m.vertex = v
}

func (m*Mesh) SetIndex(v []int32) {
	m.index = v
}

func (m*Mesh) SetRawTexture(tex uint32) {
	m.tex = tex
}

func (m *Mesh) Update() {
	m.vbo.Update(gl.Ptr(m.vertex), len(m.vertex) * 4)
	m.ebo.Update(gl.Ptr(m.index), len(m.vertex) * 4)
}

func (m*Mesh) Delete() {
	m.vbo.Delete()
	m.ebo.Delete()
	gl.DeleteVertexArrays(1, &m.vao)
}

// new mesh from Texture
func NewQuadMesh(tex *Texture2D) *Mesh{
	m := new(Mesh)
	m.tex = tex.Id

	fmt.Println("w:", tex.Width, " h:", tex.Height)

	tex.Width = 50
	tex.Height = 50

	m.vertex = []float32{
		// Pos      	 // Tex
		0.0, tex.Height, 0.0, 1.0,
		tex.Width, 0.0 , 1.0, 0.0,
		0.0, 0.0  	   , 0.0, 0.0,

		0.0, tex.Height, 0.0, 1.0,
		tex.Width, tex.Height, 1.0, 1.0,
		tex.Width, 0.0 , 1.0, 0.0,
	}
	return m
}

// new mesh from SubTexture
func NewQuadMeshSubTex(tex *SubTex) *Mesh {
	m := new(Mesh)
	m.tex = tex.Id

	h, w := tex.Height, tex.Width
	m.vertex = []float32{
		// pos 			 // tex
		0, h, tex.Min[0]/tex.Width, tex.Max[1]/tex.Height,
		w, 0, tex.Max[0]/tex.Width, tex.Min[1]/tex.Height,
		0, 0, tex.Min[0]/tex.Width, tex.Min[1]/tex.Height,

		0, h, tex.Min[0]/tex.Width, tex.Max[1]/tex.Height,
		w, h, tex.Max[0]/tex.Width, tex.Max[1]/tex.Height,
		w, 0, tex.Max[0]/tex.Width, tex.Min[1]/tex.Height,
	}
	return m
}

func NewIndexedMesh(tex *Texture2D) *Mesh {
	m := new(Mesh)
	m.tex = tex.Id

	h, w := tex.Height, tex.Width
	m.vertex = []float32{
		0,  h,  0.0, 1.0,
		w,  0,  1.0, 0.0,
		0,  0,  0.0, 0.0,
		w,  h,  1.0, 1.0,
	}
	m.index = []int32{
		0, 1, 2,
		0, 3, 1,
	}
	return m
}

// Configure VAO/VBO TODO
// 如果每个Sprite都创建一个VBO还是挺浪费的，
// 但是如果不创建新的VBO，那么怎么处理纹理坐标呢？
// 2D 场景中会出现大量模型和纹理相同的物体，仅仅
// 位置不同，比如满屏的子弹
// 或许可以通过工厂来构建mesh，这样自动把重复的mesh丢弃
// mesh 数量最终 <= 精灵的数量
var vertices = []float32{
	// Pos      // Tex
	0.0, 1.0, 0.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 0.0,

	0.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
}



