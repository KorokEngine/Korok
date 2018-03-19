package korok

import (
	"log"
	"io/ioutil"
	"reflect"

	"korok.io/korok/gfx"
	"korok.io/korok/engi"
	"korok.io/korok/game"
	"korok.io/korok/hid"
	"korok.io/korok/gfx/dbg"
	"korok.io/korok/effect"
	"korok.io/korok/hid/input"
	"korok.io/korok/math/f32"
)

const VERSION_CODE  = 2
const VERSION_NAME  = "0.2"

type Options struct {
	Title string
	Width, Height int
	Clear f32.Vec4
}

func RunScene(options *Options, sc game.Scene) {
	PushScene(sc)
	Run(options)
}

func Run(options *Options)  {
	log.Println("Game Start! " + options.Title)

	g := &game.Game{}
	G = g
	g.Init(game.Options{options.Width, options.Height})

	Entity = g.DB.EntityM

	for _, table := range g.DB.Tables {
		switch t := table.(type) {
		case *gfx.SpriteTable:
			Sprite = t
		case *gfx.MeshTable:
			Mesh = t
		case *gfx.TransformTable:
			Transform = t
		case *gfx.TextTable:
			Text = t
		case *effect.ParticleSystemTable:
			ParticleSystem = t
		case *game.TagTable:
			Tag = t
		case *game.ScriptTable:
			Script = t
		}
	}

	log.Printf("LoadBitmap table: %v", len(g.DB.Tables))
	for i, v := range g.DB.Tables {
		log.Println(i, "table - ", reflect.TypeOf(v))
	}

	hid.RegisterWindowCallback(g)
	hid.RegisterInputCallback(g)
	hid.CreateWindow(&hid.WindowOptions{
		options.Title,
		options.Width,
		options.Height,
		options.Clear,
	})
}

func PushScene(sc game.Scene) {
	game.AddScene(sc)
}

func SetDebug(enable bool) {
	if enable == false {
		dbg.SetOutput(ioutil.Discard)
	}
}
var G *game.Game

///// entity-api
var Entity *engi.EntityManager

var Script *game.ScriptTable
var Tag    *game.TagTable

///// shortcut component-api for rendering system
var Sprite 	   *gfx.SpriteTable
var Mesh       *gfx.MeshTable
var Transform  *gfx.TransformTable
var Text       *gfx.TextTable

///// particle system
var ParticleSystem *effect.ParticleSystemTable

///// input system
var Input *input.InputSystem
