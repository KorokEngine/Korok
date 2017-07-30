package anim

import (
	"korok/gfx"
	"korok/anim/spine"
)

type AnimatorComp struct {

}

func (*AnimatorComp) Play(id string) {

}

// 整个骨骼用一个Mesh表示
type SkeletonMesh struct {
	gfx.Mesh
}

func (sm*SkeletonMesh)Update(skeleton *spine.Skeleton) {
	for _, slot := range skeleton.Slots {
		//
		if attachment, ok := slot.Attachment.(*spine.RegionAttachment); ok {
			// 计算得到插件坐标
			attachment.Update(slot)
			//fmt.Println("index:", i, " name:", attachment.Name(), " update:", vert)
			// sr.draw(vert[0:], attachment.Uvs[0:])
		} else {
			// fmt.Println("index:",i , " is null")
		}
	}
}


func NewSkeleton() {

}

type AnimationSystem struct {

}

func NewAnimationSystem() *AnimationSystem {
	return &AnimationSystem{}
}

func (sys*AnimationSystem) Update(dt float32) {

}

