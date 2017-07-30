package gameplay

import (
	"korok/ecs"
	"korok/gfx"
	"korok/space"
	"korok/physics"
)

/**
	在 korok 系统中是不需要精灵的，为了符合常见的2D引擎，加上精灵
 */

type Sprite ecs.Entity

func (sprite *Sprite) SetPosition(x, y float32) {

}

func (sprite *Sprite) Get()  {

}

func (sprite *Sprite) GetComponent(name string) interface{} {
	return nil
}

func (sprite *Sprite) AddComponent(comp interface{})  {
	// todo
}

// 行为脚本
func (sprite *Sprite) AddBehaviorComp(behavior BehaviorComp) {

}

func (sprite *Sprite) GetBehaviorComp() *BehaviorComp{
	return nil
}

// 渲染组件
func (sprite *Sprite) AddRenderComp(comp gfx.RenderComp)  {

}

func (sprite *Sprite) GetRenderComp()  *gfx.RenderComp{
	return nil
}

// 空间组件
func (sprite *Sprite) AddTransformComp(transform space.Transform) {

}

func (sprite *Sprite) GetTransformComp() *space.Transform {
	return nil
}

// 物理引擎 - 刚体/碰撞
func (sprite *Sprite) AddRigidBodyComp(rigidBody physics.RigidBodyComp)  {

}

func (sprite *Sprite) GetRigidBodyComp() *physics.RigidBodyComp{
	return nil
}

func (sprite *Sprite) AddColliderComp(collider physics.ColliderComp) {
	
}

func (sprite *Sprite) GetColliderComp() *physics.ColliderComp {
	return nil
}

// go 语言中没有泛型，好像只能这么写了..



