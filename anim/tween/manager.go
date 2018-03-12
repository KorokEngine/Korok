package tween

import (
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

func (eng *Engine) NewAnimation() int {
	index := eng.active
	eng.active ++
	anim := &eng.anims[index];
	anim.Reset()
	eng.values[index] = Value{}
	eng.callbacks[index] = dummyCallback
	eng._map[index] = index
	return index
}

func (eng *Engine) SetTimeScale(sk float32) {
	eng.scale = sk
}

func (eng *Engine) Update(dt float32) {
	size := eng.active

	// 1. update
	for i := 0; i < size; i++ {

		anim := &eng.anims[i]

		f, v, done := anim.Animate(dt)
		eng.values[i] = Value{f, v, done }
	}

	// 2. callback
	for i := 0; i < size; i++ {
		if v := eng.values[i]; v.done {
			eng.callbacks[i].End(v.v)
		} else {
			eng.callbacks[i].Update(v.f, v.v)
		}
	}
	// 3. delete dead
	for i := 0; i < size; i++ {
		if eng.values[i].done {
			// find last live
			j := size - 1
			for ; j > i && eng.values[j].done; j-- { size-- }
			if j > i {
				eng.erase(i, j)
			}
			size -= 1
		}
	}
	if eng.active > size {
		eng.active = size
		eng.resize(size)
	}
}

func (eng *Engine) erase(i, j int) {
	eng.anims[i] = eng.anims[j]
	eng.values[i] = eng.values[j]
	eng.callbacks[i] = eng.callbacks[j]
}

func (eng *Engine) resize(size int) {
	eng.anims = eng.anims[:size]
	eng.values = eng.values[:size]
	eng.callbacks = eng.callbacks[:size]
}

func (eng *Engine) Start(index int) {
	eng.anims[index].state = Running
	if cb := eng.callbacks[index].Start; cb != nil {
		cb(eng.values[index].v)
	}
}

func (eng *Engine) Stop(index int) {

}

func (eng *Engine) SetDuration(index int, d float32) {
	if ii, ok := eng._map[index]; ok {
		eng.anims[ii].duration = d
	}
}

func (eng *Engine) SetValue(index int, v0, v1 float32) {
	if ii, ok := eng._map[index]; ok {
		eng.anims[ii].start = v0
		eng.anims[ii].end = v1
		eng.anims[ii].delta = v1-v0
	}
}

func (eng *Engine) SetRepeat(index int, count int) {
	if ii, ok := eng._map[index]; ok {
		eng.anims[ii].repeatCount = count
	}
}

func (eng *Engine) SetFunction(index int, fn ease.Function) {
	if ii, ok := eng._map[index]; ok {
		if fn != nil {
			eng.anims[ii].interpolator = fn
		} else {
			eng.anims[ii].interpolator = ease.Linear
		}
	}
}


func (eng *Engine) SetStartCallback(index int, cb StartCallback) {
	if ii, ok := eng._map[index]; ok {
		eng.callbacks[ii].Start = cb
	}
}

func (eng *Engine) SetUpdateCallback(index int, cb UpdateCallback) {
	if ii, ok := eng._map[index]; ok {
		eng.callbacks[ii].Update = cb
	}
}

func (eng *Engine) SetCompleteCallback(index int, cb EndCallback) {
	if ii, ok := eng._map[index]; ok {
		eng.callbacks[ii].End = cb
	}
}

func (eng *Engine) Value(index int) (f, v float32) {
	if ii, ok := eng._map[index]; ok {
		val := eng.values[ii]
		f, v = val.f, val.v
	}
	return
}

func (eng *Engine) Duration(index int) float32 {
	if ii, ok := eng._map[index]; ok {
		return eng.anims[ii].duration
	} else {
		return 0
	}
}


func (eng *Engine) Animation(index int) (anim *Animation, ok bool) {
	if ii, ok := eng._map[index]; ok {
		anim = &eng.anims[ii]
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


