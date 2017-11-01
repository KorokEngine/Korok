package gfx

import (
	"korok/engi"
)

/// SpriteComp & SpriteTable
/// Usually, sprite can be rendered with a BatchRenderer

type SpriteComp struct {
	engi.Entity
	*SubTex

	Scale float32
	Color uint32
}

type SpriteTable struct {
	_comps [1024]SpriteComp
	_index uint32
	_map   [1024]uint32
}

func (st *SpriteTable) NewComp(entity engi.Entity, tex *SubTex) (sc *SpriteComp) {
	sc = &st._comps[st._index]
	sc.SubTex = tex
	sc.Entity = entity
	st._map[entity] = st._index
	st._index ++
	return
}

func (st *SpriteTable) Comp(id uint32) *SpriteComp {
	return &st._comps[st._map[id]]
}




