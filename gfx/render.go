package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"korok/ecs"
)

/**
	处理渲染相关问题
 */
type SpriteRender struct {
	Shader *Shader
}

func NewSpriteRender(shader *Shader)  *SpriteRender{
	renderer := &SpriteRender{Shader:shader}
	return renderer
}

// 目前只支持位置变换
func (renderer *SpriteRender)Draw(comp *RenderComp)  {
	// Prepare transformations
	renderer.Shader.Use()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	ident := mgl32.Ident4()
	//ident.Mul4(mgl32.Translate3D(50, 50, 0))
	renderer.Shader.SetMatrix4("model\x00", *comp.SRT(&ident))

	//绑定纹理
	gl.BindTexture(gl.TEXTURE_2D, comp.tex)

	// 绘制
	gl.BindVertexArray(comp.vao)

	if comp.ebo > 0 {
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
	}
	gl.BindVertexArray(0)

	gl.Disable(gl.BLEND)
}

// 渲染组件
type RenderComp struct{
	ecs.Entity
	// Mesh - pointer to mesh data
	mesh *Mesh

	// Texture
	tex uint32

	// Vertex - [<x,y>, <u,v>]
	vao, vbo, ebo uint32

	// Position
	position mgl32.Vec2

	// Rotation
	rotation float32

	// Scale
	scale mgl32.Vec2

	// width & height
	width, height float32
}

func (comp *RenderComp) SetTexture(tex *Texture2D)  {
	comp.tex = tex.Id
	comp.width = tex.Width
	comp.height = tex.Height

	m := NewQuadMesh(tex)
	vao, vbo := m.Setup()
	comp.vao, comp.vbo = vao, vbo
	comp.mesh = m
}

func (comp *RenderComp) SetSubTexture(tex *SubTex) {
	comp.tex = tex.Id
	comp.width = tex.Max[0] - tex.Min[0]
	comp.height = tex.Max[1] - tex.Min[1]

	m := NewQuadMeshSubTex(tex)
	vao, vbo := m.Setup()
	comp.vao, comp.vbo = vao, vbo
	comp.mesh = m
}

func (comp *RenderComp) SetMesh(m *Mesh, setup func() (vao, vbo, ebo uint32)) {
	comp.tex = m.tex
	comp.mesh = m
	comp.vao, comp.vbo, comp.ebo = setup()
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
	renderer *SpriteRender

	comps []RenderComp
	_map  []int
	index int
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

func (th *RenderSystem) Update(dt float32) {
	for i := 1; i <= th.index; i++ {
		comp := &th.comps[i]
		th.renderer.Draw(comp)
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

func NewRenderSystem(shader *Shader) *RenderSystem {
	th := new(RenderSystem)
	th.renderer = NewSpriteRender(shader)
	th.comps = make([]RenderComp, STEP)
	th._map = make([]int, STEP)
	return th
}
