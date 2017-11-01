package particle

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok/gfx"
)

type FireCloud struct {
	cap int32

	// f32 channels
	life channel_f32
	size channel_f32

	// vector channels
	position channel_v2
	velocity channel_v2

	// color channel
	color channel_v4

	// live particle count
	live int32

	// vbo data <u, v, r, g, b, a, x, y, size, rot>
	vbo []float32
	ebo []int32

	// for render!
	id uint32
	gfx.Mesh
}

func (f *FireCloud) Initialize(cap int32) {
	f.cap = cap
	f.life = make([]float32, cap)
	f.size = make([]float32, cap)
	f.position = make([]mgl32.Vec2, cap)
	f.velocity = make([]mgl32.Vec2, cap)

	// vbo data <u,v, x,y, size, r, g, b, a>


	// init life 10 ,20
	f.life.SetRandom(cap, Range{0.5, 1})

	// init size
	f.size.SetRandom(cap, Range{5, 10})
	f.color.SetConst(cap, 1, 1, 1, 1)

	// init position = {50, 50} , velocity = {20, 20}
	//f.position.
	f.position.SetConst(cap, 50, 50)
	f.velocity.SetConst(cap, 20, 20)
}

func (f *FireCloud) Simulate(dt float32) {
	// spawn new particle
	if new := f.Rate(dt); new > 0 {
		f.NewParticle(new)
	}

	n := f.live

	// update old particle
	f.life.Sub(n, dt)

	// position integrate: p' = p + v * t
	f.position.Integrate(n, f.velocity, dt)

	// color
	f.color.Sub(n, 0, 0, 0, 2.5 * dt)

	// gc
	f.gc()
}

func (f *FireCloud) SwapErase(i, j int32) {
	f.life[i] = f.life[j]
	f.size[i] = f.size[j]
	f.position[i] = f.position[j]
	f.velocity[i] = f.velocity[j]
}

func (f *FireCloud) NewParticle(new int32) {
	start := f.live
	f.live += new

	for i := start; i < f.live; i++ {
		f.life[i] = random(0.5, 1)
		f.size[i] = random(5, 20)
		f.position[i] = mgl32.Vec2{50, 50}
		f.velocity[i] = mgl32.Vec2{20, 20}
	}
}

func (f *FireCloud) gc() {
	n := f.live
	// remove dead particle TODO!!!
	for i:= int32(0); i < n; i++{
		if f.life[i] <= 0 {
			// find last live
			j := f.live - 1
			for ; j > i && f.life[j] <= 0; j-- {
				f.live --
			}

			if j == i {
				f.live --
			} else {
				f.SwapErase(i, j)
			}
		}
	}
}

// 控制是否产生新的粒子, 返回粒子数量
func (f *FireCloud) Rate(dt float32) int32 {
	return 1
}

// 总觉的这里可以有很大的优化空间，尤其是顶点和纹理数据大量的重复
// write to a vertex buffer!!!
// <u, v, r, g, b, a, x, y, size, rot>
func (f *FireCloud) Visualize() {
	// compute vbo
	var vi int32
	for i := int32(0); i < f.live; i++ {
		vi += 10 * 4 * i

		h_size := f.size[i] / 2

		// bottom-left
		f.vbo[vi + 2] = f.position[i][0] - h_size
		f.vbo[vi + 3] = f.position[i][1] - h_size

		f.vbo[vi + 4] = f.color[i][0]
		f.vbo[vi + 5] = f.color[i][1]
		f.vbo[vi + 6] = f.color[i][2]
		f.vbo[vi + 7] = f.color[i][3]

		/*size and rotation TODO */

		// bottom-right
		vi += 10
		f.vbo[vi + 2] = f.position[i][0] + h_size
		f.vbo[vi + 3] = f.position[i][1] - h_size

		f.vbo[vi + 4] = f.color[i][0]
		f.vbo[vi + 5] = f.color[i][1]
		f.vbo[vi + 6] = f.color[i][2]
		f.vbo[vi + 7] = f.color[i][3]

		// top-right
		vi += 10
		f.vbo[vi + 2] = f.position[i][0] + h_size
		f.vbo[vi + 3] = f.position[i][1] + h_size

		f.vbo[vi + 4] = f.color[i][0]
		f.vbo[vi + 5] = f.color[i][1]
		f.vbo[vi + 6] = f.color[i][2]
		f.vbo[vi + 7] = f.color[i][3]

		// top-left
		vi += 10
		f.vbo[vi + 2] = f.position[i][0] - h_size
		f.vbo[vi + 3] = f.position[i][1] + h_size

		f.vbo[vi + 4] = f.color[i][0]
		f.vbo[vi + 5] = f.color[i][1]
		f.vbo[vi + 6] = f.color[i][2]
		f.vbo[vi + 7] = f.color[i][3]
	}

	// write mesh!
	f.Mesh.SetVertex(f.vbo)

	// upload
	f.Mesh.Setup()
}

func (f *FireCloud) WriteTexCoords() {
	for i := int32(0); i < f.live; i++ {
		vi := 4 * 10 * i

		// bottom-left
		f.vbo[vi + 0] = 0
		f.vbo[vi + 1] = 0

		// bottom-right
		vi += 10
		f.vbo[vi + 0] = 1
		f.vbo[vi + 1] = 0


		// top-right
		vi += 10
		f.vbo[vi + 0] = 1
		f.vbo[vi + 1] = 1

		// top-left
		vi += 10
		f.vbo[vi + 0] = 0
		f.vbo[vi + 1] = 1
	}
}

func (f *FireCloud) WriteIndex() {
	for i := int32(0); i < f.live; i++ {
		ei := i * 6
		bi := int32(i * 4)

		f.ebo[ei+0] = bi + 1
		f.ebo[ei+1] = bi + 2
		f.ebo[ei+2] = bi + 3

		f.ebo[ei+3] = bi + 0
		f.ebo[ei+4] = bi + 1
		f.ebo[ei+5] = bi + 3
	}
}


