package spine

import (
	"math"
)

type Attachment interface {
	Name() string
}

type RegionAttachment struct {
	name     string
	X        float32
	Y        float32
	Rotation float32
	ScaleX   float32
	ScaleY   float32
	Width    float32
	Height   float32

	RendererObject       interface{}
	RegionOffsetX        float32
	RegionOffsetY        float32
	RegionWidth          float32
	RegionHeight         float32
	RegionOriginalWidth  float32
	RegionOriginalHeight float32

	Uvs    [8]float32
	offset [8]float32
}

func (a RegionAttachment) Name() string {
	return a.name
}

func NewRegionAttachment(name string) *RegionAttachment {
	return &RegionAttachment{
		name:   name,
		ScaleX: 1,
		ScaleY: 1,
	}
}

func (r *RegionAttachment) SetUVs(u float32, v float32, u2 float32, v2 float32, rotate bool) {
	uvs := &r.Uvs
	if rotate {
		uvs[2] = u
		uvs[3] = v2
		uvs[4] = u
		uvs[5] = v
		uvs[6] = u2
		uvs[7] = v
		uvs[0] = u2
		uvs[1] = v2
	} else {
		uvs[0] = u
		uvs[1] = v2
		uvs[2] = u
		uvs[3] = v
		uvs[4] = u2
		uvs[5] = v
		uvs[6] = u2
		uvs[7] = v2
	}
}

func (r *RegionAttachment) updateOffset() {
	width := r.Width
	height := r.Height
	scaleX := r.ScaleX
	scaleY := r.ScaleY
	regionScaleX := width / r.RegionOriginalWidth * scaleX
	regionScaleY := height / r.RegionOriginalHeight * scaleY
	localX := -width/2*scaleX + r.RegionOffsetX*regionScaleX
	localY := -height/2*scaleY + r.RegionOffsetY*regionScaleY
	localX2 := localX + r.RegionWidth*regionScaleX
	localY2 := localY + r.RegionHeight*regionScaleY
	rotation := r.Rotation
	rads := float64(rotation) * math.Pi / 180
	cos := float32(math.Cos(rads))
	sin := float32(math.Sin(rads))
	x := r.X
	y := r.Y
	localXCos := localX*cos + x
	localXSin := localX * sin
	localYCos := localY*cos + y
	localYSin := localY * sin
	localX2Cos := localX2*cos + x
	localX2Sin := localX2 * sin
	localY2Cos := localY2*cos + y
	localY2Sin := localY2 * sin
	r.offset[0] = localXCos - localYSin
	r.offset[1] = localYCos + localXSin
	r.offset[2] = localXCos - localY2Sin
	r.offset[3] = localY2Cos + localXSin
	r.offset[4] = localX2Cos - localY2Sin
	r.offset[5] = localY2Cos + localX2Sin
	r.offset[6] = localX2Cos - localYSin
	r.offset[7] = localYCos + localX2Sin
}

func (r *RegionAttachment) Update(slot *Slot) (verts [8]float32) {
	bone := slot.Bone

	s := slot.Skeleton()
	x := s.X + bone.WorldX
	y := s.Y + bone.WorldY
	m00 := bone.M00
	m01 := bone.M01
	m10 := bone.M10
	m11 := bone.M11

	verts[0] = r.offset[0]*m00 + r.offset[1]*m01 + x
	verts[1] = r.offset[0]*m10 + r.offset[1]*m11 + y
	verts[2] = r.offset[2]*m00 + r.offset[3]*m01 + x
	verts[3] = r.offset[2]*m10 + r.offset[3]*m11 + y
	verts[4] = r.offset[4]*m00 + r.offset[5]*m01 + x
	verts[5] = r.offset[4]*m10 + r.offset[5]*m11 + y
	verts[6] = r.offset[6]*m00 + r.offset[7]*m01 + x
	verts[7] = r.offset[6]*m10 + r.offset[7]*m11 + y

	return
}