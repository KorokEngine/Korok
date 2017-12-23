package effect

import (
	"testing"
	"github.com/go-gl/mathgl/mgl32"
)

func TestFireSimulator(t *testing.T) {
	fire := NewFireSimulator(1024)
	fire.Initialize()

	for i := 0; i < 10; i++ {
		fire.Simulate(1.0/60)
	}


}

func TestRadiusSimulator(t *testing.T) {
	cfg := &RadiusConfig{
		Config:Config {
			Max:128,
			Duration:0,
			Life:Var{1.1, 0.4},
			Size:Range{Var{10 ,20}, Var{10, 20}},
			X:Var{0, 0}, Y:Var{0, 0},
		},
		Radius:Range{Var{50, 10}, Var{50, 10}},
		Angle:Var{0, 3.14},
	}
	radius := NewRadiusSimulator(cfg)
	radius.Initialize()

	radius.newParticle(2)

	t.Log("life:", radius.life[:2])
	t.Log("radius:", radius.radius[:2])
	t.Log("angle:", radius.angle[:2])

	for i := 0; i < 10; i++ {
		radius.Simulate(1.0/60)
	}
}

func TestGravitySimulator(t *testing.T) {
	cfg := &GravityConfig{
		Config:Config {
			Max:128,
			Duration:0,
			Life:Var{1.1, 0.4},
			Size:Range{Var{10 ,5}, Var{20, 5}},
			X:Var{0, 0}, Y:Var{0, 0},
			A: Range{Var{1, 0}, Var{0, 0}},
		},
		Velocity: [2]Var{{10, 0}, {10, 0}},
		Gravity:mgl32.Vec2{0, 90},
	}
	gravity := NewGravitySimulator(cfg)
	gravity.Initialize()

	gravity.newParticle(2)
	t.Log("life:", gravity.life[:2])
	t.Log("velocity:", gravity.velocity[:2])
	t.Log("pose:", gravity.pose[:2])
}
