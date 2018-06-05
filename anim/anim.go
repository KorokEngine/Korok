package anim

import (
	"korok.io/korok/anim/frame"
	"korok.io/korok/anim/ween"
	"korok.io/korok/gfx"
)

type AnimationSystem struct {
	*frame.SpriteEngine
	*ween.TweenEngine

	// tables
	st *gfx.SpriteTable
	xf *gfx.TransformTable
}

func NewAnimationSystem() *AnimationSystem {
	return &AnimationSystem{
		SpriteEngine: frame.NewEngine(),
		TweenEngine: ween.NewEngine(),
	}
}

func (as *AnimationSystem) RequireTable(tables []interface{}) {
	as.SpriteEngine.RequireTable(tables)

	for _, t := range tables {
		switch table := t.(type) {
		case *gfx.SpriteTable:
			as.st = table
		case *gfx.TransformTable:
			as.xf = table
		}
	}
}

func (as *AnimationSystem) Update(dt float32) {
	as.SpriteEngine.Update(dt)
	as.TweenEngine.Update(dt)
}

// set shortcut
func SetDefaultAnimationSystem(as *AnimationSystem) {
	animationSystem = as
	spriteEngine = as.SpriteEngine
	tweenEngine = as.TweenEngine
}

// shortcut
var spriteEngine *frame.SpriteEngine
var tweenEngine *ween.TweenEngine
var animationSystem *AnimationSystem