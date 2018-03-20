package effect

import (
	"korok.io/korok/gfx"
	"korok.io/korok/math/f32"
	"korok.io/korok/math"

	"log"

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

	poseStart channel_v2
	colorDelta channel_v4
	sizeDelta channel_f32
	rot channel_f32
	rotDelta channel_f32

	angle channel_f32
	angleDelta channel_f32
	radius channel_f32
	radiusDelta channel_f32

	*RadiusConfig
}

func NewRadiusSimulator(cfg *RadiusConfig) *RadiusSimulator {
	r := &RadiusSimulator{Pool:Pool{cap: cfg.Max}, RadiusConfig: cfg}
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
	r.life = r.Field(Life).(channel_f32)
	r.pose = r.Field(Position).(channel_v2)
	r.poseStart = r.Field(PositionStart).(channel_v2)
	r.color = r.Field(Color).(channel_v4)
	r.colorDelta = r.Field(ColorDelta).(channel_v4)
	r.size = r.Field(Size).(channel_f32)
	r.sizeDelta = r.Field(SizeDelta).(channel_f32)
	r.rot = r.Field(Rotation).(channel_f32)
	r.rotDelta = r.Field(RotationDelta).(channel_f32)
	r.angle = r.Field(Angle).(channel_f32)
	r.angleDelta = r.Field(AngleDelta).(channel_f32)
	r.radius = r.Field(Radius).(channel_f32)
	r.radiusDelta = r.Field(RadiusDelta).(channel_f32)

	// init controller
	r.RateController.Initialize(r.Duration, r.RadiusConfig.Rate)
}

func (r *RadiusSimulator) Simulate(dt float32) {
	if new := r.Rate(); new > 0 {
		r.newParticle(new)
	}
	n := int32(r.live)

	r.life.Sub(n, dt)
	r.angle.Integrate(n, r.angleDelta, dt)
	r.radius.Integrate(n, r.radiusDelta, dt)

	// 极坐标转换
	for i := int32(0); i < n; i ++ {
		x := float32(math.Cos(r.angle[i])) * r.radius[i]
		y := float32(math.Sin(r.angle[i])) * r.radius[i]
		r.pose[i] = f32.Vec2{x, y}
	}
	r.color.Integrate(n, r.colorDelta, dt)
	r.size.Integrate(n, r.sizeDelta, dt)
	r.rot.Integrate(n, r.rotDelta, dt)
	// recycle dead particle
	r.GC(&r.Pool)
}

func (r *RadiusSimulator) Rate() int {
	return 1
}

func (r *RadiusSimulator) newParticle(new int) {
	if (r.live + new) > r.cap {
		return
	}
	start := r.live
	r.live += new

	cfg := r.RadiusConfig

	for i := start; i < r.live; i++ {
		r.life[i] = cfg.Life.Random()
		invLife := 1/r.life[i]
		r.pose[i] = f32.Vec2{cfg.X.Random(), cfg.Y.Random()}

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
		r.color[i] = f32.Vec4{red, _g, b, a}
		r.colorDelta[i] = f32.Vec4{redd, gd, bd, ad}


		r.size[i] = cfg.Size.Start.Random()
		if cfg.Size.Start != cfg.Size.End {
			r.sizeDelta[i] = (cfg.Size.End.Random() - r.size[i]) * invLife
		}
		// rot
		r.rot[i] = cfg.Rot.Start.Random()
		if cfg.Rot.Start != cfg.Rot.End {
			r.rotDelta[i] = (cfg.Rot.End.Random() - r.rot[i]) * invLife
		}
		// start position
		r.poseStart[i] = r.pose[i]

		// radius
		r.radius[i] = cfg.Radius.Start.Random()
		if cfg.Radius.Start != cfg.Radius.End {
			r.radiusDelta[i] = (cfg.Radius.End.Random() - r.rot[i]) * invLife
		}
		// angle
		r.angle[i] = cfg.Angle.Random()
		r.angleDelta[i] = cfg.AngleDelta.Random()
	}
}

func (r *RadiusSimulator) Visualize(buf []gfx.PosTexColorVertex) {
	r.VisualController.Visualize(buf, int(r.live))
}

func (r *RadiusSimulator) Size() (live, cap int) {
	return r.live, r.cap
}


