package tween

import "korok.io/korok/anim/tween/ease"

type Animator struct {
	en *Engine
	index int
}

func (eng *Engine) NewAnimator() Animator {
	return Animator{eng, eng.NewAnimation()}
}

func (anim *Animator) SetDuration(d float32) *Animator{
	anim.en.SetDuration(anim.index, d)
	return anim
}

func (anim *Animator) SetValue(start, end float32) *Animator{
	anim.en.SetValue(anim.index, start, end)
	return anim
}

func (anim *Animator) SetRepeat(count int) *Animator {
	anim.en.SetRepeat(anim.index, count)
	return anim
}

func (anim *Animator) SetFunction(function ease.Function) *Animator {
	anim.en.SetFunction(anim.index, function)
	return anim
}

func (anim *Animator) OnStart(cb StartCallback) *Animator{
	anim.en.SetStartCallback(anim.index, cb)
	return anim
}

func (anim *Animator) OnUpdate(cb UpdateCallback) *Animator{
	anim.en.SetUpdateCallback(anim.index, cb)
	return anim
}

func (anim *Animator) OnComplete(cb EndCallback) *Animator{
	anim.en.SetCompleteCallback(anim.index, cb)
	return anim
}

func (anim *Animator) Value() (f, v float32) {
	return anim.en.Value(anim.index)
}

func (anim *Animator) Duration() float32 {
	return anim.en.Duration(anim.index)
}

func (anim *Animator) Animation() (*Animation, bool) {
	return anim.en.Animation(anim.index)
}

func (anim *Animator) Valid() bool {
	return anim.index == anim.en.anims[anim.index].index
}

func (anim *Animator) Start() {
	anim.en.Start(anim.index)
}

func (anim *Animator) Stop() {
	anim.en.Stop(anim.index)
}


