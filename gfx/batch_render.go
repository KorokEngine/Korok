package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok.io/korok/gfx/bk"

	"unsafe"
	"log"
)

// Batch Render:
// Use PosTexColorVertex struct with P4C4 format

/// A Sprite Batch TypeRender
type BatchRender struct {
	stateFlags uint64
	rgba       uint32

	// shader program
	program uint16

	// uniform handle
	umh_PJ uint16 	// Projection
	umh_S0 uint16 	// Sampler0

	// Camera
	Camera
	// batch context
	BatchContext
}

func NewBatchRender(vsh, fsh string) *BatchRender {
	br := new(BatchRender)

	// setup state
	br.stateFlags |= bk.ST_BLEND.ALPHA_NON_PREMULTIPLIED

	// setup shader
	if shId, sh := bk.R.AllocShader(vsh, fsh); shId != bk.InvalidId {
		br.program = shId
		sh.Use()

		// setup attribute
		sh.AddAttributeBinding("xyuv\x00", 0, P4C4[0])
		sh.AddAttributeBinding("rgba\x00", 0, P4C4[1])

		p := mgl32.Ortho2D(0, 480, 0, 320)
		s0 := int32(0)

		// setup uniform
		if id, _ := bk.R.AllocUniform(shId, "proj\x00", bk.UniformMat4, 1); id != bk.InvalidId {
			br.umh_PJ = id
			bk.SetUniform(id, unsafe.Pointer(&p[0]))
		}
		if id, _ := bk.R.AllocUniform(shId, "tex\x00", bk.UniformSampler, 1); id != bk.InvalidId {
			br.umh_S0 = id
			bk.SetUniform(id, unsafe.Pointer(&s0))
		}
		//bk.Touch(0)
		bk.Submit(0, shId, 0)
	}
	// setup batch context
	br.BatchContext.init()
	return br
}

func (br *BatchRender) SetCamera(camera Camera) {
	br.Camera = camera
}

// submit all batched group
func (br *BatchRender) submit(bList []Batch) {
	for i := range bList {
		b := &bList[i]

		// state
		bk.SetState(br.stateFlags, br.rgba)
		bk.SetTexture(0, br.umh_S0, b.TextureId, 0)

		// set vertex
		bk.SetVertexBuffer(0, b.VertexId, uint32(b.firstVertex), uint32(b.numVertex) )
		bk.SetIndexBuffer(b.IndexId, uint32(b.firstIndex), uint32(b.numIndex))

		// submit draw-call
		bk.Submit(0, br.program, 0)
	}
}

func (br *BatchRender) Begin(tex uint16) {
	br.BatchContext.begin(tex)
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



const MAX_BATCH_QUAD_SIZE   uint32 = 1000
const MAX_BATCH_INDEX_SIZE  uint32 = 6 * 1000
const MAX_BATCH_VERTEX_SIZE uint32 = 4 * 1000

// 管理一或多个Batch实例
// 最多可以生成 128 个 Batch 分组
// 最多可以使用 8 个 VBO 缓存
type BatchContext struct {
	// shared index
	ib *bk.IndexBuffer
	indexId   uint16
	index   []uint16

	// shared vertex
	vb [8]*bk.VertexBuffer
	vertexId [8]uint16
	vertex []PosTexColorVertex
	vertexPos uint32
	firstVertex uint32

	// state
	vbUsed 	  int
	batchUsed int
	texId     uint16

	// batch-list
	BatchList [128]Batch
}

func (bc *BatchContext)init() {
	// init shared index
	bc.index = make([]uint16, MAX_BATCH_INDEX_SIZE)
	iFormat := [6]uint16 {3, 0, 1, 3, 1, 2}
	for i := uint32(0); i < MAX_BATCH_INDEX_SIZE; i += 6 {
		copy(bc.index[i:], iFormat[:])
		iFormat[0] += 4
		iFormat[1] += 4
		iFormat[2] += 4
		iFormat[3] += 4
		iFormat[4] += 4
		iFormat[5] += 4
	}
	bc.indexId ,bc.ib = bk.R.AllocIndexBuffer(bk.Memory{unsafe.Pointer(&bc.index[0]), MAX_BATCH_INDEX_SIZE * 2})

	// init shared vertex
	bc.vertex = make([]PosTexColorVertex, MAX_BATCH_VERTEX_SIZE)
	for i := 0; i < 8; i++ {
		bc.vertexId[i], bc.vb[i] = bk.R.AllocVertexBuffer(bk.Memory{nil, MAX_BATCH_VERTEX_SIZE * 20}, 20)
	}
	bc.batchUsed = 0
	bc.vbUsed = 0
}

func (bc *BatchContext) begin(tex uint16) {
	bc.texId = tex
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
		bc.vbUsed += 1

		if bc.vbUsed >= 8 {
			log.Printf("VertexBuffer out of size: (%d, %d)", 8, bc.vbUsed)
		}
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

	batch.VertexId = bc.vertexId[bc.vbUsed]
	batch.firstVertex = uint16(bc.firstVertex)
	batch.numVertex = uint16(bc.vertexPos-bc.firstVertex)

	batch.IndexId = bc.indexId
	batch.firstIndex = uint16(batch.firstVertex/4 * 6)
	batch.numIndex = uint16(batch.numVertex/4 * 6)

	bc.batchUsed += 1
}

// upload buffer
func (bc *BatchContext) reset() {
	bc.texId = 0
	bc.firstVertex = 0
	bc.vertexPos = 0
	bc.batchUsed = 0
	bc.vbUsed = 0
}

// flushBuffer() will write and switch vertex-buffer
// we must submit batch with a end() method
func (bc *BatchContext) flushBuffer() {
	vb := bc.vb[bc.vbUsed]
	vb.Update(0, bc.vertexPos * 20, unsafe.Pointer(&bc.vertex[0]), false)
}

