package effect

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok/math"
)

// 计算通道定义：
// 此处定了 3 中类型的计算通道：float32/Vec2/Vec4
// 定义了 10 个左右的基本算法: random/const/add/sub/mul/integrate

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

