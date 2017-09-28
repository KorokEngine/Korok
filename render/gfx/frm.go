package gfx

import (
	"korok/math"
	"unsafe"
	"encoding/binary"
	"log"
	"korok/render/bx"
)

/// Vertex attribute enum
type AttribEnum uint16
const (
	ATTRIB_POSITION  AttribEnum = iota
	ATTRIB_NORMAL
	ATTRIB_TANGENT
	ATTRIB_BITANGENT
	ATTRIB_COLOR0
	ATTRIB_COLOR1
	ATTRIB_COLOR2
	ATTRIB_COLOR3
	ATTRIB_INDICES
	ATTRIB_WEIGHT
	ATTRIB_TEXCOORD0
	ATTRIB_TEXCOORD1
	ATTRIB_TEXCOORD2
	ATTRIB_TEXCOORD3
	ATTRIB_TEXCOORD4
	ATTRIB_TEXCOORD5
	ATTRIB_TEXCOORD6
	ATTRIB_TEXCOORD7

	ATTRIB_COUNT
)

/// Vertex attribute type enum
type AttribType uint16
const (
	ATTRIB_TYPE_UINT8 AttribType = iota 	// uint8
	ATTRIB_TYPE_UIN10 						// uint10,
	ATTRIB_TYPE_INT16						// int16
	ATTRIB_TYPE_HALF							// half
	ATTRIB_TYPE_FLOAT

	ATTRIB_TYPE_COUNT
)

type BindEnum uint16
const (
	BIND_IMAGE BindEnum = iota
	BIND_INDEX
	BIND_VERTEX
	BIND_TEXTURE

	BIND_COUNT
)

type Binding struct {
	idx uint16
	bType uint8

	draw struct{
		textureFlags uint32
	}
}

var StreamSize = unsafe.Sizeof(Stream{})

type Stream struct {
	startVertex uint32
	handle      VertexBufferHandle
	layout      VertexLayoutHandle
}

func (s *Stream) clear() {
	//s.startVertex = 0
	//s.handle = kInvalidHandle
	//s.layout   = kInvalidHandle
}

type RenderBind struct {
	bind [CONFIG_MAX_SAMPLERS]Binding
}

func (rb *RenderBind) clear() {
	for i := uint32(0); i < CONFIG_MAX_SAMPLERS; i++ {
		bind := rb.bind[i]
		bind.idx = 0
		bind.bType = 0
		bind.draw.textureFlags = 0
	}
}

/// Vertex declaration
type VertexLayout struct {
	hash 		uint32
	stride 		uint16
	offset 		[ATTRIB_COUNT]uint16
	attributes 	[ATTRIB_COUNT]uint16
}

func (vd *VertexLayout) begin(renderer RendererType) *VertexLayout {
	return nil
}

func (vd *VertexLayout) end() {

}

/// default: normalized=false, asInt=false
func (vd *VertexLayout) add(attrib AttribEnum, num uint8, _type AttribType, normalized, asInt bool) *VertexLayout {
	return nil
}

func (vd *VertexLayout) skip(num uint8) *VertexLayout {
	return nil
}

func (vd *VertexLayout) decode() {

}

func (vd *VertexLayout) has(attrib AttribEnum) {

}

func (vd *VertexLayout) getOffset(attrib AttribEnum) uint16 {
	return 0
}

func (vd *VertexLayout) getStride() uint16{
	return 0
}

func (vd *VertexLayout) getSize(num uint32) uint32 {
	return 0
}

func packStencil(fstencil, bstencil uint32) uint64 {
	return (uint64(bstencil)<<32) | uint64(fstencil);
}

func unpackStencil(_0or1 uint8, _stencil uint64) uint32 {
	return uint32( _stencil >> (32*_0or1))
}

type RenderDraw struct {
	indexBuffer IndexBufferHandle

	stream [CONFIG_MAX_STREAMS]Stream

	stateFlags 	uint64
	stencil 	uint64
	rgba 		uint32
	constBegin  uint32
	constEnd 	uint32
	matrix 		uint32
	startIndex 	uint32
	numIndices 	uint32
	numVertices uint32

	num 		uint16
	scissor 	uint16
	submitFlags uint8
	streamMask 	uint8
}

func (rd *RenderDraw) setStreamBit(stream uint8, handle VertexBufferHandle) bool{
	bit  := uint8(1 << stream)
	mask := rd.streamMask & ^bit
	tmp  := uint8(0)

	if handle.isValid() {
		rd.streamMask = mask | tmp
	}

	return 0 != tmp
}

func (rd *RenderDraw) clear() {
	rd.constBegin 	= 0
	rd.constEnd   	= 0
	rd.stateFlags 	= STATE_DEFAULT
	rd.stencil		= packStencil(STENCIL_DEFAULT, STENCIL_DEFAULT)
	rd.rgba 		= 0
	rd.matrix 		= 0
	rd.startIndex 	= 0
	rd.numIndices 	= UINT32_MAX
	rd.numVertices  = UINT32_MAX

	rd.num				= 1
	rd.submitFlags 		= SUBMIT_EYE_FIRST
	rd.scissor			= UINT16_MAX
	rd.streamMask 		= 0

	bx.MemZero(unsafe.Pointer(&rd.stream[0]), len(rd.stream) * int(StreamSize))

	rd.indexBuffer.idx		  = kInvalidHandle
}

type RenderItem RenderDraw

type CommandEnum uint8
const (
	Command_RendererInit CommandEnum = iota
	Command_RendererShutdownBegin
	Command_CreateVertexLayout
	Command_CreateIndexBuffer
	Command_CreateVertexBuffer

	Command_CreateDynamicIndexBuffer
	Command_UpdateDynamicIndexBuffer
	Command_CreateDynamicVertexBuffer
	Command_UpdateDynamicVertexBuffer

	Command_CreateShader
	Command_CreateProgram
	Command_CreateTexture
	Command_UpdateTexture
	Command_ResizeTexture
	Command_CreateFrameBuffer
	Command_CreateUniform
	Command_UpdateViewName
	Command_SetName
	Command_End

	Command_RendererShutdownEnd
	Command_DestroyVertexLayout
	Command_DestroyIndexBuffer
	Command_DestroyVertexBuffer
	Command_DestroyDynamicIndexBuffer
	Command_DestroyDynamicVertexBuffer
	Command_DestroyShader
	Command_DestroyProgram
	Command_DestroyTexture
	Command_DestroyFrameBuffer
	Command_DestroyUniform
	Command_ReadTexture
)

type CommandBuffer struct {
	pos 	uint32
	size 	uint32
	buffer 	[CONFIG_MAX_COMMAND_BUFFER_SIZE]uint8
}

func (buff *CommandBuffer) writeString(v string) {
	array := ([]byte)(v)
	size := uint32(len(array))
	copy(buff.buffer[buff.pos:], array)
	buff.pos += size
}

func (buff *CommandBuffer) writeBool(v bool) {
	if v {
		buff.buffer[buff.pos] = 1
	}	else {
		buff.buffer[buff.pos] = 0
	}
	buff.pos += 1
}

func (buff *CommandBuffer) writePointer(pointer unsafe.Pointer) {
	size := uint32(PointerSize)
	ptr := uintptr(pointer)

	switch size {
	case 4:
		binary.BigEndian.PutUint32(buff.buffer[buff.pos:], uint32(ptr))
	case 8:
		binary.BigEndian.PutUint64(buff.buffer[buff.pos:], uint64(ptr))
	}

	buff.pos += size
}

func (buff *CommandBuffer) writeUInt8(data uint8) {
	buff.buffer[buff.pos] = data; buff.pos += 1
}

func (buff *CommandBuffer) writeUInt16(data uint16) {
	binary.BigEndian.PutUint16(buff.buffer[buff.pos:], data); buff.pos += 2
}

func (buff *CommandBuffer) writeUInt32(data uint32) {
	binary.BigEndian.PutUint32(buff.buffer[buff.pos:], data); buff.pos += 4
}

func (buff *CommandBuffer) writeMemory(m *Memory) {
	size := uint32(MemorySize)
	ptr := (*[100000]byte)(unsafe.Pointer(m))
	copy(buff.buffer[buff.pos:], ptr[:size])
	buff.pos += size
}

func (buff *CommandBuffer) writeVertexDecl(decl *VertexLayout) {

}

func (buff *CommandBuffer) writeAttachmentArray(attachment *Attachment, num int) {

}

func (buff *CommandBuffer) writeRect(rect *Rect) {

}

func (buff *CommandBuffer) WriteType(in interface{}) {

}

func (buff *CommandBuffer) write(data interface{}, size uint32) {

}

func (buff *CommandBuffer) readPointer() unsafe.Pointer {
	size := uint32(PointerSize)
	var ptr uintptr
	switch size {
	case 4:
		ptr = uintptr(binary.BigEndian.Uint32(buff.buffer[buff.pos:]))
	case 8:
		ptr = uintptr(binary.BigEndian.Uint64(buff.buffer[buff.pos:]))
	}
	return unsafe.Pointer(ptr)
}

func (buff *CommandBuffer) readString(array []byte) {
	size := uint32(len(array))
	copy(array, buff.buffer[buff.pos: buff.pos+size])
	buff.pos += size
}

func (buff *CommandBuffer) readBool() (data bool) {
	data = buff.buffer[buff.pos] == 1; buff.pos += 1
	return
}

func (buff *CommandBuffer) readUInt8() (data uint8) {
	data = buff.buffer[buff.pos]; buff.pos += 1
	return
}

func (buff *CommandBuffer) readUInt16() (data uint16){
	data = binary.BigEndian.Uint16(buff.buffer[buff.pos:]); buff.pos += 2
	return
}

func (buff *CommandBuffer) readUInt32() (data uint32){
	data = binary.BigEndian.Uint32(buff.buffer[buff.pos:]); buff.pos += 4
	return
}

func (buff *CommandBuffer) readMemory(m *Memory) {
	size := uint32(MemorySize)
	ptr := (*[1000]byte)(unsafe.Pointer(m))
	copy(ptr[:size], buff.buffer[buff.pos:buff.pos+size])
	buff.pos += size
}

func (buff *CommandBuffer) readVertexDecl(vd *VertexLayout) {
	size := uint32(VertexLayoutSize)
	ptr := (*[1000]byte)(unsafe.Pointer(vd))
	copy(ptr[:size], buff.buffer[buff.pos:buff.pos+size])
	buff.pos += size
}

func (buff *CommandBuffer) readType(out interface{}) {

}

func (buff *CommandBuffer) read(data interface{}, size uint32) {

}

func (buff *CommandBuffer) skip(size uint32) *uint8 {
	return nil
}

func (buff *CommandBuffer) align(alignment uint32) {

}

func (buff *CommandBuffer) reset() {

}

func (buff *CommandBuffer) start() {

}

func (buff *CommandBuffer) finish() {

}

type SortKey struct {
	depth uint32
	seq   uint32
	program uint16
	view uint8
	trans uint8
}

func (sk *SortKey) encodeDraw(key1 bool) uint64{
	return 0
}

func (sk *SortKey) encodeCompute() uint64{
	return 0
}

func (sk *SortKey) decode(key uint64, viewRemap []uint8) bool {
	return false
}

func (sk *SortKey) decodeView(key uint64) uint8 {
	return 0
}

func (sk *SortKey) remapView(key uint64, viewRemap []uint8) uint64{
	return 0
}

func (sk *SortKey) reset() {

}


type Matrix4 struct {

}

func (m *Matrix4) setIdentity() {

}

type RectCache struct {
	cache []Rect
}

func (rc *RectCache) reset() {

}

func (rc *RectCache) add(x, y, width, height uint16) uint16 {
	return 0
}

type UniformBuffer struct {

}

func NewUniformBuffer() *UniformBuffer{
	return nil
}

func (ub *UniformBuffer) reset() {

}

func (ub *UniformBuffer) getPos() uint32 {
	return 0
}

func (ub *UniformBuffer) finish() {

}

func (ub *UniformBuffer) writeMarker(name string) {

}

func (ub *UniformBuffer) writeUniform(xType UniformType, handle uint16, value interface{}, num uint16) {

}

type UniformRegInfo struct {

}

type UniformRegistry struct {

}

type FreeHandle struct {
	_queue[]uint16
	num  uint16
}

func (fh *FreeHandle) isQueued(handle uint16) bool{
	return false
}

func (fh *FreeHandle) queue(handle uint16) bool{
	return false
}

func (fh *FreeHandle) reset() {

}

func (fh *FreeHandle) get(idx uint16) uint16{
	return 0
}

func (fh *FreeHandle) getNumQueued() uint16 {
	return 0
}

type Frame struct {
	key SortKey

	fb 			[CONFIG_MAX_VIEWS]FrameBufferHandle
	clear 		[CONFIG_MAX_VIEWS]Clear
	colorPalette[CONFIG_MAX_COLOR_PALETTE][4]float32
	rect 		[CONFIG_MAX_VIEWS]Rect
	scissor 	[CONFIG_MAX_VIEWS]Rect
	view 		[CONFIG_MAX_VIEWS]Matrix4
	proj 		[CONFIG_MAX_VIEWS]Matrix4
	viewFlags 	[CONFIG_MAX_VIEWS]uint8

	sortKeys 	[CONFIG_MAX_DRAW_CALLS+1]uint64
	sortValues 	[CONFIG_MAX_DRAW_CALLS+1]uint16
	renderItem  [CONFIG_MAX_DRAW_CALLS+1]RenderItem
	renderBind  [CONFIG_MAX_DRAW_CALLS+1]RenderBind

	draw 		RenderDraw
	bind 		RenderBind

	numVertices 	[CONFIG_MAX_STREAMS]uint32
	stateFlags 		uint64
	uniformBegin 	uint32
	uniformEnd 		uint32
	uniformMax 		uint32

	uniformBuffer *UniformBuffer

	num 			uint16
	numRenderItems 	uint16
	numDropped 		uint16

	rectCache 	RectCache

	iboffset	uint32
	vboffset 	uint32

	transientIb 	*TransientIndexBuffer
	transientVb 	*TransientVertexBuffer

	resolution 	Resolution
	debug 		uint32

	cmdPre 	CommandBuffer
	cmdPost CommandBuffer

	///// free handle
	freeIndexBuffer  FreeHandle
	freeVertexLayout FreeHandle
	freeVertexBuffer FreeHandle
	freeShader       FreeHandle
	freeTexture      FreeHandle
	freeFrameBuffer  FreeHandle
	freeUniform      FreeHandle

	perfStats 		Stats

	waitSubmit 		int64
	waitRender 		int64

	capture 		bool
	discard 		bool
}

func (f *Frame) create() {
	f.uniformBuffer = NewUniformBuffer()
	f.reset()
	f.start()
}

func (f *Frame) destroy() {

}

func (f *Frame) reset() {
	f.start()
	f.finish()
	f.resetFreeHandles()
}

func (f *Frame) start() {
	f.stateFlags = STATE_NONE
	f.uniformBegin = 0
	f.uniformEnd   = 0
	f.draw.clear()
	f.bind.clear()
	f.key.reset()
	f.num = 0
	f.numRenderItems = 0
	f.numDropped = 0
	f.iboffset = 0
	f.vboffset = 0
	f.cmdPre.start()
	f.cmdPost.start()
	f.uniformBuffer.reset()
	f.capture = false
	f.discard = false
}

func (f *Frame) finish() {
	f.cmdPre.finish()
	f.cmdPost.finish()

	f.uniformMax = math.UInt32_max(f.uniformMax, f.uniformBuffer.getPos())
	f.uniformBuffer.finish()

	if f.numDropped > 0 {
		log.Printf("Too many draw calls: %d, dropped %d (max: %d)",
			f.num+f.numDropped,
			f.numDropped,
			CONFIG_MAX_DRAW_CALLS)
	}
}

func (f *Frame) setState(state uint64, rgba uint32) {
	blend := uint8(((state & STATE_BLEND_MASK) >> STATE_BLEND_SHIFT) & 0xFF)
	alphaRef := uint8(((state & STATE_ALPHA_REF_MASK) >> STATE_ALPHA_REF_SHIFT) & 0xFF)
	// transparent sort order table
	// f.key.trans = "\x0\x2\x2\x3\x3\x2\x3\x2\x3\x2\x2\x2\x2\x2\x2\x2\x2\x2\x2"[]
	f.draw.stateFlags = state
	f.draw.rgba = rgba
}

func (f *Frame) setStencil(fstencil, bstencil uint32) {
	f.draw.stencil = packStencil(fstencil, bstencil)
}

func (f *Frame) setScissor(x, y, width, height uint16) uint16{
	scissor := f.rectCache.add(x, y, width, height)
	f.draw.scissor = scissor
	return scissor
}

func (f *Frame) setScissorCache(cache uint16) {
	f.draw.scissor = cache
}

func (f *Frame) setTransform(mtx interface{}, num uint16) uint32{
	f.draw.matrix = f.matrixCache.add(mtx, num)
	f.draw.num    = num
	return f.draw.matrix
}

func (f *Frame) setIndexBuffer(handle IndexBufferHandle, firstIndex, numIndices uint32) {
	f.draw.startIndex = firstIndex
	f.draw.numIndices = numIndices
	f.draw.indexBuffer = handle
}

func (f *Frame) setDynamicIndexBuffer(dib *DynamicIndexBuffer, firstIndex, numIndices uint32) {
	var indexSize uint32 = 4
	if (dib.flags & BUFFER_INDEX32) == 0 {
		indexSize = 2
	}

	f.draw.startIndex = dib.startIndex + firstIndex
	f.draw.numIndices = math.UInt32_min(numIndices, dib.size/indexSize)
	f.draw.indexBuffer = dib.handle
}

func (f *Frame) setTransientIndexBuffer(tib *TransientIndexBuffer, firstIndex, numIndices uint32) {
	f.draw.indexBuffer = tib.handle
	f.draw.startIndex = firstIndex
	f.draw.numIndices = numIndices
	f.discard = 0 == numIndices
}

func (f *Frame) setVertexBuffer(_stream uint8, handle VertexBufferHandle, firstIndex, numIndices uint32) {
	if f.draw.setStreamBit(_stream, handle) {
		stream := f.draw.stream[_stream]
		stream.startVertex = firstIndex
		stream.handle = handle
		stream.layout.idx  = kInvalidHandle
		f.numVertices[_stream] = numIndices
	}
}

func (f *Frame) setDynamicVertexBuffer(_stream uint8, dvb *DynamicVertexBuffer, firstIndex, numIndices uint32) {
	if f.draw.setStreamBit(_stream, dvb.handle) {
		stream := f.draw.stream[_stream]
		stream.startVertex = dvb.startVertex + firstIndex
		stream.handle = dvb.handle
		stream.layout = dvb.layout
		f.numVertices[_stream] = math.UInt32_min(math.UInt32_max(0, dvb.numVertices - firstIndex), numIndices)
	}
}

func (f *Frame) setTransientVertexBuffer(_stream uint8, tvb *TransientVertexBuffer, firstIndex, numIndices uint32) {
	if f.draw.setStreamBit(_stream, tvb.handle) {
		stream := f.draw.stream[_stream]
		stream.startVertex = tvb.startVertex + firstIndex
		stream.handle = tvb.handle
		stream.layout = tvb.layout
		f.numVertices[_stream] = math.UInt32_min(math.UInt32_max(0, tvb.size/uint32(tvb.stride) - firstIndex), numIndices)
	}
}

func (f *Frame) setTexture(stage uint8, sampler UniformHandle, handle TextureHandle, flags uint32) {
	bind := f.bind.bind[stage]
	bind.idx = uint16(handle)
	bind.bType = uint8(BIND_TEXTURE)

	if 0 == (flags & TEXTURE_INTERNAL_DEFAULT_SAMPLER) {
		bind.draw.textureFlags = flags
	} else {
		bind.draw.textureFlags = TEXTURE_INTERNAL_DEFAULT_SAMPLER
	}

	if sampler.isValid() {
		SetUniform(sampler, uint32(stage), 1) // TODO set sampler
	}
}

func (f *Frame) Discard() {
	f.discard = false
	f.draw.clear()
	f.stateFlags = STATE_NONE
}

func (f *Frame) Submit(id uint8, program ProgramHandle, depth int32, preserveState bool) uint32{
	if f.discard {
		f.Discard(); return uint32(f.num)
	}

	if CONFIG_MAX_DRAW_CALLS-1 <= uint32(f.num) || (0 == f.draw.numVertices && 0 == f.draw.numIndices) {
		f.numDropped ++
		return uint32(f.num)
	}

	f.uniformEnd = f.uniformBuffer.getPos()

	if program.idx == kInvalidHandle {
		f.key.program = 0
	} else {
		f.key.program = uint16(program)
	}

	f.key.view = id

	var key1 bool = false
	switch s_ctx.viewMode[id] {
	case VIEW_MODE_SEQUENTIAL:
		f.key.seq = uint32(s_ctx.seq[id])
	case VIEW_MODE_DEPTH_ASCENDING:
		f.key.depth = uint32(depth); key1 = true
	case VIEW_MODE_DEPTH_DESCENDING:
		f.key.depth = -uint32(depth); key1 = true
	}
	s_ctx.seq[id] ++

	key := f.key.encodeDraw(key1)
	f.sortKeys[f.num] = key
	f.sortValues[f.num] = f.numRenderItems
	f.num ++

	f.draw.constBegin = f.uniformBegin
	f.draw.constEnd   = f.uniformEnd
	f.draw.stateFlags |= f.stateFlags

	var numVertices uint32 = UINT32_MAX
	var idx uint32 = 0
	var streamMask uint32 = uint32(f.draw.streamMask)
	var ntz uint32 = math.UInt32_cnttz(streamMask)

	for 0 != streamMask {
		streamMask >>= ntz
		idx 	 	+= ntz
		numVertices  = math.UInt32_min(numVertices, f.numVertices[idx])

		streamMask >>= 1; idx += 1; ntz = math.UInt32_cnttz(streamMask)
	}
	f.draw.numVertices = numVertices

	f.renderItem[f.numRenderItems] 	= RenderItem(f.draw)
	f.renderBind[f.numRenderItems] 	= f.bind
	f.numRenderItems ++

	if !preserveState {
		f.draw.clear()
		f.bind.clear()
		f.uniformBegin = f.uniformEnd
		f.stateFlags   = STATE_NONE
	}

	return uint32(f.num)
}

func (f *Frame) Sort() {
	/// todo radix sort!!
}

func (f *Frame) renderFrame() {
	/// 单线程无需实现
}

// 2 = sizeOf(uint16)
func (f *Frame) getAvailTransientIndexBuffer(num uint32) uint32 {
	iboffset := f.iboffset + num * 2
	iboffset = math.UInt32_min(iboffset, uint32(CONFIG_TRANSIENT_INDEX_BUFFER_SIZE))
	_num := (iboffset - f.iboffset)/ 2
	return _num
}

func (f *Frame) allocTransientIndexBuffer(num *uint32) uint32 {
	_num := f.getAvailTransientIndexBuffer(*num)
	f.iboffset = f.iboffset + _num * 2
	*num = _num
	return f.iboffset
}

func (f *Frame) getAvailTransientVertexBuffer(num uint32, stride uint16) uint32 {
	vboffset := f.vboffset + num * uint32(stride)
	vboffset = math.UInt32_min(vboffset, uint32(CONFIG_TRANSIENT_VERTEX_BUFFER_SIZE))
	_num := (vboffset-f.vboffset)/uint32(stride)
	return _num
}

func (f *Frame) allocTransientVertexBuffer(num *uint32, stride uint16) uint32 {
	_num := f.getAvailTransientVertexBuffer(*num, stride)
	f.vboffset = f.vboffset + _num * uint32(stride)
	*num = _num
	return f.vboffset
}

//// try to resize UniformBuffer
func (f *Frame) writeUniform(_type UniformType, handle UniformHandle, value interface{}, num uint16) {
	f.uniformBuffer.writeUniform(_type, uint16(handle.idx), value, num)
}

func (f *Frame) FreeIndexBuffer(handle IndexBufferHandle) bool{
	return f.freeIndexBuffer.queue(uint16(handle.idx))
}

func (f *Frame) FreeVertexBuffer(handle VertexBufferHandle) bool{
	return f.freeVertexBuffer.queue(uint16(handle.idx))
}

func (f *Frame) freeVertexDecBuffer(handle VertexLayoutHandle) bool{
	return f.freeVertexLayout.queue(uint16(handle.idx))
}

func (f *Frame) FreeShader(handle ShaderHandle) bool{
	return f.freeShader.queue(uint16(handle.idx))
}

func (f *Frame) FreeProgram(handle ProgramHandle) bool{
	return false
}

func (f *Frame) FreeTexture(handle TextureHandle) bool{
	return false
}

func (f *Frame) FreeFrameBuffer(handle FrameBufferHandle) bool {
	return false
}

func (f *Frame) FreeUniform(handle UniformHandle) bool{

	return false
}

func (f *Frame) resetFreeHandles() {
	f.freeIndexBuffer.reset()
	f.freeVertexBuffer.reset()
	f.freeVertexLayout.reset()
	/// todo
}

///////////// static & global variable

var PointerSize = unsafe.Sizeof(unsafe.Pointer(0))
var MemorySize = unsafe.Sizeof(Memory{})
var VertexLayoutSize = unsafe.Sizeof(VertexLayout{})
