package gfx

import "korok.io/korok/gfx/bk"

// graphics context
// a wrapper for bk-api

func Init() {
	bk.Init()
	bk.Reset(480, 320)

	// Enable debug text
	bk.SetDebug(bk.DEBUG_R|bk.DEBUG_Q)
}

func Flush() {
	bk.Flush()
}

func Destroy() {
	bk.Destroy()
}
