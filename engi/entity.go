package engi

/**
	定义ECS组件: Entity/Component/System

	Entity 仅仅是一个类型为 uint32 的 id， 高8位表示出生代，低24位表示索引。
	因为 Entity 可以被不断的销毁和重建，可以使用高8位来记录在同一个索引位创建的Entity，
	这样最多可以用24位来记录Entity的id，8位记录id复用情况，如此，同一个id在循环了256次
	之后才会出现Entity重复的情况。

	三个主要接口：
	1. Create() Entity
	2. Alive(Entity)
	3. Destroy(Entity)
	增删查
 */

const IndexBits = 24
const IndexMask = (1<< IndexBits)-1

const GenerationBits = 8
const GenerationMask = (1<< GenerationBits)-1

type ComponentType uint16

type Entity uint32

func (e Entity) Index() uint32 {
	return uint32(e) & IndexMask
}

func (e Entity) Gene() uint8 {
	return uint8(e >>IndexBits) & GenerationMask
}

// 要不要用一个数组来管理所有的 Entity-Id 索引？
// 好处是可以方便的跟踪游戏整个生命周期中产生的所有对象.
// 同时需要一个 FreeList 来记录所有的对象.
type EntityManager struct {
	generation []uint8
	freelist   []uint32
	id Entity
}

func NewEntityManager() *EntityManager {
	return &EntityManager{}
}

func (em *EntityManager) New() Entity {
	var ei uint32
	var eg uint8

	if size := len(em.freelist); size > 0 {
		ei, em.freelist = em.freelist[0], em.freelist[1:]
		eg = em.generation[ei]
	} else {
		ei = uint32(len(em.generation))
		em.generation = append(em.generation, 0)
	}
	return Entity((uint32(eg) << IndexBits) | ei)
}

func (em *EntityManager) Alive(e Entity) bool {
	return em.generation[e.Index()] == e.Gene()
}

func (em *EntityManager) Destroy(e Entity) {
	ei := e.Index()
	em.generation[ei] ++
	em.freelist = append(em.freelist, ei)
}
