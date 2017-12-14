package game

import "github.com/go-gl/glfw/v3.2/glfw"

//
type FPS struct {
	scale float32
	curTime float32
	preTime float64

	dt float32
	fps int32
}

func (*FPS) SetScale(factor float32) {

}

func (f *FPS) Step() {
	time := glfw.GetTime()
	dt := time - f.preTime
	f.preTime = time
	if dt <= 0.001 {
		f.dt = 16
		f.fps = 60
	} else {
		f.dt = float32(dt)
		f.fps = int32(1/dt)
	}
}

func (*FPS) Sleep(d float32) {

}

func (*FPS) Pause() {

}

func (*FPS) Resume() {

}
