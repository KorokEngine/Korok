package effect

import (
	"korok.io/korok/gfx"
)

/**
	需要解决两个问题：
	1. 粒子系统的着色器能够直接使用Mesh着色器？如果不能，需要考虑如何切换着色器的问题，由于渲染需要从后往前，所以
		必须统一到渲染系统里面渲染，不能单独写成一个新的渲染系统。

		- 粒子系统渲染使用自定义的粒子着色器，不使用针对精灵的Mesh着色器
	2. 现在已经勾画出了整个业务流程，考虑如何将之动态化，这样可以通过配置脚本来驱动粒子系统！

		- 暂时无法实现动态化，构建大而全的系统会牺牲一部分的性能和内存
 */

//// 算子定义操作：赋值， 相加，积分

type Simulator interface {
	Initialize()

	Simulate(dt float32)

	Visualize(buf []gfx.PosTexColorVertex)

	Size() (live, cap int)
}

// Emitter的概念可以提供一种可能：
// 在这里可以通过各种各样的Emitter实现，生成不同的初始粒子位置
// 这样可以实现更丰富的例子形状
// 之前要么认为粒子都是从一个点发射出来的，要么是全屏发射的，这只是hardcode了特殊情况
// 同时通过配置多个Emitter还可以实现交叉堆叠的形状
type Emitter interface {
}

// 基于上面的想法，还可以设计出 Updater 的概念，不同的 Updater 对粒子执行不同的
// 行走路径，这会极大的增加粒子弹性
type Updater interface {

}





// ps-comp simulate
// 在仿真系统中，直接读取 PSTable 的 Comp 进行模拟仿真
type ParticleSimulateSystem struct {
	pst *ParticleSystemTable
	init bool
}

func NewSimulationSystem () *ParticleSimulateSystem {
	return &ParticleSimulateSystem{}
}

func (pss *ParticleSimulateSystem) RequireTable(tables []interface{}) {
	for _, t := range tables {
		switch table := t.(type) {
		case *ParticleSystemTable:
			pss.pst = table
		}
	}
}

// System 的生命周期中，应该安排一个 Initialize 的阶段
func (pss *ParticleSimulateSystem) Initialize() {
	et := pss.pst
	for i, n := 0, et.index; i < n; i++ {
		et.comps[i].sim.Initialize()
	}
}

func (pss *ParticleSimulateSystem) Update(dt float32) {
	if !pss.init {
		pss.Initialize()
		pss.init = true
	}
	et := pss.pst
	for i, n := 0, et.index; i < n; i++ {
		et.comps[i].sim.Simulate(dt)
	}
}
