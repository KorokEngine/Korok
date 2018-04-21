package anim

import "korok.io/korok/anim/tween"

func OfFloat(target *float32, from, to float32) tween.Animator {
	animator := tweenEngine.NewAnimator()
	animator.SetValue(from, to)
	animator.OnUpdate(func (f, v float32) {
		*target = v
	}).Start()
	return animator
}
