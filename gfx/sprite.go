package gfx

/// SpriteComp & SpriteTable
/// Usually, sprite can be rendered with a BatchRenderer

type SpriteComp struct {
	*SubTex

	Scale float32
	Color uint32
}

type SpriteTable struct {
	_comps [1024]SpriteComp
	_index uint32
	_map   [1024]int
}

func (st *SpriteTable) NewComp(tex *SubTex) (sc *SpriteComp) {
	sc = &st._comps[st._index]
	sc.SubTex = tex
	st._index ++
	return
}

func (st *SpriteTable) Comp(id uint32) *SpriteComp {
	return &st._comps[st._map[id]]
}




