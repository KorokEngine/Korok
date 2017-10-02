package gfx

import (
	"github.com/go-gl/gl/v3.2-core/gl"

	"korok/gfx/bk"
)

/**
Batch 系统设计
对所有物体进行分类：静态和动态，静态物体使用稳定的batch系统，动态物体每次重新计算batch

静态：
使用稳定的batch面临的最大问题是，batch内存的物体的可见性问题。如果batch内的物体有的可见，
有的不可见（有可能是不再视野内，也有可能是主动隐藏）. 采用 unity 方案，只减少状态切换，不减少 drawcall 的做法，
实现比较简单（否则要面临修改VBO缓存的问题）。

动态：
动态Batch每次都需要重新构建，只要实现正确的合并算法即可。

暗示：在 sortkey 中提供一个字段 batch=0 (默认情况使用动态batch)，batch=1,2..N 的情况按batch值进行
分批处理。
优化：按照空间划分可以得到不同子空间的batch，可以在batch中再做一次筛选滤掉不可见的物体。

渲染：
静态系统，按照 batch-id 找到对应的batch数据，生成渲染命令
动态系统，直接合并数据并生成渲染命令

目前Batch实现只支持格式：pos_uv_color
 */
const MaxBatch  = 2000
const QuadSize  = 32

// 使用临时内存即可
type Batch struct {
	vertex []QuadVertex
	index  []int32

	count int32

	vao uint32
	vbo, ebo bk.Buffer
}

func NewBatch() *Batch {
	b := new(Batch)
	b.vbo = bk.NewArrayBuffer(Format_POS_COLOR_UV)
	b.ebo = bk.NewIndexBuffer()
	return b
}

func (b *Batch) Add(qv [4]QuadVertex) bool{
	if b.count + QuadSize > MaxBatch {
		return false
	}

	copy(b.vertex[b.count:], qv[:])
	b.count += 1

	return true
}

//////
///////////////
///// update batch Data
func (b *Batch) Commit() {
	b.count = 3

	gl.BindVertexArray(b.vao)

	b.vbo.Update(gl.Ptr(b.vertex), int(b.count * 4 * 20))
	b.ebo.Update(gl.Ptr(b.index), int(b.count * 6 * 4))

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 20, gl.Ptr(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, gl.Ptr(8))
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 4, gl.UNSIGNED_BYTE, true, 20, gl.Ptr(16))

	gl.BindVertexArray(0)
}

//func NewBatchCommandTest(tex uint32) {
//	b := new(Batch)
//	var w, h float32 = 50, 50
//	b.vertex = []float32{
//		0,  h,  0.0, 1.0,
//		w,  0,  1.0, 0.0,
//		0,  0,  0.0, 0.0,
//		w,  h,  1.0, 1.0,
//	}
//	b.index = []int32{
//		0, 1, 2,
//		0, 3, 1,
//	}
//	b.count = 1
//	b.tex = tex
//
//	gl.GenVertexArrays(1, &b.vao)
//	gl.GenBuffers(1, &b.vbo)
//	gl.GenBuffers(1, &b.ebo)
//
//	gl.BindVertexArray(b.vao)
//	// vbo
//	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
//	gl.BufferData(gl.ARRAY_BUFFER, len(b.vertex)*4, gl.Ptr(b.vertex), gl.STATIC_DRAW)
//
//	// ebo optional
//	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
//	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(b.index)*4, gl.Ptr(b.index), gl.STATIC_DRAW)
//
//	gl.EnableVertexAttribArray(0)
//	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
//
//	gl.EnableVertexAttribArray(1)
//	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(8))
//
//	gl.BindVertexArray(0)
//}


// 在这里申请一组VBO
// 2000/20 = 100, 一次最多合并 100 个矩形， 具体多少以后再说吧
// use it as a stack
type BatchSystem struct {
	VAOs [128]uint32
	VBOs [128]uint32  		// max = 128， 共享的VBO数量, 这样最多可以合并10000个矩形应该已经足够用了
	EBOs [128]uint32

	// usage
	count int32

	// temp storage
	vertex [2000]float32
	index  [2000]int32
}

func NewBatchSystem() *BatchSystem {
	bs := new(BatchSystem)

	// init all vao & vbo
	gl.GenVertexArrays(int32(len(bs.VAOs)), &bs.VAOs[0])
	gl.GenBuffers(int32(len(bs.VBOs)), &bs.VBOs[0])
	gl.GenBuffers(int32(len(bs.EBOs)), &bs.EBOs[0])

	// init index
	n := int32(len(bs.VAOs))
	for i := int32(0); i < n; i ++ {
		vi := i * 4
		ii := i * 6
		bs.index[ii + 0] = vi + 0
		bs.index[ii + 1] = vi + 1
		bs.index[ii + 2] = vi + 2

		bs.index[ii + 3] = vi + 0
		bs.index[ii + 4] = vi + 3
		bs.index[ii + 5] = vi + 1
	}

	return bs
}
//
//func (bs *BatchSystem) NewBatch(tex uint32) *Batch {
//	b := new(Batch)
//
//	b.tex = tex
//	b.vao = bs.VAOs[bs.count]
//	b.vbo = bs.VBOs[bs.count]
//	b.ebo = bs.EBOs[bs.count]
//
//	b.vertex = bs.vertex[:]
//	b.index  = bs.index[:]
//
//	bs.count ++
//	return b
//}
//
//func (bs *BatchSystem) Reset() {
//	bs.count = 0
//}
//
//func (bs *BatchSystem) Release() {
//	bs.count = 0
//	gl.DeleteVertexArrays(int32(len(bs.VBOs)), &bs.VAOs[0])
//	gl.DeleteBuffers(int32(len(bs.VBOs)), &bs.VBOs[0])
//}
