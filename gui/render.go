package gui

import (
	"korok.io/korok/gfx"
	"unsafe"
	"github.com/go-gl/mathgl/mgl32"
	"korok.io/korok/gfx/bk"
)


type UISystem struct {
	*gfx.MeshRender
	context []*Context
	mesh gfx.Mesh
}

func NewUISystem(render *gfx.MeshRender) *UISystem {
	ui :=  &UISystem{
		MeshRender: render,
	}
	// &m.d.VtxBuffer[0]
	if id, _ := bk.R.AllocVertexBuffer(bk.Memory{nil, 4000*20}, 20); id != bk.InvalidId {
		ui.mesh.VertexId = id
	}
	// &m.d.IdxBuffer[0]
	if id, _ := bk.R.AllocIndexBuffer(bk.Memory{nil, 4000*6}); id != bk.InvalidId {
		ui.mesh.IndexId = id
	}
	return ui
}

func (ui *UISystem) RegisterContext(c *Context) {
	ui.context = append(ui.context, c)
}

// Loop!!
func (ui *UISystem) Draw(dt float32) {
	// --------------
	for _, c := range ui.context {
		dl := &c.DrawList

		// 1. update buffer
		iSize, vSize := dl.Size()
		vBuffer := (*[4*10000]gfx.PosTexColorVertex)(unsafe.Pointer(&dl.VtxBuffer[0]))[:vSize]
		iBuffer := (*[6*10000]uint16)(unsafe.Pointer(&dl.IdxBuffer[0]))[:iSize]

		ui.mesh.SetVertex(vBuffer)
		ui.mesh.SetIndex(iBuffer)
		ui.mesh.Update()

		// 2. draw command
		firstIndex := uint16(0)
		for _, cmd := range dl.Commands() {
			ui.mesh.FirstVertex = 0
			ui.mesh.NumVertex = uint16(len(vBuffer))
			ui.mesh.FirstIndex = firstIndex
			ui.mesh.NumIndex = uint16(cmd.ElemCount)
			firstIndex += uint16(cmd.ElemCount)
			ui.mesh.TextureId = cmd.TextureId

			ui.MeshRender.Draw(&ui.mesh, &mat4)
		}

		//log.Println("cmds:", dl.Commands())
		//log.Println("vbuffer:", vBuffer)
		//log.Println("ibuffer:", iBuffer)

		// reset drawlist
		dl.Clear()
	}
}

var	mat4 = mgl32.Ident4()

func (ui *UISystem) Destroy() {

}


