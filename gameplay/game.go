package gameplay

import (
	"korok/gfx"
	"korok/gfx/text"
	"korok/assets"
	"korok/space"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"korok/physics"
	"os"
	"fmt"
	"korok/anim/spine"
	timer "time"
	"korok/anim"
)

// 统一管理游戏各个子系统的创建和销毁的地方
var G *Game
var scenes = make(map[string]Scene)
var current Scene

type Game struct {
	*gfx.RenderSystem
	*space.NodeSystem
	*physics.CollisionSystem
	*anim.AnimationSystem
	tRender *text.Renderer

	State int

	Renderer *gfx.SpriteRender
	Player *gfx.RenderComp
	Ball *gfx.RenderComp
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
	//
	g.Init()
}

// destroy subsystem
func (g *Game) Destroy() {
	g.RenderSystem.Destroy()
	g.NodeSystem.Destroy()
}

func (g *Game) Init()  {
	// assets
	assets.LoadShader()

	// Shader
	// init MVP
	shader := assets.GetShader("dft")
	if shader != nil {
		shader.Use()
		//  ---- Vertex Shader
		// projection
		p := mgl32.Ortho2D(0, 480, 0, 320)
		shader.SetMatrix4("projection\x00", p)

		// model
		model := mgl32.Ident4()
		shader.SetMatrix4("model\x00", model)

		// ---- Fragment Shader
		shader.SetInteger("tex\x00", 0)
		gl.BindFragDataLocation(shader.Program, 0, gl.Str("outputColor\x00"))

		log.Println("---> init shader ok ...")
	}

	tShader := assets.GetShader("text")
	if tShader != nil {
		tShader.Use()

		p := mgl32.Ortho2D(0, 480, 0, 320)

		// vertex
		tShader.SetMatrix4("projection\x00", p)
		tShader.SetVector3f("model\x00", 50, 50, 10)

		// fragment
		tShader.SetInteger("text\x00", 0)
		gl.BindFragDataLocation(tShader.Program, 0, gl.Str("color\x00"))
	}

	g.Renderer = gfx.NewSpriteRender(shader)

	// text render
	g.tRender = text.NewTextRenderer(tShader)

	// render
	g.RenderSystem = gfx.NewRenderSystem(assets.GetShader("dft"))

	// entity graph
	g.NodeSystem = space.NewNodeSystem()

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
		label.Scale = 10
		label.Position = mgl32.Vec2{50, 50}
		label.Color = mgl32.Vec3{0, 0, 0}
		g.Label = label
	}

	// animation
	////// 1. Loading Atlas
	{
		atlasFile, err := os.Open("assets/spine/alien.atlas")
		if err != nil {
			fmt.Println(err)
			return
		}

		assets.LoadTexture("assets/spine/alien.png")
		atlas, err := spine.NewAtlas(atlasFile, &TextureLoader{})
		if err != nil {
			fmt.Println(err)
			return
		}

		////// 2. create attachment loader
		loader := &spine.AtlasAttachmentLoader{atlas}

		////// 3. Loading Skeleton
		spineFile, err := os.Open("assets/spine/alien.json")
		if err != nil {
			fmt.Println(err)
			return
		}

		spineData, err := spine.New(spineFile, 0.6, loader)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(spineData)
		skeleton := spine.NewSkeleton(spineData)

		fmt.Println("bones:", len(skeleton.Bones))
		fmt.Println("slots:", len(skeleton.Slots))

		////// 4. start animation
		skeleton.SetToSetupPose()
		skeleton.X = 250
		skeleton.Y = 10

		g.run = skeleton.FindAnimation("run")
		g.jump = skeleton.FindAnimation("jump")
		g.hit = skeleton.FindAnimation("hit")
		g.death = skeleton.FindAnimation("death")

		g.skeleton = skeleton

		/// 4.1 init renderComp
		g.skRender = anim.NewSkeletonRender(shader)
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
	g.AnimationSystem.Update(dt)


	// g.CollisionSystem.Update(dt)

	/// sync
	//////// *20
	//for b := g.B2World.GetBodyList(); b != nil; b = b.Next {
	//		if b.UserData != nil {
	//			id, ok := b.UserData.(uint32)
	//			if ok {
	//				p := b.GetPosition()
	//				comp := g.RenderSystem.GetComp(id)
	//				comp.Model = mgl32.Translate3D(p.X * 20, p.Y * 20 - comp.Height / 2, 0).Mul4(mgl32.Scale3D(comp.Width, comp.Height, 1))
	//			}
	//
	//			//fmt.Print("position:", b.GetPosition())
	//		}
	//}
	//fmt.Println("contact count:", g.B2World.GetContactCount())
	//g.RenderSystem.Update(dt)

	//comp1 := g.RenderSystem.GetComp(1)
	//m := mgl32.Translate3D(100, 50, 0).Mul4(mgl32.Scale3D(1, 1, 1))

	//g.tRender.RenderText(g.Label)

	// 作用动画
	g.death.Apply(g.skeleton, float32(time), true)

	// 更新骨骼坐标
	g.skeleton.UpdateWorldTransform()

	//
	g.skRender.Draw(g.skeleton)

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


