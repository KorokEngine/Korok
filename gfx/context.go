package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	//"github.com/go-gl/mathgl/mgl32"
	//"github.com/go-gl/mathgl/mgl32"
)

type BlendFunc struct {
	Src, Dst uint32
}

var BF_Add = BlendFunc{gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA}
var BF_Sub = BlendFunc{gl.ONE, gl.ONE_MINUS_SRC_ALPHA}

type VertexLayout []VertexAttr

type UniformLayout []Uniform

// 数据格式描述
type PipelineState struct {
	BlendFunc
	GLShader
	tex uint32

	UniformLayout
	VertexLayout
}

//// state
//// 0x000 000 - 00
//// shader - blend - depth

// 表示GPU状态和并管理GPU资源
type RenderContext struct {
	PipelineState
	BlendFunc

	// bound program
	program uint32


	VAO uint32
	// <shader | blend | depth>
	state uint32
	tex   uint32

	// user defined uniform-data
	UniformData

	// Vertex and Index buffer
	VertexBuffer Buffer
	IndexBuffer Buffer
}

func NewRenderContext(shader GLShader, tex uint32) *RenderContext {
	rc := new(RenderContext)
	rc.GLShader = shader
	rc.PipelineState.tex = tex
	rc.PipelineState.BlendFunc = BlendFunc{gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA}

	shader.Use()

	shader.SetInteger("tex\x00", 0)
	gl.BindFragDataLocation(shader.Program, 0, gl.Str("outputColor\x00"))


	m := NewQuadMesh(&Texture2D{Id:tex})
	m.Setup()

	rc.VAO = m.VAO()
	rc.SetVertexBuffer(m.VertexBuffer())
	return rc
}

func (rc *RenderContext) Require(state uint32) {
	if rc.state != state {

	}
}

func (rc *RenderContext) SetPipelineState(state PipelineState) {
	rc.PipelineState = state
}

func (rc *RenderContext) SetUniformData(uniform UniformData) {
	rc.UniformData = uniform
}

func (rc *RenderContext) SetVertexBuffer(buffer Buffer) {
	rc.VertexBuffer = buffer
}

func (rc *RenderContext) SetIndexBuffer(buffer Buffer) {
	rc.IndexBuffer = buffer
}

// uniform 有些是固定的，比如projection有些是每次更新的，需要区别对待
// 现在每次都更新Uniform TODO 优化
func (rc *RenderContext) PreDraw() {
	// blend func
	bf := rc.PipelineState.BlendFunc
	if rc.BlendFunc != bf {
		rc.BlendFunc = bf
		gl.Enable(gl.BLEND)
		gl.BlendFunc(bf.Src, bf.Dst)
	}

	// bind program
	program := rc.PipelineState.GLShader.Program
	if program != rc.program {
		rc.program = program
		gl.UseProgram(program)
	}

	// bind texture
	tex := rc.PipelineState.tex
	if tex != rc.tex {
		rc.tex = tex
		gl.BindTexture(gl.TEXTURE_2D, tex)
	}

	// bind vbo and vertex attributes //

	gl.BindBuffer(gl.ARRAY_BUFFER, rc.VertexBuffer.Id)
	//for _, vertexAttr := range rc.PipelineState.VertexLayout {
	//	log.Println(vertexAttr)
	//	gl.EnableVertexAttribArray(vertexAttr.Slot)
	//	gl.VertexAttribPointer(vertexAttr.Slot, vertexAttr.Size, vertexAttr.Type, vertexAttr.Normalized, vertexAttr.Stride, gl.Ptr(0))
	//}

	gl.BindVertexArray(rc.VAO)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	//gl.BindVertexArray(0)

	// bind uniform
	for _, uniform := range rc.PipelineState.UniformLayout {
		uniformIndex := uniform.Slot
		dataIndex := uniform.Data

		var data interface{}
		if data = rc.UniformData.Index(dataIndex); data == nil {
			continue
		}

		switch uniform.Type {
		case UniformFloatN:
			gl.Uniform1fv(int32(uniformIndex), uniform.Count, data.(*float32))
		case UniformIntN:
			gl.Uniform1iv(int32(uniformIndex), uniform.Count, data.(*int32))
		case UniformMat4:
			gl.UniformMatrix4fv(int32(uniformIndex), uniform.Count, false, data.(*float32))
		}
	}
}

//
func (rc *RenderContext) SetViewport() {

}

// 我觉得数组应该是最通用的数据类型，所以先支持下数组，如果数组无法 cover 其它情况
func (rc *RenderContext) Draw() {
	//// setup gpu-state
	//rc.Require(rc.state)
	//
	//
	rc.PreDraw()

	//
	//// 4. draw
	//if rc.IndexBuffer.Id != 0 {
	//	b := rc.IndexBuffer
	//	gl.DrawElements(b.T, b.Count, b.Type, nil)
	//} else {
	//	//b := rc.VertexBuffer
	//	gl.DrawArrays(gl.ARRAY_BUFFER, 0, 6)
	//}

	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	//gl.BindVertexArray(0)
	//gl.BindBuffer(gl.ARRAY_BUFFER, rc.VertexBuffer.Id)
	//gl.EnableVertexAttribArray(0)
	//gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	//gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

// gl.VertexAttribPointer(bs.inPosition, 2, gl.FLOAT, false, 20, gl.Ptr(0))
type VertexAttr struct {
	Data uint32 			// index of vertex data, 目前只支持一个VBO，所以此字段现在并没有特殊的用途
	Slot uint32				// slot  in shader

	Size int32
	Type uint32
	Normalized bool
	Stride int32
	Offset int32
	Pointer int32
}

type BufferDescription struct {
	Attr []VertexAttr
}

type Uniform struct {
	Data uint32 		// index of uniform data
	Slot uint32 		// slot in shader

	Size int32
	Type UniformType
	Count int32
}

type UniformType byte

/// 支持哪些类型呢?
const (
	UniformStart UniformType = iota
	UniformMat4

	UniformFloat1 	// float
	UniformFloat2   // vec2
	UniformFloat3   // vec3
	UniformFloat4   // vec4
	UniformFloatN   // float[]

	UniformInt1 	// int
	UniformInt2		// int_vec2
	UniformInt3 	// int_vec3
	UniformInt4 	// int_vec4
	UniformIntN  	// int[]
)

//
type UniformData struct {
	data [8]interface{}
}

func (ud *UniformData) AddUniform(slot uint32, value interface{}) {
	ud.data[slot] = value
}

func (ud *UniformData) Index(slot uint32) interface{} {
	return ud.data[slot]
}

//func DrawCommand(cmd *Command)  {
//	old := uint32(renderer.SortKey)
//	renderer.SortKey = cmd.Key
//
//	// 1. is shader changed
//	if (old & ShaderMask) != (uint32(cmd.Key) & ShaderMask) {
//		rs := renderer.states[cmd.T]
//		renderer.SetShader(rs.Shader)
//		rs.Shader.Prepare()
//
//		//log.Println("use shader:", rs.Shader)
//	}
//
//	// 2. is blend-func changed
//	if (old & BlendMask) != (uint32(cmd.Key) & BlendMask) {
//		rs := renderer.states[cmd.T]
//		gl.Enable(gl.BLEND)
//		renderer.SetBlendFunc(rs.Src, rs.Dst)
//
//		log.Println("use blend-func")
//	}
//
//	// 3. is texture changed
//	if (old & TextureMask) != (uint32(cmd.Key) & TextureMask) {
//		var tex uint32
//		if cmd.T == RenderType_Batch {
//			tex = cmd.Data.(*Batch).tex
//		} else {
//			tex = cmd.Data.(*RenderComp).tex
//		}
//		renderer.SetTexture(tex)
//
//		log.Println("use texture!!")
//	}
//
//	// 4. render
//	renderer.Shader.Draw(cmd.Data)
//}
//
//
