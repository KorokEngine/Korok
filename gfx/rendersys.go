package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok/space"
)

type RenderType int32

const (
	RenderType_Mesh 	RenderType = iota
	RenderType_Batch
)

// TypeRender 负责把各种各样的 RenderData 从 RenderComp 里面取出来
type TypeRender interface {
	Draw(ref []CompRef)
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

type RenderSystem struct {
	MainCamera Camera

	// cull
	C CullSystem

	// batch
	B BatchSystem

	// render for each-type render-data
	renders [8]TypeRender
}

// 在渲染系统里面，可以维护一组 Transform 缓存，
// 这样不需要访问 TransformSystem 就可以快速的使用这些数据
// 有两个地方期待这里的数据：
// 1. CullingSystem - 更新位移数据
// 2. RenderSystem 的 Matrix 缓存
// 现在问题来了，部分游戏对象是通过Batch系统绘制的，不需要 Matrix，同时希望能够直接
// 访问 SRT 数据，而不是 Matrix！
func (th *RenderSystem) UpdateTransform(transforms []space.Transform) {
	// update culling system TODO
	cs := th.C
	for _, xform := range transforms {
		cs.UpdateBounding(int32(xform.Entity), BoundingBox{})
	}
}

// register type-render
func (th *RenderSystem) RegisterTypeRender(t RenderType, render TypeRender) {
	th.renders[t] = render
}

func (th *RenderSystem) Update(dt float32) {
	// 使用筛选系统对所有的可渲染对象进行筛选, 得到一个可见对象列表
	// 如果需要实现多相机，只需要在此进行多次筛选，比如:
	// main_view = th.C.Collect(&th.MainCamera)
	// view1 = th.C.Collect(&th.SecondCamera)
	// view2 = th.C.Collect(&th.ThirdCamera)
	visibleObjects := th.C.Collect(&th.MainCamera)

	// 对当前的可见列表进行渲染
	// 可见对象的底层数据是按类型存储的，比如 Sprite/Mesh 分别存储在各自的列表里面
	// 所以应该对可见对象按类型排序，这样实际渲染的时候，对每种底层类型数据的访问，
	// 内存是连续的.
	renderObjects := visibleObjects

	// 2. sort by type
	// TODO
	// 假设已经排好序了. - 可以在此时筛选出所有的可Batch对象，执行Batch系统
	// 这样就不必把 Batch 放在 BatchRender 里面，Render 继续负责最底层的渲染
	// 典型的Batch对象，Sprite/Text/
	// 筛选batch对象可以非常迅速，采用首尾交换的遍历，可以迅速分开可以Batch的不可Batch的对象
	// 2.1 find batchable object
	unbatchObjects := renderObjects[:10]
	batchableObject := renderObjects[10:]

	// 2.2 execute batch system
	batchObjects := th.B.Batch(batchableObject)

	// 2.3 合并batch后的结果, 此时可以进行最终的绘制了！
	// 无论是 Cull 还是 Batch 只是性能优化而已， 渲染需要转化最终的渲染命令
	// 这是 RenderFeature 的事情
	renderObjects = append(unbatchObjects, batchObjects...)

	// 使用 RenderFeature 转化渲染命令
	// 渲染不仅需要待渲染对象还要得到 Transform 目前
	// 其实此时已经没有必要合并 renderObject了，既然batch后的全部都是batch对象，那么直接渲染即可。

	// 以上逻辑不对，Cull系统中保存的 id 是 Entity-Id，但是对应的确实 Comp，所以！！！！
	// 所以需要重新思考：Cull/场景管理/渲染系统集成
	//

	// 3. extract and draw
	//var N = len(refs)
	//var xType int32
	//var left, right int
	//for i := 0; i < N; i++ {
	//
	//	right = left + 1
	//	xType = refs[left].Type
	//	for right < N && refs[right].Type == xType {
	//		right ++
	//	}
	//	th.renders[xType].Draw(refs[left:right])
	//}
	// 4. debug draw
	//
}

func (th *RenderSystem) Destroy() {

}

func NewRenderSystem() *RenderSystem {
	th := new(RenderSystem)
	return th
}
