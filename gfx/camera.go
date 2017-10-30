package gfx

import "github.com/go-gl/mathgl/mgl32"

type CameraMode uint8
const (
	Perspective CameraMode = iota
	Orthographic
)

// TODO
type Camera struct {
	Eye mgl32.Vec3

}

type CameraComp struct {

}

type CameraSystem struct {

}
