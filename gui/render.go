package gui

import (
	"korok.io/korok/gfx"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/math/f32"

	"unsafe"
	"korok.io/korok/gfx/dbg"
)

type UIRenderFeature struct {
	id int
	*gfx.MeshRender
	*DrawList
	*gfx.Camera

	Buffer struct{
		iid, vid uint16
		isz, vsz int
		vertex *bk.VertexBuffer
		index *bk.IndexBuffer
	}
}

func (f *UIRenderFeature) SetDrawList(dl *DrawList) {
	f.DrawList = dl
}

func (f *UIRenderFeature) Register(rs *gfx.RenderSystem) {
	f.Camera = &rs.MainCamera
	// init render
	for _, r := range rs.RenderList {
		switch render := r.(type) {
		case *gfx.MeshRender:
			f.MeshRender = render; break
		}
	}
	// add new feature, use the index as id
	f.id = rs.Accept(f)
	f.DrawList = &gContext.DrawList
	f.Buffer.vid = bk.InvalidId
	f.Buffer.iid = bk.InvalidId
}

func (f *UIRenderFeature) Extract(v *gfx.View) {
	if dl := f.DrawList; !dl.Empty() {
		fi := uint32(f.id)<<16
		for i, cmd := range dl.Commands() {
			sid := gfx.PackSortId(cmd.zOrder, 0)
			val := fi + uint32(i)
			v.RenderNodes = append(v.RenderNodes, gfx.SortObject{SortId:sid, Value:val})
		}
	}
}

// TODO scale.. and rotatation..
func (f *UIRenderFeature) Draw(nodes gfx.RenderNodes) {
	var (
		sw, _ = f.Camera.Screen()
		x, y, w, h = f.Camera.View()
		dx = x - w/2
		dy = y - h/2
		sk = w/sw
	)

	// setup buffer
	isz, vsz := f.DrawList.Size()
	f.allocBuffer(isz, vsz)
	mesh := &gfx.Mesh{
		IndexId:f.Buffer.iid,
		VertexId:f.Buffer.vid,
	}
	mat4 := &f32.Mat4{sk, 0, 0, 0, 0, sk, 0, 0, dx, dy, 1, 0, 0, 0, 0, 1}
	commands := f.DrawList.Commands()
	for _, node := range nodes {
		index := node.Value&0xFFFF
		cmd := commands[index]

		mesh.FirstVertex = 0
		mesh.NumVertex = uint16(vsz)
		mesh.FirstIndex = cmd.FirstIndex
		mesh.NumIndex = cmd.ElemCount
		mesh.SetTexture(cmd.TextureId)

		f.MeshRender.Draw(mesh, mat4, int32(cmd.zOrder))
	}
	f.Buffer.vertex.Update(0, uint32(vsz*20), unsafe.Pointer(&f.DrawList.VtxBuffer[0]),false)
	f.Buffer.index.Update(0, uint32(isz*2), unsafe.Pointer(&f.DrawList.IdxBuffer[0]), false)
	f.DrawList.Clear()

	dbg.Hud("gui DrawList: %d, %d", isz, vsz)
}

func (f *UIRenderFeature) allocBuffer(isz, vsz int) {
	if isz > f.Buffer.isz {
		if iid := f.Buffer.iid; iid != bk.InvalidId {
			bk.R.Free(iid)
		}
		{
			isz--
			isz |= isz >> 1
			isz |= isz >> 2
			isz |= isz >> 3
			isz |= isz >> 8
			isz |= isz >> 16
			isz++
		}
		id, ib := bk.R.AllocIndexBuffer(bk.Memory{nil, uint32(isz)*2})
		f.Buffer.iid = id
		f.Buffer.isz = isz
		f.Buffer.index = ib
	}

	vid, _, vb := gfx.Context.TempVertexBuffer(vsz, 20)
	f.Buffer.vid = vid
	f.Buffer.vertex = vb
}

