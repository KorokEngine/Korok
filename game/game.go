package game

import (
	timer "time"


	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx"
	"korok/gfx/text"
	"korok/assets"
	"korok/physics"
	"korok/anim/spine"
	"korok/anim"
	"korok/audio"
)

// 统一管理游戏各个子系统的创建和销毁的地方
var G *Game
var scenes = make(map[string]Scene)
var current Scene

type Game struct {
	*gfx.RenderSystem
	*physics.CollisionSystem
	*anim.AnimationSystem
	*audio.AudioSystem

	State int

	Label *text.LabelComp

	Font uint32

	skeleton *spine.Skeleton
	run *spine.Animation
	jump *spine.Animation
	hit  *spine.Animation
	death *spine.Animation
	skRender *anim.SkeletonRender
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
	g.AnimationSystem = anim.NewAnimationSystem()

	/// Customized scene
	if current != nil {
		current.Preload()
		current.Setup(g)
	}

	// init text
	{
		assets.LoadFont("assets/font/font.png", "assets/font/font.json")
		font := assets.GetFont("assets/font/font.json")
		g.Font = font.Texture

		label := text.NewText(font)
		label.SetString("Hello, %s", "Korok!")
		//label.Scale = 10
		//label.Position = mgl32.Vec2{50, 50}
		label.Color = mgl32.Vec3{0, 0, 0}
		g.Label = label
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


type TextureLoader struct {
}

func (*TextureLoader) Load(page *spine.AtlasPage) error {
	return nil
}

func (*TextureLoader) Unload(page *spine.AtlasPage) error {
	return nil
}


