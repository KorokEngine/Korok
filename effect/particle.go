package effect

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx"
)

// particle component
type ParticleComp struct {
	engi.Entity
	sim Simulator

	tex *gfx.SubTex
	color uint32
	scale float32
}

func (ec *ParticleComp) SetSimulator(sim Simulator) {
	ec.sim = sim
}

func (ec *ParticleComp) SetTexture(tex *gfx.SubTex) {
	ec.tex = tex
}

func (ec *ParticleComp) SetColor(color uint32) {
	ec.color = color
}

func (ec *ParticleComp) SetScale(scale float32) {
	ec.scale = scale
}

func (ec *ParticleComp) Play() {

}

func (ec *ParticleComp) Pause() {

}

func (ec *ParticleComp) Resume() {

}

func (ec *ParticleComp) Stop() {

}

// component manager
type ParticleSystemTable struct {
	comps []ParticleComp
	_map   map[uint32]int
	index, cap int
}

func NewEffectTable(cap int) *ParticleSystemTable {
	return &ParticleSystemTable{
		_map:make(map[uint32]int),
		cap:cap,
	}
}

func (et *ParticleSystemTable) NewComp(entity engi.Entity) (ec *ParticleComp) {
	if size := len(et.comps); et.index >= size {
		et.comps = effectCompResize(et.comps, size + 64)
	}
	ei := entity.Index()
	if v, ok := et._map[ei]; ok {
		return &et.comps[v]
	}
	ec = &et.comps[et.index]
	ec.Entity = entity
	et._map[ei] = et.index
	et.index ++
	return
}

func (et *ParticleSystemTable) Alive(entity engi.Entity) bool {
	if v, ok := et._map[entity.Index()]; ok {
		return et.comps[v].Entity != 0
	}
	return false
}

func (et *ParticleSystemTable) Comp(entity engi.Entity) (ec *ParticleComp) {
	if v, ok := et._map[entity.Index()]; ok {
		ec = &et.comps[v]
	}
	return
}

func (et *ParticleSystemTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := et._map[ei]; ok {
		if tail := et.index -1; v != tail && tail > 0 {
			et.comps[v] = et.comps[tail]
			// remap index
			tComp := &et.comps[tail]
			ei := tComp.Entity.Index()
			et._map[ei] = v
			tComp.Entity = 0
		} else {
			et.comps[tail].Entity = 0
		}

		et.index -= 1
		delete(et._map, ei)
	}
}

func (et *ParticleSystemTable) Size() (size, cap int) {
	return et.index, et.cap
}

func effectCompResize(slice []ParticleComp, size int) []ParticleComp {
	newSlice := make([]ParticleComp, size)
	copy(newSlice, slice)
	return newSlice
}

type ParticleRenderFeature struct {
	mr *gfx.MeshRender

	et *ParticleSystemTable
	xt *gfx.TransformTable
}


// 此处初始化所有的依赖
func (erf *ParticleRenderFeature) Register(rs *gfx.RenderSystem) {
	// init render
	for _, r := range rs.RenderList {
		switch mr := r.(type) {
		case *gfx.MeshRender:
			erf.mr = mr; break
		}
	}
	// init table
	for _, t := range rs.TableList {
		switch table := t.(type){
		case *ParticleSystemTable:
			erf.et = table
		case *gfx.TransformTable:
			erf.xt = table
		}
	}
	// add new feature
	rs.Accept(erf)
}

// 此处执行渲染
// BatchRender 需要的是一组排过序的渲染对象！！！
func (erf *ParticleRenderFeature) Draw(filter []engi.Entity) {
	//xt, mt, n := erf.xt, erf.et, erf.et.index
	//mr := erf.mr
	//
	//for i := 0; i < n; i++ {
	//	mesh := &mt.comps[i]
	//	entity := mesh.Entity
	//	xform  := xt.Comp(entity)
	//
	//	// TODO transform!!
	//	mat4.Set(0, 3, xform.world.Position[0])
	//	mat4.Set(1, 3, xform.world.Position[1])
	//
	//	mr.Draw(&mesh.Mesh, &mat4)
	//}
}

//
type EffectRenderContext struct {

}

