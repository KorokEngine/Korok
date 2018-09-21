package gfx

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/engi"
	"sort"
	"korok.io/korok/gfx/dbg"
)

type RenderType int32

type Render interface {
	SetCamera(camera *Camera)
}

// 适合于渲染系统访问的表达方式.
// 其实不必这么麻烦，我们在 RenderFeature里面涉及一个 Extract 步骤，构建一个渲染列表，然后再绘制即可.
// 这个列表需要动态构建
//type renderObject struct {
//
//}

// Shared properties by most render object.
type RenderObject struct {
	engi.Entity
	Size f32.Vec2
	Center f32.Vec2
	ZOrder int16
	BatchId uint16
}

type RenderNodes []SortObject

type View struct {
	*Camera
	RenderNodes
}

// 传入参数是经过可见性系统筛选后的 Entity，这是一个很小的数组，可以
// 直接传给各个 RenderFeature 来做可见性判断.
type RenderFeature interface {
	Extract(v *View)
	Draw(nodes RenderNodes)
	Flush()
}

// 所有的Table和Render都在此管理
// 其它的 RenderFeature 在此提取依赖
// 这样的话， RenderSystem 就沦为一个管理 RenderFeature 和 Table 的地方
// 它们之间也会存在各种组合...
type RenderSystem struct {
	MainCamera Camera
	View

	// shortcut for TransformTable
	xfs *TransformTable

	// visibility test
	V VisibilitySystem

	// render-data
	TableList []interface{}

	// render
	RenderList []Render

	// feature knows how to use render-data and render
	FeatureList []RenderFeature
}

func (th *RenderSystem) RequireTable(tables []interface{}) {
	th.TableList = tables
	for _, table := range tables {
		if t, ok := table.(*TransformTable); ok {
			th.xfs = t; break
		}
	}
}

func (th *RenderSystem) Accept(rf RenderFeature) (index int){
	index = len(th.FeatureList)
	th.FeatureList = append(th.FeatureList, rf)
	return
}

// register type-render
func (th *RenderSystem) RegisterRender(t RenderType, render Render) {
	th.RenderList = append(th.RenderList, render)
}

func (th *RenderSystem) Update(dt float32) {
	// update camera
	if c := &th.MainCamera; c.follow != engi.Ghost {
		xf := th.xfs.Comp(c.follow)
		p  := xf.Position()
		dx := (p[0]-c.mat.x)*.1
		dy := (p[1]-c.mat.y)*.1
		c.MoveBy(dx, dy)
	}

	// main camera
	for _, r := range th.RenderList {
		r.SetCamera(&th.MainCamera)
	}
	if dbg.DEBUG != dbg.None {
		dbg.SetCamera(th.MainCamera.View())
	}

	// build view
	v := th.View

	// extract
	for _, f := range th.FeatureList {
		f.Extract(&v)
	}

	// sort
	var (
		nodes, n = v.RenderNodes, len(v.RenderNodes)
	)
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].SortId < nodes[j].SortId
	})

	// draw
	for i, j := 0, 0; i < n; i = j {
		fi := nodes[i].Value >>16
		j = i+1
		for j < n && nodes[j].Value>>16 == fi {
			j++
		}
		f := th.FeatureList[fi]
		f.Draw(v.RenderNodes[i:j])
	}

	// flush, release any resource
	for _, f := range  th.FeatureList {
		f.Flush()
	}


	// view reset
	th.View.RenderNodes = th.View.RenderNodes[:0]
}

func (th *RenderSystem) Destroy() {

}

func NewRenderSystem() (rs *RenderSystem) {
	rs = &RenderSystem{MainCamera:Camera{follow:engi.Ghost}}
	rs.View.Camera = &rs.MainCamera
	rs.View.RenderNodes = make([]SortObject, 0)
	rs.MainCamera.initialize()
	return
}
