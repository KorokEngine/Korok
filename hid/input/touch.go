package input

import "korok.io/korok/math/f32"

type FingerId int

//// Touch & Mouse input go here
type PointerInput struct {
	// The mouse pointer always has a pointer-id of 0
	Id FingerId

	// The position and moved amount of pointer
	MousePos, MouseDelta f32.Vec2

	used bool
}


