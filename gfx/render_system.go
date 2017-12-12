package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok.io/korok/engi"
	"korok.io/korok/gfx/bk"
)

type RenderType int32

type Render interface {
	SetCamera(camera Camera)
}

// 适合于渲染系统访问的表达方式.
// 其实不必这么麻烦，我们在 RenderFeature里面涉及一个 Extract 步骤，构建一个渲染列表，然后再绘制即可.
// 这个列表需要动态构建
type RenderObject struct {
	RenderData

	Type uint32

	// Position
	position mgl32.Vec2

	// Rotation
	rotation float32

	// Scale
	scale mgl32.Vec2
}

// 传入参数是经过可见性系统筛选后的 Entity，这是一个很小的数组，可以
// 直接传给各个 RenderFeature 来做可见性判断.
type RenderFeature interface {
	Draw(filter []engi.Entity)
}

// 所有的Table和Render都在此管理
// 其它的 RenderFeature 在此提取依赖
// 这样的话， RenderSystem 就沦为一个管理 RenderFeature 和 Table 的地方
// 它们之间也会存在各种组合...
type RenderSystem struct {
	MainCamera Camera

	// visibility test
	V VisibilitySystem

	// render-data
	TableList []interface{}

	// render
	RenderList []Render

	// feature knows how to use render-data and render
	FeatureList []RenderFeature
}

func (th *RenderSystem) Initialize() {
	bk.Init()
	bk.Reset(480, 320)

	// Enable debug text
	bk.SetDebug(bk.DEBUG_R|bk.DEBUG_Q)
}

func (th *RenderSystem) RequireTable(tables []interface{}) {
	th.TableList = tables
}

func (th *RenderSystem) Accept(rf RenderFeature) {
	th.FeatureList = append(th.FeatureList, rf)
}

func (th *RenderSystem) featureUpdate(dt float32) {
	entities := th.V.Collect(&th.MainCamera)

	for _, f := range th.FeatureList {
		f.Draw(entities)
	}
}

// register type-render
func (th *RenderSystem) RegisterRender(t RenderType, render Render) {
	th.RenderList = append(th.RenderList, render)
}

func (th *RenderSystem) Update(dt float32) {
	// main camera
	for _, r := range th.RenderList {
		r.SetCamera(th.MainCamera)
	}

	// draw
	for _, f := range th.FeatureList {
		f.Draw(nil)
	}

	bk.Flush()
}

func (th *RenderSystem) Destroy() {

}

func NewRenderSystem() *RenderSystem {
	th := new(RenderSystem)
	return th
}
