package anim

import (
	"korok.io/korok/anim/ween"
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
)

func OfFloat(target *float32, from, to float32) ween.Animator {
	animator := tweenEngine.NewAnimator()
	animator.OnUpdate(func (r bool, f float32) {
		*target = ween.F32Lerp(from, to, f)
	})
	animator.OnComplete(func(reverse bool) {
		animator.Dispose()
	})
	return animator
}

func OfVec2(target *f32.Vec2, from, to f32.Vec2) ween.Animator {
	animator := tweenEngine.NewAnimator()
	animator.OnUpdate(func(reverse bool, f float32) {
		*target = ween.Vec2Lerp(from, to, f)
	})
	animator.OnComplete(func(reverse bool) {
		animator.Dispose()
	})
	return animator
}

func OfColor(target *gfx.Color, from, to gfx.Color) ween.Animator {
	animator := tweenEngine.NewAnimator()
	animator.OnUpdate(func(reverse bool, f float32) {
		*target = ween.ColorLerp(from, to, f)
	})
	animator.OnComplete(func(reverse bool) {
		animator.Dispose()
	})
	return animator
}