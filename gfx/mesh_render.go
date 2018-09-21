package gfx

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx/bk"

	"unsafe"
)

/// Simple Mesh TypeRender
/// For simple mesh, not 3D model

type MeshRender struct {
	stateFlags uint64
	rgba       uint32

	// shader program
	program uint16

	// uniform handle
	umhProjection uint16 // Projection
	umhModel      uint16 // Model
	umhSampler0   uint16 // Sampler0
}

func NewMeshRender(vsh, fsh string) *MeshRender {
	mr := new(MeshRender)
	// blend func
	mr.stateFlags |= bk.ST_BLEND.ALPHA_PREMULTIPLIED

	// setup shader
	if id, sh := bk.R.AllocShader(vsh, fsh); id != bk.InvalidId {
		mr.program = id
		sh.Use()

		// setup attribute
		sh.AddAttributeBinding("xyuv\x00", 0, P4C4[0])
		sh.AddAttributeBinding("rgba\x00", 0, P4C4[1])

		s0 := int32(0)
		// setup uniform
		if pid, _ := bk.R.AllocUniform(id, "proj\x00", bk.UniformMat4, 1); pid != bk.InvalidId {
			mr.umhProjection = pid
		}

		if mid, _ := bk.R.AllocUniform(id, "model\x00", bk.UniformMat4, 1); mid != bk.InvalidId {
			mr.umhModel = mid
		}

		if sid,_ := bk.R.AllocUniform(id, "tex\x00", bk.UniformSampler, 1); sid != bk.InvalidId {
			mr.umhSampler0 = sid
			bk.SetUniform(sid, unsafe.Pointer(&s0))
		}

		// submit render state
		// bk.Touch(0)
		bk.Submit(0, id, 0)
	}
	return mr
}

func (mr *MeshRender) SetCamera(camera *Camera) {
	left, right, bottom, top := camera.P()
	p := f32.Ortho2D(left, right, bottom, top)

	// setup uniform
	bk.SetUniform(mr.umhProjection, unsafe.Pointer(&p[0]))
	bk.Submit(0, mr.program, 0)
}

type RenderMesh struct {
	*Mesh
	Matrix []float32
}

// extract render object
func (mr *MeshRender) Extract(visibleObjects []uint32) {

}

// draw
func (mr *MeshRender) Draw(m *Mesh, mat4 *f32.Mat4, depth int32) {
	// state
	bk.SetState(mr.stateFlags, mr.rgba)
	bk.SetTexture(0, mr.umhSampler0, m.textureId, 0)

	// set uniform - mvp
	bk.SetUniform(mr.umhModel, unsafe.Pointer(&mat4[0]))

	// set vertex
	bk.SetVertexBuffer(0, m.VertexId, uint32(m.FirstVertex), uint32(m.NumVertex))
	bk.SetIndexBuffer(m.IndexId, uint32(m.FirstIndex), uint32(m.NumIndex))
	//
	bk.Submit(0, mr.program, depth)
}
