package bk

import (
	"strings"
	"fmt"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log"
)

// 编译/配置着色器相关数据
// 可以在 Shader 中实现预定义的 Uniform
type Shader struct {
	GLShader

	// attribute layout
	// Shader 中并没有 stream 的概念，只知道 Attribute 的布局
	// 需要外部告知 attribute 和 stream 的映射关系
	// 从 Shader 的布局中只能取到所有的 Attribute 的布局
	// 但是无法知道在当前的 drawCall 中，某个 Attribute 是否 enable or disable
	// 只能从 stream 取得
	// 算法：
	// 遍历 Shader 中的 Attribute 布局，从 Attribute 中取得 stream 的映射，
	// 绑定 stream，查询 stream 中和当前 Attribute 的关联的组件是否使用,
	// 如果否则 disable，是则 enable
	// 实现：
	// 在 Shader 中定义一个 EnabledVertexAttribArrays
	// 在 stream 中也定义一个 EnabledVertexAttribArrays

	// 需要提前配置好，此处
	// 所有启动的插槽 32bit:
	EnabledVertexAttribArrays uint32
	// 插槽到 stream 的关联： slot | stream
	// slot + stream + Component
	AttribBinds [32]AttribBind
	numAttrib   uint32

	// predefined uniform
	M,V, P Uniform

	// custom uniform
	customUniforms []uint16
}

// 如果 AttribBinds 指定了一个 Stream，但是 Stream 并没有提供相应的数据(stride < Offset)
// 此时应该 disable 当前 Attribute,
func (sh *Shader) BindAttributes(R *ResManager, streams []Stream) {
	var bindStream uint16 = UINT16_MAX
	var bindStride uint16
	for i := uint32(0); i < sh.numAttrib; i++ {
		bind := sh.AttribBinds[i]
		stream := streams[bind.stream]

		if bind.stream != bindStream {
			buffer := R.vertexBuffers[stream.vertexBuffer & 0x0FFF]
			gl.BindBuffer(gl.ARRAY_BUFFER, buffer.Id)
			bindStream = bind.stream
			bindStride = (buffer.layout >> 16) & 0xFF
		}

		slot := uint32(bind.slot)
		enable := bindStride != 0

		if enable {
			gl.EnableVertexAttribArray(slot)
			xType := g_AttrType[bind.comp.Type]
			size := g_AttrType2Size[xType]
			offset := int(bind.comp.Offset)

			var norm bool
			if (bind.comp.Normalized & 0x01) != 0 {
				norm = true
			}
			if offset < int(bindStride) {
				gl.VertexAttribPointer(slot, size, xType, norm,  int32(bindStride), gl.PtrOffset(int(offset)))
			} else {
				gl.DisableVertexAttribArray(slot)
			}
		} else {
			gl.DisableVertexAttribArray(slot)
		}
	}
}

type AttribBind struct {
	slot   uint16 	// slot location
	stream uint16 	// stream index

	comp VertexComp // attribute component format
}

func (sh *Shader) AddAttributeBinding(attr string, stream uint32, comp VertexComp) {
	slot := gl.GetAttribLocation(sh.Program, gl.Str(attr))

	if (g_debug & DEBUG_R) != 0 {
		log.Printf("Bind attr: %s => %d", attr, slot)
	}

	bind := &sh.AttribBinds[sh.numAttrib]
	bind.slot = uint16(slot)
	bind.stream = uint16(stream)
	bind.comp = comp

	sh.numAttrib ++
}

func (sh *Shader) AddUniformBinding(uniform string) {

}

type GLShader struct {
	Program uint32
}

func (s *GLShader) Use()  {
	gl.UseProgram(s.Program)
}

func (s *GLShader) create(vsh, fsh string) {
	if program, err := Compile(vsh, fsh); err == nil {
		s.Program = program
		gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))
	} else {
		log.Println("Failed to alloc shader..", err)
	}
}

func (s *GLShader) Destroy() {

}

func (s *GLShader)SetFloat(name string, value float32)  {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform1f(location, value)
}

func (s *GLShader)SetInteger(name string, value int32)  {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform1i(location, value)
}

func (s *GLShader)SetMatrix4(name string, mat4 mgl32.Mat4)  {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.UniformMatrix4fv(location, 1, false, &mat4[0])
}

func (s *GLShader)SetVector2f(name string, x, y float32)  {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform2f(location, x, y)
}


func (s *GLShader)SetVector3f(name string, x, y, z float32)  {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform3f(location, x, y, z)
}

func (s *GLShader)SetVector4f(name string, x, y, z, w float32)  {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform4f(location, x, y, z, w)
}

func (s *GLShader)SetVec4fArray(name string, array []float32, count int32) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform4fv(location, count, &array[0])
}

func (s *GLShader) GetAttrLocation(attr string)  uint32{
	return uint32(gl.GetAttribLocation(s.Program, gl.Str(attr)))
}

func (s *GLShader) GetUniformLocation(uniform string) uint32 {
	return uint32(gl.GetUniformLocation(s.Program, gl.Str(uniform)))
}

func GetErrors() string {
	return ""
}

func Compile(vertex, fragment string)  (uint32, error){
	// 1. 编译顶点着色器
	vertexShader, err := compileShader(vertex, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	// 2. 编译颜色着色器
	fragmentShader, err := compileShader(fragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	// 3. 链接程序
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	// 4. 如果有错误，读取日志
	ok, desp := getProgramStatus(program)
	if !ok {
		return 0, fmt.Errorf("failed to link program %v", desp)
	}

	// 5. 删除Shader占用的资源
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

// 编译着色器小程序, 类型： gl.VertexShader or gl.fragmentShader
// 如果错误，提取错误信息并返回
func compileShader(src string, shaderType uint32) (uint32, error)  {
	shader := gl.CreateShader(shaderType)

	cstr, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, cstr, nil)
	free()
	gl.CompileShader(shader)

	ok, err := getShaderStatus(shader)
	if !ok {
		return 0, fmt.Errorf("failed to compile %v: %v", src, err)
	}
	return shader, nil
}

func getShaderStatus(shader uint32) (bool, string) {
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.TRUE {
		return true, ""
	}

	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength + 1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return false, log
}

func getProgramStatus(program uint32) (bool, string) {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status == gl.TRUE {
		return true, ""
	}

	var logLength  int32
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	return false, log
}


// 解析 attribute 和 uniform
//func (s *Shader) Setup() {
//	// 4. 解析 attribute
//	count := int32(0)
//	nameL := int32(0)
//	gl.GetProgramiv(s.Program, gl.ACTIVE_ATTRIBUTES, &count)
//	gl.GetProgramiv(s.Program, gl.ACTIVE_ATTRIBUTE_MAX_LENGTH, &nameL)
//
//	for i := int32(0); i < count; i++ {
//		var name string = strings.Repeat("\x00", int(nameL+1))
//		var attr VertexAttr
//		gl.GetActiveAttrib(s.Program, uint32(i), nameL, nil, &attr.Size, &attr.Type, gl.Str(name))
//
//		attr.slot = uint32(gl.GetAttribLocation(s.Program, gl.Str(name)))
//		s.AttributeMap[name] = attr
//	}
//
//	// 5. 解析 uniform
//	gl.GetProgramiv(s.Program, gl.ACTIVE_UNIFORMS, &count)
//	gl.GetProgramiv(s.Program, gl.ACTIVE_UNIFORM_MAX_LENGTH, &nameL)
//
//	for i := int32(0); i < count; i++ {
//		var name string = strings.Repeat("\x00", int(nameL+ 1))
//		var uniform Uniform
//		var xtype uint32
//
//		gl.GetActiveUniform(s.Program, uint32(i), nameL, nil, &uniform.Size, &xtype, gl.Str(name))
//		uniform.slot = uint32(gl.GetUniformLocation(s.Program, gl.Str(name)))
//		uniform.Type = UniformType(xtype)
//
//		s.UniformMap[name] = uniform
//	}
//}