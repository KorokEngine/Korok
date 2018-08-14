package ween

import (
	"korok.io/korok/math/ease"
)

type UpdateCallback func(reverse bool, f float32)
type EndCallback func(reverse bool)

type Callback struct {
	UpdateCallback
	EndCallback
}

type Value struct {
	f float32
}

type AnimState uint8

const (
	Waiting AnimState = iota
	Running
	Stopped
	Dispose
)

// Defines what this animation should do when it reaches the end.
type LoopType uint8
const (
	Restart LoopType = iota
	PingPong
)

// Sets how many times the animation should be repeated.
const (
	RepeatOnce = 1
	RepeatInfinite = -1
)

// 维护动画的状态数据
// 底层动画系统，使用float作为单位 0-1
type Animation struct {
	index int
	clock, duration float32
	iteration, repeatCount int
	interpolator ease.Function
	LoopType
	state struct{
		AnimState
		dirty bool
	}
	reverse bool
	delay float32
}

func (anim *Animation) Reset() {
	anim.interpolator = ease.Linear
	anim.state.AnimState = Waiting
	anim.clock = 0
	anim.LoopType = Restart
	anim.repeatCount = 0
	anim.reverse = false
}

// 动画核心算法
func (anim *Animation) Animate(dt float32) (f float32) {
	anim.clock += dt
	fr := anim.clock / anim.duration
	if fr >= 1 {
		if anim.iteration < anim.repeatCount || anim.repeatCount == RepeatInfinite {
			// Time to repeat
			anim.iteration += int(fr)
			anim.clock = 0
			if anim.LoopType == PingPong {
				anim.reverse = !anim.reverse
			}
			for;fr >= 1; { fr = fr - 1 }
		} else {
			anim.state.AnimState = Stopped
			anim.state.dirty = true
			fr = 1
		}
	}
	if anim.state.AnimState == Stopped {
		f = 1
	} else {
		f = float32(anim.interpolator(float64(fr)))
	}
	if anim.reverse {
		f = 1 - f
	}
	return
}

type TweenEngine struct {
	anims []Animation
	values []Value
	callbacks []Callback

	time, scale float32
	active, cap int
	lookup map[int]int
	uniqueId int
}

func NewEngine() *TweenEngine {
	return &TweenEngine{
		scale:     1,
		anims:     make([]Animation, 32),
		values:    make([]Value, 32),
		callbacks: make([]Callback, 32),
		lookup:    make(map[int]int),
	}
}

func (eng *TweenEngine) New() (uid int) {
	uid = eng.uniqueId; eng.uniqueId++
	index := eng.active
	eng.active ++
	anim := &eng.anims[index];
	anim.Reset()
	anim.index = uid
	eng.values[index] = Value{}
	eng.callbacks[index] = Callback{}
	eng.lookup[uid] = index
	return
}

func (eng *TweenEngine) Delete(index int) {
	if v, ok := eng.lookup[index]; ok {
		eng.anims[v].state.AnimState = Dispose
	}
}

func (eng *TweenEngine) SetTimeScale(sk float32) {
	eng.scale = sk
}

// It's not easy to maintain the callback, one many add or delete animation
// while looping. Thanks to our delayed deleting strategy that delete operation
// only mark the animation 'Disposed' and will not change the size of the
// active animation. But add new animation will change the size, since new
// animation is always appended to the end. We can easily know the real size.
// Note: we need a BETTER callback policy.
func (eng *TweenEngine) Update(dt float32) {
	size := eng.active
	// 1. update
	for i := 0; i < size; i++ {
		if anim := &eng.anims[i]; anim.state.AnimState == Running {
			f := anim.Animate(dt)
			eng.values[i] = Value{f}
		}
	}

	// 2. callback
	for i := 0; i < eng.active; i++ {
		if anim := &eng.anims[i]; anim.state.AnimState == Stopped && anim.state.dirty {
			anim.state.dirty = false
			if cb := eng.callbacks[i].EndCallback; cb != nil {
				cb(anim.reverse)
			}
		} else if anim.state.AnimState == Running {
			if cb := eng.callbacks[i].UpdateCallback; cb != nil {
				cb(anim.reverse, eng.values[i].f)
			}
		}
	}
	// 3. delete dead
	var	i, j = 0, eng.active-1
	for i <= j {
		if anim := &eng.anims[i]; anim.state.AnimState == Dispose {
			eng.lookup[eng.anims[j].index] = i
			delete(eng.lookup, anim.index)
			eng.anims[i] = eng.anims[j]
			eng.values[i] = eng.values[j]
			eng.callbacks[i] = eng.callbacks[j]
			j--
		} else {
			i++
		}
	}
	eng.active = i

	//dbg.Move(10, 300)
	//dbg.DrawStrScaled(fmt.Sprintf("animation: %d", i), .6)
}

// Play an animation, produces values that range from 0.0 to 1.0,
// during a given duration.
func (eng *TweenEngine) Forward(index int) {
	if v, ok := eng.lookup[index]; ok {
		anim := &eng.anims[v]
		anim.clock = 0
		anim.state.AnimState = Running
		anim.state.dirty = true
		anim.iteration = 0
		anim.reverse = false

		eng.values[v] = Value{}
	}
}

// Play an animation in reverse. If the animation is already running,
// stop itself and play backwards from the point. If the animation is not
// running, then it will start from the end and play backwards.
func (eng *TweenEngine) Reverse(index int) {
	if v, ok := eng.lookup[index]; ok {
		if anim := &eng.anims[v]; anim.state.AnimState == Running {
			anim.clock = anim.duration - anim.clock
			anim.reverse = !anim.reverse
		} else {
			anim.reverse = !anim.reverse
			anim.clock = 0
			anim.state.AnimState = Running
			anim.state.dirty = true
			anim.iteration = 0
		}
	}
}

// Stops running this animation.
func (eng *TweenEngine) Stop(index int) {
	if v, ok := eng.lookup[index]; ok {
		eng.anims[v].state.AnimState = Stopped
		eng.anims[v].state.dirty = true
	}
}

// Duration is the length of time this animation should last.
func (eng *TweenEngine) SetDuration(index int, d float32) {
	if v, ok := eng.lookup[index]; ok {
		eng.anims[v].duration = d
	}
}

// Repeat the animation. If playback type is forward, restart the animation
// from start, if the playback type is backward or ping-pong,
func (eng *TweenEngine) SetRepeat(index int, count int, loop LoopType) {
	if v, ok := eng.lookup[index]; ok {
		eng.anims[v].repeatCount = count
		eng.anims[v].LoopType = loop
	}
}

func (eng *TweenEngine) SetFunction(index int, fn ease.Function) {
	if v, ok := eng.lookup[index]; ok {
		if fn != nil {
			eng.anims[v].interpolator = fn
		} else {
			eng.anims[v].interpolator = ease.Linear
		}
	}
}

func (eng *TweenEngine) SetUpdateCallback(index int, cb UpdateCallback) {
	if v, ok := eng.lookup[index]; ok {
		eng.callbacks[v].UpdateCallback = cb
	}
}

func (eng *TweenEngine) SetCompleteCallback(index int, cb EndCallback) {
	if v, ok := eng.lookup[index]; ok {
		eng.callbacks[v].EndCallback = cb
	}
}

func (eng *TweenEngine) Value(index int) (f float32) {
	if v, ok := eng.lookup[index]; ok {
		f = eng.values[v].f
	}
	return
}

func (eng *TweenEngine) Duration(index int) float32 {
	if v, ok := eng.lookup[index]; ok {
		return eng.anims[v].duration
	} else {
		return 0
	}
}

