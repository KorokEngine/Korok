// bk-api provide low-level graphics api.
package bk

import (
	"korok.io/korok/math/f32"
	"unsafe"
	"log"
)

// SetDebug set debug Flags, DebugResMan enable ResourceManager's log output and
// DebugQueue enable RenderQueue's log output.
func SetDebug(debug uint32) {
	gDebug = debug
}

// Init init the bk-api.
func Init() {
	R.Init()
	gRenderQ.Init()
}

// Reset resets RenderContext's internal state, such as frame-buffer size.
func Reset(width, height uint32, pixelRatio float32) {
	gRenderQ.Reset(uint16(width), uint16(height), pixelRatio)
}

// Destroy release any resource used by the bk-api.
func Destroy() {
	gRenderQ.Destroy()
	R.Destroy()
}

// SetState set render's states for drawCall primitive.
// State flags is defined  by ST_BLEND.
func SetState(state uint64, rgba uint32) {
	gRenderQ.SetState(state, rgba)
}

// SetIndexBuffer sets index buffer for drawCall primitive.
func SetIndexBuffer(id uint16, firstIndex, num uint32) {
	gRenderQ.SetIndexBuffer(id, uint16(firstIndex), uint16(num))
}

// SetVertexBuffer sets vertex buffer for drawCall primitive.
func SetVertexBuffer(stream uint8, id uint16, firstVertex, numVertex uint32) {
	gRenderQ.SetVertexBuffer(stream, id, uint16(firstVertex), uint16(numVertex))
}

// SetSprite sets texture stages for drawCall primitive.
func SetTexture(stage uint8, sampler uint16, handle uint16, flags uint32) {
	gRenderQ.SetTexture(stage, sampler, handle, flags)
}

// SetTransform sets Model matrix.
func SetTransform(mtx *f32.Mat4) {
	gRenderQ.SetTransform(mtx)
}

// SetUniform sets shader uniform parameter for drawCall primitive.
func SetUniform(id uint16, ptr unsafe.Pointer) {
	gRenderQ.SetUniform(id, ptr)
}

// SetStencil sets stencil test state.
func SetStencil(stencil uint32) {
	gRenderQ.SetStencil(stencil)
}

// SetScissor set scissor for drawCall primitive. Return a cached ScissorRect handler.
func SetScissor(x, y, width, height uint16) uint16 {
	return gRenderQ.SetScissor(x, y, width, height)
}

// SetScissorCached set scissor rect with a cached ScissorRect handler.
func SetScissorCached(id uint16) {
	gRenderQ.SetScissorCached(id)
}

// SetViewScissor view scissor. Draw primitive outsize view will be clipped. When
// x, y, with, height are set to 0, scissor will be disabled.
func SetViewScissor(id uint8, x, y, width, height uint16) {
	gRenderQ.SetViewScissor(id, x, y, width, height)
}

func SetViewPort(id uint8, x, y, width, height uint16) {
	gRenderQ.SetViewPort(id, x, y, width, height)
}

// Set view clear Flags. rgba is Color clear value(default = 0x000000ff), depth is Depth clear value
// (default = 1.0), stencil is Stencil clear value(default = 0).
func SetViewClear(id uint8, flags uint16, rgba uint32, depth float32, stencil uint8) {
	gRenderQ.SetViewClear(id, flags, rgba, depth, stencil)
}

// SetViewTransform sets view and projection matrices, all drawCall primitives in this
// view will use these matrices.
func SetViewTransform(id uint8, view, proj *f32.Mat4, flags uint8) {
	gRenderQ.SetViewTransform(id, view, proj, flags)
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
	return gRenderQ.Submit(id, program, depth)
}

// Reset DrawCall state
func ResetDrawCall() {
	gRenderQ.drawCall.reset()
}

// Execute final draw
func Flush() int {
	return gRenderQ.Flush()
}

func Dump() {
	size := gRenderQ.drawCallNum;
	drawCall := gRenderQ.drawCallList[:size]
	log.Println("drawCall:", drawCall)
}

// Global Resources Manager!
var R *ResManager

const (
	DebugResMan uint32 = 0x000000001
	DebugQueue  uint32 = 0x000000002
)

// -- private field
var gDebug uint32
var gRenderQ *RenderQueue

func init() {
	R = NewResManager()
	// after res-manager!
	gRenderQ = NewRenderQueue(R)
}
