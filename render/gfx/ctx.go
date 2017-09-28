package gfx

import (
	"korok/math"
	"korok/render/bx"

	"log"
	"unsafe"
)

type idx uint16

func (h idx) isValid() bool{
	return h == kInvalidHandle
}

func isValid(handle idx) bool {
	return false
}

const kInvalidHandle idx = idx(UINT16_MAX)

var INVALID_HANDLE = struct {
	idx
}{kInvalidHandle}

type DynamicIndexBufferHandle struct {
	idx
}
type DynamicVertexBufferHandle struct {
	idx
}
type FrameBufferHandle struct {
	idx
}
type IndexBufferHandle struct {
	idx
}
type IndirectBufferHandle struct {
	idx
}
type OcclusionQueryHandle struct {
	idx
}
type ProgramHandle struct {
	idx
}
type ShaderHandle struct {
	idx
}
type TextureHandle struct {
	idx
}
type UniformHandle struct {
	idx
}
type VertexBufferHandle struct {
	idx
}
type VertexLayoutHandle struct {
	idx
}

/// Memory obtained by calling `bgfx.alloc`, `bgfx.copy` or `bgfx.makeRef`
type Memory struct {
	data uintptr
	size uint32
}

type Clear struct {
	index 	[8]uint8
	depth 	float32
	stencil uint8
	flags 	uint16
}

type Rect struct {
	x, y, width, height uint16
}

func (r *Rect) clear() {
	r.x = 0
	r.y = 0
	r.width = 0
	r.height = 0
}

func (r *Rect) isZero() bool{
	u64 := (*uint64)(unsafe.Pointer(r))
	return *u64 == 0
}

func (r *Rect) isZeroArea() bool {
	return 0 == r.width || 0 == r.height
}

func (r *Rect) set(x, y, width, height uint16) {
	r.x, r.y = x, y
	r.width, r.height = width, height
}

// 交集
func (r *Rect) setIntersect(a, b Rect) {
	sx := math.UInt16_max(a.x, b.x)
	sy := math.UInt16_max(a.y, b.y)

	ex := math.UInt16_min(a.x + a.width, b.x + b.width)
	ey := math.UInt16_min(a.y + a.height, b.y + b.height)

	r.x, r.y = sx, sy

	/// todo! width and height
}

func (r *Rect) intersect(a Rect) {
	r.setIntersect(*r, a)
}

/// Frame buffer texture attachment info
type Attachment struct {
	handle TextureHandle
	mip uint16
	layer uint16
}

/// Transform data
type Transform struct {
	data *float32		// Pointer to first 4*4 matrix
	num  uint16			// number of matrices
}

/// Transient index buffer
type TransientIndexBuffer struct {
	data *uint8			// Pointer to data
	size uint32
	startIndex uint32
	handle IndexBufferHandle
}

/// Transient vertex buffer
type TransientVertexBuffer struct {
	data        *uint8
	size        uint32
	startVertex uint32
	stride      uint16
	handle      VertexBufferHandle
	layout      VertexLayoutHandle
}

/// View stats
type ViewStats struct {
	name string
	view uint8
	cpuTimeElapsed uint64
	gpuTimeElapsed uint64
}

/// View mode sets draw call sort order
type ViewMode uint16
const (
	VIEW_MODE_DEFAULT ViewMode = iota
	VIEW_MODE_SEQUENTIAL
	VIEW_MODE_DEPTH_ASCENDING
	VIEW_MODE_DEPTH_DESCENDING

	VIEW_MODE_COUNT
)

type TextureFormat uint16
const (
	TEXTURE_FORMAT_RGBA16F TextureFormat = iota
	TEXTURE_FORMAT_RGBA32F
	TEXTURE_FORMAT_RGBA8
	TEXTURE_FORMAT_R5G6B5
	TEXTURE_FORMAT_RGBA4
	TEXTURE_FORMAT_RGB5A1

	TEXTURE_FORMAT_UNKNOWN
	TEXTURE_FORMAT_COUNT
)


/// Uniform type enum
type UniformType uint8
const (
	UNIFORM_TYPE_INT1 UniformType = iota
	UNIFORM_TYPE_END

	UNIFORM_TYPE_VEC4
	UNIFORM_TYPE_MAT3
	UNIFORM_TYPE_MAT4

	UNIFORM_TYPE_COUNT
)

type BackBufferRatio uint8
const (
	BBR_EQUAL 		BackBufferRatio = iota 	// Equal to back buffer
	BBR_HALF								// One half size of back buffer
	BBR_QUARTER								// One quarter size of back buffer
	BBR_EIGHTH								// One eighth size of back buffer
	BBR_SIXTEENTH 							// One sixteenth size of back buffer
	BBR_DOUBLE								// Double size of back buffer

	BBR_COUNT
)

/// max = 256
type UpdateBatch struct {
	num uint32
	keys	[256]uint32
	values 	[256]uint32
}

func (b *UpdateBatch) add(key uint32, value uint32) {

}

func (b *UpdateBatch) sort() {

}

func (b *UpdateBatch) isFull() bool{
	return false
}

func (b *UpdateBatch) reset() {

}

/// Renderer statistics data
type Stats struct {
	CpuTimeBegin uint64
	CpuTimeEnd 	 uint64
	CpuTimerFreq  uint64

	GpuTimeBegin uint64
	GpuTimeEdn   uint64
	GpuTimerFreq uint64

	WaitRender   int64
	WaitSubmit 	 int64

	NumDraw 	  uint32
	NumCompute 		uint32
	MaxGpuLatency 	uint32

	Width uint16
	Height uint16
	TextWidth uint16
	TextHeight uint16

	NumViews uint16
	ViewStats [256]ViewStats
}


/// Renderer capabilities
type Caps struct {
	RendererType RendererType

	/// Supported functionality
	supported uint64

	vendorId uint16
	deviceId uint16
	homogeneousDepth bool
	originBottomLeft bool
	numGPUs uint8

	/// GPU info
	gpu [4]struct{
		vendorId uint16
		deviceId uint16
	}

	limits struct{
		MaxDrawCalls 	uint32
		MaxBlit 	 	uint32
		MaxTextureSize 	uint32
		MaxViews 	 	uint32
		MaxFrameBuffers uint32
		MaxFBAttachments uint32
		MaxPrograms 	uint32
		MaxShaders 		uint32
		MaxTextures 	uint32
		MaxTextureSampler uint32
		MaxVertexDecls 	uint32
		MaxVertexStreams uint32
		MaxIndexBuffers uint32
		MaxVertexBuffers uint32
		MaxDynamicIndexBuffers uint32
		MaxDynamicVertexBuffers uint32
		MaxUniforms 	uint32
		MaxOcclusionQueries uint32
	}

	formats [TEXTURE_FORMAT_COUNT]uint16
}

type PreDefUniformEnum uint16
const(
	UNIFORM_VIEW_RECT PreDefUniformEnum = iota
	UNIFORM_VIEW_TEXEL
	UNIFORM_V
	UNIFORM_INV_V
	UNIFORM_P
	UNIFORM_INV_P
	UNIFORM_M
	UNIFORM_MV
	UNIFORM_MVP
	UNIFORM_ALPHA_REF

	UNIFORM_COUNT
)

type PredefinedUniform struct {
	loc uint32
	count uint16
	uType uint8
}

type Resolution struct {
	width 	uint32
	height 	uint32
	flags	uint32
}

/// Texture info
type TextureInfo struct {
	format TextureFormat
	storageSize uint32
	width 	uint16
	height 	uint16
	depth 	uint16
	numLayers uint16
	numMips uint8
	bitsPerPixel uint8
	cubeMap bool
}

/// Uniform info
type UniformInfo struct {
	name string 			// Uniform name
	uType UniformType 		// Uniform tpe
	num uint16 				// Number of elements in array
}


const kInvalidBlock = UINT64_MAX

// First-fit no-local allocator
type NonLocalAllocator struct {
	freeList []struct{
		ptr  uint64
		size uint32
	}
	used map[uint64]uint32
}

func (alloc *NonLocalAllocator) reset() {

}

func (alloc *NonLocalAllocator) add(ptr uint64, size uint32) {

}

func (alloc *NonLocalAllocator) remove() uint64 {
	return 0
}

func (alloc *NonLocalAllocator) alloc(size uint32) uint64{
	return 0
}

func (alloc *NonLocalAllocator) free(block uint64) {

}

func (alloc *NonLocalAllocator) compact() bool {
	return false
}

type RendererContext interface {
	GetRendererType() RendererType
	GetRendererName() string
	IsDeviceRemoved() bool

	/// static index and vertex buffer
	CreateIndexBuffer(handle IndexBufferHandle, mem *Memory,flags uint16)
	DestroyIndexBuffer(handle IndexBufferHandle)
	CreateVertexLayout(handle VertexLayoutHandle, layout *VertexLayout)
	DestroyVertexLayout(handle VertexLayoutHandle)
	CreateVertexBuffer(handle VertexBufferHandle, mem *Memory, layoutHandle VertexLayoutHandle, flags uint16)
	DestroyVertexBuffer(handle VertexBufferHandle)

	/// dynamic index and vertex buffer
	CreateDynamicIndexBuffer(handle IndexBufferHandle, size uint32, flags uint16)
	UpdateDynamicIndexBuffer(handle IndexBufferHandle, offset uint32, size uint32, mem *Memory)
	DestroyDynamicIndexBuffer(handle IndexBufferHandle)
	CreateDynamicVertexBuffer(handle VertexBufferHandle, size uint32, flags uint16)
	UpdateDynamicVertexBuffer(handle VertexBufferHandle, offset uint32,  size uint32, mem *Memory)
	DestroyDynamicVertexBuffer(handle VertexBufferHandle)

	/// shader and program
	CreateShader(handle ShaderHandle, mem *Memory)
	DestroyShader(handle ShaderHandle)
	CreateProgram(handle ProgramHandle, vsh ShaderHandle, fsh ShaderHandle)
	DestroyProgram(handle ProgramHandle)

	/// textures
	CreateTexture(handle TextureHandle,  mem *Memory, flags uint32, skip uint8)
	UpdateTextureBegin(handle TextureHandle, side uint8, mip uint8)
	UpdateTexture(handle TextureHandle, side, mip uint8, rect *Rect, z, depth, pitch uint16, mem *Memory)
	UpdateTextureEnd()
	ReadTexture(handle TextureHandle, data interface{}, mip uint8)
	ResizeTexture(handle TextureHandle, width, height uint16, numMips uint8)
	OverrideInternal(handle TextureHandle, ptr uintptr)
	GetInternal(handle TextureHandle) uintptr
	DestroyTexture(handle TextureHandle)

	/// frame buffer
	CreateFrameBufferAttachment(handle FrameBufferHandle, num uint8, attachment *Attachment)
	CreateFrameBuffer(handle FrameBufferHandle, nwh interface{}, width, height uint32, depthFormat TextureFormat)
	DestroyFrameBuffer(handle FrameBufferHandle)

	/// uniform
	CreateUniform(handle UniformHandle, uType UniformType, num uint16, name string)
	DestroyUniform(handle UniformHandle)
	UpdateUniform(loc uint16, data interface{}, size uint32)

	///
	SetName(handle idx, name string)

	/// submit
	Submit(render *Frame)
}

type Context struct {
	submit *Frame
	render *Frame

	tempKeys 	[CONFIG_MAX_DRAW_CALLS]uint64
	tempValues  [CONFIG_MAX_DRAW_CALLS]uint16

	fb 			[CONFIG_MAX_VIEWS]FrameBufferHandle
	clear 		[CONFIG_MAX_VIEWS]Clear
	clearColor 	[CONFIG_MAX_COLOR_PALETTE][4]float32
	rect 		[CONFIG_MAX_VIEWS]Rect
	scissor 	[CONFIG_MAX_VIEWS]Rect
	view 		[CONFIG_MAX_VIEWS]Matrix4
	proj 		[CONFIG_MAX_VIEWS]Matrix4
	viewFlags	[CONFIG_MAX_VIEWS]uint8
	seq 		[CONFIG_MAX_VIEWS]uint16
	viewMode 	[CONFIG_MAX_VIEWS]uint8

	colorPatetteDirty uint8

	resolution      Resolution
	instBufferCount int32
	frames          uint32
	debug           uint32

	renderCtx 	RendererContext
	renderMain  RendererContext
	renderNoop  RendererContext

	renderInitialized 	bool
	exit 				bool

	TextureUpdateBatch UpdateBatch
}

func (ctx *Context) init(_type RendererType) bool{
	if ctx.renderInitialized {
		log.Println("Already initialized?")
	}
	ctx.exit = false
	ctx.frames = 0
	ctx.debug = DEBUG_NONE

	ctx.submit.create()

	for i := uint32(0); i < CONFIG_MAX_VIEWS; i++ {
		ctx.resetView(uint8(i))
	}

	for i := range ctx.clearColor {
		ctx.clearColor[i][0] = 0.0
		ctx.clearColor[i][1] = 0.0
		ctx.clearColor[i][2] = 0.0
		ctx.clearColor[i][3] = 1.0
	}

	ctx.layoutRef.init()

	cmdBuf := ctx.getCommandBuffer(Command_RendererInit)
	cmdBuf.writeUInt8(uint8(_type))

	ctx.frameNoRenderWait()

	// Make sure renderer init is called from render thread.
	// g_caps is initialized and available after this point
	ctx.NextFrame(false)

	if !ctx.renderInitialized {
		ctx.getCommandBuffer(Command_RendererShutdownEnd)
		ctx.NextFrame(false)
		ctx.NextFrame(false)
		return false
	}

	// init texture format todo

	g_caps.RendererType = ctx.renderCtx.GetRendererType()
	initAttribTypeSizeTable(_type)
	dumpCaps()

	ctx.submit.transientIb = ctx.createTransientIndexBuffer(CONFIG_TRANSIENT_INDEX_BUFFER_SIZE)
	ctx.submit.transientVb = ctx.createTransientVertexBuffer(CONFIG_TRANSIENT_VERTEX_BUFFER_SIZE, nil)
	ctx.NextFrame(false)

	g_internalData.caps = GetCaps()

	return true
}

func (ctx *Context) shutdown() {
	ctx.getCommandBuffer(Command_RendererShutdownBegin)
	ctx.NextFrame(false)

	ctx.destroyTransientIndexBuffer(ctx.submit.transientIb)
	ctx.destroyTransientVertexBuffer(ctx.submit.transientVb)

	ctx.NextFrame(false)
	ctx.NextFrame(false)

	ctx.getCommandBuffer(Command_RendererShutdownEnd)
	ctx.NextFrame(false)

	ctx.dynVertexBufferAllocator.compact()
	ctx.dynIndexBufferAllocator.compact()

	if uint32(ctx.vertexLayoutHandle.GetNumHandles()) != ctx.layoutRef.vertexLayoutMap.GetNumElements() {
		log.Printf("VertexLayoutRef mismatch, num handles %d, handles in hash map %d",
			ctx.vertexLayoutHandle.GetNumHandles(),
			ctx.layoutRef.vertexLayoutMap.GetNumElements())
	}
	ctx.layoutRef.shutdown(ctx.vertexLayoutHandle)

	bx.MemZero(unsafe.Pointer(g_internalData), int(InternalDataSize))
	s_ctx = nil

	ctx.submit.destroy()

	/// memory leak checks...
}

func (ctx *Context) getCommandBuffer(cmd CommandEnum) *CommandBuffer{
	var cmdBuf *CommandBuffer
	if cmd < Command_End {
		cmdBuf = &ctx.submit.cmdPre
	} else {
		cmdBuf = &ctx.submit.cmdPost
	}
	cmdBuf.writeUInt8(uint8(cmd))
	return cmdBuf
}

func (ctx *Context) reset(width, height uint32, flags uint32) {
	ctx.resolution.width = math.UInt32_clamp(width, 1, g_caps.limits.MaxTextureSize)
	ctx.resolution.height = math.UInt32_clamp(height, 1, g_caps.limits.MaxTextureSize)
	if g_platformDataChangedSinceReset {
		ctx.resolution.flags = flags | RESET_INTERNAL_FORCE
	} else {
		ctx.resolution.flags = flags | 0
	}
	g_platformDataChangedSinceReset = false

	bx.MemFill(unsafe.Pointer(&ctx.fb), len(ctx.fb) * int(unsafe.Sizeof(ctx.fb[0])))

	for ii, num := uint16(0), ctx.textureHandle.GetNumHandles(); ii < num; ii ++ {
		textureIdx := ctx.textureHandle.GetHandleAt(ii)
		textureRef := ctx.textureRef[textureIdx]
		if textureRef.bbRatio != BBR_COUNT {
			handle := TextureHandle{idx(textureIdx)}
			ctx.resizeTexture(handle,
				uint16(ctx.resolution.width),
				uint16(ctx.resolution.height),
				textureRef.numMips)

			ctx.resolution.flags |= RESET_INTERNAL_FORCE
		}
	}
}

func (ctx *Context) setDebug(debug uint32) {
	ctx.debug = debug
}

func (ctx *Context) getPerfStats() *Stats {
	stats := ctx.submit.perfStats
	resolution := ctx.submit.resolution
	stats.Width = uint16(resolution.width)
	stats.Height = uint16(resolution.height)
	return &stats
}

func (ctx *Context) createIndexBuffer(mem *Memory, flags uint16) IndexBufferHandle {
	handle := IndexBufferHandle{idx(ctx.indexBufferHandle.Alloc())}

	if handle.isValid() {
		cmdBuf := ctx.getCommandBuffer(Command_CreateIndexBuffer)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeMemory(mem)
		cmdBuf.writeUInt16(flags)
	}

	if !handle.isValid() {
		log.Println("Failed to allocate index buffer handle")
	}

	return IndexBufferHandle(handle)
}

func (ctx *Context) destroyIndexBuffer(handle IndexBufferHandle) {
	ctx.submit.FreeIndexBuffer(handle)

	cmdBuf := ctx.getCommandBuffer(Command_DestroyIndexBuffer)
	cmdBuf.writeUInt16(uint16(handle))
}

func (ctx *Context) findVertexLayout(layout *VertexLayout) VertexLayoutHandle {
	declHandle := ctx.layoutRef.find(layout.hash)
	if !declHandle.isValid() {
		temp := VertexLayoutHandle{idx(ctx.vertexLayoutHandle.Alloc())}
		declHandle = temp
		cmdBuf := ctx.getCommandBuffer(Command_CreateVertexLayout)
		cmdBuf.writeUInt16(uint16(declHandle))
		cmdBuf.writeVertexDecl(layout)
	}

	return declHandle
}

func (ctx *Context) createVertexBuffer(mem *Memory, layout *VertexLayout, flags uint16) VertexBufferHandle{
	handle := VertexBufferHandle{ idx(ctx.vertexBufferHandle.Alloc())}

	if handle.isValid() {
		declHandle := ctx.findVertexLayout(layout)
		ctx.layoutRef.addVertexBuffer(handle, declHandle, layout.hash)
		ctx.vertexBuffers[handle.idx].stride = layout.stride

		cmdBuf := ctx.getCommandBuffer(Command_CreateVertexBuffer)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeMemory(mem)
		cmdBuf.writeUInt16(uint16(declHandle))
		cmdBuf.writeUInt16(flags)
	}

	if !handle.isValid() {
		log.Println("Failed to allocate vertex buffer handle")
	}
	return handle
}

func (ctx *Context) destroyVertexBuffer(handle VertexBufferHandle) {
	ok := ctx.submit.FreeVertexBuffer(handle)

	if !ok {
		log.Printf("Vertex buffer handle %d is already destroyed!", handle.idx)
	}

	cmdBuf := ctx.getCommandBuffer(Command_DestroyVertexBuffer)
	cmdBuf.writeUInt16(uint16(handle))
}

/// 这个方法不知道是什么地方调用的
func (ctx *Context) destroyVertexBufferInternal(handle VertexBufferHandle) {
	layoutHandle := ctx.layoutRef.releaseVertexBuffer(handle)
	if layoutHandle.isValid() {
		cmdBuf := ctx.getCommandBuffer(Command_DestroyVertexLayout)
		cmdBuf.writeUInt16(uint16(layoutHandle))

		ctx.render.freeVertexDecBuffer(layoutHandle)
	}
	ctx.vertexBufferHandle.Free(uint16(handle.idx))
}

func (ctx *Context) allocDynamicIndexBuffer(size uint32, flags uint16) uint64{
	ptr := ctx.dynIndexBufferAllocator.alloc(size)
	if ptr == kInvalidBlock {
		indexBufferHandle := IndexBufferHandle{idx(ctx.indexBufferHandle.Alloc())}
		if !indexBufferHandle.isValid() {
			return ptr
		}

		allocSize := math.UInt32_max(CONFIG_DYNAMIC_INDEX_BUFFER_SIZE, size)
		cmdBuf := ctx.getCommandBuffer(Command_CreateDynamicIndexBuffer)
		cmdBuf.WriteType(indexBufferHandle)
		cmdBuf.WriteType(allocSize)
		cmdBuf.WriteType(flags)

		ctx.dynIndexBufferAllocator.add(uint64(indexBufferHandle.idx) << 32, allocSize)
		ptr = ctx.dynIndexBufferAllocator.alloc(size)
	}
	return ptr
}

func (ctx *Context) createDynamicIndexBuffer(num uint32, flags uint16) DynamicIndexBufferHandle {
	var handle DynamicIndexBufferHandle = INVALID_HANDLE
	var indexSize uint32 = 4
	if (flags & BUFFER_INDEX32) == 0 {
		indexSize = 2
	}
	var size uint32 = bx.Align16(num*indexSize)
	var ptr  uint64
	if 0 != (flags & BUFFER_COMPUTE_WRITE) {
		indexBufferHandle := IndexBufferHandle{idx(ctx.indexBufferHandle.Alloc())}
		if !indexBufferHandle.isValid() {
			return handle
		}
		cmdBuf := ctx.getCommandBuffer(Command_CreateDynamicIndexBuffer)
		cmdBuf.writeUInt16(uint16(indexBufferHandle))
		cmdBuf.writeUInt32(size)
		cmdBuf.writeUInt16(flags)

		ptr = uint64(indexBufferHandle.idx) << 32
	} else {
		ptr = ctx.allocDynamicIndexBuffer(size, flags)
		if ptr == kInvalidBlock {
			return handle
		}
	}

	handle.idx = idx(ctx.dynamicIndexBufferHandle.Alloc())
	if !handle.isValid() {
		return handle
	}

	dib := ctx.dynamicIndexBuffers[handle.idx]
	dib.handle.idx 	= idx(ptr>>32)
	dib.offset		= uint32(ptr)
	dib.size 		= num * indexSize
	dib.startIndex	= bx.StrideAlign(dib.offset, indexSize)/indexSize
	dib.flags 		= flags

	return handle
}

func (ctx *Context) updateDynamicIndexBuffer(handle DynamicIndexBufferHandle, startIndex uint32, mem *Memory) {
	dib := ctx.dynamicIndexBuffers[handle.idx]

	if 0 != (dib.flags & BUFFER_COMPUTE_WRITE) {
		log.Printf("Can't update GPU from CPU")
	}

	indexSize := uint32(4)
	if 0 == (dib.flags & BUFFER_INDEX32) {
		indexSize = 2
	}

	if dib.size < mem.size && 0 != (dib.flags & BUFFER_ALLOW_RESIZE) {
		ctx.dynIndexBufferAllocator.free(uint64(dib.handle.idx)<<32 | uint64(dib.offset))
		ctx.dynIndexBufferAllocator.compact()

		ptr := ctx.allocDynamicIndexBuffer(mem.size, dib.flags)
		dib.handle.idx  = idx(ptr>>32)
		dib.offset 		= uint32(ptr)
		dib.size		= mem.size
		dib.startIndex	= bx.StrideAlign(dib.offset, indexSize)/indexSize
	}

	offset 	:= (dib.startIndex + startIndex)*indexSize
	size 	:= math.UInt32_min(offset +
			math.UInt32_min(bx.UInt32_stasub(dib.size, startIndex * indexSize), mem.size),
			CONFIG_DYNAMIC_INDEX_BUFFER_SIZE) - offset

	if mem.size > size {
		log.Printf("Truncating dynamic index buffer update (size %d, mem size %d).",
			size,
			mem.size)
	}

	cmdBuf := ctx.getCommandBuffer(Command_UpdateDynamicIndexBuffer)
	cmdBuf.writeUInt16(uint16(dib.handle))
	cmdBuf.writeUInt32(offset)
	cmdBuf.writeUInt32(size)
	cmdBuf.writeMemory(mem)
}

func (ctx *Context) destroyDynamicIndexBuffer(handle DynamicIndexBufferHandle) {
	ctx.freeDynamicIndexBufferHandle[ctx.numFreeDynamicIndexBufferHandles] = handle
	ctx.numFreeDynamicIndexBufferHandles++
}

func (ctx *Context) destroyDynamicIndexBufferInternal(handle DynamicIndexBufferHandle) {
	dib := ctx.dynamicIndexBuffers[handle.idx]
	if 0 != (dib.flags & BUFFER_COMPUTE_WRITE) {
		ctx.destroyIndexBuffer(dib.handle)
	} else {
		ctx.dynIndexBufferAllocator.free(uint64(dib.handle.idx)<<32 | uint64(dib.offset))
		if ctx.dynIndexBufferAllocator.compact() {
			for ptr := ctx.dynIndexBufferAllocator.remove(); 0 != ptr; ptr = ctx.dynIndexBufferAllocator.remove() {
				handle := IndexBufferHandle{idx(ptr>>32)}
				ctx.destroyIndexBuffer(handle)
			}
		}
	}
	ctx.dynamicIndexBufferHandle.Free(uint16(handle.idx))
}

func (ctx *Context) allocDynamicVertexBuffer(size uint32, flags uint16) uint64 {
	ptr := ctx.dynVertexBufferAllocator.alloc(size)
	if ptr == kInvalidBlock {
		vertexBufferHandle := VertexBufferHandle{idx(ctx.vertexBufferHandle.Alloc())}
		if !vertexBufferHandle.isValid() {
			log.Println("Failed to allocate dynamic vertex buffer handle")
			return kInvalidBlock
		}

		allocSize := math.UInt32_max(CONFIG_DYNAMIC_VERTEX_BUFFER_SIZE, size)

		cmdBuf := ctx.getCommandBuffer(Command_CreateDynamicVertexBuffer)
		cmdBuf.writeUInt16(uint16(vertexBufferHandle))
		cmdBuf.writeUInt32(allocSize)
		cmdBuf.writeUInt16(flags)

		ctx.dynVertexBufferAllocator.add(uint64(vertexBufferHandle.idx) << 32, allocSize)
		ptr = ctx.dynVertexBufferAllocator.alloc(size)
	}

	return ptr
}

func (ctx *Context) createDynamicVertexBuffer(num uint32, lyt *VertexLayout, flags uint16) DynamicVertexBufferHandle {
	var handle DynamicVertexBufferHandle = INVALID_HANDLE
	var size uint32 = uint32(bx.StrideAlign16(uint16(num)*lyt.stride, lyt.stride))
	var ptr uint64 = 0

	if 0 != (flags & BUFFER_COMPUTE_READ_WRITE) {
		vertexBufferHandle := VertexBufferHandle{idx(ctx.vertexBufferHandle.Alloc())}
		if !vertexBufferHandle.isValid() {
			return handle
		}
		cmdBuf := ctx.getCommandBuffer(Command_CreateDynamicVertexBuffer)
		cmdBuf.writeUInt16(uint16(vertexBufferHandle))
		cmdBuf.writeUInt32(size)
		cmdBuf.writeUInt16(flags)

		ptr = uint64(vertexBufferHandle.idx) << 32
	} else {
		ptr = ctx.allocDynamicVertexBuffer(size, flags)
		if ptr == kInvalidBlock {
			return handle
		}
	}

	layoutHandle := ctx.findVertexLayout(lyt)
	handle.idx = idx(ctx.dynamicVertexBufferHandle.Alloc())
	dvb := ctx.dynamicVertexBuffers[handle.idx]
	dvb.handle.idx 	= idx(ptr>>32)
	dvb.offset 		= uint32(ptr)
	dvb.size 		= num * uint32(lyt.stride)
	dvb.startVertex = bx.StrideAlign(dvb.offset, uint32(lyt.stride))/uint32(lyt.stride)
	dvb.numVertices	= num
	dvb.stride		= lyt.stride
	dvb.layout 		= layoutHandle
	dvb.flags		= flags
	ctx.layoutRef.addDynamicBuffer(handle, layoutHandle, lyt.hash)

	return handle
}

func (ctx *Context) updateDynamicVertexBuffer(handle DynamicVertexBufferHandle, startVertex uint32, mem *Memory) {
	dvb := ctx.dynamicVertexBuffers[handle.idx]

	if 0 != (dvb.flags & BUFFER_COMPUTE_WRITE) {
		log.Printf("Can't update GPU write buffer from CPU")
	}

	if dvb.size < mem.size && 0 != (dvb.flags & BUFFER_ALLOW_RESIZE) {
		ctx.dynVertexBufferAllocator.free(uint64(dvb.handle.idx)<<32 | uint64(dvb.offset))
		ctx.dynVertexBufferAllocator.compact()

		ptr := ctx.allocDynamicVertexBuffer(mem.size, dvb.flags)
		dvb.handle.idx 	= idx(ptr>>32)
		dvb.offset		= uint32(ptr)
		dvb.size 		= mem.size
		dvb.numVertices = dvb.size/uint32(dvb.stride)
		dvb.startVertex = bx.StrideAlign(dvb.offset, uint32(dvb.stride))/uint32(dvb.stride)
	}

	offset := (dvb.startVertex + startVertex) * uint32(dvb.stride)
	size 	:= math.UInt32_min(offset +
		math.UInt32_min(bx.UInt32_stasub(dvb.size, startVertex*uint32(dvb.stride)), mem.size), CONFIG_DYNAMIC_VERTEX_BUFFER_SIZE) - offset

	cmdBuf := ctx.getCommandBuffer(Command_UpdateDynamicVertexBuffer)
	cmdBuf.writeUInt16(uint16(dvb.handle))
	cmdBuf.writeUInt32(offset)
	cmdBuf.writeUInt32(size)
	cmdBuf.writeMemory(mem)
}

func (ctx *Context) destroyDynamicVertexBuffer(handle DynamicVertexBufferHandle) {
	ctx.freeDynamicVertexBufferHandle[ctx.numFreeDynamicVertexBufferHandles] = handle
	ctx.numFreeDynamicVertexBufferHandles++
}

func (ctx *Context) destroyDynamicVertexBufferInternal(handle DynamicVertexBufferHandle) {
	layoutHandle := ctx.layoutRef.releaseDynamicBuffer(handle)

	if layoutHandle.isValid() {
		cmdBuf := ctx.getCommandBuffer(Command_DestroyVertexLayout)
		cmdBuf.writeUInt16(uint16(layoutHandle))

		ctx.render.freeVertexDecBuffer(layoutHandle)
	}

	dvb := ctx.dynamicVertexBuffers[handle.idx]
	if 0 != (dvb.flags & BUFFER_COMPUTE_READ_WRITE) {
		ctx.destroyVertexBuffer(dvb.handle)
	} else {
		ctx.dynVertexBufferAllocator.free(uint64(dvb.handle.idx)<<32 | uint64(dvb.offset))
		if ctx.dynVertexBufferAllocator.compact() {
			for ptr := ctx.dynVertexBufferAllocator.remove(); 0 != ptr; ptr = ctx.dynVertexBufferAllocator.remove() {
				handle := VertexBufferHandle{idx(ptr>>32)}
				ctx.destroyVertexBuffer(handle)
			}
		}
	}
	ctx.dynamicVertexBufferHandle.Free(uint16(handle.idx))
}

func (ctx *Context) getAvailTransientIndexBuffer(num uint32) uint32 {
	return ctx.submit.getAvailTransientIndexBuffer(num)
}

func (ctx *Context) getAvailTransientVertexBuffer(num uint32, stride uint16) uint32 {
	return ctx.submit.getAvailTransientVertexBuffer(num, stride)
}

func (ctx *Context) createTransientIndexBuffer(size uint32) *TransientIndexBuffer {
	var tib *TransientIndexBuffer
	handle := IndexBufferHandle{idx(ctx.indexBufferHandle.Alloc())}

	if !handle.isValid() {
		log.Println("Failed to allocate transient index buffer handle.")
	}

	if handle.isValid() {
		cmdBuf := ctx.getCommandBuffer(Command_CreateDynamicIndexBuffer)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeUInt32(size)
		var flags uint16 = BUFFER_NONE
		cmdBuf.writeUInt16(flags)

		// var size uint32 = bx.Align16(uint32(TransientIndexBufferSize)) + bx.Align16(size)
		data := make([]uint8, size)

		tib = new(TransientIndexBuffer)
		tib.data = &data[0]
		tib.size = size
		tib.handle = handle
	}

	return tib
}

func (ctx *Context) destroyTransientIndexBuffer(tib *TransientIndexBuffer) {
	cmdBuf := ctx.getCommandBuffer(Command_DestroyDynamicIndexBuffer)
	cmdBuf.writeUInt16(uint16(tib.handle))

	ctx.submit.FreeIndexBuffer(tib.handle)
}

func (ctx *Context) allocTransientIndexBuffer(_tib *TransientIndexBuffer, num uint32) {
	offset := ctx.submit.allocTransientIndexBuffer(&num)
	tib := ctx.submit.transientIb

	_tib.data = (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(tib.data)) + uintptr(offset)))
	_tib.size = num * 2
	_tib.handle = tib.handle
	_tib.startIndex = bx.StrideAlign(offset, 2) / 2
}

func (ctx *Context) createTransientVertexBuffer(size uint32, decl *VertexLayout) *TransientVertexBuffer {
	var tvb *TransientVertexBuffer
	var handle = VertexBufferHandle{idx(ctx.vertexBufferHandle.Alloc())}

	if !handle.isValid() {
		log.Println("Failed to allocate transient vertex buffer handle")
	}

	if handle.isValid() {
		var stride uint16
		var layoutHandle VertexLayoutHandle = INVALID_HANDLE

		if nil != decl {
			layoutHandle = ctx.findVertexLayout(decl)
			ctx.layoutRef.addVertexBuffer(handle, layoutHandle, decl.hash)
			stride = decl.stride
		}

		cmdBuf := ctx.getCommandBuffer(Command_CreateDynamicVertexBuffer)
		cmdBuf.writeUInt16(uint16(handle))
		cmdBuf.writeUInt32(size)
		flags := uint16(BUFFER_NONE)
		cmdBuf.writeUInt16(flags)

		tvb := new(TransientVertexBuffer)
		data := make([]uint8, size)
		tvb.data = &data[0]
		tvb.startVertex = 0
		tvb.stride = stride
		tvb.handle = handle
		tvb.layout = layoutHandle
	}
	return tvb
}

func (ctx *Context) destroyTransientVertexBuffer(tvb *TransientVertexBuffer) {
	cmdBuf := ctx.getCommandBuffer(Command_DestroyDynamicVertexBuffer)
	cmdBuf.writeUInt16(uint16(tvb.handle))

	ctx.submit.FreeVertexBuffer(tvb.handle)
}

func (ctx *Context) allocTransientVertexBuffer(tvb *TransientVertexBuffer, num uint32, decl *VertexLayout) {
	declHandle := ctx.layoutRef.find(decl.hash)
	dvb := ctx.submit.transientVb

	if !declHandle.isValid() {
		temp := VertexLayoutHandle{idx(ctx.vertexLayoutHandle.Alloc())}
		declHandle = temp
		cmdBuf := ctx.getCommandBuffer(Command_CreateVertexLayout)
		cmdBuf.WriteType(declHandle)
		cmdBuf.WriteType(decl)
		ctx.layoutRef.add(declHandle, decl.hash)
	}

	offset := ctx.submit.allocTransientVertexBuffer(&num, decl.stride)

	tvb.data = (*uint8)(unsafe.Pointer(uintptr(tvb.data) + uintptr(offset)))
	tvb.size = num * uint32(decl.stride)
	tvb.startVertex = bx.StrideAlign(offset, uint32(decl.stride))/ uint32(decl.stride)
	tvb.stride = decl.stride
	tvb.handle = dvb.handle
	tvb.layout = declHandle
}

func (ctx *Context) setPaletteColor(index uint8, rgba [4]float32) {
	if uint32(index) >= CONFIG_MAX_COLOR_PALETTE {
		log.Printf("Color palette index out of bounds %d (max: %d).",
					index,
					CONFIG_MAX_COLOR_PALETTE)
	}

	ctx.clearColor[index] = rgba
	ctx.colorPatetteDirty = 2
}

func (ctx *Context) setViewName(id uint8, name string) {
	cmdBuf := ctx.getCommandBuffer(Command_UpdateViewName)
	cmdBuf.writeUInt8(id)

	len := uint16(len(name) + 1)
	cmdBuf.writeUInt16(len)
	cmdBuf.writeString(name)
}

/// 这里有些限制条件，不解！
func (ctx *Context) setViewRect(id uint8, x, y, width, height uint16) {
	rect := &ctx.rect[id]
	rect.x = x
	rect.y = y
	rect.width = width
	rect.height = height
}

func (ctx *Context) setViewScissor(id uint8, x, y, width, height uint16) {
	scissor := &ctx.scissor[id]
	scissor.x = x
	scissor.y = y
	scissor.width = width
	scissor.height = height
}

func (ctx *Context) setViewClear(id uint8, flags uint16, rgba uint32, depth float32, stencil uint8) {
	clear := &ctx.clear[id]
	clear.flags = flags
	clear.index[0] = uint8(rgba>>24)
	clear.index[1] = uint8(rgba>>16)
	clear.index[2] = uint8(rgba>> 8)
	clear.index[3] = uint8(rgba>> 0)
	clear.depth    = depth
	clear.stencil  = stencil
}

func (ctx *Context) setViewMode(id uint8, mode ViewMode) {
	ctx.viewMode[id] = uint8(mode)
}

func (ctx *Context) setViewFrameBuffer(id uint8, handle FrameBufferHandle) {
	ctx.fb[id] = handle
}

func (ctx *Context) setViewTransform(id uint8, view *Matrix4, proj *Matrix4, flags uint8) {
	ctx.viewFlags[id] = flags

	if view != nil {
		ctx.view[id] = *view
	} else {
		ctx.view[id].setIdentity()
	}

	if proj != nil {
		ctx.proj[id] = *proj
	} else {
		ctx.proj[id].setIdentity()
	}
}

func (ctx *Context) resetView(id uint8) {
	ctx.setViewRect(id, 0, 0, 1, 1)
	ctx.setViewScissor(id, 0, 0, 0, 0)
	ctx.setViewClear(id, CLEAR_NONE, 0, 0, 0)
	ctx.setViewMode(id, VIEW_MODE_DEFAULT)
	ctx.setViewFrameBuffer(id, INVALID_HANDLE)
	ctx.setViewTransform(id, nil, nil, VIEW_NONE)
}

func (ctx *Context) setState(state uint64, rgba uint32) {
	ctx.submit.setState(state, rgba)
}

func (ctx *Context) setStencil(fstencil, bstencil uint32) {
	ctx.submit.setStencil(fstencil, bstencil)
}

func (ctx *Context) setScissor(x, y, width, height uint16) uint16 {
	return ctx.submit.setScissor(x, y, width, height)
}

func (ctx *Context) setScissorCache(cache uint16) {
	ctx.submit.setScissorCache(cache)
}

func (ctx *Context) setTransform(mtx interface{}, num uint16) uint32{
	return 0
}

func (ctx *Context) allocTransform(transform *Transform, num uint16) uint32 {
	return 0
}

func (ctx *Context) setTransformCache(cache uint32, num uint16) {
}

func (ctx *Context) setUniform(handle UniformHandle, value interface{}, num uint16) {
	uniform := ctx.uniformRef[handle.idx]
	ctx.submit.writeUniform(uniform.uType, handle, value, num)
}

func (ctx *Context) setIndexBuffer(handle IndexBufferHandle, firstIndex, numIndices uint32) {
	ctx.submit.setIndexBuffer(handle, firstIndex, numIndices)
}

func (ctx *Context) setDynamicIndexBuffer(handle DynamicIndexBufferHandle, firstIndex, numIndices uint32) {
	ctx.submit.setDynamicIndexBuffer(&ctx.dynamicIndexBuffers[handle.idx], firstIndex, numIndices)
}


func (ctx *Context) setTransientIndexBuffer(tib *TransientIndexBuffer, firstIndex, numIndices uint32) {
	ctx.submit.setTransientIndexBuffer(tib, firstIndex, numIndices)
}

func (ctx *Context) setVertexBuffer(stream uint8, handle VertexBufferHandle, _startVertex, _numVertices uint32) {
	ctx.submit.setVertexBuffer(stream, handle, _startVertex, _numVertices)
}

func (ctx *Context) setDynamicVertexBuffer(stream uint8, handle DynamicVertexBufferHandle, _startVertex, _numVertices uint32) {
	ctx.submit.setDynamicVertexBuffer(stream, &ctx.dynamicVertexBuffers[handle.idx], _startVertex, _numVertices)
}

func (ctx *Context) setTransientVertexBuffer(stream uint8, tvb *TransientVertexBuffer, _startVertex, _numVertices uint32) {
	ctx.submit.setTransientVertexBuffer(stream, tvb, _startVertex, _numVertices)
}

func (ctx *Context) Submit(id uint8, handle ProgramHandle, depth int32, preserveState bool) uint32{
	return ctx.submit.Submit(id, handle, depth, preserveState)
}

func (ctx *Context) freeDynamicBuffers() {

}

func (ctx *Context) freeAllHandles(frame *Frame) {

}

func (ctx *Context) NextFrame(capture bool) uint32{
	/// guard or debug todo
	ctx.submit.capture = capture

	// wait for render thread to finish
	ctx.frameNoRenderWait()

	return ctx.frames
}

/// Frame => frameNoRenderWait => Swap =>
func (ctx *Context) frameNoRenderWait() {
	ctx.swap()
}

func (ctx *Context) swap() {
	ctx.freeDynamicBuffers()

	ctx.submit.resolution = ctx.resolution
	ctx.resolution.flags &= ^RESET_INTERNAL_FORCE
	ctx.submit.debug = ctx.debug
	ctx.submit.perfStats.NumViews = 0

	/// memory copy
	copy(ctx.submit.fb, ctx.fb)
	copy(ctx.submit.clear, ctx.clear)
	copy(ctx.submit.rect, ctx.rect)
	copy(ctx.submit.scissor, ctx.scissor)
	copy(ctx.submit.view, ctx.view)
	copy(ctx.submit.proj, ctx.proj)
	copy(ctx.submit.viewFlags, ctx.viewFlags)

	if ctx.colorPatetteDirty > 0 {
		ctx.colorPatetteDirty --
		copy(ctx.submit.colorPalette, ctx.clearColor)
	}
	ctx.submit.finish()

	ctx.render, ctx.submit = ctx.submit, ctx.render

	// render frame for single thread
	ctx.renderFrame(-1)

	ctx.frames ++
	ctx.submit.start()

	// memSet array !! ctx.seq =
	ctx.freeAllHandles(ctx.submit)
	ctx.submit.resetFreeHandles()
}

// default=-1
func (ctx *Context) renderFrame(msecs int32) {
	if ctx.renderInitialized {
		ctx.renderCtx.Submit(ctx.render)
	}
}

func (ctx *Context) rendererCreate(_type RendererType) RendererContext {
	return nil
}

func (ctx *Context) renderDestroy(rc RendererContext) {

}

func (ctx *Context) flushTextureUpdateBatch(buffer *CommandBuffer) {

}

/////////////// static function
func dumpCaps() {

}


/////////////// static & global field
var g_internalData *InternalData
var g_platformData *PlatformData
var g_platformDataChangedSinceReset bool

var g_uniformTypeSize = []uint32 {
	uint32(Int32Size),
	0,
	4 * uint32(Float32Size),
	3 * 3 * uint32(Float32Size),
	4 * 4 * uint32(Float32Size),
	1,
}


/// Type Size
var TransientIndexBufferSize = unsafe.Sizeof(TransientIndexBuffer{})
var InternalDataSize 		 = unsafe.Sizeof(InternalData{})
var Int32Size 				 = unsafe.Sizeof(int32(0))
var Float32Size 			 = unsafe.Sizeof(float32(0))


