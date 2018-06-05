package ween

import "korok.io/korok/math/ease"

type Animator struct {
	en *TweenEngine
	index int
}

func (eng *TweenEngine) NewAnimator() Animator {
	return Animator{eng, eng.New()}
}

func (am Animator) SetDuration(d float32) Animator{
	am.en.SetDuration(am.index, d)
	return am
}

func (am Animator) SetRepeat(count int, loop LoopType) Animator {
	am.en.SetRepeat(am.index, count, loop)
	return am
}

func (am Animator) SetFunction(function ease.Function) Animator {
	am.en.SetFunction(am.index, function)
	return am
}

func (am Animator) OnUpdate(cb UpdateCallback) Animator{
	am.en.SetUpdateCallback(am.index, cb)
	return am
}

func (am Animator) OnComplete(cb EndCallback) Animator{
	am.en.SetCompleteCallback(am.index, cb)
	return am
}

func (am Animator) Value() (f float32) {
	return am.en.Value(am.index)
}

func (am Animator) Valid() bool {
	if am.en == nil {
		return false
	}
	if _, ok := am.en.lookup[am.index]; !ok {
		return false
	}
	return true
}

func (am Animator) Forward() {
	am.en.Forward(am.index)
}

func (am Animator) Reverse() {
	am.en.Reverse(am.index)
}

func (am Animator) Stop() {
	am.en.Stop(am.index)
}

func (am Animator) Dispose() {
	am.en.Delete(am.index)
}


