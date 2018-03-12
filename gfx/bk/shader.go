package bk

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
	"korok.io/korok/math/f32"
	"log"
)

// 编译/配置着色器相关数据
// 可以在 Shader 中实现预定义的 Uniform
type Shader struct {
	GLShader

	// bind attribute
	AttrBinds [32]AttribBind
	numAttr   uint32

	// predefined uniform
	M, V, P Uniform

	// custom uniform
	customUniforms []uint16
}

// 如果 AttrBinds 指定了一个 Stream，但是 Stream 并没有提供相应的数据(stride < Offset)
// 此时应该 disable 当前 Attribute,
func (sh *Shader) BindAttributes(R *ResManager, streams []Stream) {
	var bindStream uint16 = UINT16_MAX
	var bindStride uint16
	for i := uint32(0); i < sh.numAttr; i++ {
		bind := sh.AttrBinds[i]
		stream := streams[bind.stream]

		if bind.stream != bindStream {
			buffer := R.vertexBuffers[stream.vertexBuffer&0x0FFF]
			gl.BindBuffer(gl.ARRAY_BUFFER, buffer.Id)
			bindStream = bind.stream
			bindStride = buffer.layout
		}

		slot := uint32(bind.slot)
		enable := bindStride != 0

		if enable {
			gl.EnableVertexAttribArray(slot)

			comp := bind.comp
			num := int32(comp.Num)
			xType := g_AttrType[comp.Type]
			offset := int(comp.Offset)

			var norm bool
			if (comp.Normalized & 0x01) != 0 {
				norm = true
			}
			if offset < int(bindStride) {
				gl.VertexAttribPointer(slot, num, xType, norm, int32(bindStride), gl.PtrOffset(int(offset)))
			} else {
				gl.DisableVertexAttribArray(slot)
			}
		} else {
			gl.DisableVertexAttribArray(slot)
		}
	}
}

type AttribBind struct {
	slot   uint16 // slot location
	stream uint16 // stream index

	comp VertexComp // attribute component format
}

func (sh *Shader) AddAttributeBinding(attr string, stream uint32, comp VertexComp) {
	slot := gl.GetAttribLocation(sh.Program, gl.Str(attr))

	if slot < 0 {
		log.Printf("fail to bind attribute: %s, %v", attr, comp)
	} else {
		if (g_debug & DEBUG_R) != 0 {
			log.Printf("bind attr: %s => %d", attr, slot)
		}
	}

	bind := &sh.AttrBinds[sh.numAttr]
	bind.slot = uint16(slot)
	bind.stream = uint16(stream)
	bind.comp = comp

	sh.numAttr++
}

func (sh *Shader) AddUniformBinding(uniform string) {

}

type GLShader struct {
	Program uint32
}

func (s *GLShader) Use() {
	gl.UseProgram(s.Program)
}

func (s *GLShader) Create(vsh, fsh string) error{
	if program, err := Compile(vsh, fsh); err == nil {
		s.Program = program
		gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))
		return nil
	} else {
		return err
	}
}

func (s *GLShader) Destroy() {

}

func (s *GLShader) SetFloat(name string, value float32) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform1f(location, value)
}

func (s *GLShader) SetInteger(name string, value int32) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform1i(location, value)
}

func (s *GLShader) SetMatrix4(name string, mat4 f32.Mat4) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.UniformMatrix4fv(location, 1, false, &mat4[0])
}

func (s *GLShader) SetVector2f(name string, x, y float32) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform2f(location, x, y)
}

func (s *GLShader) SetVector3f(name string, x, y, z float32) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform3f(location, x, y, z)
}

func (s *GLShader) SetVector4f(name string, x, y, z, w float32) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform4f(location, x, y, z, w)
}

func (s *GLShader) SetVec4fArray(name string, array []float32, count int32) {
	location := gl.GetUniformLocation(s.Program, gl.Str(name))
	gl.Uniform4fv(location, count, &array[0])
}

func (s *GLShader) GetAttrLocation(attr string) uint32 {
	return uint32(gl.GetAttribLocation(s.Program, gl.Str(attr)))
}

func (s *GLShader) GetUniformLocation(uniform string) uint32 {
	return uint32(gl.GetUniformLocation(s.Program, gl.Str(uniform)))
}

func GetErrors() string {
	return ""
}

func Compile(vertex, fragment string) (uint32, error) {
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
func compileShader(src string, shaderType uint32) (uint32, error) {
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

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return false, log
}

func getProgramStatus(program uint32) (bool, string) {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status == gl.TRUE {
		return true, ""
	}

	var logLength int32
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
