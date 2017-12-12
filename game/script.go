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
	_comps []ScriptComp
	_index uint32
	_map   map[int]uint32
}


func (st *ScriptTable) NewComp(entity engi.Entity) (sc *ScriptComp) {
	sc = &st._comps[st._index]
	st._map[int(entity)] = st._index
	st._index ++
	return
}


func (st *ScriptTable) Comp(entity engi.Entity) (sc *ScriptComp) {
	if v, ok := st._map[int(entity)]; ok {
		sc = &st._comps[v]
	}
	return
}

func (*ScriptTable) Delete(entity engi.Entity) {
	// todo
}


type ScriptSystem struct {
	*ScriptTable
}

func (ss *ScriptSystem) Update(dt float32) {
	N := ss.ScriptTable._index
	comps := ss.ScriptTable._comps
	for i := uint32(0); i < N; i++ {
		comps[i].Update(dt)
	}
}



