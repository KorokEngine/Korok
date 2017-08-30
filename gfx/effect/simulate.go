package effect

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok/math"
)

/**
	需要解决两个问题：
	1. 粒子系统的着色器能够直接使用Mesh着色器？如果不能，需要考虑如何切换着色器的问题，由于渲染需要从后往前，所以
		必须统一到渲染系统里面渲染，不能单独写成一个新的渲染系统。

		- 粒子系统渲染使用自定义的粒子着色器，不使用针对精灵的Mesh着色器
	2. 现在已经勾画出了整个业务流程，考虑如何将之动态化，这样可以通过配置脚本来驱动粒子系统！

		- 暂时无法实现动态化，构建大而全的系统会牺牲一部分的性能和内存
 */

//// 算子定义操作：赋值， 相加，积分

type Simulator interface {
	Initialize()

	Simulate(dt float32)

	Visualize()
}

type channel_f32 []float32

func (ch channel_f32) SetConst(n int32, v float32) {
	for i := int32(0); i < n; i++ {
		ch[i] = v
	}
}

func (ch channel_f32) SetRandom(n int32, r Range) {
	for i := int32(0); i < n; i++ {
		ch[i] = math.Random(r[0], r[1])
	}
}

func (ch channel_f32) Add(n int32, v float32) {
	for i := int32(0); i < n; i ++ {
		ch[i] += v
	}
}

func (ch channel_f32) Sub(n int32, v float32) {
	for i := int32(0); i < n; i ++ {
		ch[i] -= v
	}
}

func (ch channel_f32) Mul(n int32, v float32) {
	for i := int32(0); i < n; i++ {
		ch[i] *= v
	}
}

func (ch channel_f32) Integrate(n int32, ch1 channel_f32, dt float32) {
	for i := int32(0); i < n; i++ {
		ch[i] += ch1[i] * dt
	}
}

type channel_v2 []mgl32.Vec2

func (ch channel_v2) SetConst(n int32, x, y float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0], ch[i][1] = x, y
	}
}

func (ch channel_v2) SetRandom(n int32, xlow, xhigh float32, ylow, yhigh float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0], ch[i][1] = xlow, ylow
	}
}

func (ch channel_v2) Add(n int32, x, y float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0] += x
		ch[i][1] += y
	}
}

func (ch channel_v2) Integrate(n int32, ch1 channel_v2, dt float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0] += ch1[i][0] * dt
		ch[i][1] += ch1[i][1] * dt
	}
}

// ch = normal * m, normal = normal_vector(x, y), m = magnitude
func (ch channel_v2) radialIntegrate(n int32, xy channel_v2, m channel_f32, dt float32) {
	for i := int32(0); i < n; i++ {
		normal := [2]float32{0, 0}
		p := xy[i]

		if p[0] != 0 || p[1] != 0 {
			normalize(p[0], p[1], &normal)
		}

		ch[i][0] += m[i] * normal[0] * dt
		ch[i][1] += m[i] * normal[1] * dt
	}
}

// ch = tangent * m, normal = normal_vector(y, x), m = magnitude
func (ch channel_v2) tangentIntegrate(n int32, xy channel_v2, m channel_f32, dt float32) {
	for i := int32(0); i < n; i++ {
		tangent := [2]float32{0, 0}
		p := xy[i]

		if p[0] != 0 || p[1] != 0 {
			normalize(p[1], p[0], &tangent)
		}

		ch[i][0] += m[i] * tangent[0] * dt
		ch[i][1] += m[i] * tangent[1] * dt
	}
}

type channel_v3 []mgl32.Vec3

// maybe only color will use it
type channel_v4 []mgl32.Mat4

func (ch channel_v4) SetConst(n int32, x, y, z, w float32) {
	for i := int32(0); i < n; i ++ {
		ch[i][0], ch[i][1], ch[i][2], ch[i][3] = x, y, z, w
	}
}

func (ch channel_v4) SetRandom(n int32, x, y, z, w Range) {
	for i := int32(0); i < n; i ++ {
		ch[i][0] = math.Random(x[0], x[1])
		ch[i][1] = math.Random(y[0], y[1])
		ch[i][2] = math.Random(z[0], z[1])
		ch[i][3] = math.Random(w[0], w[1])
	}
}

func (ch channel_v4) Add(n int32, x, y, z, w float32) {
	for i := int32(0); i < n; i ++ {
		ch[i][0] += x
		ch[i][1] += y
		ch[i][2] += z
		ch[i][2] += w
	}
}

func (ch channel_v4) Sub(n int32, x, y, z, w float32) {
	for i := int32(0); i < n; i ++ {
		ch[i][0] -= x
		ch[i][1] -= y
		ch[i][2] -= z
		ch[i][2] -= w
	}
}

func (ch channel_v4) Integrate(n int32, d channel_v4, dt float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0] += d[i][0] * dt
		ch[i][1] += d[i][1] * dt
		ch[i][2] += d[i][2] * dt
		ch[i][3] += d[i][3] * dt
	}
}

func normalize(x, y float32, n *[2]float32)  {
	div := math.InvSqrt(x * x + y * y)
	n[0] = x * div
	n[1] = x * div
}

func random(low, high float32) float32 {
	return math.Random(low, high)
}

func max(a ,b float32) float32{
	if a > b {
		return a
	}
	return b
}

//// 算子局限于目前的支持的方法数量，唯一限制!
