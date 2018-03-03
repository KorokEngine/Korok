package tween

import (
	//"fmt"
	"korok.io/korok/anim/tween/ease"
)

// 设定三种状态:
// 1. Running
// 2. Stopped
// 3. Waiting
// 调用 engine.NewAnimation 之后，得到一个 Animation，默认是 Waiting 状态
// 调用 animator.start() 之后，状态设置为 Running, 主循环开始处理动画
// 调用 animator.stop() 或者动画执行结束，状态为 Stopped, 实例会被回收

type Engine struct {
	anims []Animation

	values []Value
	callbacks []Callback

	time, scale float32
	active, cap int

	_map map[int]int
}

func NewEngine() *Engine {
	return &Engine {
		scale:1,
		anims:make([]Animation, 32),
		values:make([]Value, 32),
		callbacks:make([]Callback, 32),
		_map:make(map[int]int),
	}
}

func (en *Engine) NewAnimation() int {
	index := en.active
	en.active ++
	anim := &en.anims[index];
	anim.Reset()
	en.values[index] = Value{}
	en.callbacks[index] = dummyCallback
	en._map[index] = index
	return index
}

func (en *Engine) SetTimeScale(sk float32) {
	en.scale = sk
}

func (en *Engine) Update(dt float32) {
	size := en.active

	// 1. update
	for i := 0; i < size; i++ {

		anim := &en.anims[i]

		f, v, done := anim.Animate(dt)
		en.values[i] = Value{f, v, done }
	}

	// 2. callback
	for i := 0; i < size; i++ {
		if v := en.values[i]; v.done {
			en.callbacks[i].End(v.v)
		} else {
			en.callbacks[i].Update(v.f, v.v)
		}
	}
	// 3. delete dead
	for i := 0; i < size; i++ {
		if en.values[i].done {
			// find last live
			j := size - 1
			for ; j > i && en.values[j].done; j-- { size-- }
			if j > i {
				en.erase(i, j)
			}
			size -= 1
		}
	}
	if en.active > size {
		en.active = size
		en.resize(size)
	}
}

func (en *Engine) erase(i, j int) {
	en.anims[i] = en.anims[j]
	en.values[i] = en.values[j]
	en.callbacks[i] = en.callbacks[j]
}

func (en *Engine) resize(size int) {
	en.anims = en.anims[:size]
	en.values = en.values[:size]
	en.callbacks = en.callbacks[:size]
}

func (en *Engine) Start(index int) {
	en.anims[index].state = Running
	if cb := en.callbacks[index].Start; cb != nil {
		cb(en.values[index].v)
	}
}

func (en *Engine) Stop(index int) {

}

func (en *Engine) SetDuration(index int, d float32) {
	if ii, ok := en._map[index]; ok {
		en.anims[ii].duration = d
	}
}

func (en *Engine) SetValue(index int, v0, v1 float32) {
	if ii, ok := en._map[index]; ok {
		en.anims[ii].start = v0
		en.anims[ii].end = v1
		en.anims[ii].delta = v1-v0
	}
}

func (en *Engine) SetRepeat(index int, count int) {
	if ii, ok := en._map[index]; ok {
		en.anims[ii].repeatCount = count
	}
}

func (en *Engine) SetFunction(index int, fn ease.Function) {
	if ii, ok := en._map[index]; ok {
		if fn != nil {
			en.anims[ii].interpolator = fn
		} else {
			en.anims[ii].interpolator = ease.Linear
		}
	}
}


func (en *Engine) SetStartCallback(index int, cb StartCallback) {
	if ii, ok := en._map[index]; ok {
		en.callbacks[ii].Start = cb
	}
}

func (en *Engine) SetUpdateCallback(index int, cb UpdateCallback) {
	if ii, ok := en._map[index]; ok {
		en.callbacks[ii].Update = cb
	}
}

func (en *Engine) SetCompleteCallback(index int, cb EndCallback) {
	if ii, ok := en._map[index]; ok {
		en.callbacks[ii].End = cb
	}
}

func (en *Engine) Value(index int) (f, v float32) {
	if ii, ok := en._map[index]; ok {
		val := en.values[ii]
		f, v = val.f, val.v
	}
	return
}

func (en *Engine) Duration(index int) float32 {
	if ii, ok := en._map[index]; ok {
		return en.anims[ii].duration
	} else {
		return 0
	}
}


func (en *Engine) Animation(index int) (anim *Animation, ok bool) {
	if ii, ok := en._map[index]; ok {
		anim = &en.anims[ii]
		ok = true
	}
	return
}

var dummyCallback = Callback{
	Start: func(v float32) {

	},
	Update: func(f, v float32) {

	},
	End: func(v float32) {

	},
}


