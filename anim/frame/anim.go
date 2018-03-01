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

//
type AnimationState struct {
	engi.Entity
	define, n int
	dt, rate float32
	ii int
	running bool
}

// 序列帧动画
type SpriteAnimationSystem struct {
	// 原始的帧地址
	frames []*gfx.SubTex

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

// 创建新的动画数据
// 现在 subText 还是指针，稍后会全部用 id 来索引。
// 动画资源全部存储在一个大的buffer里面，在外部使用索引引用即可.
// 采用这种设计，删除动画将会变得麻烦..
// 或者说无法删除动画，只能全部删除或者完全重新加载...
// 如何动画以组的形式存在，那么便可以避免很多问题
//
func (sas *SpriteAnimationSystem) NewAnimation(name string, frames []*gfx.SubTex, loop bool) {
	// copy frames
	start := len(sas.frames)
	sas.frames = append(sas.frames, frames...)
	end := len(sas.frames)
	// new animation
	sas.data = append(sas.data, SpriteAnimation{name, start, end, loop})
	// keep mapping
	sas.names[name] = len(sas.data)-1
}

// 返回动画定义 - 好像并没有太大的意义
func (sas *SpriteAnimationSystem) Animation(name string) *SpriteAnimation {
	if ii, ok := sas.names[name]; ok {
		return &sas.data[ii]
	}
	return nil
}

func (sas *SpriteAnimationSystem) newAnimationState() int {
	id := len(sas.states)
	sas.states = append(sas.states, AnimationState{})
	return id
}

// 返回当前 Entity 绑定的动画状态
// 新建一个动画执行器？？
func (sas *SpriteAnimationSystem) With(entity engi.Entity) Animator {
	if ii, ok := sas._map[entity]; ok {
		return Animator{sas, ii}
	} else {
		return Animator{sas, sas.newAnimationState()}
	}
}

// 指向当前动画状态的 Handle
type Animator struct {
	sas *SpriteAnimationSystem
	index int
}


// 返回一个动画的当前执行状态
func (am *Animator) State(name string) int {
	st := am.sas.states[am.index]
	return st.n // todo 计算出当前的动画状态
}

// 创建一个动画状态，并关联到 Entity
func (am *Animator) Play(name string) {

}

func (sas *SpriteAnimationSystem) Update(dt float32) {
	// update animation
	for i := range sas.states {
		seq := &sas.states[i]
		seq.dt += dt
		if seq.dt > seq.rate {
			seq.ii = seq.ii + 1
			seq.dt = 0
		}
	}

	// update sprite-component
	for _, st := range sas.states {
		comp := sas.st.Comp(st.Entity)
		anim := sas.data[st.define]

		ii := st.ii % anim.Len
		comp.SubTex = sas.frames[anim.Start+ii]
	}

	// remove dead
	
}
