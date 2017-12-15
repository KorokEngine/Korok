package game

import "korok.io/korok/engi"

/**
	游戏对象绑定脚本/行为
 */
type ScriptComp interface {
	Init(id uint32)

	Update(dt float32)

	Destroy()
}

type ScriptTable struct {
	comps []ScriptComp
	_map   map[uint32]int
	index, cap int
}

func NewScriptTable(cap int) *ScriptTable {
	return &ScriptTable{cap: cap, _map: make(map[uint32]int)}
}

func (st *ScriptTable) NewComp(entity engi.Entity) (sc ScriptComp) {
	if size := len(st.comps); st.index >= size {
		st.comps = scriptResize(st.comps, size + 64)
	}
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		return st.comps[v]
	}
	sc = st.comps[st.index]
	st._map[ei] = st.index
	st.index ++
	return
}


func (st *ScriptTable) Comp(entity engi.Entity) (sc ScriptComp) {
	if v, ok := st._map[entity.Index()]; ok {
		sc = st.comps[v]
	}
	return
}

func (st *ScriptTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		if tail := st.index -1; v != tail && tail > 0 {
			st.comps[v] = st.comps[tail]
			// remap index TODO bug fix
			//tComp := &st.comps[tail]
			//ei := tComp.Entity.Index()
			st._map[ei] = v
			//tComp.Entity = 0
		} else {
			st.comps[tail] = nil
		}

		st.index -= 1
		delete(st._map, ei)
	}
}

func scriptResize(slice []ScriptComp, size int) []ScriptComp {
	newSlice := make([]ScriptComp, size)
	copy(newSlice, slice)
	return newSlice
}

type ScriptSystem struct {
	*ScriptTable
}

func (ss *ScriptSystem) Update(dt float32) {
	N := ss.ScriptTable.index
	comps := ss.ScriptTable.comps
	for i := 0; i < N; i++ {
		comps[i].Update(dt)
	}
}



