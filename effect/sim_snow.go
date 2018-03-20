package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"

	"log"
	"korok.io/korok/math"
)

// SnowSimulator can simulate snow effect.
type SnowSimulator struct {
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

func NewSnowSimulator(cap int, w, h float32) *SnowSimulator {
	sim := SnowSimulator{Pool: Pool{cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)

	// config
	sim.Config.Duration = math.MaxFloat32
	sim.Config.Rate = 60
	sim.Config.Life = Var{10, 4}
	sim.Config.Color = f32.Vec4{1, 0, 0, 1}
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
	sim.pose = sim.Field(Position).(channel_v2)
	sim.velocity = sim.Field(Velocity).(channel_v2)
	sim.color = sim.Field(Color).(channel_v4)

	sim.RateController.Initialize(sim.Config.Duration, sim.Config.Rate)
}

func (sim *SnowSimulator) Simulate(dt float32) {
	if new := sim.Rate(dt); new > 0 {
		sim.NewParticle(new)
	}

	n := int32(sim.live)

	// update old particle
	sim.life.Sub(n, dt)

	// position integrate: p' = p + v * t
	sim.pose.Integrate(n, sim.velocity, dt)

	// GC
	sim.GC(&sim.Pool)
}


func (sim *SnowSimulator) Size() (live, cap int) {
	return int(sim.live), sim.cap
}

func (sim *SnowSimulator) NewParticle(new int) {
	if (sim.live + new) > sim.cap {
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
		sim.pose[i] = f32.Vec2{px, py}

		dx := sim.Config.Velocity[0].Random()
		dy := sim.Config.Velocity[1].Random()
		sim.velocity[i] = f32.Vec2{dx, dy}
	}
}

func (sim *SnowSimulator) Visualize(buf []gfx.PosTexColorVertex) {
	sim.VisualController.Visualize(buf, int(sim.live))
}