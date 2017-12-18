package effect

import (
	"testing"
	"unsafe"
)

func TestPool(t *testing.T) {
	p := &Pool{cap: 1024}

	p.AddChan(Life, Position)
	p.Initialize()

	//   life: cap * sizeOf(float)
	// + pose: cap * sizeOf(vec2)
	size := 1024 * 4 + 1024 * 8
	if p.Size() != size {
		t.Error("fail to compute size")
	}

	life := p.Field(Life).(channel_f32)
	pose := p.Field(Position).(channel_v2)

	if life == nil || len(life) != 1024 {
		t.Error("fail to cast life")
	}

	if pose == nil || len(pose) != 1024 {
		t.Error("fail to cast pose")
	}

	d := uintptr(unsafe.Pointer(&pose[0])) - uintptr(unsafe.Pointer(&life[0]))
	if int(d) != 1024 * 4 { // cap * sizeOf(float)
		t.Error("fail to alloc life and position")
	}
}
