package anim

import (
	"korok.io/korok/anim/ween"
	"korok.io/korok/engi"
	"korok.io/korok/math/f32"
)

type InterpolationType uint8

const (
	Linear InterpolationType = iota
)

// OfTransform animate TransformComp by the specified property.
// 1. Position.X/Position.Y
// 2. Scale.X/Scale.Y
// 3. Rotate
func OfTransform(en engi.Entity, property string) ween.Animator {
	animator := tweenEngine.NewAnimator()
	updater := transformProperty(en, property)
	animator.OnUpdate(updater)
	return animator
}

func transformProperty(en engi.Entity, p string) (fn func(r bool, v float32)) {
	switch p {
	case "Position.X":
		fn = func(r bool, v float32) {
			xf := as.xf.Comp(en)
			xf.SetPosition(f32.Vec2{v, xf.Position()[1]})
		}
	case "Position.Y":
		fn = func(r bool, v float32) {
			xf := as.xf.Comp(en)
			xf.SetPosition(f32.Vec2{xf.Position()[0], v})
		}
	case "Scale.X":
		fn = func(r bool, v float32) {
			xf := as.xf.Comp(en)
			xf.SetScale(f32.Vec2{v, xf.Position()[1]})
		}
	case "Scale.Y":
		fn = func(r bool, v float32) {
			xf := as.xf.Comp(en)
			xf.SetScale(f32.Vec2{xf.Position()[0], v})
		}
	case "Rotate":
		fn = func(r bool, v float32) {
			xf := as.xf.Comp(en)
			xf.SetRotation(v)
		}
	default:
		fn = func(r bool, v float32) {
			// empty
		}
	}
	return
}




