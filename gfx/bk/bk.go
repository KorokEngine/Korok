package bk

import (
	"github.com/go-gl/mathgl/mgl32"
)

/// Set debug Flags
func SetDebug(debug uint32) {
	g_debug = debug
}

/// Set render states for drawCall primitive
///
/// @param state State flags. Default state for primitive type is
///        triangles. See: `BGFX_STATE_DEFAULT`
///        - ``
/// @param rgba Sets blend factor used by `BGFX_STATE_BLEND_FACTOR` and
///        `BGFX_STATE_BLEND_INV_FACTOR` blend modes
func SetState(state uint64, rgba uint32) {

}

/// Set index buffer for drawCall primitive
///
/// @param handle Index buffer
/// @param firstIndex First index to buffer
/// @param numIndices Number of indices to render
func SetIndexBuffer(id uint16, firstIndex, num uint32) {

}


/// Set vertex buffer for drawCall primitive
///
/// @param stream Vertex stream
/// @param handle Vertex buffer
/// @param startVertex First vertex to render
/// @param numVertex Number of vertices to render
func SetVertexBuffer(stream uint8, id uint16, firstVertex, numVertex uint32) {

}

/// Set texture stages for drawCall primitive
///
/// @param stage Texture unit
/// @param sampler Program sampler
/// @param handle Texture handle
/// @param flags Texture sampling mode, default=uint32_max
func SetTexture(stage uint8, sampler uint16, handle uint16, flags uint32) {

}


/// Set shader uniform parameter for drawCall primitive
///
/// @param handle Uniform
/// @param value Pointer to uniform data
/// @param num Number of elements. Passing `uint16_max` will
///        use the num passed on uniform creation, default=1
func SetUniform(id uint16, value interface{}, num uint16) {

}


/// Set stencil test state
///
/// @param stencil Stencil state
func SetStencil(stencil uint32) {

}

/// Set scissor for drawCall primitive. For scissor for all primitives in
/// view see `bgfx.SetViewScissor`
///
/// @param x, y Position from  left-top corner of the window
/// @param width, height Width, Height of scissor region
/// @return Scissor cache index
func SetScissor(x, y, width, height uint16) {

}

/// Set view scissor. Draw primitive outsize view will be clipped. When
/// x, y, with, height are set to 0, scissor will be disabled.
///
/// @param id View id
/// @param x Position x from the left corner of the window
/// @param y Position y from the top corner of the window
func SetViewScissor(id uint8, x, y, with, height uint16) {

}


/// Set view clear Flags
///
/// @param rgba Color clear value, default = 0x000000ff
/// @param depth Depth clear value, default = 1.0
/// @param stencil Stencil clear value, default = 0
func SetViewClear(id uint8, flags uint16, rgba uint32, depth float32, stencil uint8) {

}

/// Set view view and projection matrices, all drawCall primitives in this
/// view will use these matrices.
///
/// @param id View id
/// @param view View matrix
/// @param proj Projection matrix. When using stero rendering this projection matrix
///				 represent projection matrix for left eye
/// @param flags View flags. default=BGFX_VIEW_STEREO
func SetViewTransform(id uint8, view, proj *mgl32.Mat4, flags uint8) {

}

/// Submit an empty primitive for rendering. Uniforms and drawCall state
/// will be applied but no geometry will be submitted.
///
/// These empty drawCall calls will sort before ordinary drawCall calls.
/// @param id View id
/// @param Number of drawCall calls
func Touch(id uint8) uint32{
	return Submit(id, 0, 0, false)
}

/// Submit primitive for rendering
///
/// @param id View id
/// @param program Program
/// @param depth Depth for sorting, default=0
/// @param preserveState Preserve internal drawCall state for next drawCall call submit, default=false
/// @return Number of drawCall calls
func Submit(id uint8, program uint16, depth int32, preserveState bool) uint32{
	return 0
}


/// Global Resources Manager!
var R *ResManager

/// private field
var g_debug uint32


