package gfx

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
type Batch struct {
	TextureId uint16
	depth     int16

	VertexId  uint16
	IndexId   uint16

	firstVertex uint16
	numVertex   uint16

	firstIndex uint16
	numIndex   uint16
}

type BatchObject interface {
	Fill(vertex []PosTexColorVertex)
	Size() int
}

type SortObject struct {
	SortId uint32
	Value  uint32
}
