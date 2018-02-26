package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

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
	umh_P  uint16 // Projection
	umh_M  uint16 // Model
	umh_S0 uint16 // Sampler0
}

func NewMeshRender(vsh, fsh string) *MeshRender {
	mr := new(MeshRender)
	// blend func
	mr.stateFlags |= bk.ST_BLEND.ALPHA_NON_PREMULTIPLIED

	// setup shader
	if id, sh := bk.R.AllocShader(vsh, fsh); id != bk.InvalidId {
		mr.program = id
		sh.Use()

		// setup attribute
		sh.AddAttributeBinding("xyuv\x00", 0, P4C4[0])
		sh.AddAttributeBinding("rgba\x00", 0, P4C4[1])

		p := mgl32.Ortho2D(0, 480, 0, 320)
		m := mgl32.Translate3D(240, 160, 0)
		s0 := int32(0)

		// setup uniform
		if pid, _ := bk.R.AllocUniform(id, "proj\x00", bk.UniformMat4, 1); pid != bk.InvalidId {
			mr.umh_P = pid
			bk.SetUniform(pid, unsafe.Pointer(&p[0]))
		}

		if mid, _ := bk.R.AllocUniform(id, "model\x00", bk.UniformMat4, 1); mid != bk.InvalidId {
			mr.umh_M = mid
			bk.SetUniform(mid, unsafe.Pointer(&m[0]))
		}

		if sid,_ := bk.R.AllocUniform(id, "tex\x00", bk.UniformSampler, 1); sid != bk.InvalidId {
			mr.umh_S0 = sid
			bk.SetUniform(sid, unsafe.Pointer(&s0))
		}

		// submit render state
		// bk.Touch(0)
		bk.Submit(0, id, 0)
	}
	return mr
}

func (mr *MeshRender) SetCamera(camera *Camera) {
	left := camera.pos.x - camera.view.w/2
	right := camera.pos.x + camera.view.w/2
	bottom := camera.pos.y - camera.view.h/2
	top := camera.pos.y + camera.view.h/2

	p := mgl32.Ortho2D(left, right, bottom, top)

	// setup uniform
	bk.SetUniform(mr.umh_P, unsafe.Pointer(&p[0]))
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
func (mr *MeshRender) Draw(m *Mesh, mat4 *mgl32.Mat4) {
	// state
	bk.SetState(mr.stateFlags, mr.rgba)
	bk.SetTexture(0, mr.umh_S0, uint16(m.TextureId), 0)

	// set uniform - mvp
	bk.SetUniform(mr.umh_M, unsafe.Pointer(&mat4[0]))

	// set vertex
	bk.SetVertexBuffer(0, m.VertexId, uint32(m.FirstVertex), uint32(m.NumVertex))
	bk.SetIndexBuffer(m.IndexId, uint32(m.FirstIndex), uint32(m.NumIndex))
	//
	bk.Submit(0, mr.program, 0)
}
