package bk

import (
	"korok.io/korok/hid/gl"

	"log"
	"unsafe"
	"errors"
)

type IndexBuffer struct {
	Id    uint32
	size  uint32
	flags uint16
}

func (ib *IndexBuffer) Create(size uint32, data unsafe.Pointer, flags uint16) error{
	ib.size = size
	ib.flags = flags

	gl.GenBuffers(1, &ib.Id)

	if 0 == ib.Id {
		return errors.New("failed to generate buffer id")
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.Id)
	if data == nil {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(size), nil, gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(size), data, gl.STATIC_DRAW)
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return nil
}

/// discard=false
func (ib *IndexBuffer) Update(offset uint32, size uint32, data unsafe.Pointer, discard bool) {
	if 0 == ib.Id {
		log.Println("updating invalid index buffer")
	}

	if discard {
		// orphan buffer
		ib.Destroy()
		ib.Create(ib.size, nil, ib.flags)
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.Id)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(offset), int(size), data)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (ib *IndexBuffer) Destroy() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &ib.Id)
}

type VertexBuffer struct {
	Id     uint32
	target uint32
	size   uint32
	layout uint16 // Stride | Offset
}

/// draw indirect >= es 3.0 or gl 4.0
func (vb *VertexBuffer) Create(size uint32, data unsafe.Pointer, layout uint16, flags uint16) error{
	vb.size = size
	vb.layout = layout
	vb.target = gl.ARRAY_BUFFER

	gl.GenBuffers(1, &vb.Id)
	if vb.Id == 0 {
		return errors.New("failed to generate buffer id")
	}
	gl.BindBuffer(vb.target, vb.Id)
	if data == nil {
		gl.BufferData(vb.target, int(size), data, gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(vb.target, int(size), data, gl.STATIC_DRAW)
	}
	gl.BindBuffer(vb.target, 0)
	return nil
}

/// discard = false
func (vb *VertexBuffer) Update(offset uint32, size uint32, data unsafe.Pointer, discard bool) {
	if vb.Id == 0 {
		log.Println("Updating invalid vertex buffer")
	}

	if discard {
		vb.Destroy()
		vb.Create(vb.size, nil, vb.layout, 0)
	}

	gl.BindBuffer(vb.target, vb.Id)
	gl.BufferSubData(vb.target, int(offset), int(size), data)
	gl.BindBuffer(vb.target, 0)
}

func (vb *VertexBuffer) Destroy() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &vb.Id)
}
