package anim

import (
	"korok.io/korok/engi"
	"korok.io/korok/anim/frame"
)

func OfSprite(entity engi.Entity) frame.Animator{
	return spriteEngine.Of(entity)
}