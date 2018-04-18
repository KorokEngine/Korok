package bk

import "korok.io/korok/hid/gl"

/// Vertex attribute type enum
type AttrType uint8

const (
	AttrInt8   AttrType = iota // byte
	AttrUInt8                  // uint8, unsigned byte
	AttrInt16                  // int16
	AttrUInt16                 // uint16
	AttrFixed                  // fixed
	AttrFloat                  // float32

	AttrCount
)

// useful defines
// <x,y, u,v, color>
var P2T2C4 = []VertexComp{
	{2, AttrFloat, 0, 0},
	{2, AttrFloat, 8, 0},
	{4, AttrUInt8, 16, 1},
}

// <x,y, u,v>
var P2T2 = []VertexComp{
	{2, AttrFloat, 0, 0},
	{2, AttrFloat, 8, 0},
}

// <x,y, color>
var P2C4 = []VertexComp{
	{2, AttrFloat, 0, 0},
	{4, AttrUInt8, 8, 1},
}

type VertexComp struct {
	Num        uint8
	Type       AttrType
	Offset     uint8
	Normalized uint8
}

func (comp *VertexComp) encode() uint16 {
	return 0
}

var g_AttrType = []uint32{
	gl.BYTE,
	gl.UNSIGNED_BYTE,
	gl.SHORT,
	gl.UNSIGNED_SHORT,
	gl.FIXED,
	gl.FLOAT,
}

var g_AttrType2Size = []int32{
	1, 1,
	2, 2,
	0, // Fixed Size ?
	4,
}
