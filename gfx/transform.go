package gfx

import (
	"korok.io/korok/math/f32"
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

type SRT struct {
	Scale f32.Vec2
	Rotation float32
	Position f32.Vec2
}

// 还可以做更细的拆分，把 Matrix 全部放到一个数组里面
// position, position, position .... position
// rotation, rotation, rotation .... rotation
// 把这两个数组申请为一块，然后分两个小组使用
type Transform struct {
	engi.Entity

	// world location
	world SRT
	// relative location to parent
	local SRT

	// graph-link
	parent     uint16
	firstChild uint16
	preSibling uint16
	nxtSibling uint16

	// pointer will cause gc check
	t *TransformTable
}

func (xf *Transform) Position() f32.Vec2 {
	return xf.local.Position
}

func (xf *Transform) Scale() f32.Vec2 {
	return xf.local.Scale
}

func (xf *Transform) Rotation() float32 {
	return xf.local.Rotation
}

func (xf *Transform) Local() SRT {
	return xf.local
}

func (xf *Transform) World() SRT {
	return xf.world
}

// Set local position relative to parent
func (xf *Transform) SetPosition(position f32.Vec2) {
	xf.local.Position = position
	// compute world position
	if xf.parent == none {
		xf.setPosition(nil, position)
	} else {
		xf.setPosition(&xf.t.comps[xf.parent].world, position)
	}
}

func (xf *Transform) MoveBy(dx, dy float32) {
	p := xf.local.Position
	p[0], p[1] = p[0]+dx, p[1]+dy
	xf.SetPosition(p)
}

// update world location: world = parent.world + self.local
func (xf *Transform) setPosition(parent *SRT, local f32.Vec2) {
	p := f32.Vec2{0, 0}
	if parent != nil {
		p = parent.Position
	}
	xf.world.Position[0] = p[0] + local[0]
	xf.world.Position[1] = p[1] + local[1]
	// all child
	for comps, child := xf.t.comps, xf.firstChild; child != none; {
		node := &comps[child]
		child = node.nxtSibling
		node.setPosition(&xf.world, node.local.Position)
	}
}

// apply scale to child
func (xf *Transform) SetScale(scale f32.Vec2) {
	xf.local.Scale = scale
	// compute world scale
	if xf.parent == none {
		xf.setScale(nil, scale)
	} else {
		xf.setScale(&xf.t.comps[xf.parent].world, scale)
	}
}

func (xf *Transform) ScaleBy(dx, dy float32) {
	sk := xf.local.Scale
	sk[0], sk[1] = sk[0]+dx, sk[1]+dy
	xf.SetScale(sk)
}

func (xf *Transform) setScale(parent *SRT, scale f32.Vec2) {
	s := f32.Vec2{1, 1}
	if parent != nil {
		s = parent.Scale
	}
	xf.world.Scale[0] = s[0] * scale[0]
	xf.world.Scale[1] = s[1] * scale[1]

	// all child
	for comps, child := xf.t.comps, xf.firstChild; child != none; {
		node := comps[child]
		child = node.nxtSibling
		node.setPosition(&xf.world, node.local.Position)
	}
}

// apply
func (xf *Transform) SetRotation(rotation float32) {
	xf.local.Rotation = rotation
	// compute world rotation
	if xf.parent == none {
		xf.setRotation(nil, rotation)
	} else {
		xf.setRotation(&xf.t.comps[xf.parent].world, rotation)
	}
}

func (xf *Transform) RotateBy(d float32) {
	r := xf.local.Rotation
	r += d
	xf.SetRotation(r)
}

func (xf *Transform) setRotation(parent *SRT, rotation float32) {
	r := float32(0)
	if parent != nil {
		r = parent.Rotation
	}
	xf.world.Rotation = r + rotation

	// all child
	for comps, child := xf.t.comps, xf.firstChild; child != none; {
		node := comps[child]
		child = node.nxtSibling
		node.setRotation(&xf.world, node.local.Rotation)
	}
}

func (xf *Transform) LinkChildren(list... *Transform) {
	for _, c := range list {
		xf.LinkChild(c)
	}
}

func (xf *Transform) LinkChild(c *Transform) {
	mp, comps := xf.t._map, xf.t.comps
	pi, ci := mp[xf.Entity.Index()], mp[c.Entity.Index()]

	if xf.firstChild == none {
		xf.firstChild = uint16(ci)
		c.parent = uint16(pi)
	} else {
		var prev uint16
		for next := xf.firstChild; next != none; {
			prev = next
			next = comps[next].nxtSibling
		}
		comps[prev].nxtSibling = uint16(ci)
		c.preSibling = prev
		c.parent = uint16(pi)
	}
}

func (xf *Transform) RemoveChild(c *Transform) {
	mp, comps := xf.t._map, xf.t.comps
	pi, ci := uint16(mp[xf.Entity.Index()]), uint16(mp[c.Entity.Index()])

	if c.parent != pi {
		return
	}

	if xf.firstChild == ci {
		xf.firstChild = c.nxtSibling
	} else {
		comps[c.preSibling].nxtSibling = c.nxtSibling
	}
	if nxt := c.nxtSibling; nxt != none {
		comps[nxt].preSibling = c.preSibling
	}
	c.parent, c.preSibling, c.nxtSibling = none, none, none
}

func (xf *Transform) FirstChild() (c *Transform) {
	if first := xf.firstChild; first != none {
		c = &xf.t.comps[first]
	}
	return
}

func (xf *Transform) Parent() (p *Transform) {
	if x := xf.parent; x != none {
		p = &xf.t.comps[x]
	}
	return
}

func (xf *Transform) Sibling() (prev, next *Transform) {
	if x := xf.preSibling; x != none {
		prev = &xf.t.comps[x]
	}
	if x := xf.nxtSibling; x != none {
		next = &xf.t.comps[x]
	}
	return
}

func (xf *Transform) reset() {
	if p := xf.parent; p != none {
		xf.t.comps[p].RemoveChild(xf)
	}
	*xf = Transform{}
}

type TransformTable struct {
 	comps []Transform
	_map  map[uint32]int
	index, cap int
}

func NewTransformTable(cap int) *TransformTable {
	return &TransformTable{
		cap:cap,
		_map:make(map[uint32]int),
		index:1, // skip first
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
	xf.Entity = entity
	xf.local.Scale = f32.Vec2{1, 1}
	xf.world.Scale = f32.Vec2{1, 1}
	xf.t = tt
	tt._map[ei] = tt.index
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
// Delete will unlink the parent-child relation
func (tt *TransformTable) Delete(entity engi.Entity) {
	ei := entity.Index()
	if v, ok := tt._map[ei]; ok {
		if tail := tt.index -1; v != tail && tail > 0 {
			tt.comps[v] = tt.comps[tail]
			tt.relink(uint16(tail), uint16(v))

			// remap index
			tComp := &tt.comps[tail]
			ei := tComp.Entity.Index()
			tt._map[ei] = v
			tt.comps[tail].reset()
		} else {
			tt.comps[tail].reset()
		}
		tt.index -= 1
		delete(tt._map, ei)
	}
}

func (tt *TransformTable) relink(old, new uint16) {
	xf := &tt.comps[old]
	// relink parent
	if p := xf.parent; p != none {
		if pxf := &tt.comps[p]; pxf.firstChild == old {
			pxf.firstChild = new
		} else {
			prev := &tt.comps[xf.preSibling]
			prev.nxtSibling = new
		}
	}
	// relink children
	if child := xf.firstChild; child != none {
		node := &tt.comps[child]
		node.parent = new
	}
}


func (tt *TransformTable) Destroy() {
	tt.comps = make([]Transform, 0)
	tt._map = make(map[uint32]int)
	tt.index = 1
}

func (tt *TransformTable) Size() (size, cap int) {
	return tt.index-1, tt.cap
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
