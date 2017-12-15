package gfx

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok.io/korok/engi"
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

const STEP = 64
const none = 0

// 还可以做更细的拆分，把 Matrix 全部放到一个数组里面
// position, position, position .... position
// rotation, rotation, rotation .... rotation
// 把这两个数组申请为一块，然后分两个小组使用
type Transform struct {
	engi.Entity
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
	_map  map[uint32]int
	index, cap int
}

func NewTransformTable(cap int) *TransformTable {
	return &TransformTable{
		cap:cap,
		_map:make(map[uint32]int),
	}
}

// Create a new TransformComp for the entity,
// Return the old one, if it already exist .
func (tt *TransformTable) NewComp(entity engi.Entity) (xf *Transform) {
	if size := len(tt.comps); tt.index >= size {
		tt.comps = transformResize(tt.comps, size + STEP)
	}
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		xf = &tt.comps[v]
		return
	}

	xf = &tt.comps[tt.index]
	tt._map[ei] = tt.index
	xf.Entity = entity
	xf.Scale = mgl32.Vec2{1, 1}
	tt.index += 1
	return
}

// Return the TransformComp or nil
func (tt *TransformTable) Comp(entity engi.Entity) (xf *Transform) {
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		xf = &tt.comps[v]
	}
	return
}

func (tt *TransformTable) Alive(entity engi.Entity) bool {
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		return tt.comps[v].Entity != 0
	}
	return false
}

// Swap erase the TransformComp if exist
func (tt *TransformTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		if tail := tt.index -1; v != tail && tail > 0 {
			tt.comps[v] = tt.comps[tail]
			// remap index
			tComp := &tt.comps[tail]
			ei := tComp.Entity.Index()
			tt._map[ei] = v
			tComp.Entity = 0
		} else {
			tt.comps[tail].Entity = 0
		}

		tt.index -= 1
		delete(tt._map, ei)
	}
}

func (tt *TransformTable) Destroy() {
	tt.comps = make([]Transform, 0)
	tt._map = make(map[uint32]int)
	tt.index = 0
}

func (tt *TransformTable) Compact() {
	//
}

func (tt *TransformTable) Size() (size, cap int) {
	return tt.index, tt.cap
}

func transformResize(slice []Transform, size int) []Transform {
	newSlice := make([]Transform, size)
	copy(newSlice, slice)
	return newSlice
}

func intResize(slice []int, size int) []int {
	newSlice := make([]int, size)
	copy(newSlice, slice)
	return newSlice
}
