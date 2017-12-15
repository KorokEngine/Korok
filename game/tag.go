package game

import "korok.io/korok/engi"

/**
标记并分类游戏对象
 */
type TagComp struct {
	Name string
}

type TagTable struct {
	comps []TagComp
	_map   map[uint32]int
	index, cap int

	d map[string][]engi.Entity
}

func (tt *TagTable) NewComp(entity engi.Entity) (tc *TagComp) {
	tc = &tt.comps[tt.index]
	tt._map[entity.Index()] = tt.index
	tt.index ++
	return
}

func (tt *TagTable) Comp(entity engi.Entity) (tc *TagComp) {
	if v, ok := tt._map[entity.Index()]; ok {
		tc = &tt.comps[v]
	}
	return
}

// todo
func (tt *TagTable) Delete(entity engi.Entity) (tc *TagComp) {
	return nil
}

func (tt *TagTable) Group(tag string) []engi.Entity {
	return nil
}

func (tt *TagTable) Size() (size, cap int) {
	return 0, 0
}


