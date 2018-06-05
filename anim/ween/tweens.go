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

// A float32 linear interpolation between a beginning and ending value.
// It use the Animator as the input.
type F32Tween struct {
	am Animator
	from, to float32
}

// Animate sets the Animator that drives the Tween.
func (t *F32Tween) Animate(am Animator) Animator {
	t.am = am
	return am
}

// Animator returns the Animator driving the Tween.
func (t *F32Tween) Animator() Animator {
	if !t.am.Valid() {
		log.Println("animator is unavailable")
	}
	return t.am
}

// Range sets the beginning and ending value of the F32Tween.
func (t *F32Tween) Range(from, to float32) *F32Tween {
	t.from = from
	t.to = to
	return t
}

// Returns the interpolated value for the current value of the given Animator.
func (t *F32Tween) Value() float32 {
	return F32Lerp(t.from, t.to, t.am.Value())
}

// A f32.Vec2 linear interpolation between a beginning and ending value.
// It use the Animator as the input.
type Vec2Tween struct {
	am Animator
	from, to f32.Vec2
}

// Animate sets the Animator that drives the Tween.
func (t *Vec2Tween) Animate(am Animator) Animator {
	t.am = am
	return am
}

// Animator returns the Animator driving the Tween.
func (t *Vec2Tween) Animator() Animator {
	return t.am
}

// Range sets the beginning and ending value of the Vec2Tween.
func (t *Vec2Tween) Range(from, to f32.Vec2) *Vec2Tween {
	t.from, t.to = from, to
	return t
}

// Returns the interpolated value for the current value of the given Animator.
func (t *Vec2Tween) Value() f32.Vec2 {
	return Vec2Lerp(t.from, t.to, t.am.Value())
}

// A gfx.Color linear interpolation between a beginning and ending value.
// It use the Animator as the input.
type ColorTween struct {
	am Animator
	from, to gfx.Color
}

// Range sets the beginning and ending value of the ColorTween.
func (t *ColorTween) Range(from, to gfx.Color) *ColorTween {
	t.from, t.to = from, to
	return t
}

// Animate sets the Animator that drives the Tween.
func (t *ColorTween) Animate(animator Animator) Animator {
	t.am = animator
	return animator
}

// Animator returns the Animator driving the Tween.
func (t *ColorTween) Animator() Animator {
	return t.am
}

// Returns the interpolated value for the current value of the given Animator.
func (t *ColorTween) Value() gfx.Color {
	return ColorLerp(t.from, t.to, t.am.Value())
}


