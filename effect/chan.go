package effect

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/math"
)

// Compute Channel：
// Types：float32/Vec2/Vec4
// Methods: random/const/add/sub/mul/integrate
//
// Compute channel is an array of values that can be computed in a for-loop.
// It's data structure is cache-friendly.
type ChanType uint8
const (
	ChanF32 ChanType = iota
	ChanV2
	ChanV4
)

type Channel_f32 []float32

func (ch Channel_f32) SetConst(n int32, v float32) {
	for i := int32(0); i < n; i++ {
		ch[i] = v
	}
}

func (ch Channel_f32) SetRandom(n int32, v Var) {
	for i := int32(0); i < n; i++ {
		ch[i] = math.Random(v.Base, v.Base+v.Var)
	}
}

func (ch Channel_f32) Add(n int32, v float32) {
	for i := int32(0); i < n; i ++ {
		ch[i] += v
	}
}

func (ch Channel_f32) Sub(n int32, v float32) {
	for i := int32(0); i < n; i ++ {
		ch[i] -= v
	}
}

func (ch Channel_f32) Mul(n int32, v float32) {
	for i := int32(0); i < n; i++ {
		ch[i] *= v
	}
}

func (ch Channel_f32) Integrate(n int32, ch1 Channel_f32, dt float32) {
	for i := int32(0); i < n; i++ {
		ch[i] += ch1[i] * dt
	}
}

type Channel_v2 []f32.Vec2

func (ch Channel_v2) SetConst(n int32, x, y float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0], ch[i][1] = x, y
	}
}

func (ch Channel_v2) SetRandom(n int32, xlow, xhigh float32, ylow, yhigh float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0], ch[i][1] = xlow, ylow
	}
}

func (ch Channel_v2) Add(n int32, x, y float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0] += x
		ch[i][1] += y
	}
}

func (ch Channel_v2) Integrate(n int32, ch1 Channel_v2, dt float32) {
	for i := int32(0); i < n; i++ {
		ch[i][0] += ch1[i][0] * dt
		ch[i][1] += ch1[i][1] * dt
	}
}

// ch = normal * m, normal = normal_vector(x, y), m = magnitude
func (ch Channel_v2) radialIntegrate(n int32, xy Channel_v2, m Channel_f32, dt float32) {
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
func (ch Channel_v2) tangentIntegrate(n int32, xy Channel_v2, m Channel_f32, dt float32) {
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

// maybe only Color will use it
type Channel_v4 []f32.Vec4

func (ch Channel_v4) SetConst(n int32, x, y, z, w float32) {
	for i := int32(0); i < n; i ++ {
		ch[i][0], ch[i][1], ch[i][2], ch[i][3] = x, y, z, w
	}
}

func (ch Channel_v4) SetRandom(n int32, x, y, z, v [4]Var) {
	for i := int32(0); i < n; i ++ {
		ch[i][0] = math.Random(v[0].Base, v[0].Base+v[0].Var)
		ch[i][1] = math.Random(v[1].Base, v[1].Base+v[1].Var)
		ch[i][2] = math.Random(v[2].Base, v[2].Base+v[2].Var)
		ch[i][3] = math.Random(v[3].Base, v[3].Base+v[3].Var)
	}
}

func (ch Channel_v4) Add(n int32, x, y, z, w float32) {
	for i := int32(0); i < n; i ++ {
		ch[i][0] += x
		ch[i][1] += y
		ch[i][2] += z
		ch[i][3] += w
	}
}

func (ch Channel_v4) Sub(n int32, x, y, z, w float32) {
	for i := int32(0); i < n; i ++ {
		ch[i][0] -= x
		ch[i][1] -= y
		ch[i][2] -= z
		ch[i][3] -= w
	}
}

func (ch Channel_v4) Integrate(n int32, d Channel_v4, dt float32) {
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
