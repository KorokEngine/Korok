package effect

import (
	"korok.io/korok/gfx"
	"korok.io/korok/math"
)

/**
	需要解决两个问题：
	1. 粒子系统的着色器能够直接使用Mesh着色器？如果不能，需要考虑如何切换着色器的问题，由于渲染需要从后往前，所以
		必须统一到渲染系统里面渲染，不能单独写成一个新的渲染系统。

		- 粒子系统渲染使用自定义的粒子着色器，不使用针对精灵的Mesh着色器
	2. 现在已经勾画出了整个业务流程，考虑如何将之动态化，这样可以通过配置脚本来驱动粒子系统！

		- 暂时无法实现动态化，构建大而全的系统会牺牲一部分的性能和内存
 */


/**
	显示一个粒子最终只需要：Position/Rotation/Scale

	但是对粒子进行建模需要更多的属性：
	1. 粒子发射频率
	2. 重力/重力模式
	3. 方向
	4. 速度/加速度
	5. 切线速度/微分
	6. 角速度
	7. 放射模式
	8. 开始放射半径/微分
	9. 结束放射半径/微分
	10. 旋转
	11. 共有的属性
		1. 生命/微分
		2. 初始/结束旋转/微分
		3. 初始/结束大小/微分
		4. 初始/结束颜色/微分
		5. 混合方程
		6. 纹理

	ParticleDesigner 提供了两种模式：
	1. Radius  - 基于极坐标建模的仿真，需要把极坐标转化为笛卡尔坐标绘制
	2. Gravity - 基于笛卡尔坐标建模的仿真，也提供了切线/法线加速度的支持

	以上基于配置的建模方式，在实现的时候需要提前申请好所有需要的变量，通常仿真
	可能只需要很少的属性配置，而大部分没有使用，这是很浪费的。
*/


//// 算子定义操作：赋值， 相加，积分


type EmitterMode int32

// 发射模式
const (
	ModeGravity EmitterMode = iota
	ModeRadius
)

// base value + var value
type Var struct {
	Base, Var float32
}

func (v Var) Used() bool {
	return v.Base != 0 || v.Var != 0
}

func (v Var) Random() float32{
	return math.Random(v.Base-v.Var/2 , v.Base+v.Var/2)
}

// range [start, end]
type Range struct {
	Start, End Var
}

func (r Range) Used() bool {
	return r.Start.Used() || r.End.Used()
}

func (r Range) HasRange() bool {
	return r.Start != r.End
}

func (r *Range) RangeInit(invLife float32) (start, d float32) {
	start = r.Start.Random()
	if r.Start != r.End {
		d = (r.End.Random() - start) * invLife
	}
	return
}

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

// 把粒子系统的各个子控制器拆成小的片段，这样
// 通过组合就可以最优化的实现一个例子仿真器
type RateController struct {
	// control emitter-rate
	accTime float32
	threshTime float32

	// lifetime
	lifeTime float32
	duration float32
	stop bool
}

func (ctr *RateController) Rate(dt float32) (n int) {
	ctr.lifeTime += dt
	if ctr.stop || ctr.lifeTime > ctr.duration {
		return
	}
	ctr.accTime += dt
	if ctr.accTime >= ctr.threshTime {
		acc := ctr.accTime
		for d := ctr.threshTime ; acc > d; {
			acc -= d; n++
		}
		ctr.accTime = acc
	}
	return
}

type LifeController struct {
	// channel ref
	life channel_f32
	live int
}

func (ctr *LifeController) gc(p *Pool) {
	for i, n := 0, ctr.live ; i < n; i++{
		if ctr.life[i] <= 0 {
			// find last live
			j := ctr.live - 1
			for ; j > i && ctr.life[j] <= 0; j-- {
				ctr.live --
				n = ctr.live
			}

			if j > i {
				p.Swap(i, j)
			}
			ctr.live --
			n = ctr.live
		}
	}
}

type VisualController struct {
	pose channel_v2
	color channel_v4
	size channel_f32
}

func (ctr *VisualController) Visualize(buf []gfx.PosTexColorVertex, live int) {
	size := ctr.size
	pose := ctr.pose

	// compute vbo
	for i := 0; i < live; i ++ {
		vi := i << 2
		h_size := size[i] / 2

		var (
			r = math.Clamp(ctr.color[i][0], 0, 1)
			g_ = math.Clamp(ctr.color[i][1], 0, 1)
			b = math.Clamp(ctr.color[i][2], 0, 1)
			a = math.Clamp(ctr.color[i][3], 0, 1)
		)

		c := uint32(a*255) << 24 + uint32(b*255) << 16 + uint32(g_*255) << 8 + uint32(r*255)


		// bottom-left
		buf[vi+0].X = pose[i][0] - h_size
		buf[vi+0].Y = pose[i][1] - h_size
		buf[vi+0].U = 0
		buf[vi+0].V = 0
		buf[vi+0].RGBA = c

		// bottom-right
		buf[vi+1].X = pose[i][0] + h_size
		buf[vi+1].Y = pose[i][1] - h_size
		buf[vi+1].U = 1
		buf[vi+1].V = 0
		buf[vi+1].RGBA = c

		// top-right
		buf[vi+2].X = pose[i][0] + h_size
		buf[vi+2].Y = pose[i][1] + h_size
		buf[vi+2].U = 1
		buf[vi+2].V = 1
		buf[vi+2].RGBA = c

		// top-left
		buf[vi+3].X = pose[i][0] - h_size
		buf[vi+3].Y = pose[i][1] + h_size
		buf[vi+3].U = 0
		buf[vi+3].V = 1
		buf[vi+3].RGBA = c
	}
}

type ColorController struct {

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
