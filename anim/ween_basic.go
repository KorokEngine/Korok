package anim

import (
	"korok.io/korok/anim/ween"
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/math/ease"
)

// A fire-and-forget pattern Animator.
type proxyAnimator struct {
	ween.Animator
	update ween.UpdateCallback
	complete ween.EndCallback
}

// Proxy the SetDuration method.
func (proxy *proxyAnimator) SetDuration(d float32) *proxyAnimator {
	proxy.Animator.SetDuration(d)
	return proxy
}

// Proxy the SetRepeat method.
func (proxy *proxyAnimator) SetRepeat(count int, loop ween.LoopType) *proxyAnimator {
	proxy.Animator.SetRepeat(count, loop)
	return proxy
}

// Proxy the SetFunction method.
func (proxy *proxyAnimator) SetFunction(function ease.Function) *proxyAnimator {
	proxy.Animator.SetFunction(function)
	return proxy
}

// Proxy the OnUpdate method. proxyAnimator uses the UpdateCallback to update values internally,
// the user UpdateCallback will be called after it.
func (proxy *proxyAnimator) OnUpdate(fn ween.UpdateCallback) *proxyAnimator{
	proxy.update = fn
	return proxy
}

// Proxy the OnComplete method. proxyAnimator uses the CompleteCallback to remove itself from
// the TweenEngine, the user CompleteCallback will be called after it.
func (proxy *proxyAnimator) OnComplete(fn ween.EndCallback) *proxyAnimator{
	proxy.complete = fn
	return proxy
}

// OfFloat returns a Animator that animates between float values.
func OfFloat(target *float32, from, to float32) *proxyAnimator {
	proxy := &proxyAnimator{ Animator: tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func (r bool, f float32) {
		*target = ween.F32Lerp(from, to, f)
		if fn := proxy.update; fn != nil {
			fn(r, f)
		}
	})
	proxy.Animator.OnComplete(func(reverse bool) {
		proxy.Dispose()
		if fn := proxy.complete; fn != nil {
			fn(reverse)
		}
	})
	return proxy
}

// OfVec2 returns a Animator that animates between f32.Vec2 values.
func OfVec2(target *f32.Vec2, from, to f32.Vec2) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(r bool, f float32) {
		*target = ween.Vec2Lerp(from, to, f)
		if fn := proxy.update; fn != nil {
			fn(r, f)
		}
	})
	proxy.Animator.OnComplete(func(reverse bool) {
		proxy.Dispose()
		if fn := proxy.complete; fn != nil {
			fn(reverse)
		}
	})
	return proxy
}

// OfColor returns a Animator that animates between gfx.Color values.
func OfColor(target *gfx.Color, from, to gfx.Color) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(r bool, f float32) {
		*target = ween.ColorLerp(from, to, f)
		if fn := proxy.update; fn != nil {
			fn(r, f)
		}
	})
	proxy.Animator.OnComplete(func(reverse bool) {
		proxy.Dispose()
		if fn := proxy.complete; fn != nil {
			fn(reverse)
		}
	})
	return proxy
}