package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"korok/ecs"
	"korok/gfx/bk"
)

// 可以在划分出几个子系统：
// 1. BatchSystem 批量合并系统
// 2. CullingSystem 不可见剔除
// 3. LayerSystem Z-Order 绘制顺序管理系统
// 4. RenderSystem - 最终调用 OpenGL API进行绘制
//
// 应该设计一个渲染上下文的概念，每个上下文具备所有需要条件
// 来渲染指定的数据流 - 这里面抽象成了 Command
// RenderComp 现在变成一个中间层的概念，在 gameplay > RenderComp > Command
// 用来协调不同渲染数据的承载
type RenderType int32

const (
	RenderType_Mesh 	RenderType = iota
	RenderType_Quad
	RenderType_Shape
	RenderType_Text
	RenderType_Batch
	RenderType_Ptl
)

// TypeRender 负责把各种各样的 RenderData 从 RenderComp 里面取出来
type TypeRender interface {
	Extract()

	Draw(d RenderData, pos, scale mgl32.Vec2, rot float32)
}


/**
	处理渲染相关问题

	关于 Render 的设计，shader, texture, func 其实属于GPU的资源，是GPU的状态
	状态应该单独控制。用 Render 来管理GPU，但是还有些数据，比如如何渲染某个图形，这类
	操作大多是用户定义的，例如 uniform 的操作！所以应该把 uniform 的操作从 Render
	中剥离.

	uniform 涉及的内容其实就是非常具体的如何绘制的问题！
 */

// 真正渲染的时候需要的数据：
// 1. VAO - 如果支持，仅此即可。否则需要下面的数据
// 2. VBO - 绑定 buffer
// 3. Shader内部属性偏移，用来执行 VertexAttribPointer 操作 (这部分数据可以放到Shader里面自动执行)

// 渲染组件
type RenderComp struct{
	ecs.Entity

	// type
	Type RenderType

	// 渲染数据
	Data RenderData

	// <visible, >
	flag uint32

	// Position
	position mgl32.Vec2

	// Rotation
	rotation float32

	// Scale
	scale mgl32.Vec2

	// width & height
	width, height float32
}

func (comp *RenderComp) SetSize(width, height float32)  {
	comp.width = width
	comp.height = height
}

func (comp *RenderComp) SetPosition(position mgl32.Vec2)  {
	comp.position = position
}

func (comp *RenderComp) SetRotation(rotation float32)  {
	comp.rotation = rotation
}

func (comp *RenderComp) SetScale(x, y float32)  {
	comp.scale[0] = x
	comp.scale[1] = y
}

// return model translation
func (comp *RenderComp) SRT(identity *mgl32.Mat4) *mgl32.Mat4{
	id := identity.Mul4(mgl32.Translate3D(comp.position[0], comp.position[1], 0))
	return &id
}

// 渲染系统
// 渲染架构 - 草稿

const STEP  = 100

type RenderSystem struct {
	comps []RenderComp
	_map  []int
	index int

	// cull
	C CullSystem

	// batch
	B *BatchSystem

	// render for each-type render-data
	renders [8]TypeRender
}

// register type-render
func (th *RenderSystem) RegisterTypeRender(t RenderType, render TypeRender) {
	th.renders[t] = render
}

func (th *RenderSystem) NewRenderComp(id uint32) *RenderComp{
	th.index += 1
	len := len(th.comps)
	if th.index >= len {
		th.comps = resize(th.comps, len + STEP)
		th._map = resizeInt(th._map, len + STEP)
	}
	comp := RenderComp{
		Entity: ecs.Entity(th.index),
	}
	th.comps[th.index] = comp
	th._map[id] = th.index
	return &th.comps[th.index]
}

func (th *RenderSystem) Size() int{
	return th.index
}

// 在跑循环的时候处理batch吗？还是提前算好再提交？？
// 其实可以算出来一个新的队列！
// 1. 进行排序，保证渲染顺序
// 2. 过滤，得到可以直接执行的 RenderCommand
// 对 renderObj 进行排序是一个耗时的过程，主要是因为 renderObj 都比较大，其实是不方便进行排序的！！！
// 或许应该先把所有的 renderObj 转化为 Command，之后再进行排序。但是由于之前的RenderObj数量较多，
// 所以会得到灯亮的 Command，之后在对这些 Command 进行合并。
// renderCmd 的抓化虽然耗时，但是可以用多线程转化！！
func (th *RenderSystem) Update(dt float32) {
	//// cap := len(th.comps)
	//
	////
	//cmds := th.ToCommand()
	//// 1. sort cmds
	//
	//// 2. batch
	//cmds = th.B.Merge(cmds)
	//

	// 1. Cull - collect visible object
	// 2. d/c
	// 2. Batch - find batch-able object
	// 3.

	// 1. cull
	refs := th.C.Cull(th.comps, Camera{})

	// 2. sort
	// TODO

	// 3. extract and draw
	for _, ref := range refs {
		comp := ref.RenderComp
		th.renders[ref.Type].Draw(comp.Data, comp.position, comp.scale, comp.rotation)
	}
}

func (th *RenderSystem) Delete(id uint32) {
	i := th._map[id]
	if i < th.index {
		th.comps[i] = th.comps[th.index]
		th._map[th.comps[i].Index()] = i
		th._map[th.index] = 0
	} else if i == th.index {
		th._map[th.index] = 0
	}
	th.index -= 1
}

func (th *RenderSystem) GetComp(id uint32)  *RenderComp{
	if th._map[id] == 0 {
		return nil
	}
	return &th.comps[th._map[id]]
}

func (th *RenderSystem) Destroy() {

}

func resize(slice []RenderComp, size int) []RenderComp {
	newSlice := make([]RenderComp, size)
	copy(newSlice, slice)
	return newSlice
}

func resizeInt(slice []int, size int) []int {
	newSlice := make([]int, size)
	copy(newSlice, slice)
	return newSlice
}

func NewRenderSystem() *RenderSystem {
	th := new(RenderSystem)
	th.comps = make([]RenderComp, STEP)
	th._map = make([]int, STEP)
	return th
}

//
//type BatchShader struct {
//	GLShader
//
//	inPosition uint32
//	inTexCoord uint32
//	inColor    uint32
//}
//
//func (bs *BatchShader) Setup() {
//	// uniform
//	bs.Use()
//	//  ---- Vertex GLShader
//	// projection
//	p := mgl32.Ortho2D(0, 480, 0, 320)
//	bs.SetMatrix4("projection\x00", p)
//
//	//// model
//	//model := mgl32.Ident4()
//	//bs.SetModel(model)
//	//
//
//	// ---- Fragment GLShader
//	bs.SetInteger("tex\x00", 0)
//	gl.BindFragDataLocation(bs.Program, 0, gl.Str("outputColor\x00"))
//
//	// in/out stream
//	//bs.inPosition = bs.GLShader.GetAttrLocation("position\x00")
//	//bs.inTexCoord = bs.GLShader.GetAttrLocation("texCoord\x00")
//	//bs.inColor    = bs.GLShader.GetAttrLocation("color\x00")
//}
//
//func (bs *BatchShader) Prepare() {
//	bs.Use()
//}
//
//func (bs *BatchShader) Draw(any interface{}) {
//	b := any.(*Batch)
//
//	if b.vao > 0{
//		gl.BindVertexArray(b.vao)
//		gl.DrawElements(gl.TRIANGLES, int32(b.count * 6), gl.UNSIGNED_INT, nil)
//		gl.BindVertexArray(0)
//	} else {
//		gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
//		gl.EnableVertexAttribArray(bs.inPosition)
//		gl.VertexAttribPointer(bs.inPosition, 2, gl.FLOAT, false, 20, gl.Ptr(0))
//		gl.VertexAttribPointer(bs.inTexCoord, 2, gl.FLOAT, false, 20, gl.Ptr(8))
//		gl.VertexAttribPointer(bs.inColor, 4, gl.UNSIGNED_BYTE, true, 20, gl.Ptr(16))
//
//		gl.DrawElements(gl.TRIANGLES, int32(b.count), gl.UNSIGNED_INT, nil)
//		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
//	}
//}

func NewTextShader(ts bk.GLShader) bk.GLShader{
	ts.Use()

	p := mgl32.Ortho2D(0, 480, 0, 320)

	// vertex
	ts.SetMatrix4("projection\x00", p)
	ts.SetVector3f("model\x00", 50, 50, 10)

	// fragment
	ts.SetInteger("text\x00", 0)
	gl.BindFragDataLocation(ts.Program, 0, gl.Str("color\x00"))

	return ts
}
