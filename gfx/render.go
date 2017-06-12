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
	quadVAO uint32
	Shader *Shader
}

func NewSpriteRender(shader *Shader)  *SpriteRender{
	renderer := &SpriteRender{Shader:shader}

	// Configure VAO/VBO
	var vertices = []float32{
		// Pos      // Tex
		0.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo);
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW);

	vertAttrib := uint32(gl.GetAttribLocation(shader.Program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(shader.Program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	gl.BindVertexArray(0)
	renderer.quadVAO = vao

	return renderer
}

// 目前只支持位置变换
func (renderer *SpriteRender)Draw(texture *Texture2D, mat4 *mgl32.Mat4)  {
	// Prepare transformations
	renderer.Shader.Use()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	renderer.Shader.SetMatrix4("model\x00", *mat4)

	//绑定纹理
	//gl.ActiveTexture(gl.TEXTURE0)
	texture.Bind();

	// 绘制
	gl.BindVertexArray(renderer.quadVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)

	gl.Disable(gl.BLEND)
}

// 渲染组件
type RenderComp struct{
	ecs.Entity
	Texture2D
	Model mgl32.Mat4
}

func (comp *RenderComp) SetTexture(tex Texture2D)  {
	comp.Texture2D = tex
}

func (comp *RenderComp) SetSize(width, height float32)  {
	comp.Texture2D.Width = width
	comp.Texture2D.Height = height
}

func (comp *RenderComp) SetPosition(position mgl32.Vec2)  {
	comp.Model = mgl32.Translate3D(position[0], position[1], 0)
}

func (comp *RenderComp) SetRotation(rotation float32)  {
	comp.Model = mgl32.HomogRotate3DZ(rotation)
}

func (comp *RenderComp) SetScale(x, y float32)  {
	comp.Model = mgl32.Scale3D(x, y, 1)
}

// 渲染系统
// 渲染架构 - 草稿

const STEP  = 100

var (
	comps []RenderComp
	_map []int
	index int
)

var (
	renderer SpriteRender
)

func init()  {
	comps = make([]RenderComp, STEP)
}

func NewRenderComp(id uint32) *RenderComp{
	index += 1
	len := len(comps)
	if index >= len {
		comps = resize(comps, len + STEP)
	}
	comp := RenderComp{
		Model: mgl32.Ident4(),
		Entity: ecs.Entity(index),
	}
	comps[index] = comp
	_map[id] = index
	return &comp
}

func Update(dt float32) {
	for _, comp := range comps {
		renderer.Draw(&comp.Texture2D, &comp.Model)
	}
}

func Delete(id uint32) {
	i := _map[id]
	if i < index {
		comps[index], comps[i] = comps[i], comps[index]
		_map[comps[i].Index()] = i
		_map[index] = 0
	} else if i == index {
		_map[index] = 0
	}
	index -= 1
}

func GetComp(id uint32)  *RenderComp{
	if _map[id] == 0 {
		return nil
	}
	return &comps[_map[id]]
}

func resize(slice []RenderComp, size int) []RenderComp {
	newSlice := make([]RenderComp, size)
	copy(newSlice, slice)
	return newSlice
}

