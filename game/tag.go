package game

import "korok.io/korok/engi"

/**
标记并分类游戏对象, 在 Tag (Name) 的基础上再加一个 Label，作为二级分类，
在游戏中，很多时候是需要这样的二级分类的。比如: enemy {bullet, ship}

http://bitsquid.blogspot.se/2015/06/allocation-adventures-1-datacomponent.html
http://bitsquid.blogspot.com/2015/06/allocation-adventures-2-arrays-of-arrays.html

关于 string 比较的问题，还需要研究：https://stackoverflow.com/questions/20232976/how-does-go-do-string-comparison

关于 Tag 系统的设计，可以学下一下
 */
type TagComp struct {
	engi.Entity
	Name, Label string
}

// TODO 如何高效的存储和查找tag数据？
type TagTable struct {
	comps []TagComp
	_map   map[uint32]int
	index, cap int

	d map[string][]engi.Entity
}

func NewTagTable(cap int) *TagTable {
	return &TagTable{cap: cap, _map: make(map[uint32]int)}
}

func (tt *TagTable) NewComp(entity engi.Entity) (tc *TagComp) {
	if size := len(tt.comps); tt.index >= size {
		tt.comps = tagResize(tt.comps, size+64)
	}
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		return &tt.comps[v]
	}
	tc = &tt.comps[tt.index]
	tc.Entity = entity
	tt._map[entity.Index()] = tt.index
	tt.index ++
	return
}

func (tt *TagTable) Alive(entity engi.Entity) bool {
	if v, ok := tt._map[entity.Index()]; ok {
		return tt.comps[v].Entity == 0
	}
	return false
}

func (tt *TagTable) Comp(entity engi.Entity) (tc *TagComp) {
	if v, ok := tt._map[entity.Index()]; ok {
		tc = &tt.comps[v]
	}
	return
}

func (tt *TagTable) Delete(entity engi.Entity) (tc *TagComp) {
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		if tail := tt.index-1; v != tail && tail > 0 {
			tt.comps[v] = tt.comps[tail]
			// remap index
			tComp := &tt.comps[tail]
			ei := tComp.Entity.Index()
			tt._map[ei] = v
			tComp.Entity = 0
		} else {
			tt.comps[tail].Entity = 0
		}
		tt.index -= 1
		delete(tt._map, ei)
	}
	return nil
}

// 删除所有属于该标签的元素..
func (tt *TagTable) DeleteTag(entity engi.Entity) {

}

// 这个效率非常低，因为Group一次要对所有的Entity做一次线性的
// 查找..
func (tt *TagTable) Group(tag string) []engi.Entity {
	list := make([]engi.Entity, 0)
	for _, comp := range tt.comps {
		if comp.Name == tag {
			list = append(list, comp.Entity)
		}
	}
	return list
}

func (tt *TagTable) Size() (size, cap int) {
	return tt.index, tt.cap
}

func tagResize(slice []TagComp, size int) []TagComp {
	newSlice := make([]TagComp, size)
	copy(newSlice, slice)
	return newSlice
}

