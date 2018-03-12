package gfx

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/engi"
)

/// 可见性系统，以一个组件的形式呈现
type BoundingBox struct {
	Min f32.Vec2
	Max f32.Vec2
}

/// 每个Entity只有一个可见性对象，如果它包含多个 RenderComp，那么求出一个最大面积
/// 作为该对象的可见面积
type BoundingComp struct {
	BoundingBox
}

func (bc *BoundingComp) SetBounding(bb *BoundingBox) {
	bc.BoundingBox = *bb
}

// 合并两个矩形
func (bc *BoundingComp) Add(bb *BoundingBox) {

}

// 减去一个矩形的贡献
func (bc *BoundingComp) Sub(bb *BoundingBox) {

}

type BoundingTable struct {
	_comp []BoundingComp
	_index uint32
	_map [1024]uint32
}

func (bt *BoundingTable) NewComp(entity engi.Entity) (bc *BoundingComp) {
	bc = &bt._comp[bt._index]
	bt._map[entity] = bt._index
	bt._index ++
	return
}

func (bt *BoundingTable) Comp(entity engi.Entity) (bc *BoundingComp) {
	bc = &bt._comp[bt._map[entity]]
	return
}


type VisibilitySystem interface {
	UpdateBounding(id int32, bb BoundingBox)
	UpdateTransform()
	Collect(camera *Camera) []engi.Entity
}

// 对于静态的对象，可以使用一些空间算法做到快速的筛选
// 对于动态的对象，最好的做法还是跑一遍 O(N) 的循环，
// 这样可以避免因维护算法带来的开销

type visibilitySystem struct {
}

// 返回相机可见的对象集合
func (*visibilitySystem) Collect(camera *Camera) []engi.Entity{
	return nil
}

// 在此注册一个 静态的游戏对象
func (*visibilitySystem) RegisterStatic(entity engi.Entity, bb BoundingBox) int32{
	return 0
}

// 注册一个 动态的游戏对象
func (*visibilitySystem) RegisterDynamic(entity engi.Entity, bb BoundingBox) int32 {
	return 0
}

// 更新游戏对象的 AABB
func (*visibilitySystem) UpdateBounding(id int32, bb BoundingBox) {

}




