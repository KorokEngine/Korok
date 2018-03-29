package game

import (
	"korok.io/korok/engi"
	"korok.io/korok/math"
	"korok.io/korok/gfx"
	"korok.io/korok/effect"
	"korok.io/korok/anim"
	"korok.io/korok/asset"
	"korok.io/korok/hid/input"
	"korok.io/korok/gfx/dbg"
	"korok.io/korok/gui"

	"log"
	"reflect"
	"fmt"
)

const (
	MaxScriptSize = 1024

	MaxSpriteSize = 64 << 10
	MaxTransformSize = 64 << 10
	MaxTextSize = 64 << 10
	MaxMeshSize = 64 << 10

	MaxParticleSize = 1024
)


type Options struct {
	W, H int
}

type Table interface{}

type DB struct {
	EntityM *engi.EntityManager
	Tables  []interface{}
}

// 统一管理游戏各个子系统的创建和销毁的地方
var G *Game

type Game struct {
	Options; FPS; DB

	// scene manager
	SceneManager

	// system
	*gfx.RenderSystem
	*gui.UISystem
	*input.InputSystem
	*effect.ParticleSimulateSystem
	*ScriptSystem
	*anim.AnimationSystem
}

func (g *Game) Camera() *gfx.Camera {
	return &g.RenderSystem.MainCamera
}

/// window callback
func (g *Game) OnCreate(ratio float32) {
	g.Create(ratio)
}

func (g *Game) OnLoop() {
	g.Update()
}

func (g *Game) OnDestroy() {
	g.Destroy()
}

/// input callback
func (g *Game) OnKeyEvent(key int, pressed bool) {
	g.InputSystem.SetKeyEvent(key, pressed)
}

func (g *Game) OnPointEvent(key int, pressed bool, x, y float32) {
	g.InputSystem.SetPointerEvent(key, pressed, x, y)
}

// init subsystem
func (g *Game) Create(ratio float32) {
	gfx.Init(ratio)

	// render system
	rs := gfx.NewRenderSystem()
	g.RenderSystem = rs

	min, max := -math.MaxFloat32, math.MaxFloat32
	// setup camera
	c := &rs.MainCamera
	c.SetViewPort(float32(g.Options.W), float32(g.Options.H))
	c.SetBound(min, max, max, min)

	//
	// set table
	rs.RequireTable(g.DB.Tables)
	// set render
	var vertex, color string

	vertex, color = asset.Shader.GetShaderStr("batch")
	batchRender := gfx.NewBatchRender(vertex, color)
	rs.RegisterRender(gfx.RenderType(0), batchRender)

	vertex, color = asset.Shader.GetShaderStr("mesh")
	meshRender := gfx.NewMeshRender(vertex, color)
	rs.RegisterRender(gfx.RenderType(1), meshRender)

	log.Println("LoadBitmap Render:", len(rs.RenderList))
	for i, v := range rs.RenderList {
		log.Println(i, " render - ", reflect.TypeOf(v))
	}

	// set feature
	srf := &gfx.SpriteRenderFeature{}
	srf.Register(rs)
	mrf := &gfx.MeshRenderFeature{}
	mrf.Register(rs)
	trf := &gfx.TextRenderFeature{}
	trf.Register(rs)

	log.Println("LoadBitmap Feature:", len(rs.FeatureList))
	for i, v := range rs.FeatureList {
		log.Println(i, " feature - ", reflect.TypeOf(v))
	}

	// gui system
	g.UISystem = gui.NewUISystem(meshRender)
	g.UISystem.RegisterContext(gui.DefaultContext())
	gui.SetScreenSize(float32(g.W), float32(g.H))

	/// init debug
	dbg.Init()

	/// input system
	g.InputSystem = input.NewInputSystem()

	/// particle-simulation system
	pss := effect.NewSimulationSystem()
	pss.RequireTable(g.DB.Tables)
	g.ParticleSimulateSystem = pss
	// set feature
	prf := &effect.ParticleRenderFeature{}
	prf.Register(rs)

	/// script system
	g.ScriptSystem = NewScriptSystem()
	g.ScriptSystem.RequireTable(g.DB.Tables)

	/// Tex2D animation system
	g.AnimationSystem = anim.NewAnimationSystem()
	g.AnimationSystem.RequireTable(g.DB.Tables)

	/// setup scene manager
	g.SceneManager.Setup(g)
}

// destroy subsystem
func (g *Game) Destroy() {
	g.RenderSystem.Destroy()
}

func (g *Game) Init(op Options) {
	g.Options = op
	g.loadTables()
}

func (g *Game) loadTables() {
	g.DB.EntityM = engi.NewEntityManager()

	// init tables
	scriptTable := NewScriptTable(MaxScriptSize)
	tagTable := &TagTable{}

	g.DB.Tables = append(g.DB.Tables, scriptTable, tagTable)

	spriteTable := gfx.NewSpriteTable(MaxSpriteSize)
	meshTable := gfx.NewMeshTable(MaxMeshSize)
	xfTable := gfx.NewTransformTable(MaxTransformSize)
	textTable := gfx.NewTextTable(MaxTextSize)

	g.DB.Tables = append(g.DB.Tables, spriteTable, meshTable, xfTable, textTable)

	psTable := effect.NewParticleSystemTable(MaxParticleSize)
	g.DB.Tables = append(g.DB.Tables, psTable)
}

func (g *Game) Input(dt float32) {

}

func (g *Game) Update() {
	// update
	dt := g.FPS.Smooth()

	// update input-system
	g.InputSystem.Frame()

	// update scene
	g.SceneManager.Update(dt)

	// update script
	g.ScriptSystem.Update(dt)

	g.InputSystem.Reset()

	//// simulation....

	// update sprite animation
	g.AnimationSystem.Update(dt)

	/// 动画更新，骨骼数据
	///g.AnimationSystem.Update(dt)

	// g.CollisionSystem.Update(dt)

	// 粒子系统更新
	g.ParticleSimulateSystem.Update(dt)

	// Render
	g.RenderSystem.Update(dt)

	// GUI
	g.UISystem.Draw(dt)

	// fps & profile
	g.DrawProfile()

	gfx.Flush()
}

func (g *Game) DrawProfile() {
	//dbg.FPS(g.FPS.fps)
	dbg.Move(5, 5)

	dbg.Color(0xFF000000)
	dbg.DrawRect(0, 0, 50, 6)

	// format: RGBA
	dbg.Color(0xFF00FF00)

	w := float32(g.fps)/60 * 50
	dbg.DrawRect(0, 0, w, 5)

	// format: RGBA
	dbg.Color(0xFF000000)

	dbg.Move(5, 10)
	dbg.DrawStrScaled(fmt.Sprintf("%d fps", g.fps), 0.6)

	dbg.NextFrame()
}

func (g *Game) Draw(dt float32) {
}
