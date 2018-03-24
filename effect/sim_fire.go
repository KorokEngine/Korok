package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/math"
)

// FireSimulator can simulate the fire effect.
type FireSimulator struct {
	Pool

	RateController
	LifeController
	VisualController

	velocity channel_v2

	// Configuration.
	Config struct{
		Duration, Rate float32
		Life Var
		Size Var
		Color f32.Vec4
		Position [2]Var
		Velocity [2]Var
	}
}

func NewFireSimulator(cap int) *FireSimulator {
	sim := FireSimulator{Pool: Pool{cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)

	// config
	sim.Config.Duration = math.MaxFloat32
	sim.Config.Rate = 10
	sim.Config.Life = Var{3, 4}
	sim.Config.Color = f32.Vec4{.76, .25, .12, 1}
	sim.Config.Size = Var{34, 10}
	sim.Config.Position[0] = Var{0, 40}
	sim.Config.Position[1] = Var{0, 20}
	sim.Config.Velocity[0] = Var{10, 70}
	sim.Config.Velocity[1] = Var{10, 40}

	return &sim
}

func (f *FireSimulator) Initialize() {
	f.Pool.Initialize()

	f.life = f.Field(Life).(channel_f32)
	f.size = f.Field(Size).(channel_f32)
	f.pose = f.Field(Position).(channel_v2)
	f.velocity = f.Field(Velocity).(channel_v2)
	f.color = f.Field(Color).(channel_v4)

	f.RateController.Initialize(f.Config.Duration, f.Config.Rate)
}

func (f *FireSimulator) Simulate(dt float32) {
	// spawn new particle
	if new := f.Rate(dt); new > 0 {
		f.newParticle(new)
	}

	n := int32(f.live)

	// update old particle
	f.life.Sub(n, dt)

	// position integrate: p' = p + v * t
	f.pose.Integrate(n, f.velocity, dt)

	// color
	f.color.Sub(n, 0, 0, 0, .3 * dt)

	// GC
	f.GC(&f.Pool)
}

func (f *FireSimulator) Size() (live, cap int) {
	return int(f.live), f.cap
}

func (f *FireSimulator) newParticle(new int) {
	if (f.live + new) > f.cap {
		return
	}

	start := f.live
	f.live += new

	for i := start; i < f.live; i++ {
		f.life[i] = f.Config.Life.Random()
		f.color[i] = f.Config.Color
		f.size[i] = f.Config.Size.Random()

		px := f.Config.Position[0].Random()
		py := f.Config.Position[1].Random()
		f.pose[i] = f32.Vec2{px, py}

		dx := f.Config.Velocity[0].Random()
		dy := f.Config.Velocity[1].Random()
		v := f32.Vec2{dx-40, float32(30+dy)}
		f.velocity[i] = v
	}
}

func (f *FireSimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	f.VisualController.Visualize(buf, tex, f.live)
}

