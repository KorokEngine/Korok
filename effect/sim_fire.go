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

	velocity Channel_v2

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
	sim := FireSimulator{Pool: Pool{Cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)
	sim.AddChan(Rotation)

	// config
	sim.Config.Duration = math.MaxFloat32
	sim.Config.Rate = 10
	sim.Config.Life = Var{3, 4}
	sim.Config.Color = f32.Vec4{1, 0, 0, 1}
	sim.Config.Size = Var{34, 10}
	sim.Config.Position[0] = Var{0, 40}
	sim.Config.Position[1] = Var{0, 20}
	sim.Config.Velocity[0] = Var{10, 70}
	sim.Config.Velocity[1] = Var{10, 40}

	return &sim
}

func (f *FireSimulator) Initialize() {
	f.Pool.Initialize()

	f.Life = f.Field(Life).(Channel_f32)
	f.ParticleSize = f.Field(Size).(Channel_f32)
	f.Position = f.Field(Position).(Channel_v2)
	f.velocity = f.Field(Velocity).(Channel_v2)
	f.Color = f.Field(Color).(Channel_v4)
	f.Rotation = f.Field(Rotation).(Channel_f32)

	f.RateController.Initialize(f.Config.Duration, f.Config.Rate)
}

func (f *FireSimulator) Simulate(dt float32) {
	// spawn new particle
	if new := f.Rate(dt); new > 0 {
		f.newParticle(new)
	}

	n := int32(f.Live)

	// update old particle
	f.Life.Sub(n, dt)

	// position integrate: p' = p + v * t
	f.Position.Integrate(n, f.velocity, dt)

	// Color
	f.Color.Sub(n, 0, 0, 0, .3 * dt)

	// GC
	f.GC(&f.Pool)
}

func (f *FireSimulator) Size() (live, cap int) {
	return int(f.Live), f.Cap
}

func (f *FireSimulator) newParticle(new int) {
	if (f.Live + new) > f.Cap {
		return
	}

	start := f.Live
	f.Live += new

	for i := start; i < f.Live; i++ {
		f.Life[i] = f.Config.Life.Random()
		f.Color[i] = f.Config.Color
		f.VisualController.ParticleSize[i] = f.Config.Size.Random()

		px := f.Config.Position[0].Random()
		py := f.Config.Position[1].Random()
		f.Position[i] = f32.Vec2{px, py}

		dx := f.Config.Velocity[0].Random()
		dy := f.Config.Velocity[1].Random()
		v := f32.Vec2{dx-40, float32(30+dy)}
		f.velocity[i] = v
	}
}

func (f *FireSimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	f.VisualController.Visualize(buf, tex, f.Live)
}

