package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
)

type CullSystem interface {
	Cull(comps []RenderComp, camera Camera) []CompRef
	UpdateBounding(id int32, bb BoundingBox)
	Collect(camera *Camera) []RenderObject
}

type GameObject interface {}

type BoundingBox struct {
	Min mgl32.Vec2
	Max mgl32.Vec2
}

// 裁剪一个
type CullObject struct {
	Comp GameObject
	BoundingBox
}

// 对于静态的对象，可以使用一些空间算法做到快速的筛选
// 对于动态的对象，最好的做法还是跑一遍 O(N) 的循环，
// 这样可以避免因维护算法带来的开销

type cullSystem struct {
	static_Objects []CullObject
	dynamicObjects []CullObject
}

// 返回相机可见的对象集合
func (*cullSystem) Collect(camera *Camera) []CullObject {
	return nil
}

// 在此注册一个 静态的游戏对象
func (*cullSystem) RegisterStatic(obj GameObject, bb BoundingBox) int32{
	return 0
}

// 注册一个 动态的游戏对象
func (*cullSystem) RegisterDynamic(obj GameObject, bb BoundingBox) int32 {
	return 0
}

// 更新游戏对象的 AABB
func (*cullSystem) UpdateBounding(id int32, bb BoundingBox) {

}




