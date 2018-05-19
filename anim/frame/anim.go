package frame

import (
	"korok.io/korok/engi"
	"korok.io/korok/gfx"
)

// implement frame-animation system

// 动画定义
type SpriteAnimation struct {
	Name string
	Start, Len int
	Loop bool
}

// SpriteAnimComp
type AnimationState struct {
	engi.Entity
	define int
	dt, rate float32
	ii int
	running bool
	once bool
}

// 序列帧动画
type Engine struct {
	// 原始的帧地址
	frames []gfx.Tex2D

	// 动画定义
	data []SpriteAnimation

	// sprite
	st *gfx.SpriteTable

	// 正在播放的动画
	states []AnimationState

	// 从名称到动画实例的映射
	names map[string]int
	_map map[engi.Entity]int
}

func NewEngine() *Engine {
	return &Engine{
		names:make(map[string]int),
		_map:make(map[engi.Entity]int),
	}
}

func (eng *Engine) RequireTable(tables []interface{}) {
	for _, t := range tables {
		if st, ok := t.(*gfx.SpriteTable); ok {
			eng.st = st; break
		}
	}
}

// 创建新的动画数据
// 现在 subText 还是指针，稍后会全部用 id 来索引。
// 动画资源全部存储在一个大的buffer里面，在外部使用索引引用即可.
// 采用这种设计，删除动画将会变得麻烦..
// 或者说无法删除动画，只能全部删除或者完全重新加载...
// 如何动画以组的形式存在，那么便可以避免很多问题
//
func (eng *Engine) NewAnimation(name string, frames []gfx.Tex2D, loop bool) {
	// copy frames
	start, size := len(eng.frames), len(frames)
	eng.frames = append(eng.frames, frames...)
	// new animation
	eng.data = append(eng.data, SpriteAnimation{name, start, size, loop})
	// keep mapping
	eng.names[name] = len(eng.data)-1
}

// 返回动画定义 - 好像并没有太大的意义
func (eng *Engine) Animation(name string) (anim *SpriteAnimation, seq []gfx.Tex2D) {
	if ii, ok := eng.names[name]; ok {
		anim = &eng.data[ii]
		seq  = eng.frames[anim.Start:anim.Start+anim.Len]
	}
	return
}

func (eng *Engine) newAnimationState(entity engi.Entity) int {
	id := len(eng.states)
	eng.states = append(eng.states, AnimationState{Entity:entity})
	return id
}

// 返回当前 Entity 绑定的动画状态
// 新建一个动画执行器？？
func (eng *Engine) Of(entity engi.Entity) Animator {
	if ii, ok := eng._map[entity]; ok {
		return Animator{eng, ii}
	} else {
		ii = eng.newAnimationState(entity)
		eng._map[entity] = ii
		return Animator{eng, ii}
	}
}

// 指向当前动画状态的 Handle
type Animator struct {
	sas *Engine
	index int
}


// 返回一个动画的当前执行状态
//func (am *Animator) State(name string) int {
//	st := am.sas.states[am.index]
//	return st.n // todo 计算出当前的动画状态
//}

// 创建一个动画状态，并关联到 Entity
func (am Animator) Play(name string) {
	am.sas.states[am.index].define = am.sas.names[name]
}

func (am Animator) Once() Animator {
	am.sas.states[am.index].once = true
	return am
}

func (am Animator) Rate(r float32) Animator {
	am.sas.states[am.index].rate = r
	return am
}

// delete from playing list
func (am Animator) Stop() {
	sz := len(am.sas.states)
	if am.index >= 0 && am.index < sz {
		if am.index != sz-1 {
			am.sas.states[am.index] = am.sas.states[sz-1]
		}
		am.sas.states = am.sas.states[:sz-1]
		am.index = -1
	}
}

func (eng *Engine) Update(dt float32) {
	// update animation
	for i := range eng.states {
		seq := &eng.states[i]
		seq.dt += dt
		if seq.dt > seq.rate {
			seq.ii = seq.ii + 1
			seq.dt = 0
		}
	}

	// update sprite-component
	for _, st := range eng.states {
		comp := eng.st.Comp(st.Entity)
		anim := eng.data[st.define]

		ii := st.ii % anim.Len
		frame := eng.frames[anim.Start+ii]
		comp.SetSprite(frame)
	}
}
