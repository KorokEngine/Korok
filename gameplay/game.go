package gameplay

import (
	"os"
	"fmt"
	timer "time"


	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"korok/gfx"
	"korok/gfx/text"
	"korok/assets"
	"korok/space"
	"korok/physics"
	"korok/anim/spine"
	"korok/anim"
	"korok/gfx/effect"
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
	*effect.ParticleSystem

	State int

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

	comp *gfx.RenderComp

	// type-render
	*gfx.MeshRender

	// render-context
	*gfx.Mesh

	*gfx.RenderContext
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

	// GLShader
	// init MVP
	shader := assets.GetShader("dft")
	//if shader != nil {
	//	shader.Use()
	//	//  ---- Vertex GLShader
	//	// projection
	//	p := mgl32.Ortho2D(0, 480, 0, 320)
	//	shader.SetMatrix4("projection\x00", p)
	//
	//	// model
	//	model := mgl32.Ident4()
	//	shader.SetMatrix4("model\x00", model)
	//
	//	// ---- Fragment GLShader
	//	shader.SetInteger("tex\x00", 0)
	//	gl.BindFragDataLocation(shader.Program, 0, gl.Str("outputColor\x00"))
	//
	//	log.Println("---> init shader ok ...")
	//}

	// mesh shader

	// register render-state

	// render
	g.RenderSystem = gfx.NewRenderSystem()

	// entity graph
	g.NodeSystem = space.NewNodeSystem()

	// physics system
	g.CollisionSystem = physics.NewCollisionSystem()

	// animation
	g.AnimationSystem = anim.NewAnimationSystem()

	// particle system
	g.ParticleSystem = effect.NewParticleSystem()

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

		//// Mesh Render

		assets.LoadTexture("assets/ball.png")

		b := assets.GetTexture("assets/ball.png")

		fmt.Println("texture", b)
		sub := b.Sub(50, 50, 300, 300)

		fmt.Println("sub texture", sub)


		//id := ecs.Create()
		//comp := g.RenderSystem.NewRenderComp(id.Index())
		////comp.SetTexture(b)
		//comp.SetSubTexture(sub)
		//comp.SetPosition(mgl32.Vec2{100, 50})

		// comp test!!
		comp := new(gfx.RenderComp)
		//comp.SetTexture(b)
		comp.SetPosition(mgl32.Vec2{50, 50})

		text := text.NewText(assets.GetFont("assets/font/font.json"))
		text.SetString("Hello")

		//m := text.Mesh()
		//
		//comp.SetMesh(m, func() (vao, vbo, ebo uint32) {
		//	return m.Handle()
		//})
		//
		//g.comp = comp
		assets.LoadTexture("assets/ball.png")
		tex := assets.GetTexture("assets/ball.png")

		g.MeshRender = gfx.NewMeshRender(*shader)
		g.Mesh = gfx.NewQuadMesh(tex)
		g.Mesh.Setup()

		g.RenderContext = gfx.NewRenderContext(*shader, tex.Id)
	}

	// sprite
	sprite := newSprite(50, 50)
	g.comp = sprite
}

func newSprite(x, y float32) *gfx.RenderComp {
	//assets.LoadTexture("assets/ball.png")
	//tex := assets.GetTexture("assets/ball.png")

	comp := new(gfx.RenderComp)
	comp.Type = gfx.RenderType_Mesh
	//comp.SetTexture(tex)
	//comp.SetPosition(mgl32.Vec2{x, y})
	//
	//comp.Sort.SetShader(1)
	//comp.Sort.SetTexture(1)
	//comp.Sort.SetBlendFunc(1)

	return comp
}

func newBatchCommand() {

	//assets.LoadTexture("assets/ball.png")
	//tex := assets.GetTexture("assets/ball.png")
	//tex.Width, tex.Height = 50, 50
	//
	//m0 := gfx.NewIndexedMesh(tex); m0.SRT(mgl32.Vec2{50, 50}, 0, mgl32.Vec2{0, 0})
	//m1 := gfx.NewIndexedMesh(tex); m1.SRT(mgl32.Vec2{110, 50}, 0, mgl32.Vec2{0, 0})
	//m2 := gfx.NewIndexedMesh(tex); m2.SRT(mgl32.Vec2{170, 50}, 0, mgl32.Vec2{0, 0})
	//
	//bs := gfx.NewBatchSystem()
	//
	//b := bs.NewBatch(tex.Id)
	//
	//b.AddVertex(m0.Vertex(), nil)
	//b.AddVertex(m1.Vertex(), nil)
	//b.AddVertex(m2.Vertex(), nil)
	//
	//c := b.Commit()
	//// only shader changes..
	//c.Key.SetShader(2)
	//c.Key.SetBlendFunc(1)
	//c.Key.SetTexture(1)
	//
	//return &c
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

	// 渲染系统更新
	//g.RenderSystem.Update(dt)

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
	//g.death.Apply(g.skeleton, float32(time), true)

	// 更新骨骼坐标
	//g.skeleton.UpdateWorldTransform()

	//
	//g.skRender.Draw(g.skeleton)

	//g.RenderSystem.Update(dt)

	g.MeshRender.Draw(g.Mesh, mgl32.Vec2{100, 50}, mgl32.Vec2{1, 1}, 0)
	// g.RenderContext.Draw()
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


