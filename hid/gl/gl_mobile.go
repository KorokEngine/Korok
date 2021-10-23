// +build !js
// +build android ios

package gl

import (
	"unsafe"

	"golang.org/x/mobile/gl"

	//"log"
	"runtime"
)

var glc gl.Context

func NeedVao() bool {
	return false
}

func InitContext(drawContext interface{}) {
	glc = drawContext.(gl.Context)
}

func Release() {
	glc = nil
	runtime.GC()
}

func GetError() uint32 {
	return uint32(glc.GetError())
}

func Viewport(x, y, width, height int32) {
	glc.Viewport(int(x), int(y), int(width), int(height))
}

func ClearColor(r, g, b, a float32) {
	glc.ClearColor(r, g, b, a)
}

func Clear(flags uint32) {
	glc.Clear(gl.Enum(flags))
}

func Disable(flag uint32) {
	glc.Disable(gl.Enum(flag))
}

func Enable(flag uint32) {
	glc.Enable(gl.Enum(flag))
}

func Scissor(x, y, w, h int32) {
	glc.Scissor(x, y, w, h)
}

func DepthMask(flag bool) {
	glc.DepthMask(flag)
}

func ColorMask(r, g, b, a bool) {
	glc.ColorMask(r, g, b, a)
}

func BlendFunc(src, dst uint32) {
	glc.BlendFunc(gl.Enum(src), gl.Enum(dst))
}

func DepthFunc(fn uint32) {
	glc.DepthFunc(gl.Enum(fn))
}

// vao

func GenVertexArrays(n int32, arrays *uint32) {
	va := glc.CreateVertexArray()
	*arrays = va.Value
}

func BindVertexArray(array uint32) {
	glc.BindVertexArray(gl.VertexArray{array})
}

func DeleteVertexArrays(n int32, array *uint32) {
	glc.DeleteVertexArray(gl.VertexArray{*array})
}

// program & shader

func CreateProgram() uint32 {
	return glc.CreateProgram().Value
}

func AttachShader(program, shader uint32) {
	glc.AttachShader(gl.Program{true, program}, gl.Shader{shader})
}

func LinkProgram(program uint32) {
	glc.LinkProgram(gl.Program{true, program})
}

func UseProgram(id uint32) {
	glc.UseProgram(gl.Program{true, id})
}

func GetProgramiv(program uint32, pname uint32, params *int32) {
	v := glc.GetProgrami(gl.Program{true, program}, gl.Enum(pname))
	*params = int32(v)
}

func GetProgramInfoLog(program uint32) string {
	return glc.GetProgramInfoLog(gl.Program{true, program})
}

func CreateShader(xtype uint32) uint32 {
	return glc.CreateShader(gl.Enum(xtype)).Value
}

func ShaderSource(shader uint32, src string) {
	glc.ShaderSource(gl.Shader{shader}, src)
}

func CompileShader(shader uint32) {
	glc.CompileShader(gl.Shader{shader})
}

func GetShaderiv(shader uint32, pname uint32, params *int32) {
	v := glc.GetShaderi(gl.Shader{shader}, gl.Enum(pname))
	*params = int32(v)
}

func GetShaderInfoLog(shader uint32) string {
	return glc.GetShaderInfoLog(gl.Shader{shader})
}

func DeleteShader(shader uint32) {
	glc.DeleteShader(gl.Shader{shader})
}

// buffers & draw

func GenBuffers(n int32, buffers *uint32) {
	buffer := glc.CreateBuffer()
	*buffers = buffer.Value
}

func BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	if data == nil {
		glc.BufferInit(gl.Enum(target), size, gl.Enum(usage))
	} else {
		d := ((*[1 << 24]byte)(data))[:size]
		glc.BufferData(gl.Enum(target), d, gl.Enum(usage))
		//log.Println("update buffer:", target, " data:", (*[20]float32)(data))
	}
}

func BufferSubData(target uint32, offset, size int, data unsafe.Pointer) {
	glc.BufferSubData(gl.Enum(target), offset, ((*[1 << 24]byte)(data))[:size])
}

func BindBuffer(target uint32, buffer uint32) {
	glc.BindBuffer(gl.Enum(target), gl.Buffer{buffer})
}

func DeleteBuffers(n int32, buffers *uint32) {
	glc.DeleteBuffer(gl.Buffer{*buffers})
}

func DrawElements(mode uint32, count int32, typ uint32, offset int) {
	glc.DrawElements(gl.Enum(mode), int(count), gl.Enum(typ), offset)
}

func DrawArrays(mode uint32, first, count int32) {
	glc.DrawArrays(gl.Enum(mode), int(first), int(count))
}

// uniform

func GetUniformLocation(program uint32, name string) int32 {
	u := glc.GetUniformLocation(gl.Program{true, program}, name)
	return u.Value
}

func Uniform1i(loc, v int32) {
	glc.Uniform1i(gl.Uniform{loc}, int(v))
}

func Uniform1iv(loc, num int32, v *int32) {
	glc.Uniform1iv(gl.Uniform{loc}, ((*[1 << 24]int32)(unsafe.Pointer(v)))[:num])
}

func Uniform1f(location int32, v0 float32) {
	glc.Uniform1f(gl.Uniform{location}, v0)
}

func Uniform2f(location int32, v0, v1 float32) {
	glc.Uniform2f(gl.Uniform{location}, v0, v1)
}

func Uniform3f(location int32, v0, v1, v2 float32) {
	glc.Uniform3f(gl.Uniform{location}, v0, v1, v2)
}

func Uniform4f(location int32, v0, v1, v2, v3 float32) {
	glc.Uniform4f(gl.Uniform{location}, v0, v1, v2, v3)
}

func Uniform1fv(loc, num int32, v *float32) {
	glc.Uniform1fv(gl.Uniform{loc}, ((*[1 << 24]float32)(unsafe.Pointer(v)))[:num])
}

func Uniform4fv(loc, num int32, v *float32) {
	glc.Uniform4fv(gl.Uniform{loc}, ((*[1 << 24]float32)(unsafe.Pointer(v)))[:num*4])
}

func UniformMatrix3fv(loc, num int32, t bool, v *float32) {
	glc.UniformMatrix3fv(gl.Uniform{loc}, ((*[1 << 24]float32)(unsafe.Pointer(v)))[:num*9])
}

func UniformMatrix4fv(loc, num int32, t bool, v *float32) {
	glc.UniformMatrix4fv(gl.Uniform{loc}, ((*[1 << 24]float32)(unsafe.Pointer(v)))[:num*16])
}

// attribute

func EnableVertexAttribArray(index uint32) {
	glc.EnableVertexAttribArray(gl.Attrib{uint(index)})
}

func VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset int) {
	glc.VertexAttribPointer(gl.Attrib{uint(index)}, int(size), gl.Enum(xtype), normalized, int(stride), offset)
}

func DisableVertexAttribArray(index uint32) {
	glc.DisableVertexAttribArray(gl.Attrib{uint(index)})
}

func GetAttribLocation(program uint32, name string) uint32 {
	return uint32(glc.GetAttribLocation(gl.Program{true, program}, name).Value)
}

// texture

func ActiveTexture(texture uint32) {
	glc.ActiveTexture(gl.Enum(texture))
}

func BindTexture(target uint32, texture uint32) {
	glc.BindTexture(gl.Enum(target), gl.Texture{texture})
}

func TexSubImage2D(target uint32, level int32, xOffset, yOffset, width, height int32, format, xtype uint32, pixels unsafe.Pointer) {
	glc.TexSubImage2D(gl.Enum(target), int(level), int(xOffset), int(xOffset), int(width), int(height), gl.Enum(format), gl.Enum(xtype), ((*[1 << 24]byte)(pixels))[:])
}

func TexImage2D(target uint32, level int32, internalFormat int32, width, height, border int32, format, xtype uint32, pixels unsafe.Pointer) {
	glc.TexImage2D(gl.Enum(target), int(level), int(internalFormat), int(width), int(height), gl.Enum(format), gl.Enum(xtype), ((*[1 << 24]byte)(pixels))[:])
}

func DeleteTextures(n int32, textures *uint32) {
	glc.DeleteTexture(gl.Texture{*textures})
}

func GenTextures(n int32, textures *uint32) {
	tex := glc.CreateTexture()
	*textures = tex.Value
}

func TexParameteri(texture uint32, pname uint32, params int32) {
	glc.TexParameteri(gl.Enum(texture), gl.Enum(pname), int(params))
}
