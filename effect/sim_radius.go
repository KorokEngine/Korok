package effect

import (
	"korok.io/korok/gfx"
	"korok.io/korok/math/f32"
	"korok.io/korok/math"

)

// RadiusConfig used to configure the RadiusSimulator.
type RadiusConfig struct {
	Config

	Radius Range
	Angle, AngleDelta    Var
}

// RadiusSimulator works as the radius mode of the Cocos2D's particle-system.
type RadiusSimulator struct {
	Pool

	LifeController
	RateController
	VisualController

	poseStart  Channel_v2
	colorDelta Channel_v4
	sizeDelta  Channel_f32
	rot        Channel_f32
	rotDelta   Channel_f32

	angle       Channel_f32
	angleDelta  Channel_f32
	radius      Channel_f32
	radiusDelta Channel_f32

	*RadiusConfig
}

func NewRadiusSimulator(cfg *RadiusConfig) *RadiusSimulator {
	r := &RadiusSimulator{Pool:Pool{Cap: cfg.Max}, RadiusConfig: cfg}
	r.Pool.AddChan(Life)
	r.Pool.AddChan(Position, PositionStart)
	r.Pool.AddChan(Color, ColorDelta)
	r.Pool.AddChan(Size, SizeDelta)
	r.Pool.AddChan(Rotation, RotationDelta)

	r.Pool.AddChan(Angle, AngleDelta)
	r.Pool.AddChan(Radius, RadiusDelta)

	return r
}

// prepare data
func (r *RadiusSimulator) Initialize() {
	r.Pool.Initialize()
	r.Life = r.Field(Life).(Channel_f32)
	r.Position = r.Field(Position).(Channel_v2)
	r.poseStart = r.Field(PositionStart).(Channel_v2)
	r.Color = r.Field(Color).(Channel_v4)
	r.colorDelta = r.Field(ColorDelta).(Channel_v4)
	r.ParticleSize = r.Field(Size).(Channel_f32)
	r.sizeDelta = r.Field(SizeDelta).(Channel_f32)
	r.Rotation = r.Field(Rotation).(Channel_f32)
	r.rotDelta = r.Field(RotationDelta).(Channel_f32)
	r.angle = r.Field(Angle).(Channel_f32)
	r.angleDelta = r.Field(AngleDelta).(Channel_f32)
	r.radius = r.Field(Radius).(Channel_f32)
	r.radiusDelta = r.Field(RadiusDelta).(Channel_f32)

	// init controller
	r.RateController.Initialize(r.Duration, r.RadiusConfig.Rate)
}

func (r *RadiusSimulator) Simulate(dt float32) {
	if new := r.RateController.Rate(dt); new > 0 {
		r.newParticle(new)
	}
	n := int32(r.Live)

	r.Life.Sub(n, dt)
	r.angle.Integrate(n, r.angleDelta, dt)
	r.radius.Integrate(n, r.radiusDelta, dt)

	// 极坐标转换
	for i := int32(0); i < n; i ++ {
		x := float32(math.Cos(r.angle[i])) * r.radius[i]
		y := float32(math.Sin(r.angle[i])) * r.radius[i]
		r.Position[i] = f32.Vec2{x, y}
	}
	r.Color.Integrate(n, r.colorDelta, dt)
	r.ParticleSize.Integrate(n, r.sizeDelta, dt)
	r.Rotation.Integrate(n, r.rotDelta, dt)
	// recycle dead particle
	r.GC(&r.Pool)
}

func (r *RadiusSimulator) newParticle(new int) {
	if (r.Live + new) > r.Cap {
		return
	}
	start := r.Live
	r.Live += new

	cfg := r.RadiusConfig

	for i := start; i < r.Live; i++ {
		r.Life[i] = cfg.Life.Random()
		invLife := 1/r.Life[i]
		r.Position[i] = f32.Vec2{cfg.X.Random(), cfg.Y.Random()}

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
		r.Color[i] = f32.Vec4{red, _g, b, a}
		r.colorDelta[i] = f32.Vec4{redd, gd, bd, ad}


		r.ParticleSize[i] = cfg.Size.Start.Random()
		if cfg.Size.Start != cfg.Size.End {
			r.sizeDelta[i] = (cfg.Size.End.Random() - r.ParticleSize[i]) * invLife
		}
		// rot
		r.Rotation[i] = cfg.Rot.Start.Random()
		if cfg.Rot.Start != cfg.Rot.End {
			r.rotDelta[i] = (cfg.Rot.End.Random() - r.Rotation[i]) * invLife
		}
		// start position
		r.poseStart[i] = r.Position[i]

		// radius
		r.radius[i] = cfg.Radius.Start.Random()
		if cfg.Radius.Start != cfg.Radius.End {
			r.radiusDelta[i] = (cfg.Radius.End.Random() - r.Rotation[i]) * invLife
		}
		// angle
		r.angle[i] = cfg.Angle.Random()
		r.angleDelta[i] = cfg.AngleDelta.Random()
	}
}

func (r *RadiusSimulator) Visualize(buf []gfx.PosTexColorVertex, tex gfx.Tex2D) {
	r.VisualController.Visualize(buf, tex, int(r.Live), r.Additive)
}

func (r *RadiusSimulator) Size() (live, cap int) {
	return r.Live, r.Cap
}


