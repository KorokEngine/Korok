package korok

import (
	"log"
	"korok/hid"
	"korok/gameplay"
	"io/ioutil"
)

const VERSION_CODE  = 1
const VERSION_NAME  = "0.1"

type Options struct {
	Title string
	Width, Height int
}

func Run(options *Options)  {
	log.Println("Game Start! " + options.Title)

	g := &gameplay.Game{}
	gameplay.G = g

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


