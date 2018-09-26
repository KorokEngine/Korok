package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/math"
)

// FireSimulator can simulate the fire effect.
type FountainSimulator struct {
	Pool

	RateController
	LifeController
	VisualController

	velocity Channel_v2
	deltaColor Channel_v4
	deltaRot Channel_f32

	// Configuration.
	Config struct{
		Duration, Rate float32
		Life Var
		Size Var
		Color TwoColor
		Fading bool
		Position [2]Var
		Angle Var
		Speed Var
		Gravity float32
		Rotation Var
		Additive bool
	}
}

func NewFountainSimulator(cap int) *FountainSimulator {
	sim := FountainSimulator{Pool: Pool{Cap: cap}}
	sim.AddChan(Life, Size)
	sim.AddChan(Position, Velocity)
	sim.AddChan(Color, ColorDelta)
	sim.AddChan(Rotation, RotationDelta)

	// config
	sim.Config.Duration = math.MaxFloat32
	sim.Config.Rate = float32(cap)/3
	sim.Config.Life = Var{3, .25}
	sim.Config.Color = TwoColor{f32.Vec4{1, 1, 1, 1}, f32.Vec4{1, 1, 1, 1}, false}
	sim.Config.Fading = false
	sim.Config.Size = Var{8, 2}
	sim.Config.Angle = Var{3.14/2, 3.14/3}
	sim.Config.Speed = Var{120, 20}
	sim.Config.Rotation = Var{0, 1}
	sim.Config.Gravity = -120

	return &sim
}

func (f *FountainSimulator) Initialize() {
	f.Pool.Initialize()

	f.Life = f.Field(Life).(Channel_f32)
	f.ParticleSize = f.Field(Size).(Channel_f32)
	f.Position = f.Field(Position).(Channel_v2)
	f.velocity = f.Field(Velocity).(Channel_v2)
	f.Color = f.Field(Color).(Channel_v4)
	f.deltaColor = f.Field(ColorDelta).(Channel_v4)
	f.Rotation = f.Field(Rotation).(Channel_f32)
	f.deltaRot = f.Field(RotationDelta).(Channel_f32)

	f.RateController.Initialize(f.Config.Duration, f.Config.Rate)
}

func (f *FountainSimulator) Simulate(dt float32) {
	// spawn new particle
	if new := f.Rate(dt); new > 0 {
		f.newParticle(new)
	}

	n := int32(f.Live)

	// update old particle
	f.Life.Sub(n, dt)

	// position integrate: p' = p + v * t
	f.Position.Integrate(n, f.velocity, dt)

	// v' = v + g * t
	f.velocity.Add(n, 0, f.Config.Gravity*dt)

	// spin
	f.Rotation.Integrate(n, f.deltaRot, dt)

	// Color
	f.Color.Integrate(n, f.deltaColor, dt)

	// GC
	f.GC(&f.Pool)
}

func (f *FountainSimulator) Size() (live, cap int) {
	return int(f.Live), f.Cap
}

func (f *FountainSimulator) newParticle(new int) {
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
		if f.Config.Fading {
			invLife := 1/f.Life[i]
			f.deltaColor[i] = f32.Vec4{
				-startColor[0] * invLife,
				-startColor[1] * invLife,
				-startColor[2] * invLife,
			}
		}

		px := f.Config.Position[0].Random()
		py := f.Config.Position[1].Random()
		f.Position[i] = f32.Vec2{px, py}

		a := f.Config.Angle.Random()
		s := f.Config.Speed.Random()
		f.velocity[i] = f32.Vec2{math.Cos(a)*s, math.Sin(a)*s}

		r := f.Config.Rotation.Random()
		f.deltaRot[i] = r
	}
}

func (f *FountainSimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	f.VisualController.Visualize(buf, tex, f.Live, f.Config.Additive)
}


