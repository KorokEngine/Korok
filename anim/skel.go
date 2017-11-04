package anim

import "korok/engi"

// 当前动画模块的性能是比较低的，需要一次重写
// 同时也会用 Comp/Table/System 架构重新组织功能

type SkeletonComp struct {

}

type SkeletonTable struct {
	_comps []SkeletonComp
	_index uint32
	_map   map[int]uint32

}

func (st *SkeletonTable) NewComp(entity engi.Entity) (sc *SkeletonComp) {
	sc = &st._comps[st._index]
	st._map[int(entity)] = st._index
	st._index ++
	return
}

func (st *SkeletonTable) Comp(entity engi.Entity) (sc *SkeletonComp) {
	if v, ok := st._map[int(entity)]; ok {
		sc = &st._comps[v]
	}
	return
}

// TODO impl
func (st *SkeletonTable) Delete(entity engi.Entity) (sc *SkeletonComp) {
	return nil
}




