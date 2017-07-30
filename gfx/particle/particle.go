package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v3.2-core/gl"

	"korok/gfx"
)

/**
	实现粒子效果

	内存中的例子数量是固定的，所有的例子放在一个数组中，
	当粒子的生命殆尽时在数组中标记一次，当新的粒子需要
	产生的时候，从数组中找到一个死亡的例子重新赋予生命。

	每次更新的时候都会新加几个粒子出来。粒子都有生命所
	以一会便会自然死亡。在绘制粒子的时候用到了BLEND函数
	这个函数给粒子添加了Glow的效果。

	粒子使用了新的着色器，着色器针对粒子有优化...
 */

type Particle struct {
	Position, Velocity mgl32.Vec2
	Color mgl32.Vec4
	Life float32
}

type ParticleGenerator struct {
	Particles []*Particle 	// 粒子数组
	amount int 		// 粒子数量
	shader *gfx.Shader 		// 粒子着色器
	texture *gfx.Texture2D 	// 粒子纹理
	vao  uint32 		// VAO对象
	lastUsedParticle int
}

func (g *ParticleGenerator) init()  {
	var vertice = []float32{
		0.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
	}
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertice)*4, gl.Ptr(vertice), gl.STATIC_DRAW)

	g.shader.Use()

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	// generator particles
	g.Particles = make([]*Particle, g.amount)
	for i := 0; i < g.amount; i++ {
		g.Particles[i] = &Particle{}
	}

	// important!
	g.vao = vao
}

func (g *ParticleGenerator) Update(dt float32, new int, p, offset, velocity mgl32.Vec2)  {
	// Add new particles
	for i := 0; i < new; i++ {
		unused := g.firstUnusedParticle()
		g.respawnParticle(g.Particles[unused], p, offset, velocity)
	}
	// update all particles
	for i := 0; i < g.amount; i++ {
		p := g.Particles[i]
		p.Life -= dt
		if p.Life > 0 {
			p.Position[0] -= p.Velocity[0] * dt
			p.Position[1] -= p.Velocity[1] * dt

			// RGB - A
			p.Color[3] -= dt*2.5
		}
	}
}

func (g *ParticleGenerator) Draw()  {
	// TODO 理解含义
	// Use additive blending to give it a 'Glow' effect
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	g.shader.Use()

	for _, v := range g.Particles {
		if v.Life > 0 {
			g.shader.SetVector2f("offset\x00", v.Position[0], v.Position[1])
			g.shader.SetVector4f("color\x00", v.Color[0], v.Color[1], v.Color[2], v.Color[3])
			g.texture.Bind()

			gl.BindVertexArray(g.vao)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
			gl.BindVertexArray(0)
		}
	}

	// reset to default!
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.BLEND)
}

// 这里面的算法有待研究
func (g *ParticleGenerator) firstUnusedParticle() int {
	// First search from last used particle, this will usually return almost instantly
	for i := g.lastUsedParticle; i < g.amount; i++ {
		if g.Particles[i].Life <= 0 {
			g.lastUsedParticle = i
			return i
		}
	}
	// Otherwise, do a linear search
	for i:= 0; i < g.lastUsedParticle; i++ {
		if g.Particles[i].Life <= 0 {
			g.lastUsedParticle = i
			return i
		}
	}
	// All particles ar taken, override the first one
	g.lastUsedParticle = 0
	return 0
}

func (g *ParticleGenerator) respawnParticle(particle *Particle, position, offset, velocity mgl32.Vec2)  {
	// random = ?
	// color = ?
	particle.Position = mgl32.Vec2{position[0] + offset[0], position[1] + offset[1]}
	particle.Color = mgl32.Vec4{1, 1, 0, 1}
	particle.Life = .5
	particle.Velocity = mgl32.Vec2{velocity[0]*0.1, velocity[1] * 0.1}
}

func NewParticleGenerator(shader *gfx.Shader, tex *gfx.Texture2D, amount int) (*ParticleGenerator) {
	g := &ParticleGenerator{
		shader: shader,
		texture: tex,
		amount: amount,
	}
	g.init()
	return g
}





