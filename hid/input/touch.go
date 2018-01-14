package input

import "github.com/go-gl/mathgl/mgl32"

type FingerId int

//// Touch & Mouse input go here
type PointerInput struct {
	// The mouse pointer always has a pointer-id of 0
	Id FingerId

	// The position and moved amount of pointer
	MousePos, MouseDelta mgl32.Vec2

	used bool
}


