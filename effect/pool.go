package effect

import (
	"unsafe"
	"korok.io/korok/math/f32"
)

// A description of channel used in pool
type ChanFiled struct {
	Type ChanType
	Name string
}

type block struct {
	ChanFiled
	stride int
	data []byte
}

var (
	Life = ChanFiled{Type:ChanF32, Name:"Life"}
	Size = ChanFiled{Type:ChanF32, Name:"ParticleSize"}
	SizeDelta = ChanFiled{Type:ChanF32, Name:"ParticleSize-delta"}

	Color = ChanFiled{Type:ChanV4, Name:"Color"}
	ColorDelta = ChanFiled{Type:ChanV4, Name:"Color-delta"}

	Position = ChanFiled{Type:ChanV2, Name:"position"}
	PositionStart = ChanFiled{Type:ChanV2, Name:"position-start"}
	Velocity = ChanFiled{Type:ChanV2, Name:"velocity"}

	Speed = ChanFiled{Type:ChanF32, Name:"speed"}
	Direction = ChanFiled{Type:ChanV2, Name:"direction"}
	RadialAcc = ChanFiled{Type:ChanF32, Name:"radial-acc"}
	TangentialAcc = ChanFiled{Type:ChanF32, Name:"tangent-acc"}

	Rotation = ChanFiled{Type:ChanF32, Name:"rotation"}
	RotationDelta = ChanFiled{Type:ChanF32, Name:"rotation-delta"}

	Angle = ChanFiled{Type:ChanF32, Name:"angle"}
	AngleDelta = ChanFiled{Type:ChanF32, Name:"angle-delta"}
	Radius = ChanFiled{Type:ChanF32, Name:"radius"}
	RadiusDelta = ChanFiled{Type:ChanF32, Name:"radius-delta"}
)

// A Pool represent a particle-pool.
type Pool struct {
	blocks []block
	chans  map[ChanFiled]int
	Cap    int
}

// AddChan adds new fields to the pool.
func (p *Pool) AddChan(fields ...ChanFiled) {
	for _, f := range fields {
		p.blocks = append(p.blocks, block{ChanFiled:f})
	}
}

// Initialize the particle pool.
func (p *Pool) Initialize() {
	p.chans = make(map[ChanFiled]int)
	size := p.Size()
	pool := make([]byte, size)

	var (
		mem = uintptr(unsafe.Pointer(&pool[0]))
		offset uintptr
		cap = p.Cap
	)

	for i, b := range p.blocks {
		stride := sizeOf(b.Type)
		p.blocks[i].stride = stride
		p.blocks[i].data = (*[1<<16]byte)(unsafe.Pointer(mem + offset))[:cap*stride]
		offset += uintptr(cap * stride)
		p.chans[b.ChanFiled] = i
	}
}

func sizeOf(t ChanType) (size int) {
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

// Size return the ParticleSize (in bytes) of the pool.
func (p *Pool) Size() (size int) {
	for _, f := range p.blocks {
		size += int(sizeOf(f.Type)) * p.Cap
	}
	return
}

// Field returns pointer of the filed in the pool.
func (p *Pool) Field(t ChanFiled) (array interface{}) {
	block := p.blocks[p.chans[t]]
	mem := unsafe.Pointer(&block.data[0])
	switch t.Type {
	case ChanF32:
		array = Channel_f32((*[1<<16]float32)(mem)[:p.Cap])
	case ChanV2:
		array = Channel_v2((*[1<<16]f32.Vec2)(mem)[:p.Cap])
	case ChanV4:
		array = Channel_v4((*[1<<16]f32.Vec4)(mem)[:p.Cap])
	}
	return
}

// Swap swap all the field defined in the pool.
func (p *Pool) Swap(dst, src int) {
	for _, b := range p.blocks {
		stride := 4 << uint(b.Type)
		i, j := dst * stride, src * stride
		copy(b.data[i:], b.data[j:j+stride])
	}
}
