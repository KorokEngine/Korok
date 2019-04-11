package hid

import (
	"syscall/js"

	"korok.io/korok/hid/gl"

	"runtime"
	"strconv"
	"time"
)

var windowCallback WindowCallback
var inputCallback InputCallback
var Keys [1024]int

func init() {
	runtime.LockOSThread()
}

func RegisterWindowCallback(callback WindowCallback) {
	windowCallback = callback
}

func RegisterInputCallback(callback InputCallback) {
	inputCallback = callback
}

func consume(event js.Value) {
	event.Call("stopPropagation")
	event.Call("preventDefault")
}

func CreateWindow(option *WindowOptions) {
	document := js.Global().Get("document")
	document.Set("title", option.Title)

	canvas := document.Call("createElement", "canvas")

	ww := js.Global().Get("innerWidth").Int()
	hh := js.Global().Get("innerHeight").Int()

	r := float32(ww) / float32(option.Width)
	wh := float32(option.Width) / float32(option.Height)
	h := int(float32(ww) / wh)
	var w int
	if hh >= h {
		w = ww
	} else {
		r = float32(hh) / float32(option.Height)
		hw := float32(option.Height) / float32(option.Width)
		w = int(float32(hh) / hw)
		h = hh
	}

	canvas.Set("width", strconv.Itoa(w))
	canvas.Set("height", strconv.Itoa(h))
	document.Get("body").Call("appendChild", canvas)
	err := gl.Init(canvas)
	if err != nil {
		js.Global().Call("alert", "Error: "+err.Error())
		return
	}
	// DEBUG

	mousedown := js.FuncOf(func(this js.Value, arg []js.Value) interface{} {
		consume(arg[0])
		rect := canvas.Call("getBoundingClientRect")
		x := arg[0].Get("clientX").Int() - rect.Get("left").Int()
		y := arg[0].Get("clientY").Int() - rect.Get("top").Int()
		button := arg[0].Get("button").Int()
		inputCallback.OnPointEvent(button, true, float32(x)/r, float32(y)/r)
		return nil
	})
	mouseup := js.FuncOf(func(this js.Value, arg []js.Value) interface{} {
		consume(arg[0])
		rect := canvas.Call("getBoundingClientRect")
		x := arg[0].Get("clientX").Int() - rect.Get("left").Int()
		y := arg[0].Get("clientY").Int() - rect.Get("top").Int()
		button := arg[0].Get("button").Int()
		inputCallback.OnPointEvent(button, false, float32(x)/r, float32(y)/r)
		return nil
	})
	keydown := js.FuncOf(func(this js.Value, arg []js.Value) interface{} {
		consume(arg[0])
		// TODO 这里需要处理特殊按键
		button := arg[0].Get("key").String()
		inputCallback.OnKeyEvent(int(button[0]), true)
		return nil
	})
	keyup := js.FuncOf(func(this js.Value, arg []js.Value) interface{} {
		consume(arg[0])
		// TODO 这里需要处理特殊按键
		button := arg[0].Get("key").String()
		inputCallback.OnKeyEvent(int(button[0]), false)
		return nil
	})

	// ========== Engine Start
	windowCallback.OnCreate(float32(option.Width), float32(option.Height), r)
	windowCallback.OnResize(int32(option.Width), int32(option.Height))

	// resize := js.FuncOf(func(this js.Value, arg []js.Value) interface{} {
	// 	consume(arg[0])

	// 	w = js.Global().Get("innerWidth").Int()

	// 	wh = float32(option.Width) / float32(option.Height)
	// 	h = int(float32(w) / wh)

	// 	canvas.Set("width", strconv.Itoa(w))
	// 	canvas.Set("height", strconv.Itoa(h))

	// 	windowCallback.OnResize(int32(option.Width), int32(option.Height))

	// 	return nil
	// })

	canvas.Call("addEventListener", "mousedown", mousedown, true)
	canvas.Call("addEventListener", "mouseup", mouseup, true)
	document.Call("addEventListener", "keydown", keydown, true)
	document.Call("addEventListener", "keyup", keyup, true)

	// js.Global().Call("addEventListener", "resize", resize, true)

	st := time.Second / 60
	ticker := time.NewTicker(st)
	for _ = range ticker.C {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		windowCallback.OnLoop()
		// window.SwapBuffers()
	}

	mousedown.Release()
	mouseup.Release()
	keydown.Release()
	keyup.Release()

	// resize.Release()

	windowCallback.OnDestroy()
}
