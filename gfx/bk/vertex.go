package bk

import "github.com/go-gl/gl/v3.2-core/gl"

/// Vertex attribute type enum
type AttrType uint8
const (
	ATTR_TYPE_UINT8 AttrType = iota // uint8
	ATTR_TYPE_UIN10                 // uint10,
	ATTR_TYPE_INT16                 // int16
	ATTR_TYPE_HALF                  // half
	ATTR_TYPE_FLOAT                 // float

	ATTR_TYPE_COUNT
)

// useful defines
// <x,y, u,v, color>
var P2T2C4 = []VertexComp {
	{2, ATTR_TYPE_FLOAT,  0, 0},
	{2, ATTR_TYPE_FLOAT,  8, 0},
	{4, ATTR_TYPE_UINT8, 16, 1},
}

// <x,y, u,v>
var P2T2 = []VertexComp {
	{2, ATTR_TYPE_FLOAT,  0, 0},
	{2, ATTR_TYPE_FLOAT,  8, 0},
}

// <x,y, color>
var P2C4 = []VertexComp {
	{2, ATTR_TYPE_FLOAT,  0, 0},
	{4, ATTR_TYPE_UINT8,  8, 1},
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

var g_AttrType = []uint32 {
	gl.BYTE,
	gl.UNSIGNED_BYTE,
	gl.SHORT,
	gl.UNSIGNED_SHORT,
	gl.FIXED,
	gl.FLOAT,
}

var g_AttrType2Size = []int32 {
	1, 1,
	2, 2,
	0, 	// Fixed Size ?
	4,
}