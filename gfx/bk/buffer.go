package gfx

import(
	"github.com/go-gl/gl/v3.2-core/gl"
	"unsafe"
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