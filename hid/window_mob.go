// +build android ios !js

package hid

import (
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"korok.io/korok/hid/gl"

	"os"
	"sync"
)

var options *WindowOptions

var (
	once     sync.Once
	widthPx  int
	heightPx int
)

var (
	windowCallback WindowCallback
	inputCallback  InputCallback
	Keys           [1024]int
)

func RegisterWindowCallback(callback WindowCallback) {
	windowCallback = callback
}

func RegisterInputCallback(callback InputCallback) {
	inputCallback = callback
}

// Mobile always full-screen.
func CreateWindow(opt *WindowOptions) {
	options = opt
	app.Main(func(a app.App) {
		var (
			glctx interface{}
			sz    size.Event
		)
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageAlive) {
				case lifecycle.CrossOn:
					onCreate()
				case lifecycle.CrossOff:
					onDestroy()
				}
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx = e.DrawContext
					onStart(e)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					glctx = nil
				}
			case size.Event:
				sz = e
				onResize(e)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}
				onPaint(e, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				onTouch(e)
			case key.Event:
				onKey(e)
			}
		}
	})
}

func onCreate() {
	widthPx = app.DisplayMetrics.WidthPx
	heightPx = app.DisplayMetrics.HeightPx
}

func onStart(e lifecycle.Event) {
	if e.DrawContext == nil {
		return
	}
	gl.InitContext(e.DrawContext)
	bg := options.Clear
	if bg[3] == 0 {
		gl.ClearColor(1, 1, 1, 1)
	} else {
		gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	}
	once.Do(func() {
		windowCallback.OnCreate(float32(widthPx), float32(heightPx), 1)
	})
	windowCallback.OnResume()
}

func onStop() {
	windowCallback.OnPause()
	gl.Release()
}

func onDestroy() {
	windowCallback.OnDestroy()
	os.Exit(0)
}

func onResize(e size.Event) {
	widthPx, heightPx = e.WidthPx, e.HeightPx
	iw, ih := int32(e.WidthPx), int32(e.HeightPx)
	windowCallback.OnResize(iw, ih)
}

func onPaint(e paint.Event, sz size.Event) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	windowCallback.OnLoop()
}

func onTouch(e touch.Event) {
	var pressed bool
	switch e.Type {
	case touch.TypeBegin:
		pressed = true
	case touch.TypeMove:
		pressed = true
	case touch.TypeEnd:
		pressed = false
	}
	inputCallback.OnPointEvent(0, pressed, e.X, e.Y)
}

func onKey(e key.Event) {
	inputCallback.OnKeyEvent(int(e.Code), e.Direction == key.DirPress)
}
