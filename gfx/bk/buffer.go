package bk

import(
	"github.com/go-gl/gl/v3.2-core/gl"
	"unsafe"
	"log"
)

// TODO 设计 format 方便导出成 attribute -layout
type Format struct {

}

var Format_POS_COLOR_UV = Format{}
var Format_POS_COLOR    = Format{}
var Format_POS_UV  		= Format{}

type Buffer struct {
	Id uint32

	F Format
	T uint32

	Type uint32
	Count int32
}

func NewArrayBuffer(format Format) Buffer {
	b := Buffer{}
	gl.GenBuffers(1, &b.Id)
	b.F = format
	b.T = gl.ARRAY_BUFFER
	b.Count = 6
	return b
}

func (b *Buffer) Update(data unsafe.Pointer, size int) {
	gl.BindBuffer(b.T, b.Id)
	gl.BufferData(b.T, size, data, gl.STATIC_DRAW)

	// TODO 检测数据的合法性!!
}

func (b *Buffer) Delete() {
	gl.DeleteBuffers(1, &b.Id)
}

func NewIndexBuffer() Buffer {
	return Buffer{
		T: gl.ELEMENT_ARRAY_BUFFER,
	}
}

type IndexBuffer struct {
	Id 		uint32
	size 	uint32
	flags 	uint16
}

func (ib *IndexBuffer) create(size uint32, data unsafe.Pointer, flags uint16) {
	ib.size = size
	ib.flags = flags

	gl.GenBuffers(1, &ib.Id)

	if 0 == ib.Id {
		log.Println("Failed to generate buffer id.")
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.Id)
	if data == nil {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(size), nil, gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(size), data, gl.STATIC_DRAW)
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

/// discard=false
func (ib *IndexBuffer) update(offset uint32, size uint32, data interface{}, discard bool) {
	if 0 == ib.Id {
		log.Println("Updating invalid index buffer.")
	}

	if discard {
		// orphan buffer
		ib.destroy()
		ib.create(ib.size, nil, ib.flags)
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.Id)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(offset), int(size), gl.Ptr(data))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (ib *IndexBuffer) destroy() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.DeleteBuffers(1, &ib.Id)
}

type VertexBuffer struct {
	Id     uint32
	target uint32
	size   uint32
	layout uint16 	// Stride | Offset
}

/// draw indirect >= es 3.0 or gl 4.0
func (vb *VertexBuffer) create(size uint32, data unsafe.Pointer, layout uint16, flags uint16) {
	vb.size = size
	vb.layout = layout
	vb.target = gl.ARRAY_BUFFER

	gl.GenBuffers(1, &vb.Id)
	if vb.Id == 0 {
		log.Println("Failed to generate buffer id")
	}
	gl.BindBuffer(vb.target, vb.Id)
	if data == nil {
		gl.BufferData(vb.target, int(size), data, gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(vb.target, int(size), data, gl.STATIC_DRAW)
	}
	gl.BindBuffer(vb.target, 0)
}

/// discard = false
func (vb *VertexBuffer) update(offset uint32, size uint32, data unsafe.Pointer, discard bool) {
	if vb.Id == 0 {
		log.Println("Updating invalid vertex buffer")
	}

	if discard {
		vb.destroy()
		vb.create(vb.size, nil, vb.layout, 0)
	}

	gl.BindBuffer(vb.target, vb.Id)
	gl.BufferSubData(vb.target, int(offset), int(size), data)
	gl.BindBuffer(vb.target, 0)
}

func (vb *VertexBuffer) destroy() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0 )
	gl.DeleteBuffers(1, &vb.Id)
}


