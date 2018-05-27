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

	poseStart channel_v2
	colorDelta channel_v4
	sizeDelta channel_f32
	rot channel_f32
	rotDelta channel_f32

	velocity      channel_v2
	radialAcc     channel_f32
	tangentialAcc channel_f32

	//
	gravity f32.Vec2

	// config
	*GravityConfig
}

func NewGravitySimulator(cfg *GravityConfig) *GravitySimulator {
	g := &GravitySimulator{GravityConfig: cfg}; g.cap = cfg.Max

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

	g.life = g.Field(Life).(channel_f32)
	g.pose = g.Field(Position).(channel_v2)
	g.poseStart = g.Field(PositionStart).(channel_v2)
	g.color = g.Field(Color).(channel_v4)
	g.colorDelta = g.Field(ColorDelta).(channel_v4)
	g.size = g.Field(Size).(channel_f32)
	g.sizeDelta = g.Field(SizeDelta).(channel_f32)
	g.rot = g.Field(Rotation).(channel_f32)
	g.rotDelta = g.Field(RotationDelta).(channel_f32)

	g.velocity = g.Field(Velocity).(channel_v2)
	g.radialAcc = g.Field(RadialAcc).(channel_f32)
	g.tangentialAcc = g.Field(TangentialAcc).(channel_f32)

	// init const
	g.gravity = g.GravityConfig.Gravity

	// emit rate
	g.RateController.Initialize(g.GravityConfig.Duration, g.GravityConfig.Rate)
}

func (g *GravitySimulator) Simulate(dt float32) {
	if new := g.RateController.Rate(dt); new > 0 {
		g.newParticle(new)
	}

	n := int32(g.live)

	g.life.Sub(n, dt)

	// gravity model integrate
	g.velocity.radialIntegrate(n, g.pose, g.radialAcc, dt)
	g.velocity.tangentIntegrate(n, g.pose, g.tangentialAcc, dt)

	g.velocity.Add(n, g.gravity[0]*dt, g.gravity[1]*dt)

	// position
	g.pose.Integrate(n, g.velocity, dt)

	// color
	g.color.Integrate(n, g.colorDelta, dt)

	// size
	g.size.Integrate(n, g.sizeDelta, dt)

	// angle
	g.rot.Integrate(n, g.rotDelta, dt)

	// recycle dead
	g.GC(&g.Pool)
}

func (g *GravitySimulator) newParticle(new int) {
	if (g.live + new) > g.cap {
		return
	}
	start := g.live
	g.live += new

	cfg := g.GravityConfig
	for i := start; i < g.live; i++ {
		g.life[i] = math.Random(cfg.Life.Base, cfg.Life.Base+cfg.Life.Var)
		invLife := 1/g.life[i]

		g.pose[i] = f32.Vec2{cfg.X.Random(), cfg.Y.Random()}
		// color
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
		g.color[i] = f32.Vec4{red, _g, b, a}
		g.colorDelta[i] = f32.Vec4{redd, gd, bd, ad}

		g.size[i], g.sizeDelta[i] = cfg.Size.RangeInit(invLife)
		// rot
		g.rot[i], g.rotDelta[i] = cfg.Rot.RangeInit(invLife)

		// start position
		g.poseStart[i] = g.pose[i]

		// gravity
		g.radialAcc[i] = cfg.RadialAcc.Random()
		g.tangentialAcc[i] = cfg.TangentialAcc.Random()

		// velocity = speed * direction
		a, s := cfg.Angel.Random(), cfg.Speed.Random()
		g.velocity[i] = f32.Vec2{math.Cos(a)*s, math.Sin(a)*s}
	}
}

func (g *GravitySimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	g.VisualController.Visualize(buf, tex, g.live)
}

func (r *GravitySimulator) Size() (live, cap int) {
	return r.live, r.cap
}
