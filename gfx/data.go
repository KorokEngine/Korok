package gfx

import (
	"korok/gfx/bk"
	"unsafe"
)

type RenderData interface{
	Type() int32
}

type CompRef struct {
	Type int32
	*Transform
	*SpriteComp
}

// format <x,y,u,v rgba>
var P4C4 = []bk.VertexComp{
	{4, bk.ATTR_TYPE_FLOAT, 0, 0},
	{4, bk.ATTR_TYPE_UINT8, 16, 1},
}

// vertex struct
type PosTexColorVertex struct {
	X, Y, U, V float32
	RGBA       uint32
}

//
var PosTexColorVertexSize = unsafe.Sizeof(PosTexColorVertex{})





