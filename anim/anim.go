package anim

import (
	"korok.io/korok/anim/frame"
	"korok.io/korok/anim/tween"
)

type AnimationSystem struct {
	SpriteEngine *frame.Engine
	TweenEngine *tween.Engine
}

func NewAnimationSystem() *AnimationSystem {
	var (
		se = frame.NewEngine()
		te = tween.NewEngine()
	)
	spriteEngine = se
	tweenEngine = te
	return &AnimationSystem{
		SpriteEngine:se,
		TweenEngine:te,
	}
}

func (as *AnimationSystem) RequireTable(tables []interface{}) {
	as.SpriteEngine.RequireTable(tables)
}

func (as *AnimationSystem) Update(dt float32) {
	as.SpriteEngine.Update(dt)
	as.TweenEngine.Update(dt)
}

// shortcut
var spriteEngine *frame.Engine
var tweenEngine *tween.Engine