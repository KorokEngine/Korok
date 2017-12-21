package gfx

import "korok.io/korok/gfx/bk"

// graphics context
// a wrapper for bk-api

func Init() {
	bk.Init()
	bk.Reset(480, 320)

	// Enable debug text
	bk.SetDebug(bk.DEBUG_R|bk.DEBUG_Q)
}

func Flush() {
	bk.Flush()
}

func Destroy() {
	bk.Destroy()
}

// 目前各个 RenderFeature 都是自己管理 VBO/IBO，但是对于一些系统，比如
// Batch/ParticleSystem(2D中的大部分元素)，都是可以复用VBO的，顶点数据
// 需要每帧动态生成，如此可以把这些需要动态申请的Buffer在此管理起来，对应的
// CPU 数据可以在 StackAllocator 上申请，一帧之后就自动释放。
type context struct {
	Stack StackAllocator
}

// global shared
var Context *context = &context{}