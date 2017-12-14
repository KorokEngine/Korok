package gfx

import "korok.io/korok/gfx/bk"

// graphics context
// a wrapper for bk-api

func Init() {
	bk.Init()
}

func Flush() {
	bk.Flush()
}

func Destroy() {
	bk.Destroy()
}
