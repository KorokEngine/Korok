package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/math"
)

// FireSimulator can simulate the fire effect.
type ExplosionSimulator struct {
	Pool

	RateController
	LifeController
	VisualController

	velocity Channel_v2
	deltaColor Channel_v4

	// Configuration.
	Config struct{
		Duration, Rate float32
		Life Var
		Size Var
		Color TwoColor
		Position [2]Var
		Angle Var
		Speed Var
		Additive bool
	}
}

func NewExplosionSimulator(cap int) *ExplosionSimulator {
	sim := ExplosionSimulator{Pool: Pool{Cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color)
	sim.AddChan(Rotation)
	sim.AddChan(ColorDelta)

	// config
	sim.Config.Duration = .1
	sim.Config.Rate = float32(cap)/sim.Config.Duration
	sim.Config.Life = Var{3, 1}
	sim.Config.Color = TwoColor{f32.Vec4{1, 1, 0, 1}, f32.Vec4{.973, .349, 0, 1}, true}
	sim.Config.Size = Var{15, 10}
	sim.Config.Angle = Var{0, 6.28}
	sim.Config.Speed = Var{100, 40}

	return &sim
}

func (f *ExplosionSimulator) Initialize() {
	f.Pool.Initialize()

	f.Life = f.Field(Life).(Channel_f32)
	f.ParticleSize = f.Field(Size).(Channel_f32)
	f.Position = f.Field(Position).(Channel_v2)
	f.velocity = f.Field(Velocity).(Channel_v2)
	f.Color = f.Field(Color).(Channel_v4)
	f.deltaColor = f.Field(ColorDelta).(Channel_v4)
	f.Rotation = f.Field(Rotation).(Channel_f32)

	f.RateController.Initialize(f.Config.Duration, f.Config.Rate)
}

func (f *ExplosionSimulator) Simulate(dt float32) {
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
	f.Color.Integrate(n, f.deltaColor, dt)

	// GC
	f.GC(&f.Pool)
}

func (f *ExplosionSimulator) Size() (live, cap int) {
	return int(f.Live), f.Cap
}

func (f *ExplosionSimulator) newParticle(new int) {
	if (f.Live + new) > f.Cap {
		return
	}

	start := f.Live
	f.Live += new

	for i := start; i < f.Live; i++ {
		f.Life[i] = f.Config.Life.Random()
		f.ParticleSize[i] = f.Config.Size.Random()
		startColor := f.Config.Color.Random()
		f.Color[i] = startColor
		invLife := 1/f.Life[i]
		f.deltaColor[i] = f32.Vec4{
			-startColor[0] * invLife,
			-startColor[1] * invLife,
			-startColor[2] * invLife,
			-startColor[3] * invLife,
		}

		px := f.Config.Position[0].Random()
		py := f.Config.Position[1].Random()
		f.Position[i] = f32.Vec2{px, py}

		a := f.Config.Angle.Random()
		s := f.Config.Speed.Random()
		f.velocity[i] = f32.Vec2{math.Cos(a)*s, math.Sin(a)*s}
		f.Rotation[i] = a
	}
}

func (f *ExplosionSimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	f.VisualController.Visualize(buf, tex, f.Live, f.Config.Additive)
}

