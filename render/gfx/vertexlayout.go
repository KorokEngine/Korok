package gfx


/// Vertex attribute enum
type AttribEnum uint16
const (
	ATTRIB_POSITION  AttribEnum = iota
	ATTRIB_NORMAL
	ATTRIB_TANGENT
	ATTRIB_BITANGENT
	ATTRIB_COLOR0
	ATTRIB_COLOR1
	ATTRIB_COLOR2
	ATTRIB_COLOR3
	ATTRIB_INDICES
	ATTRIB_WEIGHT
	ATTRIB_TEXCOORD0
	ATTRIB_TEXCOORD1
	ATTRIB_TEXCOORD2
	ATTRIB_TEXCOORD3
	ATTRIB_TEXCOORD4
	ATTRIB_TEXCOORD5
	ATTRIB_TEXCOORD6
	ATTRIB_TEXCOORD7

	ATTRIB_COUNT
)

/// Vertex attribute type enum
type AttribType uint16
const (
	ATTRIB_TYPE_UINT8 AttribType = iota 	// uint8
	ATTRIB_TYPE_UIN10 						// uint10,
	ATTRIB_TYPE_INT16						// int16
	ATTRIB_TYPE_HALF							// half
	ATTRIB_TYPE_FLOAT

	ATTRIB_TYPE_COUNT
)


/// Vertex declaration
type VertexLayout struct {
	hash 		uint32
	stride 		uint16
	offset 		[ATTRIB_COUNT]uint16
	attributes 	[ATTRIB_COUNT]uint16
}

func (vd *VertexLayout) begin(renderer RendererType) *VertexLayout {
	return nil
}

func (vd *VertexLayout) end() {

}

/// default: normalized=false, asInt=false
func (vd *VertexLayout) add(attrib AttribEnum, num uint8, _type AttribType, normalized, asInt bool) *VertexLayout {
	return nil
}

func (vd *VertexLayout) skip(num uint8) *VertexLayout {
	return nil
}

func (vd *VertexLayout) decode() {

}

func (vd *VertexLayout) has(attrib AttribEnum) {

}

func (vd *VertexLayout) getOffset(attrib AttribEnum) uint16 {
	return 0
}

func (vd *VertexLayout) getStride() uint16{
	return 0
}

func (vd *VertexLayout) getSize(num uint32) uint32 {
	return 0
}


func initAttribTypeSizeTable(_type RendererType) {

}

//////////////// static & global filed
