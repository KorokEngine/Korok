package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/math"
)

// SnowSimulator can simulate snow effect.
type SnowSimulator struct {
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

func NewSnowSimulator(cap int, w, h float32) *SnowSimulator {
	sim := SnowSimulator{Pool: Pool{Cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)
	sim.AddChan(Rotation)

	// config
	sim.Config.Duration = math.MaxFloat32
	sim.Config.Rate = 60
	sim.Config.Life = Var{10, 4}
	sim.Config.Color = f32.Vec4{1, 0, 0, 1}
	sim.Config.Size = Var{6, 6}
	sim.Config.Position[0] = Var{0, w}
	sim.Config.Position[1] = Var{h/2, 0}
	sim.Config.Velocity[0] = Var{-10,  20}
	sim.Config.Velocity[1] = Var{-50, 20}

	return &sim
}

func (sim *SnowSimulator) Initialize() {
	sim.Pool.Initialize()

	sim.Life = sim.Field(Life).(Channel_f32)
	sim.ParticleSize = sim.Field(Size).(Channel_f32)
	sim.Position = sim.Field(Position).(Channel_v2)
	sim.velocity = sim.Field(Velocity).(Channel_v2)
	sim.Color = sim.Field(Color).(Channel_v4)
	sim.Rotation = sim.Field(Rotation).(Channel_f32)

	sim.RateController.Initialize(sim.Config.Duration, sim.Config.Rate)
}

func (sim *SnowSimulator) Simulate(dt float32) {
	if new := sim.Rate(dt); new > 0 {
		sim.NewParticle(new)
	}

	n := int32(sim.Live)

	// update old particle
	sim.Life.Sub(n, dt)

	// position integrate: p' = p + v * t
	sim.Position.Integrate(n, sim.velocity, dt)

	// GC
	sim.GC(&sim.Pool)
}


func (sim *SnowSimulator) Size() (live, cap int) {
	return int(sim.Live), sim.Cap
}

func (sim *SnowSimulator) NewParticle(new int) {
	if (sim.Live + new) > sim.Cap {
		return
	}
	start := sim.Live
	sim.Live += new

	for i := start; i < sim.Live; i++ {
		sim.Life[i] = sim.Config.Life.Random()
		sim.Color[i] = sim.Config.Color
		sim.ParticleSize[i] = sim.Config.Size.Random()

		f := sim.ParticleSize[i]/(sim.Config.Size.Base+sim.Config.Size.Var)
		sim.Color[i][3] = f

		px := sim.Config.Position[0].Random()
		py := sim.Config.Position[1].Random()
		sim.Position[i] = f32.Vec2{px, py}

		dx := sim.Config.Velocity[0].Random()
		dy := sim.Config.Velocity[1].Random()
		sim.velocity[i] = f32.Vec2{dx, dy}
	}
}

func (sim *SnowSimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	sim.VisualController.Visualize(buf, tex, int(sim.Live), false)
}