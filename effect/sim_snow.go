package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"

	"log"
)

type SnowSimulator struct {
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

	Config struct{
		Life Var
		Size Var
		Color f32.Vec4
		Position [2]Var
		Velocity [2]Var
	}

	warm float32
}

func NewSnowSimulator(cap int, w, h float32) *SnowSimulator {
	sim := SnowSimulator{Pool: Pool{cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)

	// config
	sim.Config.Life = Var{10, 4}
	sim.Config.Color = f32.Vec4{1, 1, 1, 1}
	sim.Config.Size = Var{6, 6}
	sim.Config.Position[0] = Var{0, w}
	sim.Config.Position[1] = Var{h, 0}
	sim.Config.Velocity[0] = Var{-10,  20}
	sim.Config.Velocity[1] = Var{-50, 20}

	return &sim
}

func (sim *SnowSimulator) Initialize() {
	sim.Pool.Initialize()

	sim.life = sim.Field(Life).(channel_f32)
	sim.size = sim.Field(Size).(channel_f32)

	log.Println("init life:",len(sim.life))

	sim.position = sim.Field(Position).(channel_v2)
	sim.velocity = sim.Field(Velocity).(channel_v2)

	sim.color = sim.Field(Color).(channel_v4)

	cap := int32(sim.cap)
	// init life 10 ,20
	sim.life.SetRandom(cap, Var{10, 12})

	log.Println("set life:", sim.life[0])

	// init size
	sim.size.SetRandom(cap, Var{5, 10})
	sim.color.SetConst(cap, 1, 1, 1, 1)

	//sim.position.
	sim.position.SetConst(cap, 100, 50)
	sim.velocity.SetConst(cap, 20, 20)
}

func (sim *SnowSimulator) Simulate(dt float32) {
	// spawn new particle
	if new := sim.Rate(dt); new > 0 {
		sim.NewParticle(new)
	}

	n := sim.live

	// update old particle
	sim.life.Sub(n, dt)

	// position integrate: p' = p + v * t
	sim.position.Integrate(n, sim.velocity, dt)

	// gc
	sim.gc()
}


func (sim *SnowSimulator) Size() (live, cap int) {
	return int(sim.live), sim.cap
}

func (sim *SnowSimulator) SwapErase(i, j int32) {
	if sim.life[j] > 0 {
		sim.life[i] = sim.life[j]
		sim.size[i] = sim.size[j]
		sim.position[i] = sim.position[j]
		sim.velocity[i] = sim.velocity[j]
		sim.color[i] = sim.color[j]
	}
}


func (sim *SnowSimulator) NewParticle(new int32) {
	if (sim.live + new) > int32(sim.cap) {
		log.Println("pool overflow...")
		return
	}
	start := sim.live
	sim.live += new

	for i := start; i < sim.live; i++ {
		sim.life[i] = sim.Config.Life.Random()
		sim.color[i] = sim.Config.Color
		sim.size[i] = sim.Config.Size.Random()

		f := sim.size[i]/(sim.Config.Size.Base+sim.Config.Size.Var)
		sim.color[i][3] = f

		px := sim.Config.Position[0].Random()
		py := sim.Config.Position[1].Random()
		sim.position[i] = f32.Vec2{px, py}

		dx := sim.Config.Velocity[0].Random()
		dy := sim.Config.Velocity[1].Random()
		sim.velocity[i] = f32.Vec2{dx, dy}
	}
}


func (sim *SnowSimulator) gc() {
	for i, n := int32(0), sim.live ; i < n; i++{
		if sim.life[i] <= 0 {
			// find last live
			j := sim.live - 1
			for ; j > i && sim.life[j] <= 0; j-- {
				sim.live --
				n = sim.live
			}

			if j > i {
				sim.SwapErase(i, j)
			}
			sim.live --
			n = sim.live
		}
	}
}

func (sim *SnowSimulator) Rate(dt float32) int32 {
	return 1
}


func (sim *SnowSimulator) Visualize(buf []gfx.PosTexColorVertex) {
	size := sim.size
	pose := sim.position
	live := int(sim.live)

	// compute vbo
	for i := 0; i < live; i ++ {
		vi := i << 2
		h_size := size[i] / 2

		var (
			r = sim.color[i][0]
			g = sim.color[i][1]
			b = sim.color[i][2]
			a = sim.color[i][3]
		)

		// normalize color
		c := uint32(a*255) << 24 + uint32(b*255) << 16 + uint32(g*255) << 8 + uint32(r*255)

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