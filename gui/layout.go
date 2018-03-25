package gui

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/math"

	"log"
)

// GUI coordinate system
// (0, 0)
//  +-------------------------+ (w, 0)
//  |  (x,y)                  |
//  |   +------> W            |
//  |   |                     |
//  |   |                     |
//  | H V                     |
//  +-------------------------+
//(0, h)                       (w, h)
// cursor is our layout-system

// 改变布局方向，方便计算从右到左
// 或者从下往上的布局
type Direction int

const (
	Left2Right Direction = iota
	Right2Left
	Top2Bottom
	Bottom2Top
)

// 当 W = 0, H = 0 的时候，按 WrapContent 的方式绘制
type Element struct{
	id ID
	// 相对偏移 和 大小
	Bound
	// Margin
	Margin
}

type Property struct {
	// margin
	MarginLeft float32
	MarginRight float32
	MarginTop float32
	MarginBottom float32

	// Gravity
	GravityH float32
	GravityV float32

	// Element Size
	Width float32
	Height float32
}

// UI绘制边界
type Bound struct {
	X, Y float32
	W, H float32
}

type Margin struct {
	Top, Left, Bottom, Right float32
}

type Gravity struct {
	X, Y float32
}

type DirtyFlag uint32

const (
	FlagSize DirtyFlag = 1 << iota
	FlagMargin
	FlagGravity
)

// Shadow of current ui-element
type cursor struct {
	Bound
	Margin
	Gravity Gravity
	owner ID
	Flag DirtyFlag // dirty flag
}

func (c *cursor) Reset()  {
	c.Margin = Margin{}
}

// set flag
func (c *cursor) SetMargin(top, left, right, bottom float32) *cursor{
	c.Flag |= FlagMargin
	c.Margin = Margin{top, left, bottom, right}
	return c
}

func (c *cursor) SetSize(w, h float32) *cursor{
	c.Flag |= FlagSize
	c.Bound.W = w
	c.Bound.H = h
	return c
}

func (c *cursor) SetGravity(x, y float32) *cursor{
	c.Flag |= FlagGravity
	c.Gravity.X = x
	c.Gravity.Y = y
	return c
}

func (c *cursor) To(id ID) {
	c.owner = id
}

func (b *Bound) Offset(x, y float32) *Bound {
	b.X, b.Y = x, y
	return b
}

func (b *Bound) Size(w, h float32) {
	b.W, b.H = w, h
}

func (b *Bound) SizeAuto() {
	b.W, b.H = 0, 0
}

func (b *Bound) InRange(p f32.Vec2) bool{
	if p[0] < b.X || p[0] > (b.X + b.W) {
		return false
	}
	if p[1] < b.Y || p[1] > (b.Y + b.H) {
		return false
	}
	return true
}

type LayoutManager struct {
	Horizontal, Vertical Direction
	Cursor               cursor
	Align
	// ui bound 是一直存储的，记录一些持久化的数据
	uiElements           []Element // element uiElements

	// group 是 fifo 的结构,记录动态的数据
	groupStack           []Group // groupStack uiElements

	// header of group stack
	hGroup *Group

	// default ui-element spacing
	spacing float32
}

func (lyt *LayoutManager) Initialize(style *Style) {
	// init size, todo resize 会导致指针错误
	lyt.uiElements = make([]Element, 0, 32)
	lyt.groupStack = make([]Group, 0, 8)
	lyt.spacing = style.Spacing

	// Create a default layout
	bb := lyt.NewElement(0)
	ii := len(lyt.groupStack)
	lyt.groupStack = append(lyt.groupStack, Group{LayoutType:LinearOverLay, Element: bb})
	lyt.hGroup = &lyt.groupStack[ii]
}

func (lyt *LayoutManager) SetDefaultLayoutSize(w, h float32) {
	if dft, ok := lyt.Element(0); ok {
		dft.W, dft.H = w, h
	}
}

// 创建新的Layout
func (lyt *LayoutManager) NewElement(id ID) *Element {
	ii := len(lyt.uiElements)
	lyt.uiElements = append(lyt.uiElements, Element{id:id})
	return &lyt.uiElements[ii]
}

// 找出前一帧保存的大小
func (lyt *LayoutManager) Element(id ID) (bb *Element, ok bool) {
	if size := len(lyt.uiElements); size > int(id) {
		if bb = &lyt.uiElements[id]; bb.id == id {
			ok = true; return
		}
	}
	// Linear Search
	for i := range lyt.uiElements {
		if bb = &lyt.uiElements[i]; bb.id == id {
			ok = true; break
		}
	}
	return
}

func (lyt *LayoutManager) Dump()  {
	log.Println("dump elemnts:", lyt.uiElements)
	log.Println("dump group:", lyt.groupStack)
}

func (lyt *LayoutManager) Reset() {
	lyt.uiElements = lyt.uiElements[:0]
}

// Cursor Operation
func (lyt *LayoutManager) Move(x, y float32) *LayoutManager {
	lyt.Cursor.X, lyt.Cursor.Y = x, y
	return lyt
}

func (lyt *LayoutManager) BoundOf(id ID) (bb Element, ok bool) {
	if size := len(lyt.uiElements); size > int(id) {
		if bb = lyt.uiElements[id]; bb.id == id {
			ok = true
		}
	}
	// 否则进行线性查找, 找出UI边界
	return
}

func (lyt *LayoutManager) Offset(dx, dy float32) *LayoutManager {
	lyt.Cursor.X += dx
	lyt.Cursor.Y += dy
	return lyt
}

func (lyt *LayoutManager) SetGravity(x, y float32) *LayoutManager {
	lyt.hGroup.Gravity.X = math.Clamp(x, 0, 1)
	lyt.hGroup.Gravity.Y = math.Clamp(y, 0, 1)
	return lyt
}

func (lyt *LayoutManager) SetSize(w, h float32) *LayoutManager {
	lyt.hGroup.Bound.W = w
	lyt.hGroup.Bound.H = h
	lyt.hGroup.hasSize = true
	return lyt
}

func (lyt *LayoutManager) SetPadding(top, left, right, bottom float32) *LayoutManager{
	lyt.hGroup.Padding = Padding{left, right, top, bottom}
	return lyt
}

func (lyt *LayoutManager) BeginLayout(id ID, xtype LayoutType) (elem *Element, ok bool) {
	// layout element
	elem, ok = lyt.BeginElement(id)

	// do layout
	ii := len(lyt.groupStack)

	// group-stack has a default parent
	// so it's safe to index
	parent := &lyt.groupStack[ii-1]
	lyt.groupStack = append(lyt.groupStack, Group{LayoutType:xtype, Element: elem})
	lyt.hGroup = &lyt.groupStack[ii]

	// stash cursor state
	parent.Cursor.X = lyt.Cursor.X
	parent.Cursor.Y = lyt.Cursor.Y

	// group's (x, y) is absolute coordinate
	// f(x, y) = group(x, y) + cursor(x, y)
	g := lyt.hGroup
	g.X, g.Y = g.X + parent.X, g.Y+parent.Y
	//g.X, g.Y = g.X+lyt.Cursor.X , g.Y+lyt.Cursor.Y

	//log.Println("group's elem:", g.Element)

	// reset cursor
	lyt.Cursor.X, lyt.Cursor.Y = 0, 0
	return
}

// 此处 EndLayout 有点问题，首先要End之前的Element，这样可以得到当前Layout的大小，然后
// 这是和之前 BeginGroup 配对的操作。
// 之后还应该再End一次，
// Group 其实相当于一个元素，对于它自己的Group还要再End一次，这样计算得到
// 1. 结束本次布局的Layout
// 2. 结束自己的父类的Layout
// PopLayout, resume parent's state
func (lyt *LayoutManager) EndLayout() {
	// 1. Set size if not set explicitly, end spacing
	size := lyt.hGroup.Size
	size.W += lyt.spacing
	size.H += lyt.spacing

	if !lyt.hGroup.hasSize || lyt.hGroup.W == 0 {
		lyt.hGroup.W = size.W
	}
	if !lyt.hGroup.hasSize || lyt.hGroup.H == 0 {
		lyt.hGroup.H = size.H
	}

	// 2. return to parent
	if size := len(lyt.groupStack); size > 1 {
		lyt.groupStack = lyt.groupStack[:size-1]
		lyt.hGroup = &lyt.groupStack[size-2]
	}

	g := lyt.hGroup
	lyt.Cursor.X, lyt.Cursor.Y = g.Cursor.X, g.Cursor.Y

	// 3. end layout, remove default spacing
	elem := &Element{Bound:Bound{0, 0, size.W, size.H}}
	lyt.EndElement(elem)

	// 3. 清除当前布局的参数
	lyt.Cursor.Reset()
}

func (lyt *LayoutManager) BeginElement(id ID) (elem *Element, ok bool){
	if elem, ok = lyt.Element(id); !ok {
		elem = lyt.NewElement(id)
	} else {
		// 计算偏移
		elem.X = lyt.Cursor.X + lyt.spacing
		elem.Y = lyt.Cursor.Y + lyt.spacing

		// Gravity
		var (
			group = lyt.hGroup
			gravity = group.Gravity
			extra = struct {
				W, H float32
			}{lyt.spacing*2, lyt.spacing*2}
		)

		// Each element's property
		if lyt.Cursor.owner == id {
			// 计算 Margin 和 偏移
			if lyt.Cursor.Flag & FlagMargin != 0 {
				elem.Margin = lyt.Cursor.Margin
				elem.X += elem.Left
				elem.Y += elem.Top

				extra.W += elem.Left + elem.Right
				extra.H += elem.Top + elem.Bottom
			}

			// 计算大小
			if lyt.Cursor.Flag & FlagSize != 0 {
				elem.Bound.W = lyt.Cursor.W
				elem.Bound.H = lyt.Cursor.H
			}

			// Overlap group's gravity
			if lyt.Cursor.Flag & FlagGravity != 0 {
				gravity = lyt.Cursor.Gravity
			}

			// 清空标记
			lyt.Cursor.owner = -1
			lyt.Cursor.Flag = 0
		}

		switch group.LayoutType {
		case LinearHorizontal:
			elem.Y += (group.H - elem.H - extra.H) * gravity.Y
		case LinearVertical:
			elem.X += (group.W - elem.W - extra.W) * gravity.X
		case LinearOverLay:
			elem.Y += (group.H - elem.H - extra.H) * gravity.Y
			elem.X += (group.W - elem.W - extra.W) * gravity.X
		}
	}
	return
}

func (lyt *LayoutManager) EndElement(elem *Element) {
	lyt.Advance(elem)
	lyt.Extend(elem)
}

// 重新计算父容器的大小
// size + margin = BoundingBox
func (lyt *LayoutManager) Extend(elem *Element) {
	var (
		g  = lyt.hGroup
		dx = elem.W + elem.Left + elem.Right + lyt.spacing //+ lyt.spacing
		dy = elem.H + elem.Top + elem.Bottom + lyt.spacing //+ lyt.spacing
	)

	switch g.LayoutType {
	case LinearHorizontal:
		// 水平加之，高度取最大
		g.Size.W += dx
		g.Size.H = math.Max(g.Size.H, dy)
	case LinearVertical:
		// 高度加之，水平取最大
		g.Size.W = math.Max(g.Size.W, dx)
		g.Size.H += dy
	case LinearOverLay:
		// 重叠, 取高或者宽的最大值
		g.Size.W = math.Max(g.Size.W, dx)
		g.Size.H = math.Max(g.Size.H, dy)
	}
}

// 重新计算父容器的光标位置
func (lyt *LayoutManager) Advance(elem *Element) {
	var (
		g, c  = lyt.hGroup, &lyt.Cursor
		dx = elem.W + elem.Left + elem.Right + lyt.spacing// + lyt.spacing
		dy = elem.H + elem.Top + elem.Bottom + lyt.spacing// + lyt.spacing
	)

	switch g.LayoutType {
	case LinearHorizontal:
		// 水平步进，前进一个控件宽度
		c.X += dx
	case LinearVertical:
		// 垂直步进，前进一个控件高度
		c.Y += dy
	case LinearOverLay:
		// 保持原来的位置不变..
	}
}

// Q. 当前 Group 的 X，Y, W, H 应该和 Group 的Cursor区分开来

type Flag uint32

type Group struct {
	LayoutType; Flag
	*Element
	// 仅用来缓存...
	Cursor struct{X, Y float32}
	Offset struct{X, Y float32}
	Padding

	// 当前帧布局的计算变量
	Size struct{W, H float32}
	Gravity struct{X, Y float32}

	// true if group has a predefined size
	hasSize bool
}