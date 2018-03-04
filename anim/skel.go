package anim

import "korok.io/korok/engi"



type SkeletonComp struct {
	//animations map[string]*AnimationClip
	//playingAnims []PlayingAnimation
	//
	////
	//Blender AnimationBlender
}

// 播放一个动画
func (skel *SkeletonComp) Play(name string) {

}

func (skel *SkeletonComp) IsPlaying(name string) bool {
	return false
}

// 添加一个动画到播放组
//func (skel *SkeletonComp) Add(clip AnimationClip) {
//
//}

// 过度到一个新的动画
func (skel *SkeletonComp) CrossFade(name string, timeSpan float32) {
	skel.Blend(name, 1.0, timeSpan)
}

// 设定动画的混合参数
func (skel *SkeletonComp) Blend(name string, weight, time float32) {

}

type SkeletonTable struct {
	_comps []SkeletonComp
	_index uint32
	_map   map[int]uint32

}

func (st *SkeletonTable) NewComp(entity engi.Entity) (sc *SkeletonComp) {
	sc = &st._comps[st._index]
	st._map[int(entity)] = st._index
	st._index ++
	return
}

func (st *SkeletonTable) Comp(entity engi.Entity) (sc *SkeletonComp) {
	if v, ok := st._map[int(entity)]; ok {
		sc = &st._comps[v]
	}
	return
}

// TODO impl
func (st *SkeletonTable) Delete(entity engi.Entity) (sc *SkeletonComp) {
	return nil
}

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



