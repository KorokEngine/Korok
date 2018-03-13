// bk-api provide low-level graphics api.
package bk

import (
	"korok.io/korok/math/f32"
	"unsafe"
)

// SetDebug set debug Flags, DEBUG_R enable ResourceManager's log output and
// DEBUG_Q enable RenderQueue's log output.
func SetDebug(debug uint32) {
	g_debug = debug
}

// Init init the bk-api.
func Init() {
	R.Init()
	g_renderQ.Init()
}

// Reset resets RenderContext's internal state, such as frame-buffer size.
func Reset(width, height uint32) {
	g_renderQ.Reset(uint16(width), uint16(height))
}

// Destroy release any resource used by the bk-api.
func Destroy() {
	g_renderQ.Destroy()
	R.Destroy()
}

// SetState set render's states for drawCall primitive.
// State flags is defined  by ST_BLEND.
func SetState(state uint64, rgba uint32) {
	g_renderQ.SetState(state, rgba)
}

// SetIndexBuffer sets index buffer for drawCall primitive.
func SetIndexBuffer(id uint16, firstIndex, num uint32) {
	g_renderQ.SetIndexBuffer(id, uint16(firstIndex), uint16(num))
}

// SetVertexBuffer sets vertex buffer for drawCall primitive.
func SetVertexBuffer(stream uint8, id uint16, firstVertex, numVertex uint32) {
	g_renderQ.SetVertexBuffer(stream, id, uint16(firstVertex), uint16(numVertex))
}

// SetTexture sets texture stages for drawCall primitive.
func SetTexture(stage uint8, sampler uint16, handle uint16, flags uint32) {
	g_renderQ.SetTexture(stage, sampler, handle, flags)
}

// SetTransform sets Model matrix.
func SetTransform(mtx *f32.Mat4) {
	g_renderQ.SetTransform(mtx)
}

// SetUniform sets shader uniform parameter for drawCall primitive
func SetUniform(id uint16, ptr unsafe.Pointer) {
	g_renderQ.SetUniform(id, ptr)
}

// SetStencil sets stencil test state.
func SetStencil(stencil uint32) {
	g_renderQ.SetStencil(stencil)
}

// SetScissor set scissor for drawCall primitive.
func SetScissor(x, y, width, height uint16) {
	g_renderQ.SetScissor(x, y, width, height)
}

// SetViewScissor view scissor. Draw primitive outsize view will be clipped. When
// x, y, with, height are set to 0, scissor will be disabled.
func SetViewScissor(id uint8, x, y, width, height uint16) {
	g_renderQ.SetViewScissor(id, x, y, width, height)
}

func SetViewPort(id uint8, x, y, width, height uint16) {
	g_renderQ.SetViewPort(id, x, y, width, height)
}

// Set view clear Flags. rgba is Color clear value(default = 0x000000ff), depth is Depth clear value
// (default = 1.0), stencil is Stencil clear value(default = 0).
func SetViewClear(id uint8, flags uint16, rgba uint32, depth float32, stencil uint8) {
	g_renderQ.SetViewClear(id, flags, rgba, depth, stencil)
}

// SetViewTransform sets view and projection matrices, all drawCall primitives in this
// view will use these matrices.
func SetViewTransform(id uint8, view, proj *f32.Mat4, flags uint8) {
	g_renderQ.SetViewTransform(id, view, proj, flags)
}

// Submit an empty primitive for rendering. Uniforms and drawCall state
// will be applied but no geometry will be submitted. These empty drawCall
// calls will sort before ordinary drawCall calls.
func Touch(id uint8) uint32 {
	return Submit(id, InvalidId, 0)
}

// Submit primitive for rendering. Default depth is zero.
// Returns Number of drawCall calls.
func Submit(id uint8, program uint16, depth int32) uint32 {
	return g_renderQ.Submit(id, program, depth)
}

// Reset DrawCall state
func ResetDrawCall() {
	g_renderQ.drawCall.reset()
}

// Execute final draw
func Flush() uint32 {
	return g_renderQ.Flush()
}

// Global Resources Manager!
var R *ResManager

const (
	DEBUG_R uint32 = 0x000000001
	DEBUG_Q uint32 = 0x000000002
)

// -- private field
var g_debug uint32
var g_renderQ *RenderQueue

func init() {
	R = NewResManager()
	// after res-manager!
	g_renderQ = NewRenderQueue(R)
}
