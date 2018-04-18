package auto

import (
	"korok.io/korok/math"
	"korok.io/korok/gui"

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
// Options is our layout-system

// 当 W = 0, H = 0 的时候，按 WrapContent 的方式绘制
type Element struct{
	id gui.ID
	// 相对偏移 和 大小
	gui.Rect
	// margin
	margin
}

type margin struct {
	Top, Left, Bottom, Right float32
}

type gravity struct {
	X, Y float32
}

type DirtyFlag uint32

const (
	FlagSize DirtyFlag = 1 << iota
	FlagMargin
	FlagGravity
)

// Shadow of current ui-element
type Options struct {
	gui.Rect
	margin
	gravity
	Flag    DirtyFlag // dirty flag
}

// set flag
func (p *Options) Margin(top, left, right, bottom float32) *Options {
	p.Flag |= FlagMargin
	p.margin = margin{top, left, bottom, right}
	return p
}

func (p *Options) Size(w, h float32) *Options {
	p.Flag |= FlagSize
	p.Rect.W = w
	p.Rect.H = h
	return p
}

func (p *Options) Gravity(x, y float32) *Options {
	p.Flag |= FlagGravity
	p.gravity.X = x
	p.gravity.Y = y
	return p
}

type layout struct {
	Cursor struct{X, Y float32}
	
	// ui bound 是一直存储的，记录一些持久化的数据
	uiElements           []Element // element uiElements

	// group 是 fifo 的结构,记录动态的数据
	groupStack           []Group // groupStack uiElements

	// header of group stack
	hGroup *Group

	// default ui-element spacing
	spacing float32
}

func (lyt *layout) Initialize(style *gui.Theme) {
	// init size, todo resize 会导致指针错误
	lyt.uiElements = make([]Element, 0, 32)
	lyt.groupStack = make([]Group, 0, 8)
	lyt.spacing = style.Spacing

	// Create a default layout
	bb := lyt.NewElement(0)
	ii := len(lyt.groupStack)
	lyt.groupStack = append(lyt.groupStack, Group{LayoutType: OverLay, Element: bb})
	lyt.hGroup = &lyt.groupStack[ii]
}

func (lyt *layout) SetDefaultLayoutSize(w, h float32) {
	if dft, ok := lyt.Element(0); ok {
		dft.W, dft.H = w, h
	}
}

// 创建新的Layout
func (lyt *layout) NewElement(id gui.ID) *Element {
	ii := len(lyt.uiElements)
	lyt.uiElements = append(lyt.uiElements, Element{id:id})
	return &lyt.uiElements[ii]
}

// 找出前一帧保存的大小
func (lyt *layout) Element(id gui.ID) (bb *Element, ok bool) {
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

func (lyt *layout) Dump()  {
	log.Println("dump elemnts:", lyt.uiElements)
	log.Println("dump group:", lyt.groupStack)
}

func (lyt *layout) Reset() {
	lyt.uiElements = lyt.uiElements[:0]
}

// Options Operation
func (lyt *layout) Move(x, y float32) *layout {
	lyt.Cursor.X, lyt.Cursor.Y = x, y
	return lyt
}

func (lyt *layout) BoundOf(id gui.ID) (bb Element, ok bool) {
	if size := len(lyt.uiElements); size > int(id) {
		if bb = lyt.uiElements[id]; bb.id == id {
			ok = true
		}
	}
	// 否则进行线性查找, 找出UI边界
	return
}

func (lyt *layout) Offset(dx, dy float32) *layout {
	lyt.Cursor.X += dx
	lyt.Cursor.Y += dy
	return lyt
}

func (lyt *layout) SetGravity(x, y float32) *layout {
	lyt.hGroup.Gravity.X = math.Clamp(x, 0, 1)
	lyt.hGroup.Gravity.Y = math.Clamp(y, 0, 1)
	return lyt
}

func (lyt *layout) SetSize(w, h float32) *layout {
	lyt.hGroup.SetSize(w, h)
	return lyt
}

func (lyt *layout) SetPadding(top, left, right, bottom float32) *layout {
	lyt.hGroup.Padding = gui.Padding{left, right, top, bottom}
	return lyt
}

func (lyt *layout) BeginLayout(id gui.ID, opt *Options, xtype LayoutType) (elem *Element, ok bool) {
	// layout element
	elem, ok = lyt.BeginElement(id, opt)

	// do layout
	ii := len(lyt.groupStack)

	// group-stack has a default parent
	// so it's safe to index
	parent := &lyt.groupStack[ii-1]
	lyt.groupStack = append(lyt.groupStack, Group{LayoutType:xtype, Element: elem})
	lyt.hGroup = &lyt.groupStack[ii]

	// stash Options state
	parent.Cursor.X = lyt.Cursor.X
	parent.Cursor.Y = lyt.Cursor.Y

	// reset Options
	lyt.Cursor.X, lyt.Cursor.Y = 0, 0
	return
}

// PopLayout, resume parent's state
func (lyt *layout) EndLayout() {
	// 1. Set size with external spacing
	v := lyt.hGroup
	if v.fixedWidth > 0 {
		v.W = v.fixedWidth
	} else {
		v.W = v.Size.W
	}
	if v.fixedHeight > 0 {
		v.H = v.fixedHeight
	} else {
		v.H = v.Size.H
	}
	v.W += lyt.spacing
	v.H += lyt.spacing

	// 2. return to parent
	if size := len(lyt.groupStack); size > 1 {
		lyt.groupStack = lyt.groupStack[:size-1]
		lyt.hGroup = &lyt.groupStack[size-2]
	}

	g := lyt.hGroup
	lyt.Cursor.X, lyt.Cursor.Y = g.Cursor.X, g.Cursor.Y

	// 3. end layout
	elem := &Element{Rect:gui.Rect{0, 0, v.W, v.H}}
	lyt.EndElement(elem)
}

func (lyt *layout) BeginElement(id gui.ID, opt *Options) (elem *Element, ok bool){
	if elem, ok = lyt.Element(id); !ok {
		elem = lyt.NewElement(id)
	} else {
		// 计算偏移
		elem.X = lyt.Cursor.X + lyt.spacing
		elem.Y = lyt.Cursor.Y + lyt.spacing

		// gravity
		var (
			group = lyt.hGroup
			gravity = group.Gravity
			extra = struct {W, H float32} {lyt.spacing*2, lyt.spacing*2}
		)

		// Each element's property
		if opt != nil && opt.Flag != 0 {
			// 计算 margin 和 偏移
			if opt.Flag & FlagMargin != 0 {
				elem.margin = opt.margin
				elem.X += elem.Left
				elem.Y += elem.Top

				extra.W += elem.Left + elem.Right
				extra.H += elem.Top + elem.Bottom
			}

			// 计算大小
			if opt.Flag & FlagSize != 0 {
				if opt.W > 0 {
					elem.Rect.W = opt.W
				}
				if opt.H > 0 {
					elem.Rect.H = opt.H
				}
			}

			// Overlap group's gravity
			if opt.Flag & FlagGravity != 0 {
				gravity = opt.gravity
			}

			// Clear flag
			opt.Flag = 0
		}

		switch group.LayoutType {
		case Horizontal:
			elem.Y += (group.H - elem.H - extra.H) * gravity.Y
		case Vertical:
			elem.X += (group.W - elem.W - extra.W) * gravity.X
		case OverLay:
			elem.Y += (group.H - elem.H - extra.H) * gravity.Y
			elem.X += (group.W - elem.W - extra.W) * gravity.X
		}

		elem.X += lyt.hGroup.X
		elem.Y += lyt.hGroup.Y
	}
	return
}

func (lyt *layout) EndElement(elem *Element) {
	if  elem == nil {
		log.Println("====> err nil")
	}
	lyt.Advance(elem)
	lyt.Extend(elem)
}

// 重新计算父容器的大小
// size + margin = BoundingBox
func (lyt *layout) Extend(elem *Element) {
	var (
		g  = lyt.hGroup
		dx = elem.W + elem.Left + elem.Right + lyt.spacing //+ lyt.spacing
		dy = elem.H + elem.Top + elem.Bottom + lyt.spacing //+ lyt.spacing
	)

	switch g.LayoutType {
	case Horizontal:
		// 水平加之，高度取最大
		g.Size.W += dx
		g.Size.H = math.Max(g.Size.H, dy)
	case Vertical:
		// 高度加之，水平取最大
		g.Size.W = math.Max(g.Size.W, dx)
		g.Size.H += dy
	case OverLay:
		// 重叠, 取高或者宽的最大值
		g.Size.W = math.Max(g.Size.W, dx)
		g.Size.H = math.Max(g.Size.H, dy)
	}
}

// 重新计算父容器的光标位置
func (lyt *layout) Advance(elem *Element) {
	var (
		g, c  = lyt.hGroup, &lyt.Cursor
		dx = elem.W + elem.Left + elem.Right + lyt.spacing// + lyt.spacing
		dy = elem.H + elem.Top + elem.Bottom + lyt.spacing// + lyt.spacing
	)

	switch g.LayoutType {
	case Horizontal:
		// 水平步进，前进一个控件宽度
		c.X += dx
	case Vertical:
		// 垂直步进，前进一个控件高度
		c.Y += dy
	case OverLay:
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
	gui.Padding

	// 当前帧布局的计算变量
	Size struct{W, H float32}
	Gravity struct{X, Y float32}

	// true if group has a predefined size
	fixedWidth float32
	fixedHeight float32
}

func (g *Group) SetGravity(x, y float32) {
	g.Gravity.X = math.Clamp(x, 0, 1)
	g.Gravity.Y = math.Clamp(y, 0, 1)
}

func (g *Group) SetPadding(top, left, right, bottom float32) {
	g.Padding = gui.Padding{left, right, top, bottom}
}

func (g *Group) SetSize(w, h float32) {
	g.fixedWidth = w
	g.fixedHeight = h
}