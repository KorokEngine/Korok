package game

import "korok/engi"

type TagComp struct {
	Name string
}

type TagTable struct {
	_comps []TagComp
	_index uint32
	_map   map[int]uint32
}

func (tt *TagTable) NewComp(entity engi.Entity) (tc *TagComp) {
	tc = &tt._comps[tt._index]
	tt._map[int(entity)] = tt._index
	tt._index ++
	return
}

func (tt *TagTable) Comp(entity engi.Entity) (tc *TagComp) {
	if v, ok := tt._map[int(entity)]; ok {
		tc = &tt._comps[v]
	}
	return
}

// todo
func (tt *TagTable) Delete(entity engi.Entity) (tc *TagComp) {
	return nil
}

type TagSystem struct {
	*TagTable
	// any ?
}


