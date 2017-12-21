package effect

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok.io/korok/gfx"
	"log"
	"korok.io/korok/gfx/dbg"
	"fmt"
)

type FireSimulator struct {
	Pool

	// f32 channels
	life channel_f32
	size channel_f32

	// vector channels
	position channel_v2
	velocity channel_v2

	// color channel
	color channel_v4

	// live particle count
	live int32

}

// TODO
func NewFireSimulator(cap int) Simulator {
	sim := FireSimulator{Pool: Pool{cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)
	return &sim
}

func (f *FireSimulator) Initialize() {
	f.Pool.Initialize()

	f.life = f.Field(Life).(channel_f32)
	f.size = f.Field(Size).(channel_f32)

	log.Println("init life:",len(f.life))

	f.position = f.Field(Position).(channel_v2)
	f.velocity = f.Field(Velocity).(channel_v2)

	f.color = f.Field(Color).(channel_v4)

	cap := int32(f.cap)
	// init life 10 ,20
	f.life.SetRandom(cap, Range{10, 12})

	log.Println("set life:", f.life[0])

	// init size
	f.size.SetRandom(cap, Range{5, 10})
	f.color.SetConst(cap, 1, 1, 1, 1)

	// init position = {50, 50} , velocity = {20, 20}
	//f.position.
	f.position.SetConst(cap, 100, 50)
	f.velocity.SetConst(cap, 20, 20)
}

func (f *FireSimulator) Simulate(dt float32) {
	// spawn new particle
	if new := f.Rate(dt); new > 0 {
		f.NewParticle(new)
	}

	n := f.live

	// update old particle
	f.life.Sub(n, dt)

	// position integrate: p' = p + v * t
	f.position.Integrate(n, f.velocity, dt)

	// color
	f.color.Sub(n, 0, 0, 0, .3 * dt)

	// gc
	f.gc()
}

func (f *FireSimulator) Size() (live, cap int) {
	return int(f.live), f.cap
}

func (f *FireSimulator) SwapErase(i, j int32) {
	if f.life[j] > 0 {
		f.life[i] = f.life[j]
		f.size[i] = f.size[j]
		f.position[i] = f.position[j]
		f.velocity[i] = f.velocity[j]
		f.color[i] = f.color[j]
	}
}

var vcount int32

func (f *FireSimulator) NewParticle(new int32) {
	if (f.live + new) > int32(f.cap) {
		log.Println("pool overflow...")
		return
	}
	start := f.live
	f.live += new

	for i := start; i < f.live; i++ {
		f.life[i] = random(3, 4)
		f.color[i] = mgl32.Vec4{1, 1, 1, 1}
		f.size[i] = random(5, 20)
		f.position[i] = mgl32.Vec2{240, 50}
		dx := random(10, 70)
		dy := random(10, 40)

		v := mgl32.Vec2{dx-40, float32(30+dy)}
		f.velocity[i] = v
	}

	vcount += new

	if vcount == 5000 {

	}
}


func (f *FireSimulator) gc() {
	for i, n := int32(0), f.live ; i < n; i++{
		if f.life[i] <= 0 {
			// find last live
			j := f.live - 1
			for ; j > i && f.life[j] <= 0; j-- {
				f.live --
				n = f.live
			}

			if j > i {
				f.SwapErase(i, j)
			}
			f.live --
			n = f.live
		}
	}
}

// 控制是否产生新的粒子, 返回粒子数量
// 每 0.5 秒产生一个粒子
var time float32
func (f *FireSimulator) Rate(dt float32) int32 {
	time += dt
	if time > 0.1 {
		time = 0
		return 2
	}
	return 0
}

// 总觉的这里可以有很大的优化空间，尤其是顶点和纹理数据大量的重复
// write to a vertex buffer!!!
// <u, v, r, g, b, a, x, y, size, rot>
func (f *FireSimulator) Visualize(buf []gfx.PosTexColorVertex) {
	size := f.size
	pose := f.position
	live := int(f.live)

	var zcc int
	var lcc int
	// compute vbo
	for i := 0; i < live; i ++ {
		vi := i << 2
		h_size := size[i] / 2
		h_size = 10
		alpha := f.color[i][3]
		if alpha <= 0 {
			alpha = 0
			zcc ++
		}

		if f.life[i] <= 0 {
			lcc ++
		}
		//alpha = 1
		c := 0x00ffffff + (uint32(0xff * alpha) << 24)

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

	dbg.Move(10, 280)
	dbg.DrawStrScaled(fmt.Sprintf("zcc: %d lcc: %d", zcc, lcc), 0.6)
}

