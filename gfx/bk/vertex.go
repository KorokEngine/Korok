package bk

import "korok.io/korok/hid/gl"

/// Vertex attribute type enum
type AttrType uint8

const (
	ATTR_TYPE_INT8   AttrType = iota // byte
	ATTR_TYPE_UINT8                  // uint8, unsigned byte
	ATTR_TYPE_INT16                  // int16
	ATTR_TYPE_UINT16                 // uint16
	ATTR_TYPE_FIXED                  // fixed
	ATTR_TYPE_FLOAT                  // float32

	ATTR_TYPE_COUNT
)

// useful defines
// <x,y, u,v, color>
var P2T2C4 = []VertexComp{
	{2, ATTR_TYPE_FLOAT, 0, 0},
	{2, ATTR_TYPE_FLOAT, 8, 0},
	{4, ATTR_TYPE_UINT8, 16, 1},
}

// <x,y, u,v>
var P2T2 = []VertexComp{
	{2, ATTR_TYPE_FLOAT, 0, 0},
	{2, ATTR_TYPE_FLOAT, 8, 0},
}

// <x,y, color>
var P2C4 = []VertexComp{
	{2, ATTR_TYPE_FLOAT, 0, 0},
	{4, ATTR_TYPE_UINT8, 8, 1},
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
