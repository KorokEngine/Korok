package frame

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx"
)

// Sprite Animation Component
type FlipbookComp struct {
	engi.Entity
	define string
	dt, rate float32
	ii int
	running bool
	loop bool
}

func (fb *FlipbookComp) Play(name string) {
	fb.define = name
	fb.running = true
}

func (fb *FlipbookComp) Stop() {
	fb.running = false
}

func (fb *FlipbookComp) SetAnimation(name string) {
	fb.define = name
}

func (fb *FlipbookComp) Loop() bool {
	return fb.loop
}

func (fb *FlipbookComp) SetLoop(v bool) {
	fb.loop = v
}

func (fb *FlipbookComp) Rate() float32 {
	return fb.rate
}

func (fb *FlipbookComp) SetRate(r float32) {
	fb.rate = r
}

// Sprite Animation Table
type FlipbookTable struct {
	comps []FlipbookComp
	_map   map[uint32]int
	index, cap int
}

func NewFlipbookTable(cap int) *FlipbookTable {
	return &FlipbookTable{
		cap:cap,
		_map:make(map[uint32]int),
	}
}

func (t *FlipbookTable) NewComp(entity engi.Entity) (am *FlipbookComp) {
	if size := len(t.comps); t.index >= size {
		t.comps = tableResize(t.comps, size + gfx.STEP)
	}
	ei := entity.Index()
	if v, ok := t._map[ei]; ok {
		am = &t.comps[v]
		return
	}
	am = &t.comps[t.index]
	am.Entity = entity
	t._map[ei] = t.index
	t.index ++
	return
}

func (t *FlipbookTable) Alive(entity engi.Entity) bool {
	ei := entity.Index()
	if v, ok := t._map[ei]; ok {
		return t.comps[v].Entity != 0
	}
	return false
}

func (t *FlipbookTable) Comp(entity engi.Entity) (sc *FlipbookComp) {
	ei := entity.Index()
	if v, ok := t._map[ei]; ok {
		sc = &t.comps[v]
	}
	return
}

func (t *FlipbookTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := t._map[ei]; ok {
		if tail := t.index -1; v != tail && tail > 0 {
			t.comps[v] = t.comps[tail]
			// remap index
			tComp := &t.comps[tail]
			ei := tComp.Entity.Index()
			t._map[ei] = v
			tComp.Entity = 0
		} else {
			t.comps[tail].Entity = 0
		}

		t.index -= 1
		delete(t._map, ei)
	}
}

func (t *FlipbookTable) Size() (size, cap int) {
	return t.index, t.cap
}

func (t *FlipbookTable) Destroy() {
	t.comps = make([]FlipbookComp, 0)
	t._map = make(map[uint32]int)
	t.index = 0
}

func tableResize(slice []FlipbookComp, size int) []FlipbookComp {
	newSlice := make([]FlipbookComp, size)
	copy(newSlice, slice)
	return newSlice
}




