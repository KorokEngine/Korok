package game

import (
	"time"
)

type FPS struct {
	startTime time.Time
	curTime time.Duration
	preTime time.Duration

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
	now := time.Since(fps.startTime)
	dt := float32(now-fps.preTime)/float32(time.Second)
	fps.preTime = now
	fps.dt = float32(dt)
	fps.fps = int32(1/dt)

	return fps.dt
}

func (fps *FPS) Smooth() float32 {
	now := time.Since(fps.startTime)
	dt := float32(now-fps.preTime)/float32(time.Second)
	fps.preTime = now

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
