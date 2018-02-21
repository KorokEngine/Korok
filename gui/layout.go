package gui

import (
	"github.com/go-gl/mathgl/mgl32"
	"korok.io/korok/engi/math"

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

// UI绘制边界
type Bound struct {
	X, Y float32
	W, H float32
}

type Margin struct {
	Top, Left, Bottom, Right float32
}

// Shadow of current ui-element
type Cursor struct {
	Bound
	Margin
	Flag uint32 // dirty flag
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

func (b *Bound) InRange(p mgl32.Vec2) bool{
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
	Cursor               Cursor
	Align
	// ui bound 是一直存储的，记录一些持久化的数据
	uiElements           []Element // element uiElements

	// group 是 fifo 的结构,记录动态的数据
	groupStack           []Group // groupStack uiElements

	// header of group stack
	hGroup *Group
}

func (lyt *LayoutManager) Initialize() {
	// init size, todo resize 会导致指针错误
	lyt.uiElements = make([]Element, 0, 32)
	lyt.groupStack = make([]Group, 0, 8)

	// Create a default layout
	bb := lyt.NewElement(0)
	ii := len(lyt.groupStack)
	lyt.groupStack = append(lyt.groupStack, Group{LayoutType:LinearOverLay, Element: bb})
	lyt.hGroup = &lyt.groupStack[ii]
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

	lyt.Cursor.Margin.Left = dx
	lyt.Cursor.Margin.Top  = dy
	return lyt
}

func (lyt *LayoutManager) Size(w, h float32) *LayoutManager {
	lyt.Cursor.W, lyt.Cursor.H = w, h
	return lyt
}

func (lyt *LayoutManager) SizeAuto() *LayoutManager {
	lyt.Cursor.W, lyt.Cursor.H = 0, 0
	return lyt
}

func (lyt *LayoutManager) Flow(h, v Direction) {
	lyt.Horizontal, lyt.Vertical = h, v
}

// AutoLayout System
func (lyt *LayoutManager) NewLayout(id ID, xtype LayoutType) {
	bb := lyt.NewElement(id)
	ii := len(lyt.groupStack)
	// 默认的布局没有 parent
	parent := &lyt.groupStack[ii-1]

	lyt.groupStack = append(lyt.groupStack, Group{LayoutType:xtype, Element: bb})
	lyt.hGroup = &lyt.groupStack[ii]

	// 启动布局之后，应该重置光标, 此时光标的位置全部相对于布局
	// 保存父容器的当前光标位置， 当前好像没设计这个变量
	parent.Cursor.X = lyt.Cursor.X
	parent.Cursor.Y = lyt.Cursor.Y

	if parent.id == 0 {
		log.Println("parent cursor:", parent.Cursor)
	}

	// 在 EndLayout的时候再恢复Cursor位置

	// 此时已经可以得出 当前Layout的位置了
	// 好像光标的位置就是 Group 的位置
	// 那么其实完全不需要 Cursor的存在！！
	//
	// Q. 目前的布局使计算父容器的坐标的时候，仅仅考虑父容器的坐标
	// 其实父容器还有自己的父容器，需要把所有父容器的本身相对其父容器的偏移加起来才是对的！！
	g := lyt.hGroup

	// 先加上父父容器
	g.X, g.Y = parent.X, parent.Y
	// 再加上当前父容器的偏移, 这样得到的容器坐标就是绝对的坐标
	g.X, g.Y = g.X+lyt.Cursor.X , g.Y+lyt.Cursor.Y

	// 记录光标的偏移...
	// EndLayout的时候还需要这个数据..
	g.Offset.X = lyt.Cursor.Left
	g.Offset.Y = lyt.Cursor.Top

	lyt.Cursor.X, lyt.Cursor.Y = 0, 0
}

func (lyt *LayoutManager) FindLayout(id ID) (bb *Element, ok bool) {
	return lyt.Element(id)
}

// set as current layout
// 这个方法的逻辑和 NewLayout 都是重复的
func (lyt *LayoutManager) PushLayout(xtype LayoutType, bb *Element) {
	ii := len(lyt.groupStack)
	// 默认的布局没有 parent
	// 此处可以安全的取出 parent，因为初始化的时候会先创建出 parent 布局..
	parent := &lyt.groupStack[ii-1]
	lyt.groupStack = append(lyt.groupStack, Group{LayoutType:xtype, Element: bb})
	lyt.hGroup = &lyt.groupStack[ii]

	// 启动布局之后，应该重置光标, 此时光标的位置全部相对于布局
	// 保存父容器的当前光标位置， 当前好像没设计这个变量
	parent.Cursor.X = lyt.Cursor.X
	parent.Cursor.Y = lyt.Cursor.Y

	// 在 EndLayout的时候再恢复Cursor位置

	// 此时已经可以得出 当前Layout的位置了
	// 好像光标的位置就是 Group 的位置
	// 那么其实完全不需要 Cursor的存在！！
	//
	// Q. 目前的布局使计算父容器的坐标的时候，仅仅考虑父容器的坐标
	// 其实父容器还有自己的父容器，需要把所有父容器的本身相对其父容器的偏移加起来才是对的！！
	g := lyt.hGroup

	// 先加上父父容器
	g.X, g.Y = parent.X, parent.Y

	// 再加上当前父容器的偏移, 这样得到的容器坐标就是绝对的坐标
	g.X, g.Y = g.X+lyt.Cursor.X , g.Y+lyt.Cursor.Y


	lyt.Cursor.X, lyt.Cursor.Y = 0, 0
}

// 重新计算父容器的大小
// size + margin = BoundingBox
func (lyt *LayoutManager) Extend(elem *Element) {
	if g := lyt.hGroup; g != nil {
		switch g.LayoutType {
		case LinearHorizontal:
			// 水平加之，高度取最大
			g.Size.W += elem.W + elem.Left
			g.Size.H = math.Max(g.Size.H, elem.H+elem.Top)
		case LinearVertical:
			// 高度加之，水平取最大
			g.Size.W = math.Max(g.Size.W, elem.W+elem.Left)
			g.Size.H += elem.H + elem.Left
		case LinearOverLay:
			// 重叠, 取高或者宽的最大值
			g.Size.W = math.Max(g.Size.W, elem.W+elem.Left)
			g.Size.H = math.Max(g.Size.H, elem.H+elem.Top)
		}
	}
}

// 重新计算父容器的光标位置
func (lyt *LayoutManager) Advance(elem *Element) {
	if g := lyt.hGroup; g != nil  {
		cursor := &lyt.Cursor
		switch g.LayoutType {
		case LinearHorizontal:
			// 水平步进，前进一个控件宽度
			cursor.X += elem.W
		case LinearVertical:
			// 垂直步进，前进一个控件高度
			cursor.Y += elem.H
		case LinearOverLay:
			// 保持原来的位置不变..
		}
	}
}

func (lyt *LayoutManager) NextBound(sq int) Element {
	return lyt.uiElements[sq]
}

// 需要重新调整光标位置，
// 如果上一层有 Layout，则应该按照该Layout的布局方式移动光标
// 如果没有，则回到Group开始的位置
// 可以在 Group 里面插入一个 默认的 Group 这样永远有一个 Group 存在
func (lyt *LayoutManager) EndLayout() {
	// 1. 计算并更新Layout的大小
	// g := lyt.hGroup
	// 在计算 UIElement 的时候已经计算过了（或许应该放在此处计算）
	size := lyt.hGroup.Size
	// 记录这一帧的值，清空计数器
	// 或者应该把计数器和帧值区分开来
	lyt.hGroup.W = size.W
	lyt.hGroup.H = size.H

	// 2. Layout出栈, 此时可以保证此处总是大于1，否则应该报错
	if size := len(lyt.groupStack); size > 1 {
		lyt.groupStack = lyt.groupStack[:size-1]
		lyt.hGroup = &lyt.groupStack[size-2]
	}

	// 3. 处理父容器的布局,
	// 先回到 Group 的位置（因为进入子容器的时候重置了Cursor，所以要先回复）, 然后才可以布局
	g := lyt.hGroup
	lyt.Cursor.X, lyt.Cursor.Y = g.Cursor.X, g.Cursor.Y

	elem := &Element{Bound:Bound{0, 0, size.W, size.H}}
	// 2. 扩展父容器大小
	lyt.Extend(elem)

	// 2. 接下来直接用 Advance 步进光标
	lyt.Advance(elem)
}

// Q. 当前 Group 的 X，Y, W, H 应该和 Group 的Cursor区分开来

type Flag uint32

type Group struct {
	LayoutType; Flag
	*Element
	// 仅用来缓存...
	Cursor struct{X, Y float32}
	Offset struct{X, Y float32}

	// 当前帧布局的计算变量
	Size struct{W, H float32}
}

// 布局算法：
// 开始布局, 查找 UIElement 的大小，如果
// 有大小则计算布局
//    布局Group + Cursor + Layout-Weight
// 否则创建一个新的UIElement, 如果是 Group
// 则在EndLayout的时候计算出代销，如果Element
// 则直接算出Layout