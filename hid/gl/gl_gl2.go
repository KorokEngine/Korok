// +build windows,!js
// +build !android,!ios,!js

package gl

import (
	"github.com/go-gl/gl/v2.1/gl"
	"unsafe"
)

func Init() error {
	return gl.Init()
}

func NeedVao() bool {
	return false
}

func GetError() uint32 {
	return gl.GetError()
}

func Viewport(x, y, width, height int32) {
	gl.Viewport(x, y, width, height)
}

func ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func Clear(flags uint32) {
	gl.Clear(flags)
}

func Disable(flag uint32) {
	gl.Disable(flag)
}

func Enable(flag uint32) {
	gl.Enable(flag)
}

func Scissor(x, y, w, h int32) {
	gl.Scissor(x, y, w, h)
}

func DepthMask(flag bool) {
	gl.DepthMask(flag)
}

func ColorMask(r, g, b, a bool) {
	gl.ColorMask(r, g, b, a)
}

func BlendFunc(src, dst uint32) {
	gl.BlendFunc(src, dst)
}

func DepthFunc(fn uint32) {
	gl.DepthFunc(fn)
}

// vao

func GenVertexArrays(n int32, arrays *uint32) {
	gl.GenVertexArrays(n, arrays)
}

func BindVertexArray(array uint32) {
	gl.BindVertexArray(array)
}

func DeleteVertexArrays(n int32, array *uint32) {
	gl.DeleteVertexArrays(n, array)
}


// program & shader

func CreateProgram() uint32{
	return gl.CreateProgram()
}

func AttachShader(program, shader uint32) {
	gl.AttachShader(program, shader)
}

func LinkProgram(program uint32) {
	gl.LinkProgram(program)
}

func UseProgram(id uint32) {
	gl.UseProgram(id)
}

func GetProgramiv(program uint32, pname uint32, params *int32) {
	gl.GetProgramiv(program, pname, params)
}

// TODO 原来的实现中 buf = loglength + 1 的，需要测试这种情况...
func GetProgramInfoLog(program uint32) string {
	var logLength int32
	GetProgramiv(program, INFO_LOG_LENGTH, &logLength)

	if logLength == 0 {
		return ""
	}
	buf := make([]uint8, logLength)
	gl.GetProgramInfoLog(program, logLength, nil, &buf[0])
	return string(buf)
}

func CreateShader(xtype uint32) uint32{
	return gl.CreateShader(xtype)
}

func ShaderSource(shader uint32, src string) {
	cstr, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, cstr, nil)
	free()
}

func CompileShader(shader uint32) {
	gl.CompileShader(shader)
}

func GetShaderiv(shader uint32, pname uint32, params *int32) {
	gl.GetShaderiv(shader, pname, params)
}

func GetShaderInfoLog(shader uint32) string {
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	buf := make([]uint8, logLength)
	gl.GetShaderInfoLog(shader, logLength, nil, &buf[0])
	return string(buf)
}

func DeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}

// buffers & draw

func GenBuffers(n int32, buffers *uint32) {
	gl.GenBuffers(n, buffers)
}

func BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	gl.BufferData(target, size, data, usage)
}

func BufferSubData(target uint32, offset, size int, data unsafe.Pointer)  {
	gl.BufferSubData(target, offset, size, data)
}

func BindBuffer(target uint32, buffer uint32) {
	gl.BindBuffer(target, buffer)
}

func DeleteBuffers(n int32, buffers *uint32) {
	gl.DeleteBuffers(n, buffers)
}

func DrawElements(mode uint32, count int32, typ uint32, offset int) {
	gl.DrawElements(mode, count, typ, gl.PtrOffset(offset))
}

func DrawArrays(mode uint32, first, count int32) {
	gl.DrawArrays(mode, first, count)
}

// uniform

func GetUniformLocation(program uint32, name string) int32{
	return gl.GetUniformLocation(program, gl.Str(name))
}

func Uniform1i(loc, v int32) {
	gl.Uniform1i(loc, v)
}

func Uniform1iv(loc, num int32, v *int32) {
	gl.Uniform1iv(loc, num, v)
}

func Uniform1f(location int32, v0 float32) {
	gl.Uniform1f(location, v0)
}

func Uniform2f(location int32, v0, v1 float32) {
	gl.Uniform2f(location, v0, v1)
}

func Uniform3f(location int32, v0, v1, v2 float32) {
	gl.Uniform3f(location, v0, v1, v2)
}

func Uniform4f(location int32, v0, v1, v2, v3 float32) {
	gl.Uniform4f(location, v0, v1, v2, v3)
}

func Uniform1fv(loc, num int32, v *float32) {
	gl.Uniform1fv(loc, num, v)
}

func Uniform4fv(loc, num int32, v *float32){
	gl.Uniform4fv(loc, num, v)
}

func UniformMatrix3fv(loc, num int32, t bool, v *float32) {
	gl.UniformMatrix3fv(loc, num, t, v)
}

func UniformMatrix4fv(loc, num int32, t bool, v *float32) {
	gl.UniformMatrix4fv(loc, num, t, v)
}

// attribute

func EnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(index)
}

func VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset int) {
	gl.VertexAttribPointer(index, size, xtype, normalized, stride, gl.PtrOffset(offset))
}

func DisableVertexAttribArray(index uint32) {
	gl.DisableVertexAttribArray(index)
}

func GetAttribLocation(program uint32, name string) int32 {
	return gl.GetAttribLocation(program, gl.Str(name))
}

// texture

func ActiveTexture(texture uint32) {
	gl.ActiveTexture(texture)
}

func BindTexture(target uint32, texture uint32) {
	gl.BindTexture(target, texture)
}

func TexSubImage2D(target uint32, level int32, xOffset, yOffset, width, height int32, format, xtype uint32, pixels unsafe.Pointer) {
	gl.TexSubImage2D(target, level, xOffset, yOffset, width, height, format, xtype, pixels)
}

func TexImage2D(target uint32, level int32, internalFormat int32, width, height, border int32, format, xtype uint32, pixels unsafe.Pointer) {
	gl.TexImage2D(target, level, internalFormat, width, height, border, format, xtype, pixels)
}

func DeleteTextures(n int32, textures *uint32) {
	gl.DeleteTextures(n, textures)
}

func GenTextures(n int32, textures *uint32) {
	gl.GenTextures(n, textures)
}

func TexParameteri(texture uint32, pname uint32, params int32) {
	gl.TexParameteri(texture, pname, params)
}
