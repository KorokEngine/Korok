//+build !android,!ios

package hid

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"korok.io/korok/hid/gl"

	"fmt"
	"runtime"
	"log"
)

var windowCallback WindowCallback
var inputCallback InputCallback
var Keys  [1024]int

func init()  {
	runtime.LockOSThread()
}

func RegisterWindowCallback(callback WindowCallback) {
	windowCallback = callback
}

func RegisterInputCallback(callback InputCallback) {
	inputCallback = callback
}

func CreateWindow(option *WindowOptions)  {
	fmt.Println(glfw.GetVersionString())

	// 初始化 glfw
	err := glfw.Init()
	defer glfw.Terminate()

	if err != nil {
		log.Fatal(err)
	}

	// 设置API版本兼容, 最低支持：3.2
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	monitor := glfw.GetPrimaryMonitor()
	mode := &glfw.VidMode{
		Width:       1,
		Height:      1,
		RedBits:     8,
		GreenBits:   8,
		BlueBits:    8,
		RefreshRate: 60,
	}
	if monitor != nil {
		mode = monitor.GetVideoMode()
	}

	if option.FullScreen {
		option.Width = mode.Width
		option.Height = mode.Height
		glfw.WindowHint(glfw.Decorated, 0)
	} else {
		monitor = nil
	}
	if option.NoTitleBar {
		glfw.WindowHint(glfw.Visible, glfw.False)
	}

	// 创建窗口
	window, err := glfw.CreateWindow(option.Width, option.Height, option.Title, monitor, nil)
	if err != nil {
		fmt.Println("fail window:", err)
		return
	}
	defer window.Destroy()

	if !option.FullScreen {
		window.SetPos((mode.Width-option.Width)/2, (mode.Height-option.Height)/2)
	}



	// make the window's context current
	window.MakeContextCurrent()

	if option.NoVsync {
		glfw.SwapInterval(0)
	} else {
		glfw.SwapInterval(1)
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// Handle input callback
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}

		//if key >= 0 && key < 1024 {
		//	if action == glfw.Press {
		//		Keys[key] = gl.TRUE
		//	} else if action == glfw.Release {
		//		Keys[key] = gl.FALSE
		//	}
		//}
		if inputCallback != nil {
			if action == glfw.Press {
				inputCallback.OnKeyEvent(int(key), true)
			} else if action == glfw.Release {
				inputCallback.OnKeyEvent(int(key), false)
			}
		}
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if inputCallback != nil {
			x, y := w.GetCursorPos()
			pb := int(button)
			if action == glfw.Press {
				inputCallback.OnPointEvent(pb, true, float32(x), float32(y))
			} else {
				inputCallback.OnPointEvent(pb, false, float32(x), float32(y))
			}
		}
	})

	window.SetIconifyCallback(func(w *glfw.Window, iconified bool) {
		if iconified {
			windowCallback.OnPause()
		} else {
			windowCallback.OnResume()
		}
	})
	window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		windowCallback.OnFocusChanged(focused)
	})

	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		windowCallback.OnResize(int32(width), int32(height))
	})

	// init openGL
	// golang 版本必须调用此init方法来加载本地OpenGL指针，C原声不需要
	// 见go-gl文档:https://github.com/go-gl/gl
	if err := gl.Init(); err != nil {
		fmt.Print(err)
		return
	}

	// 读取本机的 OpenGL 版本
	//version := gl.GoStr(gl.GetString(gl.VERSION))
	//fmt.Println("OpenGL version", version)

	// viewport size
	w, h := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))

	// DEBUG
	// ========== Engine Start
	windowCallback.OnCreate(float32(w)/float32(option.Width))
	windowCallback.OnResize(int32(option.Width), int32(option.Height))

	// ========== Engine End
	// 全局配置
	//gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	if option.Clear[3] == 0 {
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	} else {
		gl.ClearColor(option.Clear[0], option.Clear[1], option.Clear[2], option.Clear[3])
	}

	// 如果窗口没有关闭，那么应该持续当前的循环
	// main loop...
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		windowCallback.OnLoop()

		// swap buffer
		window.SwapBuffers()

		// poll event
		glfw.PollEvents()

		// cursor should be update every frame!!
		if inputCallback != nil {
			x, y := window.GetCursorPos()
			inputCallback.OnPointEvent(-1000, false, float32(x), float32(y))
		}
	}
	windowCallback.OnDestroy()
}

