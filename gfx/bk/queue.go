package bk

import (
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"unsafe"
)

type Rect struct {
	x, y uint16
	w, h uint16
}

func (r *Rect) clear() {
	r.x, r.y = 0, 0
	r.w, r.h = 0, 0
}

func (r *Rect) isZero() bool {
	u64 := (*uint64)(unsafe.Pointer(r))
	return *u64 == 0
}

type Stream struct {
	vertexBuffer uint16
	vertexFormat uint16 // Offset | Stride， not used now!!

	firstVertex uint16
	numVertex   uint16
}

type RenderDraw struct {
	indexBuffer   uint16
	vertexBuffers [2]Stream
	textures      [2]uint16

	// index params
	firstIndex, num uint16

	// uniform range
	uniformBegin uint16
	uniformEnd   uint16

	// stencil and scissor
	stencil uint32
	scissor Rect

	// required renderer state
	state uint64
}

func (rd *RenderDraw) reset() {

}

// ~ 8000 draw call
const MAX_QUEUE_SIZE = 8 << 10

type RenderQueue struct {
	// render list
	sortKey    [MAX_QUEUE_SIZE]uint64
	sortValues [MAX_QUEUE_SIZE]uint16

	drawCallList [MAX_QUEUE_SIZE]RenderDraw
	drawCallNum  uint16

	sk SortKey

	// per-drawCall state cache
	drawCall RenderDraw

	uniformBegin uint16
	uniformEnd   uint16

	// per-frame state
	viewports [4]Rect
	scissors  [4]Rect
	clears    [4]struct {
		index   [8]uint8
		rgba    uint32
		depth   float32
		stencil uint8
		flags   uint16
	}

	// per-frame data flow
	rm *ResManager
	ub *UniformBuffer

	// render context
	ctx *RenderContext
}

func NewRenderQueue(m *ResManager) *RenderQueue {
	ub := NewUniformBuffer()
	rc := NewRenderContext(m, ub)
	return &RenderQueue{
		ctx: rc,
		rm:  m,
		ub:  ub,
	}
}

func (rq *RenderQueue) Init() {
	rq.ctx.Init()
}

// reset frame-buffer size
func (rq *RenderQueue) Reset(w, h uint16) {

}

func (rq *RenderQueue) Destroy() {
	//
}

func (rq *RenderQueue) SetState(state uint64, rgba uint32) {
	rq.drawCall.state = state
}

func (rq *RenderQueue) SetIndexBuffer(id uint16, firstIndex, num uint16) {
	rq.drawCall.indexBuffer = id & ID_TYPE_MASK
	rq.drawCall.firstIndex = firstIndex
	rq.drawCall.num = num
}

func (rq *RenderQueue) SetVertexBuffer(stream uint8, id uint16, firstVertex, numVertex uint16) {
	if stream < 0 || stream >= 2 {
		log.Printf("Not support stream location: %d", stream)
		return
	}

	vbStream := &rq.drawCall.vertexBuffers[stream]
	vbStream.vertexBuffer = id & ID_TYPE_MASK
	vbStream.vertexFormat = InvalidId
	vbStream.firstVertex = firstVertex
	vbStream.numVertex = numVertex
}

func (rq *RenderQueue) SetTexture(stage uint8, samplerId uint16, texId uint16, flags uint32) {
	if stage < 0 || stage >= 2 {
		log.Printf("Not suppor texture location: %d", stage)
		return
	}

	rq.drawCall.textures[stage] = texId & ID_TYPE_MASK
}

// 复制简单数据的时候（比如：Sampler），采用赋值的方式可能更快 TODO
func (rq *RenderQueue) SetUniform(id uint16, ptr unsafe.Pointer) {
	if ok, um := rq.rm.Uniform(id); ok {
		opCode := Uniform_encode(um.Type, um.Slot, um.Size, um.Count)
		rq.ub.WriteUInt32(opCode)
		rq.ub.Copy(ptr, uint32(um.Size)*uint32(um.Count))
	}
}

// Transform 是 uniform 之一，2D 世界可以省略
func (rq *RenderQueue) SetTransform(mtx *mgl32.Mat4) {
	// TODO impl
}

func (rq *RenderQueue) SetStencil(stencil uint32) {
	rq.drawCall.stencil = stencil
}

func (rq *RenderQueue) SetScissor(x, y, width, height uint16) {
	r := &rq.drawCall.scissor
	r.x, r.y = x, y
	r.w, r.h = width, height
}

/// View Related Setting
func (rq *RenderQueue) SetViewScissor(id uint8, x, y, with, height uint16) {
	if id < 0 || id >= 4 {
		log.Printf("Not support view id: %d", id)
		return
	}
	rq.scissors[id] = Rect{x, y, with, height}
}

func (rq *RenderQueue) SetViewPort(id uint8, x, y, width, height uint16) {
	if id < 0 || id >= 4 {
		log.Printf("Not support view id: %d", id)
		return
	}
	rq.viewports[id] = Rect{x, y, width, height}
}

func (rq *RenderQueue) SetViewClear(id uint8, flags uint16, rgba uint32, depth float32, stencil uint8) {
	if id < 0 || id >= 4 {
		log.Printf("Not support view id: %d", id)
		return
	}
	clear := &rq.clears[id]
	clear.flags = flags
	clear.rgba = rgba
	clear.depth = depth
	clear.stencil = stencil
}

func (rq *RenderQueue) SetViewTransform(id uint8, view, proj *mgl32.Mat4, flags uint8) {

}

func (rq *RenderQueue) Submit(id uint8, program uint16, depth int32) uint32 {
	// uniform range
	rq.uniformEnd = uint16(rq.ub.GetPos())

//	log.Println("submit uniform end:", rq.uniformEnd)

	// encode sort-key
	sk := &rq.sk
	sk.Layer = 0

	sk.Shader = program & ID_TYPE_MASK // trip type
	sk.Blend = 0
	sk.Texture = rq.drawCall.textures[0]

	rq.sortKey[rq.drawCallNum] = rq.sk.Encode()
	rq.sortValues[rq.drawCallNum] = rq.drawCallNum

	// copy data
	rq.drawCall.uniformBegin = rq.uniformBegin
	rq.drawCall.uniformEnd = rq.uniformEnd

	rq.drawCallList[rq.drawCallNum] = rq.drawCall
	rq.drawCallNum++

	// reset state
	rq.drawCall.reset()
	rq.uniformBegin = uint16(rq.ub.GetPos())

	// return frame Num
	return 0
}

/// 执行最终的绘制
func (rq *RenderQueue) Flush() uint32 {
	//
	////// real draw !!!

	sortKeys := rq.sortKey[:rq.drawCallNum]
	sortValues := rq.sortValues[:rq.drawCallNum]
	drawList := rq.drawCallList[:rq.drawCallNum]

	rq.ctx.Draw(sortKeys, sortValues, drawList)

	rq.drawCallNum = 0
	rq.uniformBegin = 0
	rq.uniformEnd = 0
	rq.ub.Reset()

	return uint32(rq.drawCallNum)
}
