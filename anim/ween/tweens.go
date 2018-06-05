package ween

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"log"
)

func U8Lerp(from, to uint8, f float32) uint8 {
	v1, v2 := float32(from), float32(to)
	return uint8(v1+(v2-v1)*f)
}

func U16Lerp(from, to uint16, f float32) uint16 {
	v1, v2 := float32(from), float32(to)
	return uint16(v1+(v2-v1)*f)
}

func IntLerp(from, to int, f float32) int {
	return from + int(float32(to-from)*f)
}

func F32Lerp(from, to float32, f float32) float32 {
	return from + (to-from)*f
}

func Vec2Lerp(from, to f32.Vec2, f float32) f32.Vec2 {
	v2 := f32.Vec2{
		F32Lerp(from[0], to[0], f),
		F32Lerp(from[1], to[1], f),
	}
	return v2
}

func ColorLerp(from, to gfx.Color, f float32) gfx.Color {
	c := gfx.Color{
		R: U8Lerp(from.R, to.R, f),
		G: U8Lerp(from.G, to.G, f),
		B: U8Lerp(from.B, to.B, f),
		A: U8Lerp(from.A, to.A, f),
	}
	return c
}

type F32Tween struct {
	am Animator
	from, to float32
}

func (t *F32Tween) Animate(am Animator) Animator {
	t.am = am
	return am
}

func (t *F32Tween) Animator() Animator {
	if !t.am.Valid() {
		log.Println("animator is unavailable")
	}
	return t.am
}

func (t *F32Tween) Initialize(from, to float32, engine *TweenEngine) {
	t.from = from
	t.to = to
	t.am = engine.NewAnimator()
}

func (t *F32Tween) Range(from, to float32) *F32Tween {
	t.from = from
	t.to = to
	return t
}

func (t *F32Tween) Value() float32 {
	return F32Lerp(t.from, t.to, t.am.Value())
}

type Vec2Tween struct {
	Animator
	from, to f32.Vec2
}

func (t *Vec2Tween) Initialize(from, to f32.Vec2, engine *TweenEngine) {
	t.from = from
	t.to = to
	t.Animator = engine.NewAnimator()
}

func (t *Vec2Tween) Value() f32.Vec2 {
	return Vec2Lerp(t.from, t.to, t.Animator.Value())
}

type ColorTween struct {
	Animator
	from, to gfx.Color
}

func (t *ColorTween) Initialize(from, to gfx.Color, engine *TweenEngine) {
	t.from = from
	t.to = to
	t.Animator = engine.NewAnimator()
}

func (t *ColorTween) Value() gfx.Color {
	f := t.Animator.Value()
	return gfx.Color{
		R: U8Lerp(t.from.R, t.to.R, f),
		G: U8Lerp(t.from.G, t.to.G, f),
		B: U8Lerp(t.from.B, t.to.B, f),
		A: U8Lerp(t.from.A, t.to.A, f),
	}
}


