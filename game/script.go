package game

import "korok.io/korok/engi"

/**
	游戏对象绑定脚本/行为
 */

type Script interface {
	Init()

	Update(dt float32)

	Destroy()
}

type ScriptComp struct {
	engi.Entity
	Script
}

func (sc *ScriptComp) SetScript(script Script) {
	sc.Script = script
}

type ScriptTable struct {
	comps []ScriptComp
	_map   map[uint32]int
	index, cap int
}

func NewScriptTable(cap int) *ScriptTable {
	return &ScriptTable{cap: cap, _map: make(map[uint32]int)}
}

func (st *ScriptTable) NewComp(entity engi.Entity, script Script) (sc *ScriptComp) {
	if size := len(st.comps); st.index >= size {
		st.comps = scriptResize(st.comps, size + 64)
	}
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		return &st.comps[v]
	}
	sc = &st.comps[st.index]
	sc.Entity = entity
	sc.Script = script
	st._map[ei] = st.index
	st.index ++
	return
}

func (st *ScriptTable) Alive(entity engi.Entity) bool {
	if v, ok := st._map[entity.Index()]; ok {
		return st.comps[v].Entity == 0
	}
	return false
}

func (st *ScriptTable) Comp(entity engi.Entity) (sc *ScriptComp) {
	if v, ok := st._map[entity.Index()]; ok {
		sc = &st.comps[v]
	}
	return
}

func (st *ScriptTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := st._map[ei]; ok {
		if tail := st.index -1; v != tail && tail > 0 {
			st.comps[v] = st.comps[tail]
			// remap index
			tComp := &st.comps[tail]
			ei := tComp.Entity.Index()
			st._map[ei] = v
			tComp.Entity = 0
		} else {
			st.comps[tail].Entity = 0
		}

		st.index -= 1
		delete(st._map, ei)
	}
}

func (st *ScriptTable) Size() (size, cap int) {
	return st.index, st.cap
}

func scriptResize(slice []ScriptComp, size int) []ScriptComp {
	newSlice := make([]ScriptComp, size)
	copy(newSlice, slice)
	return newSlice
}

type ScriptSystem struct {
	*ScriptTable
}

func NewScriptSystem() *ScriptSystem {
	return &ScriptSystem{}
}
func (ss *ScriptSystem) RequireTable(tables []interface{}) {
	for _, t := range tables {
		switch table := t.(type) {
		case *ScriptTable:
			ss.ScriptTable = table
		}
	}
}

func (ss *ScriptSystem) Update(dt float32) {
	N := ss.ScriptTable.index
	comps := ss.ScriptTable.comps
	for i := 0; i < N; i++ {
		if script := comps[i].Script; script != nil {
			script.Update(dt)
		}
	}
}



