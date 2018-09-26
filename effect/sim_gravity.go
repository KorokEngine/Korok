package effect

import (
	"korok.io/korok/math"
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
)

type Config struct {
	Max int

	Duration float32
	Rate float32 	// Number of particles per-second

	Life  Var
	X, Y  Var
	Size  Range
	Rot   Range

	R, G, B, A Range
}

// GravityConfig used to configure the GravitySimulator.
type GravityConfig struct {
	Config

	// gravity
	Gravity f32.Vec2

	// speed and direction
	Speed Var
	Angel Var

	// Radial acceleration
	RadialAcc Var

	// tangent acceleration
	TangentialAcc Var

	RotationIsDir bool
}

// GravitySimulator works as the gravity mode of Cocos2D's particle-system.
type GravitySimulator struct {
	Pool

	RateController
	LifeController
	VisualController

	poseStart  Channel_v2
	colorDelta Channel_v4
	sizeDelta  Channel_f32
	rot        Channel_f32
	rotDelta   Channel_f32

	velocity      Channel_v2
	radialAcc     Channel_f32
	tangentialAcc Channel_f32

	//
	gravity f32.Vec2

	// config
	*GravityConfig
}

func NewGravitySimulator(cfg *GravityConfig) *GravitySimulator {
	g := &GravitySimulator{GravityConfig: cfg}; g.Cap = cfg.Max

	g.Pool.AddChan(Life)
	g.Pool.AddChan(Position, PositionStart)
	g.Pool.AddChan(Color, ColorDelta)
	g.Pool.AddChan(Size, SizeDelta)
	g.Pool.AddChan(Rotation, RotationDelta)

	g.Pool.AddChan(Velocity)
	g.Pool.AddChan(RadialAcc)
	g.Pool.AddChan(TangentialAcc)

	return g
}

// alloc a big block
func (g *GravitySimulator) Initialize() {
	g.Pool.Initialize()

	g.Life = g.Field(Life).(Channel_f32)
	g.Position = g.Field(Position).(Channel_v2)
	g.poseStart = g.Field(PositionStart).(Channel_v2)
	g.Color = g.Field(Color).(Channel_v4)
	g.colorDelta = g.Field(ColorDelta).(Channel_v4)
	g.ParticleSize = g.Field(Size).(Channel_f32)
	g.sizeDelta = g.Field(SizeDelta).(Channel_f32)
	g.Rotation = g.Field(Rotation).(Channel_f32)
	g.rotDelta = g.Field(RotationDelta).(Channel_f32)

	g.velocity = g.Field(Velocity).(Channel_v2)
	g.radialAcc = g.Field(RadialAcc).(Channel_f32)
	g.tangentialAcc = g.Field(TangentialAcc).(Channel_f32)

	// init const
	g.gravity = g.GravityConfig.Gravity

	// emit rate
	g.RateController.Initialize(g.GravityConfig.Duration, g.GravityConfig.Rate)
}

func (g *GravitySimulator) Simulate(dt float32) {
	if new := g.RateController.Rate(dt); new > 0 {
		g.newParticle(new)
	}

	n := int32(g.Live)

	g.Life.Sub(n, dt)

	// gravity model integrate
	g.velocity.radialIntegrate(n, g.Position, g.radialAcc, dt)
	g.velocity.tangentIntegrate(n, g.Position, g.tangentialAcc, dt)

	g.velocity.Add(n, g.gravity[0]*dt, g.gravity[1]*dt)

	// position
	g.Position.Integrate(n, g.velocity, dt)

	// Color
	g.Color.Integrate(n, g.colorDelta, dt)

	// ParticleSize
	g.ParticleSize.Integrate(n, g.sizeDelta, dt)

	// angle
	g.Rotation.Integrate(n, g.rotDelta, dt)

	// recycle dead
	g.GC(&g.Pool)
}

func (g *GravitySimulator) newParticle(new int) {
	if (g.Live + new) > g.Cap {
		return
	}
	start := g.Live
	g.Live += new

	cfg := g.GravityConfig
	for i := start; i < g.Live; i++ {
		g.Life[i] = math.Random(cfg.Life.Base, cfg.Life.Base+cfg.Life.Var)
		invLife := 1/g.Life[i]

		g.Position[i] = f32.Vec2{cfg.X.Random(), cfg.Y.Random()}
		// Color
		var red, _g, b, a  float32 = 0, 0, 0, 1
		var redd, gd, bd, ad float32

		if cfg.R.Used() {
			red, redd = cfg.R.RangeInit(invLife)
		}
		if cfg.G.Used() {
			_g, gd = cfg.G.RangeInit(invLife)
		}
		if cfg.B.Used() {
			b, bd = cfg.B.RangeInit(invLife)
		}
		if cfg.A.Used() {
			a, ad = cfg.A.RangeInit(invLife)
		}
		g.Color[i] = f32.Vec4{red, _g, b, a}
		g.colorDelta[i] = f32.Vec4{redd, gd, bd, ad}

		g.ParticleSize[i], g.sizeDelta[i] = cfg.Size.RangeInit(invLife)
		// rot
		g.Rotation[i], g.rotDelta[i] = cfg.Rot.RangeInit(invLife)

		// start position
		g.poseStart[i] = g.Position[i]

		// gravity
		g.radialAcc[i] = cfg.RadialAcc.Random()
		g.tangentialAcc[i] = cfg.TangentialAcc.Random()

		// velocity = speed * direction
		a, s := cfg.Angel.Random(), cfg.Speed.Random()
		g.velocity[i] = f32.Vec2{math.Cos(a)*s, math.Sin(a)*s}
	}
}

func (g *GravitySimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	g.VisualController.Visualize(buf, tex, g.Live)
}

func (r *GravitySimulator) Size() (live, cap int) {
	return r.Live, r.Cap
}
