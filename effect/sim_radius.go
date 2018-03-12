package effect

import (
	"korok.io/korok/gfx"
	"korok.io/korok/math/f32"

	"math"
	"log"

)

type RadiusConfig struct {
	Config

	Radius Range
	Angle, AngleDelta    Var
}

// Radius 模式既以极坐标表示的物理模拟
// 方便模拟粒子旋转放射的场景
type RadiusSimulator struct {
	Pool
	*RadiusConfig

	// channel ref
	life channel_f32

	pose channel_v2
	poseStart channel_v2
	color channel_v4
	colorDelta channel_v4
	size channel_f32
	sizeDelta channel_f32
	rot channel_f32
	rotDelta channel_f32

	angle channel_f32
	angleDelta channel_f32
	radius channel_f32
	radiusDelta channel_f32

	live int
}

func NewRadiusSimulator(cfg *RadiusConfig) *RadiusSimulator {
	r := &RadiusSimulator{Pool:Pool{cap: 1024}, RadiusConfig: cfg}
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
		x := float32(math.Cos(float64(r.angle[i]))) * r.radius[i]
		y := float32(math.Sin(float64(r.angle[i]))) * r.radius[i]
		r.pose[i] = f32.Vec2{x, y}
	}
	r.color.Integrate(n, r.colorDelta, dt)
	r.size.Integrate(n, r.sizeDelta, dt)
	r.rot.Integrate(n, r.rotDelta, dt)
	// recycle dead particle
	r.gc()
}

func (r *RadiusSimulator) Rate() int {
	return 1
}

func (r *RadiusSimulator) newParticle(new int) {
	if (r.live + new) > r.cap {
		log.Println("pool overflow...")
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
		var _r, g, b, a  float32 = 1, 1, 1, 1
		var redd, gd, bd, ad float32
		if cfg.R.Used() {
			_r = cfg.R.Start.Random()
			if cfg.R.HasRange() {
				redd = (cfg.R.End.Random() - _r) * invLife
			}
		}
		if cfg.G.Used() {
			g = cfg.G.Start.Random()
			if cfg.G.HasRange() {
				gd = (cfg.G.End.Random() - g) * invLife
			}
		}
		if cfg.B.Used() {
			b = cfg.B.Start.Random()
			if cfg.B.HasRange() {
				bd = (cfg.B.End.Random() - b) * invLife
			}
		}
		if cfg.A.Used() {
			a = cfg.A.Start.Random()
			if cfg.A.HasRange() {
				ad = (cfg.A.End.Random() - a) * invLife
			}
		}
		r.color[i] = f32.Vec4{_r, g, b, a}
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

func (r *RadiusSimulator) gc() {
	for i, n := 0, r.live ; i < n; i++{
		if r.life[i] <= 0 {
			// find last live
			j := r.live - 1
			for ; j > i && r.life[j] <= 0; j-- {
				r.live --
				n = r.live
			}

			if j > i {
				r.SwapErase(i, j)
			}
			r.live --
			n = r.live
		}
	}
}

func (r *RadiusSimulator) SwapErase(i, j int) {
	r.life[i] = r.life[j]
	r.pose[i] = r.pose[j]
	r.poseStart[i] = r.poseStart[j]
	r.color[i] = r.color[j]
	r.colorDelta[i] = r.colorDelta[j]
	r.size[i] = r.size[j]
	r.sizeDelta[i] = r.sizeDelta[j]
	r.rot[i] = r.rot[j]
	r.rotDelta[i] = r.rotDelta[j]
	r.angle[i] = r.angle[i]
	r.angleDelta[i] = r.angleDelta[j]
	r.radius[i] = r.radius[j]
	r.radiusDelta[i] = r.radiusDelta[i]
}

func (r *RadiusSimulator) Visualize(buf []gfx.PosTexColorVertex) {
	size := r.size
	pose := r.pose
	live := int(r.live)

	var zcc int
	var lcc int
	// compute vbo
	for i := 0; i < live; i ++ {
		vi := i << 2
		h_size := size[i] / 2
		alpha := r.color[i][3]
		if alpha <= 0 {
			alpha = 0
			zcc ++
		}

		if r.life[i] <= 0 {
			lcc ++
		}
		c := 0x00ffffff + (uint32(0xff*alpha) << 24)

		// bottom-left
		buf[vi+0].X = pose[i][0] - h_size
		buf[vi+0].Y = pose[i][1] - h_size
		buf[vi+0].U = 0
		buf[vi+0].V = 0
		buf[vi+0].RGBA = c

		// bottom-right
		buf[vi+1].X = pose[i][0] + h_size
		buf[vi+1].Y = pose[i][1] - h_size
		buf[vi+1].U = 1
		buf[vi+1].V = 0
		buf[vi+1].RGBA = c

		// top-right
		buf[vi+2].X = pose[i][0] + h_size
		buf[vi+2].Y = pose[i][1] + h_size
		buf[vi+2].U = 1
		buf[vi+2].V = 1
		buf[vi+2].RGBA = c

		// top-left
 		buf[vi+3].X = pose[i][0] - h_size
		buf[vi+3].Y = pose[i][1] + h_size
		buf[vi+3].U = 0
		buf[vi+3].V = 1
		buf[vi+3].RGBA = c
	}
}

func (r *RadiusSimulator) Size() (live, cap int) {
	return r.live, r.cap
}


