package game

import "github.com/go-gl/glfw/v3.2/glfw"

type FPS struct {
	scale float32
	curTime float32
	preTime float64

	dt float32
	fps int32
	pause bool
}

func (fps *FPS) SetScale(factor float32) {
	fps.scale = factor
}

func (fps *FPS) Step() float32 {
	time := glfw.GetTime()
	dt := time - fps.preTime
	fps.preTime = time
	fps.dt = float32(dt)
	fps.fps = int32(1/dt)

	return fps.dt
}

func (fps *FPS) Smooth() float32 {
	time := glfw.GetTime()
	dt := time - fps.preTime
	fps.preTime = time

	predt := fps.dt
	sdt := predt * .8 + float32(dt*.2)
	fps.dt = sdt
	fps.fps = int32(1/sdt)

	return fps.dt
}

func (fps *FPS) Pause() {
	fps.pause = true
}

func (fps *FPS) Resume() {
	fps.pause = false
}
