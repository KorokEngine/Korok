package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"log"
)

type Mesh struct {
	// render handle
	vao, vbo, ebo uint32

	// vertex <x,y,u,v>
	vertex []float32
	index  []uint32

	//
	tex uint32
}

func (m*Mesh) Setup() (vao, vbo uint32) {
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.BindVertexArray(m.vao)

	// vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertex)*4, gl.Ptr(m.vertex), gl.STATIC_DRAW)

	// ebo optional
	if m.index != nil {
		gl.GenBuffers(1, &m.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.index)*4, gl.Ptr(m.index), gl.STATIC_DRAW)
	}

	if m.index != nil && m.ebo == 0 {
		log.Fatal("ebo setup err!!!")
	}

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	return m.vao, m.vbo
}

func (m*Mesh) Handle() (vao, vbo, ebo uint32) {
	return m.vao, m.vbo, m.ebo
}

func (m*Mesh) SetVertex(v []float32) {
	m.vertex = v
}

func (m*Mesh) Update() {
	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	// gl.buffersub!!
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertex)*4, gl.Ptr(&m.vertex[0]), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
}

func (m*Mesh) Delete() {
	gl.DeleteBuffers(1, &m.vbo)
	gl.DeleteBuffers(1, &m.ebo)
	gl.DeleteVertexArrays(1, &m.vao)
}

// new mesh from Texture
func NewQuadMesh(tex *Texture2D) *Mesh{
	m := new(Mesh)
	m.tex = tex.Id
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
		0.0, h, tex.Min[0]/tex.Width, tex.Max[1]/tex.Height,
		w, 0.0 , tex.Max[0]/tex.Width, tex.Min[1]/tex.Height,
		0.0, 0.0	   , tex.Min[0]/tex.Width, tex.Min[1]/tex.Height,

		0.0, h, tex.Min[0]/tex.Width, tex.Max[1]/tex.Height,
		w, h, tex.Max[0]/tex.Width, tex.Max[1]/tex.Height,
		w, 0.0 , tex.Max[0]/tex.Width, tex.Min[1]/tex.Height,
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
	m.index = []uint32{
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

