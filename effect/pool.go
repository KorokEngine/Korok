package effect

import (
	"unsafe"
	"github.com/go-gl/mathgl/mgl32"
)

// A description of channel used in pool
type ChanFiled struct {
	Type ChanType
	Name string
}

var (
	Life = ChanFiled{Type:ChanF32, Name:"life"}
	Size = ChanFiled{Type:ChanF32, Name:"size"}

	Color = ChanFiled{Type:ChanV4, Name:"color"}

	Position = ChanFiled{Type:ChanV2, Name:"position"}
	Velocity = ChanFiled{Type:ChanV2, Name:"velocity"}

	Direction = ChanFiled{Type:ChanV2}
	Rotation = ChanFiled{Type:ChanF32}
	// and todo more
)

// A Pool represent a particle-pool
type Pool struct {
	fields []ChanFiled
	chans  map[ChanFiled]uintptr
	cap int
}

// 通过通道描述，可以确定要申请多大的内存，生成那些通道
func (p *Pool) AddChan(fields ...ChanFiled) {
	p.fields = append(p.fields, fields...)
}

func (p *Pool) Initialize() {
	p.chans = make(map[ChanFiled]uintptr)
	size := p.Size()
	block := make([]byte, size)
	mem := uintptr(unsafe.Pointer(&block[0]))

	var offset uintptr
	var len = uintptr(p.cap)
	for _, f := range p.fields {
		p.chans[f] = mem + offset
		offset += len * sizeOf(f.Type)
	}
}

func sizeOf(t ChanType) (size uintptr) {
	switch t {
	case ChanF32:
		size = 4
	case ChanV2:
		size = 8
	case ChanV4:
		size = 16
	}
	return
}

func (p *Pool) Size() (size int) {
	for _, f := range p.fields {
		size += int(sizeOf(f.Type)) * p.cap
	}
	return
}

func (p *Pool) Field(t ChanFiled) (array interface{}) {
	mem := unsafe.Pointer(p.chans[t])
	switch t.Type {
	case ChanF32:
		array = channel_f32((*[1<<16]float32)(unsafe.Pointer(mem))[:p.cap])
	case ChanV2:
		array = channel_v2((*[1<<16]mgl32.Vec2)(unsafe.Pointer(mem))[:p.cap])
	case ChanV4:
		array = channel_v4((*[1<<16]mgl32.Vec4)(unsafe.Pointer(mem))[:p.cap])
	}
	return
}

// TODO
func (p *Pool) Swap(src, dst int) {

}
