package effect

import (
	"testing"
	"korok.io/korok/math/f32"
)

func TestChannel_f32(t *testing.T) {
	p := Channel_f32(make([]float32, 1024))
	n  := 10

	p.Add(int32(n), 8)

	for i := 0; i < n; i++ {
		if p[i] != 8 {
			t.Error("err: chanf32.add")
		}
	}

	p.Sub(int32(n), 4)

	for i := 0; i < n; i++ {
		if p[i] != 4 {
			t.Error("err: chanf32.add")
		}
	}

	p.Mul(10, 2)
	for i := 0; i < n; i++ {
		if p[i] != 8 {
			t.Error("err: chanf32.Mul")
		}
	}

	v := Channel_f32(make([]float32, 1024))
	v.SetConst(int32(n), 1.0/60)

	p.Integrate(10, v, 2)

	// p' = p + v*dt
	pp := float32(8 + 1.0/60 * 2)
	for i := 0; i < n; i++ {
		if p[i] != pp {
			t.Error("err: chanf32.integrate")
		}
	}
}

func TestChannel_v2(t *testing.T) {
	p := Channel_v2(make([]f32.Vec2, 1024))
	n := 10

	p.SetConst(int32(n), 4, 8)
	for i := 0; i < n; i++ {
		if p[i][0] != 4 || p[i][1] != 8 {
			t.Error("err: chanv2.SetConst")
		}
	}

	p.Add(int32(n), 4, 8)
	for i := 0; i < n; i++ {
		if p[i][0] != 8 || p[i][1] != 16 {
			t.Error("err: chanv2.Add")
		}
	}

	v := Channel_v2(make([]f32.Vec2, 1024))
	v.SetConst(int32(n), 16, 32)
	p.Integrate(int32(n), v, 1.0/60)

	p1, p2 := float32(8 + 16 * 1.0/60), float32(16 + 32 * 1.0/60)
	for i := 0; i < n; i++ {
		if p[i][0] != p1 || p[i][1] != p2 {
			t.Error("err: chanv2.Integrate")
		}
	}

}
