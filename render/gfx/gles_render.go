package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"korok/render/bx"
	"korok/math"
	"log"
	"unsafe"
)

type PrimInfo struct {
	xType 	uint32
	min 	uint32
	div		uint32
	sub 	uint32
}

type Blend struct {
	src uint32
	dst uint32
	factor 	bool
}

type TextureFormatInfo struct {
	internalFmt 	uint32
	internalFmtSrgb uint32
	fmt 			uint32
	xType 			uint32
	supported 		bool
}

type VendorId struct {
	name string
	id uint16
}

type Workaround struct {
	detachShader bool
}

func (*Workaround) reset() {

}

type RendererContextGL struct {
	numWindows 	uint16
	windows	[CONFIG_MAX_FRAME_BUFFERS]FrameBufferHandle

	uniformReg 		UniformRegistry

	fbh 	  FrameBufferHandle
	fbDiscard uint16

	resolution 		Resolution
	vao 			uint32

	// >= gl es 3.0
	vaoSupport 		bool
	samplerObjectSupport bool

	hash 	uint64

	readPixelsFmt 	uint32
	backBufferFbo 	uint32

	glctx	gl.Context
	needPresent bool

	vendor 		string
	renderer 	string
	version 	string
	glsVersion 	string

	workaround 	Workaround
}


func (ctx *RendererContextGL) init() bool {
	ctx.fbh.idx = kInvalidHandle
	var anyType interface{} // TODO interface size!!
	bx.MemZero(unsafe.Pointer(&ctx.uniforms[0]), len(ctx.uniforms) * int(unsafe.Sizeof(anyType)))
	ctx.resolution = Resolution{}

	ctx.SetRenderContextSize(DEFAULT_WIDTH, DEFAULT_HEIGHT)

	ctx.vendor 	   = GetGLString(gl.VENDOR)
	ctx.renderer   = GetGLString(gl.RENDERER)
	ctx.version    = GetGLString(gl.VERSION)
	ctx.glsVersion = GetGLString(gl.SHADING_LANGUAGE_VERSION)

	for i := range s_vendorIds {
		if v := &s_vendorIds[i]; v.name == ctx.vendor {
			g_caps.vendorId = v.id; break
		}
	}

	ctx.workaround.reset()

	var numCmpFormat int32
	gl.GetIntegerv(gl.NUM_COMPRESSED_TEXTURE_FORMATS, &numCmpFormat)
	log.Printf("gl.NUM_COMPRESSED_TEXTURE_FORMATS %d", numCmpFormat)

	/// not support compressed format!!

	/// Initial binary shader hash depends on driver version!!

	/// Initial extensions!

	if opengl es {
		ctx.workaround.detachShader = false
	}

	/// 计算纹理格式是否支持
	/// 目前仅支持简单的格式
	setTextureFormat(TEXTURE_FORMAT_RGBA16F, gl.RGBA, gl.RGBA, gl.HALF_FLOAT)
	setTextureFormat(TEXTURE_FORMAT_RGBA32F, gl.RGBA, gl.RGBA, gl.FLOAT)
	// internalFormat and format must match:
	// https://www.khronos.org/opengles/sdk/docs/man/xhtml/glTexImage2D.xml
	setTextureFormat(TEXTURE_FORMAT_RGBA8,  gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE)
	setTextureFormat(TEXTURE_FORMAT_R5G6B5, gl.RGB,  gl.RGB,  gl.UNSIGNED_SHORT_5_6_5_REV)
	setTextureFormat(TEXTURE_FORMAT_RGBA4,  gl.RGBA, gl.RGBA, gl.UNSIGNED_SHORT_4_4_4_4_REV)
	setTextureFormat(TEXTURE_FORMAT_RGB5A1, gl.RGBA, gl.RGBA, gl.UNSIGNED_SHORT_1_5_5_5_REV)

	// 一下初始化 g_caps 的一些变量！

	// Other!
	if CONFIG_RENDERER_OPENGLES >= 30 {
		ctx.vaoSupport = true
	} else {
		ctx.vaoSupport = false
	}

	if ctx.vaoSupport {
		gl.GenVertexArrays(1, &ctx.vao)
	}

	/// sampler object/ shadow sampler/program binary/texture swizzle
	/// depth texture/ timer query/occlusion/ato/
	/// srgb
	//// texture address
	//// 以上皆不支持！！

	ctx.needPresent = false

	return true
}

func (ctx *RendererContextGL) shutdown() {
	if ctx.vaoSupport {
		gl.BindVertexArray(0)
		gl.DeleteVertexArrays(1, &ctx.vao)
		ctx.vao = 0
	}

	ctx.invalidateCache()
}

func (ctx *RendererContextGL) getRendererType() RendererType {
	return RENDERER_TYPE_OPENGL_ES
}

func (ctx *RendererContextGL) getRendererName() string {
	return "OpenGL ES"
}
func (ctx *RendererContextGL) IsDeviceRemoved() bool {
	return false
}

func (ctx *RendererContextGL) SetRenderContextSize(width, height uint16) {

}

func (ctx *RendererContextGL) invalidateCache() {

}


/// static index and vertex buffer
func (ctx *RendererContextGL) CreateIndexBuffer(handle IndexBufferHandle, mem *Memory, flags uint16) {
	ctx.indexBuffers[handle.idx].create(mem.size, mem.data, flags)
}

func (ctx *RendererContextGL) DestroyIndexBuffer(handle IndexBufferHandle) {
	ctx.indexBuffers[handle.idx].destroy()
}

func (ctx *RendererContextGL) CreateVertexLayout(handle VertexLayoutHandle, layout *VertexLayout) {
	ctx.vertexLayouts[handle.idx] = *layout
	log.Println(layout)
}

func (ctx *RendererContextGL) DestroyVertexLayout(handle VertexLayoutHandle) {
	// nothing!
}

func (ctx *RendererContextGL) CreateVertexBuffer(handle VertexBufferHandle, mem *Memory, declHandle VertexLayoutHandle, flags uint16) {
	ctx.vertexBuffers[handle.idx].create(mem.size, mem.data, declHandle, flags)
}

func (ctx *RendererContextGL) DestroyVertexBuffer(handle VertexBufferHandle) {
	ctx.vertexBuffers[handle.idx].destroy()
}

/// dynamic index and vertex buffer
func (ctx *RendererContextGL) CreateDynamicIndexBuffer(handle IndexBufferHandle, size uint32, flags uint16) {
	ctx.indexBuffers[handle.idx].create(size, nil, flags)
}

func (ctx *RendererContextGL) UpdateDynamicIndexBuffer(handle IndexBufferHandle, offset uint32, size uint32, mem *Memory) {
	ctx.indexBuffers[handle.idx].update(offset, math.UInt32_min(size, mem.size), mem.data, false)
}

func (ctx *RendererContextGL) DestroyDynamicIndexBuffer(_handle IndexBufferHandle) {
	ctx.indexBuffers[_handle.idx].destroy()
}

func (ctx *RendererContextGL) CreateDynamicVertexBuffer(handle VertexBufferHandle, size uint32, flags uint16) {
	decl := VertexLayoutHandle(INVALID_HANDLE)
	ctx.vertexBuffers[handle.idx].create(size, nil, decl, flags)
}

func (ctx *RendererContextGL) UpdateDynamicVertexBuffer(handle VertexBufferHandle, offset uint32,  size uint32, mem *Memory) {
	ctx.vertexBuffers[handle.idx].update(offset, math.UInt32_min(size, mem.size), mem.data, false)
}

func (ctx *RendererContextGL) DestroyDynamicVertexBuffer(handle VertexBufferHandle) {
	ctx.vertexBuffers[handle.idx].destroy()
}

/// shader and program
func (ctx *RendererContextGL) CreateShader(handle ShaderHandle, mem *Memory) {
	ctx.shaders[handle.idx].create(mem)
}

func (ctx *RendererContextGL) DestroyShader(handle ShaderHandle) {
	ctx.shaders[handle.idx].destroy()
}

func (ctx *RendererContextGL) CreateProgram(handle ProgramHandle, vsh ShaderHandle, fsh ShaderHandle) {
	ctx.programs[handle.idx].create(&ctx.shaders[vsh.idx], &ctx.shaders[fsh.idx])
}

func (ctx *RendererContextGL) DestroyProgram(handle ProgramHandle) {
	ctx.programs[handle.idx].destroy()
}

/// textures
func (ctx *RendererContextGL) CreateTexture(handle TextureHandle,  mem *Memory, flags uint32, skip uint8) {
	ctx.textures[handle.idx].create(mem, flags, skip)
}

func (ctx *RendererContextGL) UpdateTextureBegin(_handle TextureHandle, _side uint8, _mip uint8) {

}

func (ctx *RendererContextGL) UpdateTexture(handle TextureHandle, side, mip uint8, rect Rect, z, depth, pitch uint16, mem *Memory) {
	ctx.textures[handle.idx].update(side, mip, rect, z, depth, pitch, mem)
}

func (ctx *RendererContextGL) UpdateTextureEnd() {

}

func (ctx *RendererContextGL) ReadTexture(_handle TextureHandle, _data interface{}, _mip uint8) {
	// no impl
}

func (ctx *RendererContextGL) ResizeTexture(handle TextureHandle, width, height uint16, numMips uint8) {
	// todo no impl
}

func (ctx *RendererContextGL) OverrideInternal(handle TextureHandle, ptr uintptr) {
	ctx.textures[handle.idx].overrideInternal(ptr)
}

func (ctx *RendererContextGL) GetInternal(handle TextureHandle) uintptr {
	return uintptr(ctx.textures[handle.idx].id)
}

func (ctx *RendererContextGL) DestroyTexture(handle TextureHandle) {
	ctx.textures[handle.idx].destroy()
}

/// frame buffer
func (ctx *RendererContextGL) CreateFrameBufferAttachment(handle FrameBufferHandle, num uint8, attachment []Attachment) {
	ctx.frameBuffers[handle.idx].createWithAttachment(num, attachment)
}

func (ctx *RendererContextGL) CreateFrameBuffer(handle FrameBufferHandle, nwh interface{}, width, height uint32, depthFormat TextureFormat) {
	ctx.numWindows ++
	denseIdx := ctx.numWindows
	ctx.windows[denseIdx] = handle
	ctx.frameBuffers[handle.idx].create(denseIdx, nwh, width, height, depthFormat)	// todo
}

func (ctx *RendererContextGL) DestroyFrameBuffer(handle FrameBufferHandle) {
	denseIdx := ctx.frameBuffers[handle.idx].destroy()
	if denseIdx != UINT16_MAX {
		ctx.numWindows --
		if ctx.numWindows > 1 {
			handle := ctx.windows[ctx.numWindows]
			ctx.windows[denseIdx] = handle
			ctx.frameBuffers[handle.idx].denseIdx = denseIdx
		}
	}
}

/// uniform
func (ctx *RendererContextGL) CreateUniform(handle UniformHandle, uType UniformType, num uint16, name string) {
	if nil != ctx.uniforms[handle.idx] {
		// todo bx free
	}
	// todo
}

func (ctx *RendererContextGL) DestroyUniform(_handle UniformHandle) {

}

func (ctx *RendererContextGL) UpdateViewName(id uint8, name string) {
	s_viewName[id] = name
}

/// memory copy!!
func (ctx *RendererContextGL) UpdateUniform(loc uint16, data interface{}, size uint32) {
	ctx.uniforms[loc] = data
}

// 最终方法
func (ctx *RendererContextGL) Submit(render *Frame) {
	// 2.是否关闭 VAO
	if ctx.numWindows > 1 && ctx.vaoSupport {
		ctx.vaoSupport = false
		gl.BindVertexArray(0)
		gl.DeleteVertexArrays(1, &ctx.vao)
		ctx.vao = 0
	}

	// TODO
	ctx.glctx.makeCurrent(nil)

	// 3. 绑定 VAO 和 FrameBuffer
	if defaultVao := ctx.vao; 0 != defaultVao {
		gl.BindVertexArray(defaultVao)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, ctx.backBufferFbo)

	// 4. 更新分辨率
	ctx.updateResolution(&render.resolution)

	/// time query ignore!!

	// 5. 更新 TransientBuffer 的数据
	if 0 < render.iboffset {
		ib := render.transientIb
		ctx.indexBuffers[ib.handle.idx].update(0, render.iboffset, ib.data, true)
	}
	if 0 < render.vboffset {
		vb := render.transientVb
		ctx.vertexBuffers[vb.handle.idx].update(0, render.vboffset, vb.data, true)
	}

	// 6. 排序
	render.Sort()

	// 7. 初始化一些临时变量
	currentState := RenderDraw{}
	currentState.clear()
	currentState.stateFlags = STATE_NONE
	currentState.stencil = packStencil(STENCIL_NONE, STENCIL_NONE)

	currentBind := RenderBind{}
	currentBind.clear()

	// hmd 直接忽略！！
	viewState := ViewState{}
	viewState.reset(render)

	programIdx := kInvalidHandle
	key := SortKey{}
	view := UINT16_MAX
	fbh := FrameBufferHandle{idx(CONFIG_MAX_FRAME_BUFFERS)}

	// blitState 直接忽略
	resolutionHeight := render.resolution.height
	blendFactor := uint32(0)

	primIndex := uint8(uint64(0) >> STATE_PT_SHIFT)
	prim := s_primInfo[primIndex]

	var viewHasScissor bool
	viewScissorRect := Rect{}
	viewScissorRect.clear()
	var discardFlags uint16 = CLEAR_NONE

	statsNumPrimsSubmitted := [len(s_primInfo)]uint32{}
	statsNumPrimsRendered  := [len(s_primInfo)]uint32{}
	statsNumIndices := uint32(0)
	statsKeyType 	:= [2]uint32{}

	eye := uint8(0)

	// 8. Render!!

	// msaa not supported! gl.BindFrameBuffer(gl.FrameBuffer, ctx.msaaBackBufferFbo)
	var viewRestart bool
	var restartState uint8
	viewState.rect = render.rect[0]

	var numItems = uint32(render.num)
	for item, restartItem:= uint32(0), numItems; item < numItems || restartItem < numItems; {
		encodedKey := render.sortKeys[item]

		viewChanged := uint16(key.view) != view || item == numItems

		itemIdx := render.sortValues[item]
		renderItem := render.renderItem[itemIdx]
		renderBind := render.renderBind[itemIdx]
		item ++

		if viewChanged {
			if 1 == restartState {
				restartState = 2
				item = restartItem
				restartItem = numItems
				view = UINT16_MAX
				continue
			}

			view = uint16(key.view)
			programIdx = kInvalidHandle

			if render.fb[view].idx != fbh.idx {
				fbh = render.fb[0]
				resolutionHeight = render.resolution.height
				resolutionHeight = ctx.setFrameBuffer(fbh, resolutionHeight, discardFlags)
			}

			viewRestart = VIEW_STEREO == (render.viewFlags[view] & VIEW_STEREO)
			viewRestart &= hmdEnabled

			if viewRestart {
				if 0 == restartState {
					restartState = 1
					restartItem = item - 1
				}

				eye = (restartState - 1) & 1
				restartState &= 1
			} else {
				eye  = 0
			}

			if item > 1 {
				// profiler.end()
			}

			viewState.rect = render.rect[view]
			if viewRestart {
				/// 此处设置 ViewPort！ 比较诡异 TODO
				viewState.rect.x = uint16(eye) * (viewState.rect.width+1)/2
				viewState.rect.width /= 2
			}

			scissorRect := render.scissor[view]
			viewHasScissor = !scissorRect.isZero()
			if viewHasScissor {
				viewScissorRect = scissorRect
			} else {
				viewScissorRect = viewState.rect
			}

			gl.Viewport(int32(viewState.rect.x),
				int32(uint16(resolutionHeight) - viewState.rect.height - viewState.rect.y),
				int32(viewState.rect.width),
				int32(viewState.rect.height))

			clear := render.clear[view]
			discardFlags = clear.flags & uint16(CLEAR_DISCARD_MASK)

			if CLEAR_NONE != (clear.flags & CLEAR_MASK) {
				ctx.clearQuad(clearQuad, viewState.rect, clear, resolutionHeight, render.colorPalette[:])
			}

			gl.Disable(gl.STENCIL_TEST)
			gl.Disable(gl.DEPTH_TEST)
			gl.DepthFunc(gl.LESS)
			gl.Enable(gl.CULL_FACE)
			gl.Disable(gl.BLEND)

			/// gl.submitBlit // TODO
		}

		/// if isCompute {} not support!
		resetState := viewChanged || wasCompute

		draw := RenderDraw(renderItem)

		/// occlusion query not supported!

		/// 求取变化的状态位
		newFlags := draw.stateFlags
		changedFlags := currentState.stateFlags ^ draw.stateFlags // TODO golang 异或
		currentState.stateFlags = newFlags

		newStencil := draw.stencil
		changedStencil := currentState.stencil ^ draw.stencil
		currentState.stencil = newStencil

		if resetState {
			currentState.clear()
			currentState.scissor = !draw.scissor
			changedFlags = STATE_MASK
			changedStencil = packStencil(STENCIL_MASK, STENCIL_MASK)
			currentState.stateFlags = newFlags
			currentState.stencil = newStencil

			currentBind.clear()
		}

		/// 计算 Scissor！！
		scissor := draw.scissor
		if currentState.scissor != scissor {
			currentState.scissor = scissor

			if UINT16_MAX == scissor {
				if viewHasScissor {
					gl.Enable(gl.SCISSOR_TEST)
					gl.Scissor(int32(viewScissorRect.x),
						int32(uint16(resolutionHeight) - viewScissorRect.height - viewScissorRect.y),
						int32(viewScissorRect.width),
						int32(viewScissorRect.height))
				}
			} else {
				scissorRect := Rect{}
				scissorRect.setIntersect(viewScissorRect, render.rectCache.cache[scissor])

				if scissorRect.isZeroArea() {
					continue
				}

				gl.Enable(gl.SCISSOR_TEST)
				gl.Scissor(int32(scissorRect.x),
					int32(resolutionHeight - scissorRect.height - scissorRect.y),
					int32(scissorRect.width),
					int32(scissorRect.height))
			}
		}

		if 0 != changedStencil {
			if 0 != newStencil {
				gl.Enable(gl.STENCIL_TEST)

				//// stencil not supported!!!
			}
		}

		if 0 != (0 |
			STATE_CULL_MASK |
			STATE_DEPTH_WRITE |
			STATE_DEPTH_TEST_MASK |
			STATE_RGB_WRITE |
			STATE_ALPHA_WRITE |
			STATE_BLEND_MASK |
			STATE_BLEND_EQUATION_MASK |
			STATE_ALPHA_REF_MASK |
			STATE_PT_MASK |
			STATE_POINT_SIZE_MASK |
			STATE_MSAA |
			STATE_LINEAA |
			STATE_CONSERVATIVE_RASTER) & changedFlags {

			if changedFlags & STATE_CULL_MASK != 0 {
				// not supported!!
			}

			if changedFlags & STATE_DEPTH_WRITE != 0 {
				gl.DepthMask(STATE_DEPTH_WRITE & newFlags != 0)
			}

			if changedFlags & STATE_DEPTH_TEST_MASK != 0 {
				_func := (newFlags & STATE_DEPTH_TEST_MASK) >> STATE_DEPTH_TEST_SHIFT

				if _func != 0 {
					gl.Enable(gl.DEPTH_TEST)
					gl.DepthFunc(s_cmpFunc[_func])
				} else {
					if newFlags & STATE_DEPTH_WRITE != 0 {
						gl.Enable(gl.DEPTH_TEST)
						gl.DepthFunc(gl.ALWAYS)
					} else {
						gl.Disable(gl.DEPTH_TEST)
					}
				}
			}

			if STATE_ALPHA_REF_MASK & changedFlags != 0 {
				ref := (newFlags & STATE_ALPHA_REF_MASK) >> STATE_ALPHA_REF_SHIFT
				viewState.alphaRef = float32(float64(ref)/255)
			}

			///// if Enable OpenGL 这里是 OpenGL 独有的配置吗？？
			if (STATE_PT_POINTS | STATE_POINT_SIZE_MASK) & changedFlags != 0{
				pointSize := math.UInt32_max(1, uint32((newFlags & STATE_POINT_SIZE_MASK) >> STATE_POINT_SIZE_SHIFT))
				gl.PointSize(float32(pointSize))
			}

			if STATE_MSAA & changedFlags != 0 {
				// not supported!
			}

			if STATE_LINEAA & changedFlags != 0 {
				// not supported!
			}

			if STATE_CONSERVATIVE_RASTER & changedFlags != 0 {
				// not supported!
			}
			///// end opengl

			if (STATE_ALPHA_WRITE|STATE_RGB_WRITE) & changedFlags != 0 {
				alpha := (newFlags & STATE_ALPHA_WRITE) != 0
				rgb   := (newFlags & STATE_RGB_WRITE) != 0
				gl.ColorMask(rgb, rgb, rgb, alpha)
			}

			/// 所谓 blend independent 可以实现顺序无关的 alpha 混合
			/// http://www.openglsuperbible.com/2013/08/20/is-order-independent-transparency-really-necessary/
			/// Alpha to coverage
			/// 好像是MSAA里面的技术，不懂

			if ((STATE_BLEND_MASK | STATE_BLEND_EQUATION_MASK) & newFlags) != 0 || (blendFactor != draw.rgba) {
				enabled := (STATE_BLEND_MASK & newFlags) != 0

				blend := uint32(newFlags & STATE_BLEND_MASK) >> STATE_BLEND_SHIFT
				srcRGB := (blend    ) & 0xFF
				dstRGB := (blend>> 4) & 0xFF
				srcA   := (blend>> 8) & 0xFF
				dstA   := (blend>>12) & 0xFF

				equ    := uint32(newFlags & STATE_BLEND_EQUATION_MASK) >> STATE_BLEND_EQUATION_SHIFT
				equRGB := (equ    ) & 0x7
				equA   := (equ>> 3) & 0x7

				numRt := ctx.getNumRt()


			}

			////////// 以上都是各种状态的绑定

			pt := newFlags & STATE_PT_MASK
			primIndex = uint8(pt >> STATE_PT_SHIFT)
			prim = s_primInfo[primIndex]
		} /// End state change

		var programChanged bool
		var constantsChanged = draw.constBegin < draw.constEnd
		var bindAttribs bool

		/// update uniform
		rendererUpdateUniforms(ctx, render.uniformBuffer, draw.constBegin, draw.constEnd)

		/// update program
		if key.program != uint16(programIdx) {
			programIdx = idx(key.program)
			var id uint32
			if programIdx != kInvalidHandle {
				id = ctx.programs[programIdx].id
			}

			gl.UseProgram(id)
			programChanged = true
			constantsChanged = true
			bindAttribs = true
		}

		if programIdx != kInvalidHandle {
			program := ctx.programs[programIdx]

			if constantsChanged && program.constantBuffer!= nil {
				ctx.commit(program.constantBuffer)
			}

			viewState.setPredefined(ctx, view, eye, render, draw)

			{
				for stage := 0; stage < CONFIG_MAX_SAMPLERS; stage ++ {
					bind := renderBind.bind[stage]
					current := currentBind.bind[stage]

					if  current.idx != bind.idx ||
						current.bType != bind.bType ||
						current.draw.textureFlags != bind.draw.textureFlags ||
						programChanged {
							if bind.idx != uint16(kInvalidHandle) {
								switch bind.bType {
								case BIND_TEXTURE:
									texture := ctx.textures[bind.idx]
									texture.commit(stage, bind.draw.textureFlags, render.colorPalette[:])
								case BIND_INDEX:
									buffer := ctx.indexBuffers[bind.idx]
									/// >= gl es 3.0
									gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, stage, buffer.id)
								case BIND_VERTEX:
									buffer := ctx.vertexBuffers[bind.idx]
									gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, stage, buffer.id)
								}
							}
					}

					current = bind
				}
			}

			{
				/// TODO 不同这里的意思
				var diffStreamHandles bool
				for idx, streamMask := 0, draw.streamMask
			}
		}

	}

}

func (ctx *RendererContextGL) setFrameBuffer(handle FrameBufferHandle, height uint32, flags uint16) uint32 {
	return 0
}

func (ctx *RendererContextGL) getNumRt() uint32 {
	return 0
}

func (ctx *RendererContextGL) commit(buffer *UniformBuffer) {

}

func rendererUpdateUniforms(context *RendererContext, buffer *UniformBuffer, begin, end uint32) {

}

func rendererCreate() *RendererContext {
	if s_renderGL == nil {
		s_renderGL = new(RendererContextGL)
	}
	return s_renderGL
}

func rendererDestroy() {
	s_renderGL.shutdown()
	s_renderGL = nil
}
/////////////// static & global function
func GetGLString(name uint32) string{
	// gl.GetString(name) todo
	return "<unknown>"
}

func GetGLStringHash(name uint32) uint32 {
	return 0
}

func debugToString(enum uint32) string {
	switch enum {
	case gl.DEBUG_SOURCE_API:
		return "API"
	case gl.DEBUG_SOURCE_WINDOW_SYSTEM:
		return "WinSys"
	case gl.DEBUG_SOURCE_SHADER_COMPILER:
		return "Shader"
	case gl.DEBUG_SOURCE_THIRD_PARTY:
		return "3rdparty"
	case gl.DEBUG_SOURCE_APPLICATION:
		return "Application"
	case gl.DEBUG_SOURCE_OTHER:
		return "Other"
	case gl.DEBUG_TYPE_ERROR:
		return "Error"
	case gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR:
		return "Deprecated behavior"
	case gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:
		return "Undefined behavior"
	case gl.DEBUG_TYPE_PORTABILITY:
		return "Portability"
	case gl.DEBUG_TYPE_PERFORMANCE:
		return "Performance"
	case gl.DEBUG_TYPE_OTHER:
		return "Other"
	case gl.DEBUG_SEVERITY_HIGH:
		return "High"
	case gl.DEBUG_SEVERITY_MEDIUM:
		return "Medium"
	case gl.DEBUG_SEVERITY_LOW:
		return "Low"
	case gl.DEBUG_SEVERITY_NOTIFICATION:
		return "SPAM"
	}
	return "<unknown>"
}

func glGet(name uint32) (result int32) {
	gl.GetIntegerv(name, &result)
	err := gl.GetError()
	if err != gl.NO_ERROR {
		result = 0
		log.Printf("glGetIntegerv(0x%04x, ...) failed with GL error: 0x%04x.", name, err)
	}
	return
}

/// default: xType= gl.ZERO
func setTextureFormat(format TextureFormat, internalFmt, fmt uint32, xType uint32) {
	tfi := &s_textureFormat[format]
	tfi.internalFmt = internalFmt
	tfi.fmt 		= fmt
	tfi.xType 		= xType
}

func flushGLError() {
	for err := gl.GetError(); err != 0; err = gl.GetError() {
	}
}

/// not support: 3D texture / 2D array / compressed image
func texSubImage(target uint32, level int32, xoffset, yoffset, width, height int32, format, xType uint32, data interface{}) {
	gl.TexSubImage2D(target, level, xoffset, yoffset, width, height, format, xType, gl.Ptr(data))
}

func texImage(target uint32, level, internalFmt int32, width, height int32, border int32, format, xType uint32, data interface{}) {
	gl.TexImage2D(target, level, internalFmt, width, height, border, format, xType, gl.Ptr(data))
}

func initTestTexture(format TextureFormat, srgb, mipmap, array bool, dim int32) {

}

/// default: srgb = false, mipAutogen = false, array = false, dim = 16
func isTextureFormatValid(format TextureFormat, srgb, mipAutogen, array bool, dim int32) bool{
	return false
}

/// default: dim = 16
func isImageFormatValid(format TextureFormat, dim int32) bool{
	return false
}

/// default: srgb = false, writeOnly = false, dim = 16
func isFrameBufferFormatValid(format TextureFormat, srgb, writeOnly bool, dim int32) bool {
	return false
}

func getFilters(flags uint32, hasMips bool, magFilter, minFilter uint32) {

}

func frameBufferValidate() {
	complete := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if complete != 0 {
		log.Printf("glCheckFramebufferStatus failed 0x%08x: %s", complete, glEnumName(complete))
	}
}

func glEnumName(enum uint32) string {
	return ""
}

/////////////// static & global field
var s_renderGL *RendererContextGL
var s_viewName [CONFIG_MAX_VIEWS]string

var s_primInfo = []PrimInfo {
	{gl.TRIANGLES, 		3, 3, 0},
	{gl.TRIANGLE_STRIP, 	3, 1, 2},
	{gl.LINES, 			2, 2, 0},
	{gl.LINE_STRIP, 		2, 1, 1},
	{gl.POINTS, 			1, 1, 0},
}
var s_primName = []string {
	"TriList",
	"TriStrip",
	"Line",
	"LineStrip",
	"Point",
}

var s_attribName = []string {
	"a_position",
	"a_normal",
	"a_tangent",
	"a_bitangent",
	"a_color0",
	"a_color1",
	"a_color2",
	"a_color3",
	"a_indices",
	"a_weight",
	"a_texcoord0",
	"a_texcoord1",
	"a_texcoord2",
	"a_texcoord3",
	"a_texcoord4",
	"a_texcoord5",
	"a_texcoord6",
	"a_texcoord7",
}

var s_instanceDataName = []string {
	"i_data0",
	"i_data1",
	"i_data2",
	"i_data3",
	"i_data4",
}

var s_access = []uint32 {
	gl.READ_ONLY,
	gl.WRITE_ONLY,
	gl.READ_WRITE,
}	// Size = ACCESS_COUNT


var s_attribType = []uint32 {
	gl.UNSIGNED_BYTE,			// uin8
	gl.UNSIGNED_INT_10_10_10_2,	// uint10
	gl.SHORT,					// int16
	gl.HALF_FLOAT,				// half
	gl.FLOAT,					// float
}	// Size = ATTRIB_COUNT

var s_blendFactor = []Blend {
	{0,								0, 							false}, 	// ignored
	{gl.ZERO,						gl.ZERO, 					false}, 	// ZERO
	{gl.ONE,						gl.ONE, 					false}, 	// ONE
	{gl.SRC_COLOR,					gl.SRC_COLOR, 				false}, 	// SRC_COLOR
	{gl.ONE_MINUS_SRC_COLOR,		gl.ONE_MINUS_SRC_COLOR,		false}, 	// INV_SRC_COLOR
	{gl.SRC_ALPHA,					gl.SRC_ALPHA, 				false}, 	// SRC_ALPHA
	{gl.ONE_MINUS_SRC_ALPHA,		gl.ONE_MINUS_SRC_ALPHA,		false}, 	// INV_SRC_ALPHA
	{gl.DST_ALPHA,					gl.DST_ALPHA, 				false}, 	// DST_ALPHA
	{gl.ONE_MINUS_DST_ALPHA,		gl.ONE_MINUS_DST_ALPHA,		false}, 	// INV_DST_ALPHA
	{gl.DST_COLOR,					gl.DST_COLOR, 				false}, 	// DST_COLOR
	{gl.ONE_MINUS_DST_COLOR,		gl.ONE_MINUS_DST_COLOR,		false}, 	// INV_DST_COLOR
	{gl.SRC_ALPHA_SATURATE,			gl.ONE, 					false}, 	// SRC_ALPHA_SAT
	{gl.CONSTANT_COLOR,				gl.CONSTANT_COLOR, 			 true}, 	// FACTOR
	{gl.ONE_MINUS_CONSTANT_COLOR,	gl.ONE_MINUS_CONSTANT_COLOR,true}, 	// INV_FACTOR
}

var s_blendEquation = []uint32 {
	gl.FUNC_ADD,
	gl.FUNC_SUBTRACT,
	gl.FUNC_REVERSE_SUBTRACT,
	gl.MIN,
	gl.MAX,
}

var s_cmpFunc 		= []uint32 {
	0, 	// ignored
	gl.LESS,
	gl.LEQUAL,
	gl.EQUAL,
	gl.GEQUAL,
	gl.GREATER,
	gl.NOTEQUAL,
	gl.NEVER,
	gl.ALWAYS,
}

var s_stencilOp 	= []uint32 {
	gl.ZERO,
	gl.KEEP,
	gl.REPLACE,
	gl.INCR_WRAP,
	gl.INCR,
	gl.DECR_WRAP,
	gl.DECR,
	gl.INVERT,
}

var s_stencilFace	= []uint32 {
	gl.FRONT_AND_BACK,
	gl.FRONT,
	gl.BACK,
}

var s_textureAddress = []uint32 {
	gl.REPEAT,
	gl.MIRRORED_REPEAT,
	gl.CLAMP_TO_EDGE,
	gl.CLAMP_TO_BORDER,
}

var s_textureFilterMag = []uint32 {
	gl.LINEAR,
	gl.NEAREST,
	gl.LINEAR,
}

var s_textureFilterMin = [][3]uint32 {
	{gl.LINEAR, 	gl.LINEAR_MIPMAP_LINEAR, 	gl.LINEAR_MIPMAP_NEAREST},
	{gl.NEAREST, 	gl.NEAREST_MIPMAP_LINEAR, 	gl.NEAREST_MIPMAP_NEAREST},
	{gl.LINEAR, 	gl.LINEAR_MIPMAP_LINEAR, 	gl.LINEAR_MIPMAP_NEAREST},
}

/// todo supported texture format!!
var s_textureFormat = []TextureFormatInfo {

}

var s_textureFilter = [TEXTURE_FORMAT_COUNT+1]bool{

}

var s_rboFormat 	= []uint32 {

}

var s_imageFormat 	= []uint32 {

}

var s_vendorIds 	= []VendorId {
	{ "NVIDIA Corporation",           PCI_ID_NVIDIA },
	{ "Advanced Micro Devices, Inc.", PCI_ID_AMD    },
	{ "Intel",                        PCI_ID_INTEL  },
}


