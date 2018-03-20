package effect

import (
	"korok.io/korok/gfx"
	"korok.io/korok/math"
)

// Var define a variable value between [Base-Var/2, Base+Var/2].
type Var struct {
	Base, Var float32
}

// Used returns whether the value is empty.
func (v Var) Used() bool {
	return v.Base != 0 || v.Var != 0
}

// Random returns a value between [Base-Var/2, Base+Var/2].
func (v Var) Random() float32{
	return math.Random(v.Base-v.Var/2 , v.Base+v.Var/2)
}

// Range define a range between [Start, End].
type Range struct {
	Start, End Var
}

// Used returns whether the value is empty.
func (r Range) Used() bool {
	return r.Start.Used() || r.End.Used()
}

// HasRange returns whether Start and End is the same.
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

// Simulator define how a particle-system works.
type Simulator interface {
	// Initialize the particle simulator.
	Initialize()

	// Run the simulator with delta time.
	Simulate(dt float32)

	Stop()

	Play()

	// Write the result to vertex-buffer.
	Visualize(buf []gfx.PosTexColorVertex)

	// Return the size of the simulator.
	Size() (live, cap int)
}

// TODO:
// Emitter的概念可以提供一种可能：
// 在这里可以通过各种各样的Emitter实现，生成不同的初始粒子位置
// 这样可以实现更丰富的例子形状
// 之前要么认为粒子都是从一个点发射出来的，要么是全屏发射的，这只是hardcode了特殊情况
// 同时通过配置多个Emitter还可以实现交叉堆叠的形状
type Emitter interface {
}

// TODO:
// 基于上面的想法，还可以设计出 Updater 的概念，不同的 Updater 对粒子执行不同的
// 行走路径，这会极大的增加粒子弹性
type Updater interface {

}

// RateController is a helper struct to manage the EmitterRate.
type RateController struct {
	// control emitter-rate
	accTime float32
	threshTime float32

	// lifetime
	lifeTime float32
	duration float32
	stop bool
}

// Initialize init RateController with duration and emitter-rate.
func (ctr *RateController) Initialize(du, rate float32) {
	ctr.duration = du
	if rate == 0 {
		ctr.threshTime = 1.0/60
	} else {
		ctr.threshTime = 1.0/rate
	}
}

// Rate returns new particles that should spawn.
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

func (ctr *RateController) Stop() {
	ctr.stop = true
}

func (ctr *RateController) Play() {
	ctr.stop = false
}

// LifeController is a helper struct to manage the Life of particles.
type LifeController struct {
	// channel ref
	life channel_f32
	live int
}

// GC removes dead particles from the Pool.
func (ctr *LifeController) GC(p *Pool) {
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

// VisualController is a helper struct to write simulation result to vertex-buffer.
type VisualController struct {
	pose channel_v2
	color channel_v4
	size channel_f32
}

// Visualize write the live particles to vertex-buffer.
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


// ParticleSimulateSystem is the system that manage ParticleComp's simulation.
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
