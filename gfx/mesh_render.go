package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx/bk"
	"unsafe"
)

/// Simple Mesh TypeRender
/// For simple mesh, not 3D model

// <x,y,u,v rgba>
var P4C4 = []bk.VertexComp{
	{4, bk.ATTR_TYPE_FLOAT, 0, 0},
	{4, bk.ATTR_TYPE_UINT8, 16, 1},
}

type MeshRender struct {
	stateFlags uint64
	rgba       uint32

	// shader program
	program uint16

	// uniform handle
	umh_P  uint16 // Projection
	umh_M  uint16 // Model
	umh_S0 uint16 // Sampler0
}

func NewMeshRender(vsh, fsh string) *MeshRender {
	mr := new(MeshRender)
	// blend func
	mr.stateFlags |= bk.ST_BLEND.ADDITIVE

	// setup shader
	if id, sh := bk.R.AllocShader(vsh, fsh); id != bk.InvalidId {
		mr.program = id

		// setup attribute
		sh.AddAttributeBinding("xyuv", 0, P4C4[0])
		sh.AddAttributeBinding("rgba", 0, P4C4[1])

		// setup uniform
		mr.umh_P, _ = bk.R.AllocUniform(id, "proj\x00", bk.UniformMat4, 1)
		mr.umh_M, _ = bk.R.AllocUniform(id, "model\x00", bk.UniformMat4, 1)

		// TODO
		sh.SetInteger("tex\x00", 0)
		gl.BindFragDataLocation(sh.Program, 0, gl.Str("outputColor\x00"))

	}
	return mr
}

// extract render object
func (mr *MeshRender) Extract() {

}

// draw
func (mr *MeshRender) Draw(d RenderData, pos, scale mgl32.Vec2, rot float32) {
	m := d.(*Mesh)

	// state
	bk.SetState(mr.stateFlags, mr.rgba)
	bk.SetTexture(0, mr.umh_S0, uint16(m.TextureId), 0)

	// set uniform - mvp
	proj := mgl32.Ortho2D(0, 480, 0, 320)
	model := mgl32.Translate3D(pos[0], pos[1], 0)

	bk.SetUniform(mr.umh_P, unsafe.Pointer(&proj[0]), 16)
	bk.SetUniform(mr.umh_M, unsafe.Pointer(&model[0]), 16)

	// set vertex
	bk.SetVertexBuffer(0, m.VertexId, 0, 0)
	bk.SetIndexBuffer(m.IndexId, 0, 0)

	bk.Submit(0, mr.program, 0)
}
