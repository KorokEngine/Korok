package effect

import (
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"korok.io/korok/engi/math"
	"korok.io/korok/gfx"
)

// name convention: r = red, d_r = derivative of r with respect to time
// integrate: r' = r + d_r * t
//

/**
duration: life of system

life: life of particle
x, y: position of particle
size: size of particle
rot : self-rotation of particle
rgba: color of particle

delta size = size_var/life
delta rgba = rgba_var/life

对于配置有个简单的约定，几乎所有的参数都是 VAR 类型，
Var { Base, Var float32 }
这样可以产生一个随机值。对于有明确范围的参数，比如颜色
的是在 [0, 1], 那么它的参数是 Range 类型，我们会自动
根据life计算出变化值。对于没有范围的参数，往往会提供
delta值（既变化率）用来计算变化
 */
type Config struct {
	Max int

	Duration float32

	Life  Var
	X, Y  Var
	Size  Range
	Rot   Range

	R, G, B, A Range
}

type GravityConfig struct {
	Config

	// gravity
	Gravity mgl32.Vec2

	// speed and d
	Velocity [2]Var

	// Radial acceleration
	RadialAcc Var

	// tangent acceleration
	TangentialAcc Var

	RotationIsDir bool
}

type GravitySimulator struct {
	Pool
	// config
	*GravityConfig

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

	velocity      channel_v2
	radialAcc     channel_f32
	tangentialAcc channel_f32

	//
	gravity mgl32.Vec2

	live int
}

func NewGravitySimulator(cfg *GravityConfig) *GravitySimulator {
	g := &GravitySimulator{Pool:Pool{cap: cfg.Max}, GravityConfig: cfg}
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

// alloc a big block ?
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
}

func (g *GravitySimulator) Simulate(dt float32) {
	if new := g.Rate(); new > 0 {
		g.newParticle(new)
	}

	n := int32(g.live)

	g.life.Sub(n, dt)

	// gravity model integrate
	//g.velocity.radialIntegrate(n, g.pose, g.radialAcc, dt)
	//g.velocity.tangentIntegrate(n, g.pose, g.tangentialAcc, dt)
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
	g.gc()
}


func (g *GravitySimulator) Rate() int {
	return 1
}

func (g *GravitySimulator) newParticle(new int) {
	if (g.live + new) > g.cap {
		log.Println("pool overflow...")
		return
	}
	start := g.live
	g.live += new

	cfg := g.GravityConfig
	for i := start; i < g.live; i++ {
		g.life[i] = rdm(cfg.Life)
		invLife := 1/g.life[i]

		g.pose[i] = mgl32.Vec2{cfg.X.Random(), cfg.Y.Random()}
		// color
		var red, _g, b, a  float32 = 1, 1, 1, 1
		var redd, gd, bd, ad float32
		if cfg.R.Used() {
			red = cfg.R.Start.Random()
			if cfg.R.HasRange() {
				redd = (cfg.R.End.Random() - red) * invLife
			}
		}
		if cfg.G.Used() {
			_g = cfg.G.Start.Random()
			if cfg.G.HasRange() {
				gd = (cfg.G.End.Random() - _g) * invLife
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
		g.color[i] = mgl32.Vec4{red, _g, b, a}
		g.colorDelta[i] = mgl32.Vec4{redd, gd, bd, ad}

		g.size[i] = cfg.Size.Start.Random()
		if cfg.Size.Start != cfg.Size.End {
			g.sizeDelta[i] = (cfg.Size.End.Random() - g.size[i]) * invLife
		}
		// rot
		g.rot[i] = cfg.Rot.Start.Random()
		if cfg.Rot.Start != cfg.Rot.End {
			g.rotDelta[i] = (cfg.Rot.End.Random() - g.rot[i]) * invLife
		}
		// start position
		g.poseStart[i] = g.pose[i]

		// gravity
		g.radialAcc[i] = cfg.RadialAcc.Random()
		g.tangentialAcc[i] = cfg.TangentialAcc.Random()
		g.velocity[i] = mgl32.Vec2{cfg.Velocity[0].Random(), cfg.Velocity[1].Random()}
	}
}

func rdm(p Var) float32 {
	return math.Random(p.Base, p.Base + p.Var)
}

func (g *GravitySimulator) gc() {
	for i, n := 0, g.live ; i < n; i++{
		if g.life[i] <= 0 {
			// find last live
			j := g.live - 1
			for ; j > i && g.life[j] <= 0; j-- {
				g.live --
				n = g.live
			}

			if j > i {
				g.SwapErase(i, j)
			}
			g.live --
			n = g.live
		}
	}
}

func (g *GravitySimulator) SwapErase(i, j int) {
	g.life[i] = g.life[j]
	g.pose[i] = g.pose[j]
	g.poseStart[i] = g.poseStart[j]
	g.color[i] = g.color[j]
	g.colorDelta[i] = g.colorDelta[j]
	g.size[i] = g.size[j]
	g.sizeDelta[i] = g.sizeDelta[j]
	g.rot[i] = g.rot[j]
	g.rotDelta[i] = g.rotDelta[j]
	g.velocity[i] = g.velocity[j]
	g.radialAcc[i] = g.radialAcc[j]
	g.tangentialAcc[i] = g.tangentialAcc[j]
}


func (r *GravitySimulator) Visualize(buf []gfx.PosTexColorVertex) {
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

func (r *GravitySimulator) Size() (live, cap int) {
	return r.live, r.cap
}
