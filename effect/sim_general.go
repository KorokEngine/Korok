package effect

import "github.com/go-gl/mathgl/mgl32"


/**
	显示一个粒子最终只需要：Position/Rotation/Scale

	但是对粒子进行建模需要更多的属性：
	1. 粒子发射频率
	2. 重力/重力模式
	3. 方向
	4. 速度/加速度
	5. 切线速度/微分
	6. 角速度
	7. 放射模式
	8. 开始放射半径/微分
	9. 结束放射半径/微分
	10. 旋转
	11. 共有的属性
		1. 生命/微分
		2. 初始/结束旋转/微分
		3. 初始/结束大小/微分
		4. 初始/结束颜色/微分
		5. 混合方程
		6. 纹理

	最好能够解决粒子建模的问题，而不是硬编码各种组合
*/

// name convention: r = red, d_r = derivative of r with respect to time
// integrate: r' = r + d_r * t
//

type EmitterMode int32

// 发射模式
const (
	Gravity EmitterMode = iota
	Radius
)

type Range [2]float32

// base value + var value
type Var struct {
	B, Var float32
}

type Config struct {
	//
	MaxParticle uint32

	Angle, D_Angle float32
	Duration float32

	// TODO blend function
	// start_r = r[0] + var_r[0] * random
	// end_r = r[1] + var_r[1] * random
	R, VAR_R Range
	G, VAR_G Range
	B, VAR_B Range
	A, VAR_A Range

	// size
	Size, D_Size Range

	// position
	Position, D_Position mgl32.Vec2

	Spin, D_Spin float32

	// mode
	Mode EmitterMode

	// life span
	Life, D_Life float32

	// emission rate = total_particle / life
	EmissionRate float32
}

type GravityConfig struct {
	Config

	// gravity
	Gravity mgl32.Vec2

	// speed and d
	Speed, D_Speed float32

	// Radial acceleration
	RadialAccel, D_RadialAccel float32

	// tangent acceleration
	TangentialAccel, D_TangentialAccel float32

	RotationIsDir bool
}

type RadiusConfig struct {
	Config

	// min, max radius and d
	Radius, D_Radius Range

	//
	RotationPerSecond, D_RotationPerSecond float32
}

type GravitySimulator struct {
	// config
	C *GravityConfig

	// live particles
	live int32

	//	time to live
	life float32

	// color channel - if not used, it's a big waste
	r, d_r channel_f32
	g, d_g channel_f32
	b, d_b channel_f32
	a, d_a channel_f32

	// size channel
	size, d_size channel_f32

	// angle channel
	angle, d_angle channel_f32

	// gravity channel
	xy, d_xy channel_v2

	// acceleration include radial, tangent, gravity
	dd_radial, dd_tangent channel_f32
	g_x, g_y float32
}

// alloc a big block ?
func (g *GravitySimulator) Initialize(cap int32, c *Config) {
	// memory
	g.r, g.d_r = make([]float32, cap), make([]float32, cap)
	g.g, g.d_g = make([]float32, cap), make([]float32, cap)
	g.b, g.d_b = make([]float32, cap), make([]float32, cap)
	g.a, g.d_a = make([]float32, cap), make([]float32, cap)

	g.size, g.d_size = make([]float32, cap), make([]float32, cap)
	g.angle, g.d_angle = make([]float32, cap), make([]float32, cap)

	g.xy, g.d_xy = make([]mgl32.Vec2, cap), make([]mgl32.Vec2, cap)
	g.dd_radial = make([]float32,cap)
	g.dd_tangent = make([]float32, cap)

	// init color ! TODO 颜色初始化错误！
	g.r.SetRandom(cap, c.R); g.d_r.SetRandom(cap, c.VAR_R)
	g.g.SetRandom(cap, c.G); g.d_g.SetRandom(cap, c.VAR_G)
	g.b.SetRandom(cap, c.B); g.d_b.SetRandom(cap, c.VAR_B)
	g.a.SetRandom(cap, c.A); g.d_a.SetRandom(cap, c.VAR_A)

	// size
	g.size.SetRandom(cap, c.Size)

	// Haha, 其实这里面的初始化方法我都没有看的太明白...
}

func (g *GravitySimulator) Simulate(dt float32) {
	n := g.live;

	// color
	g.r.Integrate(n, g.d_r, dt)
	g.g.Integrate(n, g.d_g, dt)
	g.b.Integrate(n, g.d_b, dt)
	g.a.Integrate(n, g.d_a, dt)

	// size
	g.size.Integrate(n, g.d_size, dt)

	// angle
	g.angle.Integrate(n, g.d_angle, dt)

	// gravity model integrate
	g.d_xy.radialIntegrate(n, g.xy, g.dd_radial, dt)
	g.d_xy.tangentIntegrate(n, g.xy, g.dd_tangent, dt)
	g.d_xy.Add(n, g.g_x * dt, g.g_y * dt)

	// xy' = xy + v * t
	g.xy.Integrate(n, g.d_xy, dt)
}


// 随机值算方式：f(v) = base + var * random
// random: [0, 1]
// 变化计算方式：start = f(v); end = f(v)
// delta := (end - start) / life
// 算法：分别随机出初始值和结束值，然后除以lifetime得到变化量
// 更新：利用变化量积分
func (g *GravitySimulator) AddParticle(n int32) {
	start := g.live
	end   := g.live + n

	// life
	for i := start; i < end; i++ {
		// TODO 分析life的变化率
		// g.life[i] = random()
	}

	// position
	for i := start; i < end; i++ {

	}

	// color - RGBA
	for i := start; i < end; i++ {

	}


}

func NewGeneralSimulator() Simulator {
	return nil
}