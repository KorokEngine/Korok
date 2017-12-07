package korok

import (
	"log"
	"io/ioutil"

	"korok/gfx"
	"korok/engi"
	"korok/game"
	"korok/hid"
	"korok/gfx/dbg"
	"korok/particle"
	"korok/anim"
	"korok/physics"
	"reflect"
)

const VERSION_CODE  = 1
const VERSION_NAME  = "0.1"

type Options struct {
	Title string
	Width, Height int
}

func Run(options *Options)  {
	log.Println("Game Start! " + options.Title)

	g := &game.Game{}
	G = g
	g.Init()

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
		case *particle.ParticleSystemTable:
			ParticleSystem = t
		case *anim.SkeletonTable:
			Skeleton = t
		case *physics.RigidBodyTable:
			RigidBody = t
		case *physics.ColliderTable:
			Collider = t
		case *game.TagTable:
			Tag = t
		case *game.ScriptTable:
			Script = t
		}
	}

	log.Printf("Load table: %v", len(g.DB.Tables))
	for i, v := range g.DB.Tables {
		log.Println(i, "table - ", reflect.TypeOf(v))
	}

	hid.RegisterWindowCallback(g)
	hid.CreateWindow(&hid.WindowOptions{
		options.Title,
		options.Width,
		options.Height,
	})
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
var ParticleSystem *particle.ParticleSystemTable

///// animation
var Skeleton       *anim.SkeletonTable

///// physics
var RigidBody *physics.RigidBodyTable
var Collider  *physics.ColliderTable


