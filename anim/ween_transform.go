package anim

import (
	"korok.io/korok/anim/ween"
	"korok.io/korok/engi"
	"korok.io/korok/math/f32"
)

// Convenient methods that uses to animate the Transform Component.

// Move the Entity to given value.
func Move(e engi.Entity, from, to f32.Vec2) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(r bool, f float32) {
		animationSystem.xf.Comp(e).SetPosition(ween.Vec2Lerp(from, to, f))
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

// Move the 'x' to given value.
func MoveX(e engi.Entity, from, to float32) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(r bool, f float32) {
		xf := animationSystem.xf.Comp(e)
		x  := ween.F32Lerp(from, to, f)
		xf.SetPosition(f32.Vec2{x, xf.Position()[1]})
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

// Move the 'y' to given value.
func MoveY(e engi.Entity, from, to float32) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(r bool, f float32) {
		xf := animationSystem.xf.Comp(e)
		y  := ween.F32Lerp(from, to, f)
		xf.SetPosition(f32.Vec2{xf.Position()[0], y})
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

// Scale the Entity to the given value.
func Scale(e engi.Entity, from, to f32.Vec2) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(reverse bool, f float32) {
		animationSystem.xf.Comp(e).SetScale(ween.Vec2Lerp(from, to, f))
		if fn := proxy.update; fn != nil {
			fn(reverse, f)
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

// Scale the 'x' to the given value.
func ScaleX(e engi.Entity, from, to float32) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(reverse bool, f float32) {
		xf := animationSystem.xf.Comp(e)
		x  := ween.F32Lerp(from, to, f)
		xf.SetScale(f32.Vec2{x, xf.Scale()[1]})
		if fn := proxy.update; fn != nil {
			fn(reverse, f)
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

// Scale the 'x' to the given value.
func ScaleY(e engi.Entity, from, to float32) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(reverse bool, f float32) {
		xf := animationSystem.xf.Comp(e)
		y  := ween.F32Lerp(from, to, f)
		xf.SetScale(f32.Vec2{xf.Scale()[0], y})
		if fn := proxy.update; fn != nil {
			fn(reverse, f)
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

// Rotate the Entity to given value.
func Rotate(e engi.Entity, from, to float32) *proxyAnimator {
	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(reverse bool, f float32) {
		animationSystem.xf.Comp(e).SetRotation(ween.F32Lerp(from, to, f))
		if fn := proxy.update; fn != nil {
			fn(reverse, f)
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


