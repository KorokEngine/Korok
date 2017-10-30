package gfx

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok/ecs"
)

/**
默认情况下，每一个游戏对象都会有个 Transform 组件
为了保证内存紧凑，删除操作采用尾部复制的方式，但是
这样会导致组件的索引变动，通过  _map[]int 来维护
索引变化，这样便可以一直使用 EntityId 来访问组件。

Entity的销毁可以在每帧结束的时候扫描各个系统实现，
这样的效率会更高些。

节点系统会自然的包含组织(父子)关系，此处通过 Transform 组件
来实现类似的效果。
*/

const STEP = 100
const none = 0

// 还可以做更细的拆分，把 Matrix 全部放到一个数组里面
// position, position, position .... position
// rotation, rotation, rotation .... rotation
// 把这两个数组申请为一块，然后分两个小组使用
type Transform struct {
	ecs.Entity
	Position, Scale mgl32.Vec2
	Rotation        float32

	localPosition mgl32.Vec2
	localScale mgl32.Vec2
	localRotation float32

	Parent    uint32
	FirstChild uint32

	PrevSibling uint32
	NextSibling uint32
}

func (xf *Transform) SetPosition(position mgl32.Vec2) {
	xf.localPosition = position
	// compute world position
	if xf.Parent == none {
		xf.setPosition(mgl32.Vec2{0, 0}, position)
	} else {
		xf.setPosition(nodeSystem.comps[xf.Parent].Position, position)
	}
}

func (xf *Transform) setPosition(parent, local mgl32.Vec2) {
	xf.Position[0] = parent[0] + local[0]
	xf.Position[1] = parent[1] + local[1]

	// all child
	for child := xf.FirstChild; child != none; {
		node := nodeSystem.comps[child]
		child = node.NextSibling
		node.setPosition(xf.Position, node.localPosition)
	}
}

// apply scale to child
func (xf *Transform) SetScale(scale float32) {
	xf.Scale = mgl32.Vec2{scale, scale}
}

// apply
func (xf *Transform) SetRotation(rotation float32) {
	xf.Rotation = rotation
}

func (xf *Transform) AddChild(c *Transform) {
	if xf.FirstChild == none {
		xf.FirstChild = c.Index()
		c.Parent = xf.Index()
	} else {
		var prev uint32
		for next := xf.FirstChild; next != none; {
			prev = next
			next = nodeSystem.comps[next].NextSibling
		}
		(&nodeSystem.comps[prev]).NextSibling = c.Index()
		c.PrevSibling = prev
		c.Parent = xf.Index()
	}
}

var nodeSystem *TransformTable

type TransformTable struct {
 	comps []Transform
	_map []int
	index, capacity int
}

func (th *TransformTable) NewTransform(id uint32) *Transform {
	th.index += 1
	len := len(th.comps)
	if th.index >= len {
		th.comps = resize(th.comps, len + STEP)
	}
	comp := Transform{
		Entity: ecs.Entity(th.index),
		Scale:mgl32.Vec2{1, 1},
	}
	th.comps[th.index] = comp
	th._map[id] = th.index
	return nil
}

func (th *TransformTable) Get(id uint32) *Transform {
	return &th.comps[id]
}

func (th *TransformTable) Delete(id uint32) {
	//
}

func (th *TransformTable) Destroy() {
	//
}

func NewNodeSystem() *TransformTable {
	nodeSystem = new(TransformTable)
	nodeSystem.comps = make([]Transform, nodeSystem.capacity)
	return nodeSystem
}

func NewTransform(int uint32) *Transform {
	return nil
}

func resize(slice []Transform, size int) []Transform {
	newSlice := make([]Transform, size)
	copy(newSlice, slice)
	return newSlice
}
