package engi

/**
	定义ECS组件: Entity/Component/System

	Entity 仅仅是一个类型为 uint 的 id， 后8位表示类型，前24位表示索引
	一个Entity可以包含多个Comp所以不可能拿来表示Comp类型..
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

type ComponentType uint16

type Entity uint32

func (e Entity) Index() uint32 {
	return uint32(e) & INDEX_MASK
}

func (e Entity) Generation() uint32 {
	return uint32(e >> INDEX_BITS) & GENERATION_MASK
}

func (e Entity) Alive() bool {
	return e != 0
}

// Components
type Component interface {
}

// System
type System interface {
}

// 要不要用一个数组来管理所有的 Entity-Id 索引？
// 好处是可以方便的跟踪游戏整个生命周期中产生的所有对象.
// 同时需要一个 FreeList 来记录所有的对象.
type EntityManager struct {
	id Entity
}

func NewEntityManager() *EntityManager {
	return &EntityManager{}
}

func (em *EntityManager) New() Entity {
	em.id ++
	return em.id
}
