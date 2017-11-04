package game

import (
	timer "time"

	"github.com/go-gl/glfw/v3.2/glfw"

	"korok/gfx"
	"korok/assets"
	"korok/physics"
	"korok/anim"
	"korok/audio"
	"korok/engi"
)

type Table interface {}

type DB struct {
	EntityM *engi.EntityManager
	Tables []Table
}

func (db DB) Add(t interface{}) {
	db.Tables = append(db.Tables, t)
}

// 统一管理游戏各个子系统的创建和销毁的地方
var G *Game
var scenes = make(map[string]Scene)
var current Scene

type Game struct {
	DB

	*gfx.RenderSystem
	*physics.CollisionSystem
	*anim.AnimationSystem
	*audio.AudioSystem

	State int
}
/// window callback
func (g *Game) OnCreate() {
	g.Create()
}

func (g *Game) OnLoop() {
	g.Update()
}

func (g *Game) OnDestroy() {
	g.Destroy()
}

func AddScene(scene Scene)  {
	scenes[scene.Name()] = scene
	current = scene
}

// init subsystem
func (g *Game) Create() {
	g.Init()
}

// destroy subsystem
func (g *Game) Destroy() {
	g.RenderSystem.Destroy()
}

func (g *Game) Init()  {
	// assets
	assets.LoadShader()

	// render
	g.RenderSystem = gfx.NewRenderSystem()

	// physics system
	g.CollisionSystem = physics.NewCollisionSystem()

	// animation

	/// Customized scene
	if current != nil {
		current.Preload()
		current.Setup(g)
	}
}


func (g *Game) Input(dt float32)  {
	if current != nil {
	}
}

var previousTime float64
func (g *Game) Update()  {
	// update
	time :=  glfw.GetTime()
	elapsed := time - previousTime
	previousTime = time

	dt := float32(elapsed)
	if current != nil {
		current.Update(dt)
	}

	if dt < 0.0166 {
		timer.Sleep(timer.Duration(0.0166 - dt) * timer.Second)
	}

	//// simulation....

	/// 动画更新，骨骼数据
	///g.AnimationSystem.Update(dt)

	// g.CollisionSystem.Update(dt)

	// 粒子系统更新
	//g.ParticleSystem.Update(dt)

}

func (g *Game) Draw(dt float32)  {
}
