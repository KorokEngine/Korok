package gfx

// 因为大部分图元应该都是不动的，那么可以
// 以相机为
//  cull with camera

type CompRef struct {
	SortKey
	Type int32
	*RenderComp
}

type CullSystem interface {
	Cull(comps []RenderComp, camera Camera) []CompRef
}

type cullSystem struct {

}

// 输入一组ID，给出一组ID
func (cs *cullSystem) Cull(comps []RenderComp, camera *Camera) []CompRef {
	refs := make([]CompRef, len(comps))
	for i := range comps {
		refs[i].SortKey = comps[i].Sort
		refs[i].RenderComp = &comps[i]
	}
	return refs
}
