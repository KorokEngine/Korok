package particle

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v3.2-core/gl"
)

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

type ParticleComp struct {
	C *Config
	Simulator
}

type ParticleSystem struct {
	comps []ParticleComp
}

func NewParticleSystem() *ParticleSystem {
	return new(ParticleSystem)
}

func (p *ParticleSystem) Update(dt float32) {
	for _, comp := range p.comps {
		comp.Simulator.Simulate(dt)
	}
}

func (p *ParticleSystem) NewParticleComp(id uint32, c *Config) *ParticleComp {
	comp := new(ParticleComp)
	comp.C = c
	comp.Simulator = nil // TODO!
	return nil
}

var vertex = []float32{
		0,  1,  0.0, 1.0,
		1,  0,  1.0, 0.0,
		0,  0,  0.0, 0.0,
		1,  1,  1.0, 1.0,
	}

type ParticleGroup struct {
	vao, vbo, ebo uint32

	vertex []float32
	index  []int32
	ii     []int32

	count int32
}

func NewParticleGroup(){
	b := new(ParticleGroup)

	b.vertex = []float32{
	// <x, y, size, rot> <r, g, b, a>  <index>
		50,  20,  40, 0,  0,  0,  0, 1, 0,
		50,  20,  40, 0,  0,  0,  0, 1, 1,
		50,  20,  40, 0,  0,  0,  0, 1, 2,
		50,  20,  40, 0,  0,  0,  0, 1, 3,
	}
	b.ii = []int32{

	}
	b.index = []int32{
		0, 1, 2,
		0, 3, 1,
	}
	b.count = 6
	//b.tex = tex

	gl.GenVertexArrays(1, &b.vao)
	gl.GenBuffers(1, &b.vbo)
	gl.GenBuffers(1, &b.ebo)

	gl.BindVertexArray(b.vao)
	// vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(b.vertex)*4, gl.Ptr(b.vertex), gl.STATIC_DRAW)

	// ebo optional
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(b.index)*4, gl.Ptr(b.index), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 9*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 9*4, gl.PtrOffset(16))


	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 9*4, gl.PtrOffset(32))

	gl.BindVertexArray(0)
}

