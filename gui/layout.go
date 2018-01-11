package gui

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
type Bound struct{
	X, Y float32
	W, H float32
}

type Layout struct {
	Horizontal, Vertical Direction
	Bound
	Align
	stack []Bound
}

// Stack Operation
func (lyt *Layout) Begin(bb *Bound) {
	if bb != nil {
		lyt.stack = append(lyt.stack, *bb)
	} else {
		lyt.stack = append(lyt.stack, Bound{0, 0, 0, 0})
	}
	lyt.Align = AlignCenter
	lyt.Horizontal, lyt.Vertical = Left2Right, Top2Bottom
}

func (lyt *Layout) Push(bb *Bound) (id int) {
	id = len(lyt.stack)
	lyt.stack = append(lyt.stack, *bb)
	return
}

func (lyt *Layout) Pop() (bb Bound) {
	lyt.stack = lyt.stack[:len(lyt.stack)-1]
	return
}

func (lyt *Layout) Reset() {
	lyt.stack = lyt.stack[:0]
}

// Bound Operation
func (lyt *Layout) Move(x, y float32) *Layout {
	lyt.X, lyt.Y = x, y
	return lyt
}

func (lyt *Layout) BelowOf(id int) *Layout {
	if id < len(lyt.stack) && id >= 0{
		bb := &lyt.stack[id]
		lyt.Bound.Y = bb.Y + bb.H
	}
	return lyt
}

func (lyt *Layout) TopOf(id int) *Layout {
	if id < len(lyt.stack) && id >= 0 {
		bb := &lyt.stack[id]
		lyt.Bound.Y = bb.Y - lyt.Bound.H
	}
	return lyt
}

func (lyt *Layout) LeftOf(id int) *Layout {
	if id < len(lyt.stack) && id >= 0 {
		bb := &lyt.stack[id]
		lyt.Bound.X = bb.X - lyt.Bound.W
	}
	return lyt
}

func (lyt *Layout) AlignTopOf(id int) *Layout {
	if id < len(lyt.stack) && id >= 0 {
		bb := &lyt.stack[id]
		lyt.Bound.Y = bb.Y
	}
	return lyt
}

func (lyt *Layout) AlignLeftOf(id int) *Layout {
	if id < len(lyt.stack) && id >= 0 {
		bb := &lyt.stack[id]
		lyt.Bound.X = bb.X
	}
	return lyt
}

func (lyt *Layout) RightOf(id int) *Layout {
	if id < len(lyt.stack) && id >= 0 {
		bb := &lyt.stack[id]
		lyt.Bound.X = bb.X+bb.W
	}
	return lyt
}

func (lyt *Layout) RatioH(id int, num, i, used int) *Layout {
	if id < len(lyt.stack) && id >= 0 {
		bb := &lyt.stack[id]
		col := bb.W/float32(num)
		lyt.Bound.X = bb.X + col * float32(i)
		lyt.Bound.W = col * float32(used)
	}
	return lyt
}

func (lyt *Layout) RatioV(id int, num, i, used int) *Layout {
	if id < len(lyt.stack) && id >= 0 {
		bb := &lyt.stack[id]
		col := bb.H/float32(num)
		lyt.Bound.Y = bb.Y + col * float32(i)
		lyt.Bound.H = col * float32(used)
	}
	return lyt
}

func (lyt *Layout) Offset(dx, dy float32) *Layout {
	lyt.Bound.X += dx
	lyt.Bound.Y += dy
	return lyt
}

func (lyt *Layout) Size(w, h float32) *Layout{
	lyt.W, lyt.H = w, h
	return lyt
}

func (lyt *Layout) SizeAuto() *Layout {
	lyt.W, lyt.H = 0, 0
	return lyt
}

func (lyt *Layout) Flow(h, v Direction) {
	lyt.Horizontal, lyt.Vertical = h, v
}



