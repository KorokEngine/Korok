package effect

import (
	"testing"
	"unsafe"
)

func TestPool(t *testing.T) {
	p := &Pool{Cap: 1024}

	p.AddChan(Life, Position)
	p.Initialize()

	//   Life: Cap * sizeOf(float)
	// + Position: Cap * sizeOf(vec2)
	size := 1024 * 4 + 1024 * 8
	if p.Size() != size {
		t.Error("fail to compute ParticleSize")
	}

	life := p.Field(Life).(Channel_f32)
	pose := p.Field(Position).(Channel_v2)

	if life == nil || len(life) != 1024 {
		t.Error("fail to cast Life")
	}

	if pose == nil || len(pose) != 1024 {
		t.Error("fail to cast Position")
	}

	d := uintptr(unsafe.Pointer(&pose[0])) - uintptr(unsafe.Pointer(&life[0]))
	if int(d) != 1024 * 4 { // Cap * sizeOf(float)
		t.Error("fail to alloc Life and position")
	}
}
