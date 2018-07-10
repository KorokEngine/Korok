package game

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx"
	"korok.io/korok/effect"
	"korok.io/korok/anim"
	"korok.io/korok/anim/frame"
	"korok.io/korok/asset"
	"korok.io/korok/hid/input"
	"korok.io/korok/gfx/dbg"
	"korok.io/korok/gui"
	"korok.io/korok/audio"

	"log"
	"reflect"
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

func (g *Game) OnResize(w, h int32) {
	g.setGameSize(float32(w), float32(h))
}

func (g *Game) setGameSize(w, h float32) {
	// setup camera
	camera := &g.MainCamera
	camera.SetViewPort(w, h)
	camera.MoveTo(w/2, h/2)

	// gui real screen size
	gui.SetScreenSize(w, h)

	// dbg screen size
	dbg.SetScreenSize(w, h)
}

// init subsystem
func (g *Game) Create(ratio float32) {
	g.FPS.initialize()
	gfx.Init(ratio)
	audio.Init()

	// render system
	rs := gfx.NewRenderSystem()
	g.RenderSystem = rs

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

	// set feature
	srf := &gfx.SpriteRenderFeature{}
	srf.Register(rs)
	mrf := &gfx.MeshRenderFeature{}
	mrf.Register(rs)
	trf := &gfx.TextRenderFeature{}
	trf.Register(rs)

	// gui system
	ui := &gui.UIRenderFeature{}
	ui.Register(rs)

	/// init debug
	dbg.Init(g.Options.W, g.Options.H)

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
	anim.SetDefaultAnimationSystem(g.AnimationSystem)

	// audio system

	/// setup scene manager
	g.SceneManager.Setup(g)

	log.Println("Load Feature:", len(rs.FeatureList))
	for i, v := range rs.FeatureList {
		log.Println(i, " feature - ", reflect.TypeOf(v))
	}
}

// destroy subsystem
func (g *Game) Destroy() {
	g.RenderSystem.Destroy()
	audio.Destroy()
}

func (g *Game) Init() {
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

	spriteAnimTable := frame.NewFlipbookTable(MaxSpriteSize)
	g.DB.Tables = append(g.DB.Tables, spriteAnimTable)
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

	// fps & profile
	g.DrawProfile()

	//bk.Dump()
	audio.AdvanceFrame()

	// flush drawCall
	num := gfx.Flush()

	// drawCall = all-drawCall - camera-drawCall
	dc := num - len(g.RenderSystem.RenderList)
	dbg.LogFPS(int(g.fps), dc)
}

func (g *Game) DrawProfile() {
	// Advance frame
	dbg.NextFrame()
}

func (g *Game) Draw(dt float32) {
}
