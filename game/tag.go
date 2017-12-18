package game

import "korok.io/korok/engi"

/**
标记并分类游戏对象, 在 Tag (Name) 的基础上再加一个 Label，作为二级分类，
在游戏中，很多时候是需要这样的二级分类的。比如: enemy {bullet, ship}
 */
type TagComp struct {
	Name, Label string
}

// TODO 如何高效的存储和查找tag数据？
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


