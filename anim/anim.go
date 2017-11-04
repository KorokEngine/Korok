package anim

import (
)

// 骨骼动画系统
type SkeletonSystem struct {
	ST *SkeletonTable
}

func (sk *SkeletonSystem) Update(dt float32) {

}

/**
func (sr*SkeletonRender)Draw(skeleton *spine.Skeleton) {
	for _, slot := range skeleton.Slots {
		//
		if attachment, ok := slot.Attachment.(*spine.RegionAttachment); ok {
			// 计算得到插件坐标
			vert := attachment.Update(slot)
			//fmt.Println("index:", i, " name:", attachment.Name(), " update:", vert)
			sr.draw(vert[0:], attachment.Uvs[0:])
		} else {
			// fmt.Println("index:",i , " is null")
		}
	}
}
*/

