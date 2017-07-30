package ecs

/**
	定义ECS组件: Entity/Component/System

	Entity 仅仅是一个类型为 uint 的 id， 后8位表示类型，前24位表示索引
	三个主要接口：
	1. Create() Entity
	2. Alive(Entity)
	3. Destroy(Entity)
	增删查
 */

const INDEX_BITS  = 24
const INDEX_MASK  = (1<<INDEX_BITS)-1

const GENERATION_BITS  = 8
const GENERATION_MASK  = (1<<GENERATION_BITS)-1

type Entity uint32

func (entity Entity)Index() uint32 {
	return uint32(entity) & INDEX_MASK
}

func (entity Entity)Generation() uint32 {
	return uint32(entity >> INDEX_BITS) & GENERATION_MASK
}

// Components
type Component interface {
}

// System
type System interface {
}

var id Entity
// Entity Manager
func Create() Entity {
	id ++
	return id
}

func Alive(e Entity) bool {
	return false
}

func Destroy(e Entity) {

}
