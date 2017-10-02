package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx/bk"
)

/// Simple Mesh TypeRender
type MeshRender struct {
	pipeline PipelineState
	C RenderContext
}

func NewMeshRender(shader bk.GLShader) *MeshRender {
	mr := new(MeshRender)
	// blend func
	mr.pipeline.BlendFunc = BF_Add
	//
	//// setup shader
	mr.pipeline.GLShader = shader
	shader.Use()
	//
	//// ---- Fragment GLShader
	shader.SetInteger("tex\x00", 0)
	gl.BindFragDataLocation(shader.Program, 0, gl.Str("outputColor\x00"))
	//
	//// vertex layout
	pos := VertexAttr {
		Data: 0,
		Slot: 0,

		Size: 4,
		Type: gl.FLOAT,
		Normalized: false,
		Stride: 16,
		Offset: 0,
		Pointer: 0,
	}
	mr.pipeline.VertexLayout = append(mr.pipeline.VertexLayout, pos)

	// uniform layout
	p := bk.Uniform{
		Data: 0, 		// index of uniform data
		Slot: shader.GetUniformLocation("projection\x00"), 		// slot in shader
		Type: bk.UniformMat4,
		Count: 1,
	}

	m := bk.Uniform{
		Data: 1, 		// index of uniform data
		Slot: shader.GetUniformLocation("model\x00"), 			// slot in shader
		Type: bk.UniformMat4,
		Count: 1,
	}
	mr.pipeline.UniformLayout = append(mr.pipeline.UniformLayout, p, m)
	return mr
}

func (mr *MeshRender) Draw(d RenderData, pos, scale mgl32.Vec2, rot float32) {
	m := d.(*Mesh)
	//
	mr.pipeline.tex = m.tex
	//
	mr.C.SetPipelineState(mr.pipeline)
	mr.C.SetVertexBuffer(m.VertexBuffer())
	mr.C.VAO = m.VAO()
	mr.C.SetIndexBuffer(m.IndexBuffer())
	//

	proj := mgl32.Ortho2D(0, 480, 0, 320)
	mr.C.UniformData.AddUniform(0, &proj[0])
	model := mgl32.Translate3D(pos[0], pos[1], 0)
	mr.C.UniformData.AddUniform(1, &model[0])

	mr.C.Draw()
}
