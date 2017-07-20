package spine

import (
	"math"
)

var BoneYDown = false

type BoneData struct {
	name     string
	parent   *BoneData
	Length   float32
	x        float32
	y        float32
	rotation float32
	scaleX   float32
	scaleY   float32
}

func NewBoneData(name string, parent *BoneData) *BoneData {
	boneData := new(BoneData)
	boneData.name = name
	boneData.parent = parent
	boneData.scaleX = 1
	boneData.scaleY = 1
	return boneData
}

type Bone struct {
	name          string
	Data          *BoneData
	parent        *Bone
	X             float32
	Y             float32
	Rotation      float32
	ScaleX        float32
	ScaleY        float32
	M00           float32
	M01           float32
	M10           float32
	M11           float32
	WorldX        float32
	WorldY        float32
	WorldRotation float32
	WorldScaleX   float32
	WorldScaleY   float32
}

func NewBone(boneData *BoneData, parent *Bone) *Bone {
	bone := new(Bone)
	bone.name = boneData.name
	bone.Data = boneData
	bone.parent = parent
	bone.ScaleX = 1
	bone.ScaleY = 1
	bone.WorldScaleX = 1
	bone.WorldScaleY = 1
	bone.SetToSetupPose()
	return bone
}

func (b *Bone) SetToSetupPose() {
	data := b.Data
	b.X = data.x
	b.Y = data.y
	b.Rotation = data.rotation
	b.ScaleX = data.scaleX
	b.ScaleY = data.scaleY
}

func (b *Bone) UpdateWorldTransform(flipX, flipY bool) {
	parent := b.parent
	if parent != nil {
		b.WorldX = b.X*parent.M00 + b.Y*parent.M01 + parent.WorldX
		b.WorldY = b.X*parent.M10 + b.Y*parent.M11 + parent.WorldY
		b.WorldScaleX = parent.WorldScaleX * b.ScaleX
		b.WorldScaleY = parent.WorldScaleY * b.ScaleY
		b.WorldRotation = parent.WorldRotation + b.Rotation
	} else {
		b.WorldX = b.X
		b.WorldY = b.Y
		b.WorldScaleX = b.ScaleX
		b.WorldScaleY = b.ScaleY
		b.WorldRotation = b.Rotation
	}
	radians := float64(b.WorldRotation) * math.Pi / 180.0
	cos := float32(math.Cos(radians))
	sin := float32(math.Sin(radians))
	b.M00 = cos * b.WorldScaleX
	b.M10 = sin * b.WorldScaleX
	b.M01 = -sin * b.WorldScaleY
	b.M11 = cos * b.WorldScaleY
	if flipX {
		b.M00 = -b.M00
		b.M01 = -b.M01
	}
	if flipY {
		b.M10 = -b.M10
		b.M11 = -b.M11
	}
	if BoneYDown {
		b.M10 = -b.M10
		b.M11 = -b.M11
	}
}
