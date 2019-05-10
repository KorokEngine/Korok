// +build js

package gl

import (
	"unsafe"

	"syscall/js"

	"korok.io/korok/webav"
)

var (
	gl            *webav.Context
	programMap    = make(map[uint32]js.Value)
	programCount  uint32
	shaderMap     = make(map[uint32]js.Value)
	shaderCount   uint32
	bufferMap     = make(map[uint32]js.Value)
	bufferCount   uint32
	locationMap   = make(map[int32]js.Value)
	locationCount int32
	textureMap    = make(map[uint32]js.Value)
	textureCount  uint32
)

type Slice struct {
	Addr uintptr
	Len  int
	Cap  int
}

func cstr2str(x string) string {
	return x[0 : len(x)-1]
}

func bytes2uint8s(x []byte) []uint8 {
	s := make([]uint8, len(x))
	for i, v := range x {
		s[i] = uint8(v)
	}
	return s
}

func Init(canvas js.Value) error {
	attrs := webav.DefaultAttributes()
	attrs.Alpha = false

	var err error
	gl, err = webav.NewContext(canvas, attrs)

	return err
}

func NeedVao() bool {
	// webgl 2.0才支持
	return false
}

func GetError() uint32 {
	return uint32(gl.GetError())
}

func Viewport(x, y, width, height int32) {
	gl.Viewport(int(x), int(y), int(width), int(height))
}

func ClearColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func Clear(flags uint32) {
	gl.Clear(int(flags))
}

func Disable(flag uint32) {
	gl.Disable(int(flag))
}

func Enable(flag uint32) {
	gl.Enable(int(flag))
}

func Scissor(x, y, w, h int32) {
	gl.Scissor(int(x), int(y), int(w), int(h))
}

func DepthMask(flag bool) {
	gl.DepthMask(flag)
}

func ColorMask(r, g, b, a bool) {
	gl.ColorMask(r, g, b, a)
}

func BlendFunc(src, dst uint32) {
	gl.BlendFunc(int(src), int(dst))
}

func DepthFunc(fn uint32) {
	gl.DepthFunc(int(fn))
}

// vao webgl 2.0才支持

func GenVertexArrays(n int32, arrays *uint32) {
}

func BindVertexArray(array uint32) {
}

func DeleteVertexArrays(n int32, array *uint32) {
}

// program & shader

func CreateProgram() uint32 {
	x := gl.CreateProgram()
	if x == js.Null() {
		return 0
	}
	programCount++
	programMap[programCount] = x
	return programCount
}

func AttachShader(program, shader uint32) {
	gl.AttachShader(programMap[program], shaderMap[shader])
}

func LinkProgram(program uint32) {
	gl.LinkProgram(programMap[program])
}

func UseProgram(id uint32) {
	gl.UseProgram(programMap[id])
}

// TODO 目前只支持bool的获取
func GetProgramiv(program uint32, pname uint32, params *int32) {
	*params = 0
	if gl.GetProgramParameterb(programMap[program], int(pname)) {
		*params = 1
	}
}

// TODO 原来的实现中 buf = loglength + 1 的，需要测试这种情况...
func GetProgramInfoLog(program uint32) string {
	// var logLength int32
	// GetProgramiv(program, INFO_LOG_LENGTH, &logLength)

	// if logLength == 0 {
	// 	return ""
	// }

	return gl.GetProgramInfoLog(programMap[program])
}

func CreateShader(xtype uint32) uint32 {
	x := gl.CreateShader(int(xtype))
	if x == js.Null() {
		return 0
	}
	shaderCount++
	shaderMap[shaderCount] = x
	return shaderCount
}

func ShaderSource(shader uint32, src string) {
	// cstr, free := gl.Strs(src)
	cstr := cstr2str(src)
	gl.ShaderSource(shaderMap[shader], cstr)
	// free()
}

func CompileShader(shader uint32) {
	gl.CompileShader(shaderMap[shader])
}

// TODO 目前只支持bool的获取
func GetShaderiv(shader uint32, pname uint32, params *int32) {
	*params = 0
	if gl.GetShaderParameter(shaderMap[shader], int(pname)).Bool() {
		*params = 1
	}
}

func GetShaderInfoLog(shader uint32) string {
	// var logLength int32
	// GetShaderiv(shader, uint32(gl.INFO_LOG_LENGTH), &logLength)
	// if logLength == 0 {
	// 	return ""
	// }
	return gl.GetShaderInfoLog(shaderMap[shader])
}

func DeleteShader(shader uint32) {
	gl.DeleteShader(shaderMap[shader])
	delete(shaderMap, shader)
}

// buffers & draw

func GenBuffers(n int32, buffers *uint32) {
	x := gl.CreateBuffer()
	if x == js.Null() {
		*buffers = 0
		return
	}
	bufferCount++
	bufferMap[bufferCount] = x
	*buffers = bufferCount
}

func BufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	if nil == data {
		gl.BufferData(int(target), size, int(usage))
		return
	}
	sl := &Slice{Addr: uintptr(data), Len: size, Cap: size}
	b := *(*[]uint8)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.BufferData(int(target), ta, int(usage))
	ta.Release()
}

func BufferSubData(target uint32, offset, size int, data unsafe.Pointer) {
	sl := &Slice{Addr: uintptr(data), Len: size, Cap: size}
	b := *(*[]uint8)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.BufferSubData(int(target), int(offset), ta)
	ta.Release()
}

func BindBuffer(target uint32, buffer uint32) {
	gl.BindBuffer(int(target), bufferMap[buffer])
}

func DeleteBuffers(n int32, buffers *uint32) {
	gl.DeleteBuffer(bufferMap[*buffers])
	delete(bufferMap, *buffers)
}

func DrawElements(mode uint32, count int32, typ uint32, offset int) {
	gl.DrawElements(int(mode), int(count), int(typ), int(offset))
}

func DrawArrays(mode uint32, first, count int32) {
	gl.DrawArrays(int(mode), int(first), int(count))
}

// uniform

func GetUniformLocation(program uint32, name string) int32 {
	x := gl.GetUniformLocation(programMap[program], cstr2str(name))
	if x == js.Null() {
		return 0
	}
	locationCount++
	locationMap[locationCount] = x
	return locationCount
}

func Uniform1i(loc, v int32) {
	gl.Uniform1i(locationMap[loc], int(v))
}

func Uniform1iv(loc, num int32, v *int32) {
	sl := &Slice{Addr: uintptr(unsafe.Pointer(v)), Len: int(num), Cap: int(num)}
	b := *(*[]int32)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.Uniform1iv(locationMap[loc], ta)
	ta.Release()
}

func Uniform1f(location int32, v0 float32) {
	gl.Uniform1f(locationMap[location], v0)
}

func Uniform2f(location int32, v0, v1 float32) {
	gl.Uniform2f(locationMap[location], v0, v1)
}

func Uniform3f(location int32, v0, v1, v2 float32) {
	gl.Uniform3f(locationMap[location], v0, v1, v2)
}

func Uniform4f(location int32, v0, v1, v2, v3 float32) {
	gl.Uniform4f(locationMap[location], v0, v1, v2, v3)
}

func Uniform1fv(loc, num int32, v *float32) {
	sl := &Slice{Addr: uintptr(unsafe.Pointer(v)), Len: int(num), Cap: int(num)}
	b := *(*[]float32)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.Uniform1fv(locationMap[loc], ta)
	ta.Release()
}

func Uniform4fv(loc, num int32, v *float32) {
	sl := &Slice{Addr: uintptr(unsafe.Pointer(v)), Len: int(num * 4), Cap: int(num * 4)}
	b := *(*[]float32)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.Uniform4fv(locationMap[loc], ta)
	ta.Release()
}

func UniformMatrix3fv(loc, num int32, t bool, v *float32) {
	sl := &Slice{Addr: uintptr(unsafe.Pointer(v)), Len: int(num * 9), Cap: int(num * 9)}
	b := *(*[]float32)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.UniformMatrix3fv(locationMap[loc], t, ta)
	ta.Release()
}

func UniformMatrix4fv(loc, num int32, t bool, v *float32) {
	sl := &Slice{Addr: uintptr(unsafe.Pointer(v)), Len: int(num * 16), Cap: int(num * 16)}
	b := *(*[]float32)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.UniformMatrix4fv(locationMap[loc], t, ta)
	ta.Release()
}

// attribute

func EnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(int(index))
}

func VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset int) {
	gl.VertexAttribPointer(int(index), int(size), int(xtype), normalized, int(stride), offset)
}

func DisableVertexAttribArray(index uint32) {
	gl.DisableVertexAttribArray(int(index))
}

func GetAttribLocation(program uint32, name string) int32 {
	return int32(gl.GetAttribLocation(programMap[program], cstr2str(name)))
}

// texture

func ActiveTexture(texture uint32) {
	gl.ActiveTexture(int(texture))
}

func BindTexture(target uint32, texture uint32) {
	gl.BindTexture(int(target), textureMap[texture])
}

// 目前支持Type:UNSIGNED_BYTE, Format:RGBA,RGB,LUMINANCE_ALPHA,LUMINANCE,ALPHA
func TexSubImage2D(target uint32, level int32, xOffset, yOffset, width, height int32, format, xtype uint32, pixels unsafe.Pointer) {
	var sizePerpix int32
	switch format {
	case RGBA:
		sizePerpix = 4
	case RGB:
		sizePerpix = 3
	case LUMINANCE_ALPHA:
		sizePerpix = 2
	case LUMINANCE:
		sizePerpix = 1
	case ALPHA:
		sizePerpix = 1
	}
	sl := &Slice{Addr: uintptr(pixels), Len: int(sizePerpix * width * height), Cap: int(sizePerpix * width * height)}
	b := *(*[]uint8)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.TexSubImage2D(int(target), int(level), int(xOffset), int(yOffset), int(format), int(xtype), ta)
	ta.Release()
}

func TexImage2D(target uint32, level int32, internalFormat int32, width, height, border int32, format, xtype uint32, pixels unsafe.Pointer) {
	// b := ((*[1 << 24]byte)(pixels))[:]
	// bytes2uint8s(b)
	var sizePerpix int32
	switch format {
	case RGBA:
		sizePerpix = 4
	case RGB:
		sizePerpix = 3
	case LUMINANCE_ALPHA:
		sizePerpix = 2
	case LUMINANCE:
		sizePerpix = 1
	case ALPHA:
		sizePerpix = 1
	}
	sl := &Slice{Addr: uintptr(pixels), Len: int(sizePerpix * width * height), Cap: int(sizePerpix * width * height)}
	b := *(*[]uint8)(unsafe.Pointer(sl))

	ta := js.TypedArrayOf(b)
	gl.TexImage2D(int(target), int(level), int(internalFormat), int(width), int(height), int(border), int(format), int(xtype), ta)
	ta.Release()
}

func DeleteTextures(n int32, textures *uint32) {
	gl.DeleteTexture(textureMap[*textures])
	delete(textureMap, *textures)
}

func GenTextures(n int32, textures *uint32) {
	x := gl.CreateTexture()
	if x == js.Null() {
		*textures = 0
		return
	}
	textureCount++
	textureMap[textureCount] = x
	*textures = textureCount
}

func TexParameteri(texture uint32, pname uint32, params int32) {
	gl.TexParameteri(int(texture), int(pname), int(params))
}
