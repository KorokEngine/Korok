package game

import (
	"time"
)

type FPS struct {
	startTime time.Time
	preTime time.Time

	dt, scale float32
	realdt float32
	fps int32
	pause bool
}

func (fps *FPS) initialize() {
	fps.startTime = time.Now()
}

func (fps *FPS) SetScale(factor float32) {
	fps.scale = factor
}

func (fps *FPS) Step() float32 {
	now := time.Now()
	du := now.Sub(fps.preTime); fps.preTime = now
	dt := float32(du)/float32(time.Second)

	fps.dt = float32(dt)
	fps.fps = int32(1/dt)
	return fps.dt
}

func (fps *FPS) Smooth() float32 {
	now := time.Now()
	du := now.Sub(fps.preTime); fps.preTime = now

	var dt float32
	if du < 3*time.Second {
		dt = float32(du)/float32(time.Second)
	} else {
		dt = 1.0/60
	}

	predt := fps.dt
	sdt := predt * .8 + float32(dt*.2)
	fps.dt = sdt
	fps.fps = int32(1/sdt)
	fps.realdt = dt

	return fps.dt
}

func (fps *FPS) Pause() {
	fps.pause = true
}

func (fps *FPS) Resume() {
	fps.pause = false
}
