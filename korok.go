package korok

import (
	"log"
	"io/ioutil"

	"korok/gfx"
	"korok/engi"
	"korok/game"
	"korok/hid"
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

	hid.RegisterWindowCallback(g)
	hid.CreateWindow(&hid.WindowOptions{
		options.Title,
		options.Width,
		options.Height,
	})
}

func SetDebug(enable bool) {
	if enable == false {
		log.SetOutput(ioutil.Discard)
	}
}
var G *game.Game

///// entity-api
var Entity *engi.EntityManager

///// shortcut component-api for rendering system
var Sprite 	   *gfx.SpriteTable
var Mesh       *gfx.MeshTable
var Transform  *gfx.TransformTable


