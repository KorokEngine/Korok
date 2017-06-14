package gameplay

import (
	"korok/gfx"
)

var scenes = make(map[string]Scene)
var current Scene

type Game struct {
	State int
	Keys  [1024]int
	Renderer *gfx.SpriteRender
}

func AddScene(scene Scene)  {
	scenes[scene.Name()] = scene
	current = scene
}

func (g *Game) Init()  {
	current.Preload()
}

func (g *Game) Input(dt float32)  {

}

func (g *Game) Update(dt float32)  {
	current.Update(dt)
}

func (g *Game) Draw(dt float32)  {
	gfx.Update(dt)
}




