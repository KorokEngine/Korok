package gfx

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx/bk"

	"unsafe"
	"log"
)

// Batch Render:
// Use PosTexColorVertex struct with P4C4 format

/// A Tex2D Batch TypeRender
type BatchRender struct {
	stateFlags uint64
	rgba       uint32

	// shader program
	program uint16

	// uniform handle
	umhProjection uint16 // Projection
	umhSampler0   uint16 // Sampler0

	// batch context
	BatchContext
}

func NewBatchRender(vsh, fsh string) *BatchRender {
	br := new(BatchRender)

	// setup state
	br.stateFlags |= bk.ST_BLEND.ALPHA_PREMULTIPLIED

	// setup shader
	if shId, sh := bk.R.AllocShader(vsh, fsh); shId != bk.InvalidId {
		br.program = shId
		sh.Use()

		// setup attribute
		sh.AddAttributeBinding("xyuv\x00", 0, P4C4[0])
		sh.AddAttributeBinding("rgba\x00", 0, P4C4[1])

		s0 := int32(0)
		// setup uniform
		if id, _ := bk.R.AllocUniform(shId, "proj\x00", bk.UniformMat4, 1); id != bk.InvalidId {
			br.umhProjection = id
		}
		if id, _ := bk.R.AllocUniform(shId, "tex\x00", bk.UniformSampler, 1); id != bk.InvalidId {
			br.umhSampler0 = id
			bk.SetUniform(id, unsafe.Pointer(&s0))
		}
		//bk.Touch(0)
		bk.Submit(0, shId, 0)
	}
	// setup batch context
	br.BatchContext.init()
	return br
}

func (br *BatchRender) SetCamera(camera *Camera) {
	left := camera.pos.x - camera.view.w/2
	right := camera.pos.x + camera.view.w/2
	bottom := camera.pos.y - camera.view.h/2
	top := camera.pos.y + camera.view.h/2

	p := f32.Ortho2D(left, right, bottom, top)

	// setup uniform
	bk.SetUniform(br.umhProjection, unsafe.Pointer(&p[0]))
	bk.Submit(0, br.program, 0)
}

// submit all batched group
func (br *BatchRender) submit(bList []Batch) {
	for i := range bList {
		b := &bList[i]

		// state
		bk.SetState(br.stateFlags, br.rgba)
		bk.SetTexture(0, br.umhSampler0, b.TextureId, 0)

		// set vertex
		bk.SetVertexBuffer(0, b.VertexId, uint32(b.firstVertex), uint32(b.numVertex) )
		bk.SetIndexBuffer(b.IndexId, uint32(b.firstIndex), uint32(b.numIndex))

		// submit draw-call
		bk.Submit(0, br.program, int32(b.depth))
	}
}

func (br *BatchRender) Begin(tex uint16, depth int16) {
	br.BatchContext.begin(tex, depth)
}

func (br *BatchRender) Draw(b BatchObject) {
	br.BatchContext.drawComp(b)
}

func (br *BatchRender) End() {
	br.BatchContext.end()
}

func (br *BatchRender) Flush() (num int) {
	bc := &br.BatchContext
	num = bc.batchUsed

	// flush unclosed vertex buffer
	if bc.vertexPos > 0 {
		bc.flushBuffer()
	}

	// submit
	br.submit(bc.BatchList[:bc.batchUsed])

	// reset batch state
	bc.reset()

	return
}

// 目前采用提前申请好大块空间的方式，会导致大量的内存浪费
// 之后可以把vbo管理起来，按需使用

// ~ 640k per-batch, 32k vertex, 8k quad
const MAX_BATCH_QUAD_SIZE   = uint32(8<<10)
const MAX_BATCH_VERTEX_SIZE = 4 * MAX_BATCH_QUAD_SIZE

// 管理一或多个Batch实例
// 最多可以生成 128 个 Batch 分组
// 最多可以使用 8 个 VBO 缓存
type BatchContext struct {
	vertex []PosTexColorVertex
	vertexPos uint32
	firstVertex uint32

	// state
	batchUsed int
	texId     uint16
	depth     int16

	// batch-list
	BatchList [128]Batch
}

func (bc *BatchContext)init() {
	// init shared vertex
	bc.vertex = make([]PosTexColorVertex, MAX_BATCH_VERTEX_SIZE)
	bc.batchUsed = 0
}

func (bc *BatchContext) begin(tex uint16, depth int16) {
	bc.texId = tex
	bc.depth = depth
	bc.firstVertex = bc.vertexPos
}

// 计算世界坐标并保存到 Batch 结构
//
//   3 ---- 2
//   | `    |
//   |   `  |
//   0------1

func (bc *BatchContext) drawComp(b BatchObject) {
	step := uint32(b.Size())

	if bc.vertexPos + step > MAX_BATCH_VERTEX_SIZE {
		bc.flushBuffer()
		bc.end()

		bc.vertexPos = 0
		bc.firstVertex = 0
	}

	buf := bc.vertex[bc.vertexPos:bc.vertexPos+step]
	bc.vertexPos = bc.vertexPos+step
	b.Fill(buf)
}

// commit a batch
func (bc *BatchContext) end() {
	if bc.batchUsed >= 128 {
		log.Printf("Batch List out of size:(%d, %d) ", 128, bc.batchUsed)
	}

	batch := &bc.BatchList[bc.batchUsed]
	batch.TextureId = bc.texId
	batch.depth = bc.depth

	batch.VertexId = bk.InvalidId
	batch.firstVertex = 0 //uint16(bc.firstVertex)
	batch.numVertex = uint16(bc.vertexPos-bc.firstVertex)
	batch.firstIndex = uint16(bc.firstVertex/4 * 6)
	batch.numIndex = uint16(batch.numVertex/4 * 6)

	bc.batchUsed += 1
}

// upload buffer
func (bc *BatchContext) reset() {
	bc.texId = 0
	bc.firstVertex = 0
	bc.vertexPos = 0
	bc.batchUsed = 0
}

// flushBuffer() will write and switch vertex-buffer
// we must submit batch with a end() method
func (bc *BatchContext) flushBuffer() {
	var (
		reqSize = int(bc.vertexPos)
		stride = 20
	)

	iid, _ := Context.SharedIndexBuffer()
	vid, _, vb := Context.TempVertexBuffer(reqSize, stride)

	// flush vertex-buffer
	vb.Update(0, bc.vertexPos * uint32(stride), unsafe.Pointer(&bc.vertex[0]), false)

	// backward rewrite vertex-buffer id
	for i := bc.batchUsed; i >= 0; i-- {
		if b := &bc.BatchList[i]; b.VertexId == bk.InvalidId {
			b.VertexId = vid
			b.IndexId  = iid
		}
	}
}

