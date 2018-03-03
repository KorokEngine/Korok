package tween

import (
	"korok.io/korok/anim/tween/ease"
)

type StartCallback func(v float32)
type UpdateCallback func(f, v float32)
type EndCallback func(v float32)

type Callback struct {
	Start StartCallback
	Update UpdateCallback
	End EndCallback
}

type Value struct {
	f, v float32
	done bool
}

type AnimState uint8

const (
	Running AnimState = iota
	Stopped
	Waiting
)

const RepeatInfinite = -1

// 维护动画的状态数据
// 底层动画系统，使用float作为单位 0-1
type Animation struct {
	index int
	start, end, delta float32
	playTime, duration float32
	iteration, repeatCount int

	state AnimState

	interpolator ease.Function
}

func (anim *Animation) Reset() {
	anim.interpolator = ease.Linear
}

// 动画核心算法
func (anim *Animation) Animate(dt float32) (f, v float32, done bool) {
	anim.playTime += dt
	fr := anim.playTime / anim.duration

	if fr >= 1 {
		if anim.iteration < anim.repeatCount || anim.repeatCount == RepeatInfinite {
			// Time to repeat
			anim.iteration += int(f)
			anim.playTime = 0
			for;fr >= 1; { fr = fr - 1 }
		} else {
			done = true
			fr = 1
		}
	}
	f, v = anim.animateValue(fr)
	return
}

func (anim *Animation) animateValue(f float32) (fraction, v float32) {
	f = float32(anim.interpolator(float64(f)))
	fraction = f

	switch f {
	case 0:
		v = anim.start
	case 1:
		v = anim.end
	default:
		v = anim.start + anim.delta * f
	}
	return
}
