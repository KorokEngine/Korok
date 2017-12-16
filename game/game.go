package game

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx"
	"korok.io/korok/particle"
	"korok.io/korok/anim"
	"korok.io/korok/physics"
	"korok.io/korok/assets"
	"korok.io/korok/hid/input"
	"korok.io/korok/gfx/dbg"

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
)

type Table interface{}

type DB struct {
	EntityM *engi.EntityManager
	Tables  []interface{}
}

// 统一管理游戏各个子系统的创建和销毁的地方
var G *Game
var scenes = make(map[string]Scene)
var current Scene

type Game struct {
	FPS
	DB

	*gfx.RenderSystem
	*input.InputSystem
	*ScriptSystem
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

/// input callback
func (g *Game) OnKeyEvent(key int, pressed bool) {
	g.InputSystem.SetKeyEvent(key, pressed)
}

func AddScene(scene Scene) {
	scenes[scene.Name()] = scene
	current = scene
}

// init subsystem
func (g *Game) Create() {
	gfx.Init()
	// render system
	rs := &gfx.RenderSystem{}
	g.RenderSystem = rs

	//
	// set table
	rs.RequireTable(g.DB.Tables)
	// set render
	var vertex, color string

	vertex, color = assets.Shader.GetShaderStr("batch")
	batchRender := gfx.NewBatchRender(vertex, color)
	rs.RegisterRender(gfx.RenderType(0), batchRender)

	vertex, color = assets.Shader.GetShaderStr("mesh")
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

	/// init debug
	dbg.Init()

	/// input system
	g.InputSystem = input.NewInputSystem()
	g.ScriptSystem = NewScriptSystem()
	g.ScriptSystem.RequireTable(g.DB.Tables)

	/// Customized scene
	if current != nil {
		current.Preload()
		current.Setup(g)
	}
}

// destroy subsystem
func (g *Game) Destroy() {
	g.RenderSystem.Destroy()
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

	psTable := &particle.ParticleSystemTable{}
	g.DB.Tables = append(g.DB.Tables, psTable)

	skTable := &anim.SkeletonTable{}
	g.DB.Tables = append(g.DB.Tables, skTable)

	rigidTable := &physics.RigidBodyTable{}
	colliderTable :=& physics.ColliderTable{}
	g.DB.Tables = append(g.DB.Tables, rigidTable, colliderTable)
}

func (g *Game) Input(dt float32) {
	if current != nil {
	}
}

func (g *Game) Update() {
	// update
	g.FPS.Step()

	// update input-system
	g.InputSystem.Frame()

	dt := g.FPS.dt
	if current != nil {
		current.Update(dt)
	}
	// update script
	g.ScriptSystem.Update(dt)

	g.InputSystem.Reset()

	//// simulation....

	/// 动画更新，骨骼数据
	///g.AnimationSystem.Update(dt)

	// g.CollisionSystem.Update(dt)

	// 粒子系统更新
	//g.ParticleSystem.Update(dt)

	// Render
	g.RenderSystem.Update(dt)

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
