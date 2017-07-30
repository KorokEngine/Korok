package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

type Mesh struct {
	// render handle
	vao, vbo, ebo uint32

	// vertex <x,y,u,v>
	vertex []float32
}

func (m*Mesh) Setup() (vao, vbo uint32) {
	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)

	gl.BindVertexArray(m.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertex)*4, gl.Ptr(&m.vertex[0]), gl.STATIC_DRAW)

	gl.BindVertexArray(0)
	return m.vao, m.vbo
}

func (m*Mesh) Handle() (vao, vbo uint32) {
	return m.vao, m.vbo
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



