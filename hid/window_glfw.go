package hid

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.2-core/gl"

	"korok/gameplay"
)

var game = &gameplay.Game{}

func init()  {
	runtime.LockOSThread()
}

func CreateWindow(option *WindowOptions)  {
	fmt.Println(glfw.GetVersionString())

	// 初始化 glfw
	err := glfw.Init()
	defer glfw.Terminate()

	if err != nil {
		fmt.Println(err)
		return
	}

	// 设置API版本兼容, 最低支持：3.2
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// 创建窗口
	window, err := glfw.CreateWindow(option.Width, option.Height, option.Title, nil, nil)
	if err != nil {
		fmt.Println("fail window:", err)
		return
	}
	defer window.Destroy()

	// make the window's context current
	window.MakeContextCurrent()

	// Handle input callback
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}

		if key >= 0 && key < 1024 {
			if action == glfw.Press {
				game.Keys[key] = gl.TRUE
			} else if action == glfw.Release {
				game.Keys[key] = gl.FALSE
			}
		}
	})

	// init openGL
	// golang 版本必须调用此init方法来加载本地OpenGL指针，C原声不需要
	// 见go-gl文档:https://github.com/go-gl/gl
	if err := gl.Init(); err != nil {
		fmt.Print(err)
		return
	}

	// 读取本机的 OpenGL 版本
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// viewport size
	w, h := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))

	// DEBUG
	// ========== Engine Start

	game.Init()

	// ========== Engine End
	// 全局配置
	//gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	previousTime := glfw.GetTime()

	// 如果窗口没有关闭，那么应该持续当前的循环
	// main loop...
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// update
		time :=  glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		dt := float32(elapsed)
		game.Input(dt)
		game.Update(dt)
		game.Draw(dt)

		// swap buffer
		window.SwapBuffers()

		// poll event
		glfw.PollEvents()

	}
}

