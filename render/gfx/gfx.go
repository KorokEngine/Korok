package gfx

import (
	"log"
	"korok/math"
)

type RendererType uint16
const (
	RENDERER_TYPE_NOOP RendererType = iota  // No rendering
	RENDERER_TYPE_OPENGL_ES					// OpenGL ES 2.0+
	RENDERER_TYPE_VULKAN					// Vulkan

	RENDERER_TYPE_COUNT						// Max Enum count
)

/// Memory release callback
type ReleaseFn func (ptr interface{}, userData interface{})

func Init(_type RendererType, _vendorId uint16, _deviceId uint16) bool{
	if s_ctx != nil {
		log.Println("bgfx is already initialized.")
		return false
	}

	log.Println("Init...")

	if s_ctx.init(_type) {
		log.Println("Init complete.")
		return true
	}

	return false
}

func ShutDown() {
	log.Println("Shutdown...")

	ctx := s_ctx	// it's going to be null inside shutdown
	ctx.shutdown()

	s_ctx = nil
}

/// Reset graphic settings and back-buffer size
///
/// @param width Back-buffer width
/// @param height Back-buffer height
/// @param Flags TODO
func Reset(width, height uint32, flags uint32) {
	s_ctx.reset(width, height, flags)
}

/// Do frame rendering.
///
/// @param capture Capture frame with graphics debugger, default=false
/// @return Current frame number
func NextFrame(capture bool) uint32{
	return s_ctx.NextFrame(capture)
}

/// Returns current renderer backend API type
func GetRenderType() RendererType {
	return RENDERER_TYPE_COUNT
}

/// Return performance counters
func GetStats() *Stats{
	return s_ctx.getPerfStats()
}

func GetCaps() *Caps{
	return g_caps
}

/// Set debug Flags
func SetDebug(debug uint32) {
	s_ctx.setDebug(debug)
}

/// Create static index buffer
/// @param mem Index buffer data
/// @param Flags Buffer creation Flags, default = BGFX_BUFFER_NONE
func CreateIndexBuffer(mem *Memory, flags uint16) IndexBufferHandle {
	return s_ctx.createIndexBuffer(mem, flags)
}

/// Destroy static index buffer
func DestroyIndexBuffer(handle IndexBufferHandle) {
	s_ctx.destroyIndexBuffer(handle)
}

/// Create static vertex buffer
///
/// @param mem Vertex buffer data
/// @param layout Vertex declaration
/// @param Flags Buffer creation Flags, default= BGFX_BUFFER_NONE
func CreateVertexBuffer(mem *Memory, layout *VertexLayout, flags uint16) VertexBufferHandle{
	return s_ctx.createVertexBuffer(mem, layout, flags)
}

/// Destroy static vertex buffer
func DestroyVertexBuffer(handle VertexBufferHandle) {
	s_ctx.destroyVertexBuffer(handle)
}

/// Set view rectangle. Draw primitive outsize view will be clipped
///
/// @param id View id
/// @param x Position x from the left corner of the window
/// @param y Position y from the top corner of the window
func SetViewRect(id uint8, x, y, width, height uint16) {
	s_ctx.setViewRect(id, x, y, width, height)
}

/// Set view scissor. Draw primitive outsize view will be clipped. When
/// x, y, with, height are set to 0, scissor will be disabled.
///
/// @param id View id
/// @param x Position x from the left corner of the window
/// @param y Position y from the top corner of the window
func SetViewScissor(id uint8, x, y, with, height uint16) {
	s_ctx.setViewScissor(id, x, y, with, height)
}

/// Set view clear Flags
///
/// @param rgba Color clear value, default = 0x000000ff
/// @param depth Depth clear value, default = 1.0
/// @param stencil Stencil clear value, default = 0
func SetViewClear(id uint8, flags uint16, rgba uint32, depth float32, stencil uint8) {
	s_ctx.setViewClear(id, flags, rgba, depth, stencil)
}

/// Set view sort mode
///
/// @param id View id
/// @param mode View sort mode
func SetViewMode(id uint8, mode ViewMode) {
	s_ctx.setViewMode(id, mode)
}

/// Set view frame buffer
///
/// @param id View id
/// @param handle Frame buffer handle
func SetViewFrameBuffer(id uint8, handle FrameBufferHandle) {

}

/// Set view view and projection matrices, all draw primitives in this
/// view will use these matrices.
///
/// @param id View id
/// @param view View matrix
/// @param projL Projection matrix. When using stero rendering this projection matrix
///				 represent projection matrix for left eye
/// @param flags View flags. default=BGFX_VIEW_STEREO
/// @param projR Projection matrix for right eye in stereo mode. default=nil
func SetViewTransform(id uint8, view, projL *Matrix4, flags uint8) {
	s_ctx.setViewTransform(id, view, projL, flags)
}

/// Reset all view settings to default
func ResetView(id uint8) {
	s_ctx.resetView(id)
}

/// Set render states for draw primitive
///
/// @param state State flags. Default state for primitive type is
///        triangles. See: `BGFX_STATE_DEFAULT`
///        - ``
/// @param rgba Sets blend factor used by `BGFX_STATE_BLEND_FACTOR` and
///        `BGFX_STATE_BLEND_INV_FACTOR` blend modes
func SetState(state uint64, rgba uint32) {
	s_ctx.setState(state, rgba)
}

/// Set stencil test state
///
/// @param fstencil Front stencil state
/// @param bstencil Back stencil state. If back is set to `BGFX_STENCIL_NONE`
///        fstencil is applied to both front and back facing primitives
func SetStencil(fstencil uint32, bstencil uint32) {
	s_ctx.setStencil(fstencil, bstencil)
}

/// Set scissor for draw primitive. For scissor for all primitives in
/// view see `bgfx.SetViewScissor`
///
/// @param x, y Position from  left-top corner of the window
/// @param width, height Width, Height of scissor region
/// @return Scissor cache index
func SetScissor(x, y, width, height uint16) uint16{
	return s_ctx.setScissor(x, y, width, height)
}

/// Set scissor from cache for draw primitive
///
/// @param cache Index in scissor cache. Passing uint16_max unset primitive
///        scissor and primitive will use view scissor instead
func SetScissorCached(cache uint16) {
	s_ctx.setScissorCache(cache)
}

/// Set model matrix for draw primitive. If it's not called model will
/// be rendered with identity model matrix
///
/// @param mtx Pointer to first matrix in array
/// @param num Number of matrices in array. default=1
/// @return index into matrix cache in case the same model matrix has
/// 	    has to used for other draw primitive call
func SetTransform(mtx interface{}, num uint16) uint32{
	return s_ctx.setTransform(mtx, num)
}

/// Reserve `num` matrices in internal matrix cache
///
/// @param transform Pointer to `Transform` structure
/// @param num Number of matrices
/// @return index into matrix cache
func AllocTransform(transform *Transform, num uint16) uint32 {
	return s_ctx.allocTransofrm(transform, num)
}

/// Set model matrix from matrix cache for draw primitive
///
/// @param cache Index in matrix cache
/// @param num Number of matrix from cache, default=1
func SetTransformCache(cache uint32, num uint16) {
	s_ctx.setTransformCache(cache, num)
}

/// Set shader uniform parameter for draw primitive
///
/// @param handle Uniform
/// @param value Pointer to uniform data
/// @param num Number of elements. Passing `uint16_max` will
///        use the num passed on uniform creation, default=1
func SetUniform(handle UniformHandle, value interface{}, num uint16) {
	s_ctx.setUniform(handle, value, num)
}

/// Set index buffer for draw primitive
///
/// @param handle Index buffer
//func SetIndexBuffer(handle IndexBufferHandle) {
//
//}

/// Set index buffer for draw primitive
///
/// @param handle Index buffer
/// @param firstIndex First index to buffer
/// @param numIndices Number of indices to render
func SetIndexBuffer(handle IndexBufferHandle, firstIndex, numIndices uint32) {
	s_ctx.setIndexBuffer(handle, firstIndex, numIndices)
}

/// Set index buffer for draw primitive
func SetDynamicIndexBuffer(handle DynamicIndexBufferHandle, firstIndex, numIndices uint32) {
	s_ctx.setDynamicIndexBuffer(handle, firstIndex, numIndices)
}

/// Set index buffer for draw primitive
func SetTransientIndexBuffer(tib *TransientIndexBuffer, firstIndex, numIndices uint32) {
	s_ctx.setTransientIndexBuffer(tib, firstIndex, numIndices)
}

/// Set vertex buffer for draw primitive
///
/// @param stream Vertex stream
/// @param handle Vertex buffer
/// @param startVertex First vertex to render
/// @param numVertex Number of vertices to render
func SetVertexBuffer(stream uint8, handle VertexBufferHandle, startVertex, numVertices uint32) {
	s_ctx.setVertexBuffer(stream, handle, startVertex, numVertices)
}

/// Set vertex buffer for draw primitive
///
/// @param stream Vertex stream
/// @param handle Vertex buffer
/// @param startVertex First vertex to render
/// @param numVertex Number of vertices to render
func SetDynamicVertexBuffer(stream uint8, handle DynamicVertexBufferHandle, startVertex, numVertices uint32) {
	s_ctx.setDynamicVertexBuffer(stream, handle, startVertex, numVertices)
}

/// Set vertex buffer for draw primitive
///
/// @param stream Vertex stream
/// @param handle Vertex buffer
/// @param startVertex First vertex to render
/// @param numVertex Number of vertices to render
func SetTransientVertexBuffer(stream uint8, tvb *TransientVertexBuffer, startVertex, numVertices uint32) {
	s_ctx.setTransientVertexBuffer(stream, tvb, startVertex, numVertices)
}

/// Set texture stages for draw primitive
///
/// @param stage Texture unit
/// @param sampler Program sampler
/// @param handle Texture handle
/// @param flags Texture sampling mode, default=uint32_max
func SetTexture(stage uint8, sampler UniformHandle, handle TextureHandle, flags uint32) {

}

/// Submit an empty primitive for rendering. Uniforms and draw state
/// will be applied but no geometry will be submitted.
///
/// These empty draw calls will sort before ordinary draw calls.
/// @param id View id
/// @param Number of draw calls
func Touch(id uint8) uint32{
	return Submit(id, INVALID_HANDLE, 0, false)
}

/// Submit primitive for rendering
///
/// @param id View id
/// @param program Program
/// @param depth Depth for sorting, default=0
/// @param preserveState Preserve internal draw state for next draw call submit, default=false
/// @return Number of draw calls
func Submit(id uint8, program ProgramHandle, depth int32, preserveState bool) uint32{
	return s_ctx.Submit(id, program, depth, preserveState)
}

/////// Resource: Shader, Program, Uniform

func CreateShader(mem *Memory) ShaderHandle{
	return s_ctx.createShader(mem)
}

func GetShaderUniforms(handle ShaderHandle) []UniformHandle {
	return s_ctx.getShaderUniforms(handle)
}

func DestroyShader(handle ShaderHandle) {
	s_ctx.destroyShader(handle)
}

/// default: destroyShader=false
func CreateProgram(vsh, fsh ShaderHandle, destroyShader bool) ProgramHandle {
	return s_ctx.createProgram(vsh, fsh, destroyShader)
}

/// Destroy program
func DestroyProgram(handle ProgramHandle) {
	s_ctx.destroyProgram(handle)
}

/// Create shader uniform
///
/// @param name Uniform name in shader
/// @param uType Type of uniform
/// @param num Number of elements in array
///
/// default:num=1
func CreateUniform(name string, uType UniformType, num uint16) UniformHandle {
	return s_ctx.createUniform(name, uType, num)
}

/// Retrieve uniform info
///
/// @param handle Handle to uniform object
/// @param info Return uniform info
func GetUniformInfo(handle UniformHandle) *UniformInfo {
	return s_ctx.getUniformInfo(handle)
}

/// Destroy shader uniform
func DestroyUniform(handle UniformHandle) {
	s_ctx.destroyUniform(handle)
}

/// Calculate amount of memory required for texture
///
/// @param info Resulting texture info structure
/// @param width, height
/// @param depth Depth dimension of volume texture
/// @param cubeMap Indicates that texture contains cube-map
/// @param hasMips Indicates that texture contains full mip-map chain
/// @param numLayers Number of layers in texture array
/// @param format Texture Format
func CalcTextureSize(info *TextureInfo, width, height, depth uint16, cubeMap, hasMips bool, numLayers uint16, format TextureFormat) {
	// TODO bimg.imageGetSize
}

/// Create texture from memory buffer
///
/// @param mem Texture data
/// @param flags Texture sampling mode
/// @param skip Skip top level mips when parsing texture
/// @param info Return parsed texture information
///
/// default: flags=TEXTURE_NONE, skip=0, info=nil
func CreateTexture(mem *Memory, flags uint32, skip uint8, info *TextureInfo) TextureHandle{
	return s_ctx.createTexture(mem, flags, skip, info, 0)
}

/// Create 2D texture
/// @param width, height
/// @param format Texture format
/// @param flags Texture sampling mode
/// @param mem Texture data
///
/// default: flags=TEXTURE_NONE, mem=nil
func CreateTexture2D(width, height uint16, format TextureFormat, flags uint32, mem *Memory) TextureHandle{
	if width < 0 || height < 0 {
		log.Printf("Invalid texture size (width %d, height %d).", width, height)
	}
	return createTexture2D(BBR_COUNT, width, height, format, flags, mem)
}

/// Create Texture with size based on back buffer ratio.
///
/// @param ratio Texture size in respect to back-buffer size
/// @param format Texture format
/// @param flags Texture sampling mode
///
/// default: flags=TEXTURE_NONE
func CreateTexture2DScaled(ratio BackBufferRatio, format TextureFormat, flags uint32) TextureHandle {
	if ratio > BBR_COUNT {
		log.Println("Invalid back buffer ratio.")
	}
	return createTexture2D(ratio, 0, 0, format, flags, nil)
}

func createTexture2D(ratio BackBufferRatio, width, height uint16, format TextureFormat, flags uint32, mem *Memory) TextureHandle{
	if ratio != BBR_COUNT {
		width = uint16(s_ctx.resolution.width)
		height = uint16(s_ctx.resolution.height)
		getTextureSizeFromRatio(ratio, &width, &height)
	}

	/// 把数据写入 Memo!! 以便用来创建 Texture !!
	mem := new(Memory)
	return s_ctx.createTexture(mem, flags, 0, nil, ratio)
}

func getTextureSizeFromRatio(ratio BackBufferRatio, width, height *uint16) {
	switch ratio {
	case BBR_HALF:
		*width, *height = *width/2, *height/2
	case BBR_QUARTER:
		*width, *height = *width/4, *height/4
	case BBR_EIGHTH:
		*width, *height = *width/8, *height/8
	case BBR_SIXTEENTH:
		*width, *height = *width/16, *height/16
	case BBR_DOUBLE:
		*width, *height = *width*2, *height*2
	}
	*width = math.UInt16_max(1, *width)
	*height = math.UInt16_max(1, *height)
}

/// default: pitch=UINT16_MAX
func UpdateTexture2D(handle TextureHandle, layer uint16, mip uint8, x, y, width, height uint16, mem *Memory, pitch uint16) {
	if mem == nil {
		log.Println("mem can't be nil")
	}
	if width != 0 && height != 0 {
		s_ctx.updateTexture(handle, 0, mip, x, y, layer, width, height, 1, pitch, mem)
	}
}

/// Read back texture content
/// @param handle Texture handle
/// @param data Destination buffer
/// @param mip Mip level

/// default: mip=0
func ReadTexture(handle TextureHandle, data interface{}, mip uint8) (frameNum uint32) {
	log.Fatal("ReadTexture not impl")
	return
}

/// Destroy texture
func DestroyTexture(handle TextureHandle) {
	s_ctx.destroyTexture(handle)
}

/// Create frame buffer (simple)
///
/// @param width, height Frame buffer Size
/// @param format Texture format
/// @param textureFlags Texture sampling mode
///
/// default: textureFlags=TEXTURE_U_CLAMP|TEXTURE_V_CLAMP
func CreateFrameBuffer(width, height uint16, format TextureFormat, textureFlags uint32) FrameBufferHandle{
	texHandle := CreateTexture2D(width, height, format, textureFlags, nil)
	return CreateFrameBufferFromTexture(1, []TextureHandle{texHandle}, true)
}

/// Create frame buffer with size based on back-buffer ratio. Frame buffer will maintain ratio
/// if back buffer resolution changes.
///
/// @param ratio Frame buffer size in respect to back-buffer size.
/// @param format Texture format
/// @param textureFlags Texture sampling mode
///
/// default: textureFlags=TEXTURE_U_CLAMP|TEXTURE_V_CLAMP
func CreateFrameBufferScaled(ratio BackBufferRatio, format TextureFormat, textureFlags uint32) FrameBufferHandle{
	if ratio > BBR_COUNT {
		log.Println("Invalid back buffer ratio")
	}

	texHandle := CreateTexture2DScaled(ratio, format, textureFlags)
	return CreateFrameBufferFromTexture(1, []TextureHandle{texHandle}, true)
}

/// Create MRT frame buffer from texture handles (simple)
///
/// @param num Number of texture attachments
/// @param handles Texture attachments
/// @param destroyTexture If true, textures will be destroyed when frame buffer is destroyed
///
/// default: destroyTexture=false
func CreateFrameBufferFromTexture(num uint8, handles []TextureHandle, destroyTexture bool) FrameBufferHandle{
	attachment := [CONFIG_MAX_FB_ATTACHMENTS]Attachment{}
	for i := 0; i < int(num); i++ {
		at := attachment[i]
		at.handle = handles[i]
		at.mip = 0
		at.layer = 0
	}
	return CreateFrameBufferFromAttachment(num, attachment[:], destroyTexture)
}

/// Create MRT frame buffer from texture handles with specific layer and mip level
///
/// @param num Number of texture attachments
/// @param handles Texture attachments
/// @param destroyTexture If true, textures will be destroyed when frame buffer is destroyed
///
/// default: destroyTexture=false
func CreateFrameBufferFromAttachment(num uint8, attachment []Attachment, destroyTexture bool) FrameBufferHandle{
	if num == 0 {
		log.Println("Number of frame buffer attachment can't be 0.")
	}
	if uint32(num) > CONFIG_MAX_FB_ATTACHMENTS {
		log.Printf("Number of frame buffer attachments is larger than allowed %d (max: %d)",
					num,
					CONFIG_MAX_FB_ATTACHMENTS)
	}
	if attachment == nil {
		log.Println("attachement can't be nil")
	}
	return s_ctx.createFrameBuffer(num, attachment, destroyTexture)
}

/// Obtain texture handle of frame buffer attachment
///
/// @param handle Frame buffer handle
/// @param attachment Frame buffer attachment index
///
/// default: attachment=0
func GetTexture(handle FrameBufferHandle, attachment uint8) TextureHandle{
	return s_ctx.getTexture(handle, attachment)
}

/// Destroy frame buffer
func DestroyFrameBuffer(handle FrameBufferHandle) {
	s_ctx.destroyFrameBuffer(handle)
}

/////// static & global var

var s_ctx *Context
var g_caps *Caps