package gui

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok.io/korok/gfx"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/hid/input"

	"log"
)

type EventType uint8
const (
	EventNone EventType = iota
	EventDown
	EventUp
	EventStartDrag
	EventEndDrag
	EventDragging
)

// 一个Context维护UI层逻辑:
// 1. 一个 DrawList，负责生成顶点
// 2. 共享的 Style，指针引用，多个 Context 之间可以共享
// 3. State 维护，负责维护当前的 UI 状态，比如动画，按钮，列表...
// 窗口在 Context 之上维护，默认的Context可以随意的绘制
// 在窗口内绘制的UI会受到窗口的管理.
type Context struct {
	DrawList
	Layout LayoutManager
	*Style

	// ui global state
	state struct{
		hot, active int
		mouseX, mouseY, mouseDown int

		// drag state
		draggingPointer int
		draggingStart mgl32.Vec2

		isLastEventPointerType bool
		pointerCapture int
	}

	// sqNum should be same for  layout and drawing
	sqNum int
}

func NewContext(style *Style) *Context {
	c := &Context{
		Style:style,
	}
	c.state.draggingPointer = -1
	c.state.isLastEventPointerType = false
	c.state.pointerCapture = -1
	c.DrawList.Initialize()
	c.Layout.Initialize()
	return c
}

// Widgets: Text
// 渲染阶段，取出大小，计算位置开始渲染
// Layout: 1 2 3 4 ....
// Draw:   ^ ^ ^ ^ ....
// 布局阶段仅计算出大小
//
func (ctx *Context) Text(id ID, text string, style *TextStyle)  *Bound {
	lm := &ctx.Layout
	if bb, ok := lm.Element(id); ok {
		g := lm.hGroup


		// x, y = margin + group + group-offset
		// todo 暂时不处理 margin 的问题，现在 只能在 bound 的xy字段记录
		// todo top/left margin 无法记录 bottom/right 的margin
		x,y := Gui2Game(g.X + lm.Bound.X, g.Y + lm.Bound.Y)


		var (
			font = style.Font
			fontSize = style.Size
			color = style.Color
			wrapWidth = bb.W + 10
		)
		ctx.DrawList.AddText(mgl32.Vec2{x, y}, text, font, fontSize, color, wrapWidth)

		// 每绘制完一个元素都要偏移一下光标
		ctx.Layout.Advance(bb)
		lm.Extend(bb)

		return bb
	} else {
		bb := lm.NewElement(id)
		size := ctx.CalcTextSize(text,0, style.Font, style.Size)

		// 加入 Cursor的偏移位置
		// 这样不太好
		bb.W, bb.H = size[0], size[1]
		bb.X, bb.Y = lm.Bound.Offset.X, lm.Bound.Offset.Y

		lm.Extend(bb)
		return bb
	}
}

// 1. LayoutManager pass - 计算出各个控件的大小..
func (ctx *Context) TextLayout(id ID, text string, style *TextStyle) {

	// 计算大小
	size := ctx.CalcTextSize(text,0, style.Font, style.Size)
	bb := ctx.Layout.Bound
	bb.W, bb.H = size[0], size[1]

	// 计算绝对位置
	//g := ctx.Layout.Group()
	//bb.X, bb.Y = g.X, g.Y
}

func (ctx *Context) CalcTextSize(text string, wrapWidth float32, font gfx.FontSystem, fontSize float32) mgl32.Vec2 {
	fr := &FontRender{
		font: font,
		fontSize:fontSize,
	}
	return fr.CalculateTextSize1(text)
}

// Widgets: InputEditor
func (ctx *Context) InputText(hint string, lyt LayoutManager, style *InputStyle) {

}

// Widget: Image
func (ctx *Context) Image(texId uint16, uv mgl32.Vec4, style *ImageStyle) (id int) {
	bound := ctx.Layout.Bound
	min := mgl32.Vec2{bound.X, bound.Y}
	if bound.W == 0 {
		if ok, tex := bk.R.Texture(texId); ok {
			bound.W, bound.H = tex.Width, tex.Height
		}
	}
	max := min.Add(mgl32.Vec2{bound.W, bound.H})
	var color uint32
	if style != nil {
		color = style.TintColor
	} else {
		color = ctx.Style.Image.TintColor
	}
	min[0], min[1] = Gui2Game(min[0], min[1])
	max[0], max[1] = Gui2Game(max[0], max[1])
	ctx.DrawList.AddImage(texId, min, max, mgl32.Vec2{uv[0], uv[1]}, mgl32.Vec2{uv[2], uv[3]}, color)
	return
}

// Widget: Button
func (ctx *Context) Button(text string, style *ButtonStyle) (id int, event EventType) {
	// Render Button
	var color uint32
	var rounding float32
	if style == nil {
		style = &ctx.Style.Button
		color = ctx.Style.Button.Color
		rounding = ctx.Style.Button.Rounding
	} else {
		color = style.Color
		rounding = style.Rounding
	}

	textStyle := ctx.Style.Text
	textSize := ctx.CalcTextSize(text, 0, textStyle.Font, textStyle.Size)
	bound    := ctx.Layout.Bound

	x, y := bound.X, bound.Y
	w, h := textSize[0]+20, textSize[1]+20

	// Check Event
	event = ctx.CheckEvent(&Bound{0, x, y, w, h}, false)

	// Render Frame
	ctx.renderFrame(x, y, w, h, color, rounding)

	// Render Text
	x = bound.X + w/2 - textSize[0]/2
	y = bound.Y + h/2 - textSize[1]/2
	ctx.renderTextClipped(text, &Bound{0, x, y, 0, 0}, &style.TextStyle)

	// push bound todo refactor button bounds
	//bound.W, bound.H = w, h
	// id = ctx.Layout.Push(&bound)
	return
}

func (ctx *Context) renderTextClipped(text string, bb *Bound, style *TextStyle) {
	x, y := Gui2Game(bb.X, bb.Y)
	font := ctx.Style.Text.Font
	if bb.W == 0 {
		ctx.DrawList.AddText(mgl32.Vec2{x, y}, text, font, 12, 0xFF000000, 0)
	} else {
		ctx.DrawList.AddText(mgl32.Vec2{x, y}, text, font, 12, 0xFF000000, bb.W)
	}
}

// 偷师 flat-ui 中的设计，把空间的前景和背景分离，背景单独根据事件来变化..
// 在 Android 中，Widget的前景和背景都可以根据控件状态发生变化
// 但是在大部分UI中，比如 Text/Image 只会改变背景的状态
// 偷懒的自定义UI，不做任何状态的改变... 所以说呢, 我们也采用偷懒的做法呗。。
func (ctx *Context) EventBackground(event EventType) {

}

// 现在只检测一个点, 通常是鼠标的左键或者是多点触控时的第一个手指的位置
// 这样可以记录当前控件的状态...
// 如何根据状态绘制？？
func (ctx *Context) CheckEvent(bound *Bound, checkDragOnly bool) EventType {
	event := EventNone

	// log.Println("check scope bouds:", bound, " point:",  input.PointerPosition(0).MousePos)

	if p := input.PointerPosition(0); bound.InRange(p.MousePos) || ctx.state.pointerCapture == 123{
		btn := input.PointerButton(0)
		id := int(0) // // todo 设计ID系统，记录每个按键的位置..

		// todo finger id!!!
		// 进入此处的条件为：有正确的手指头，手指头说明已经处于 drag 状态！！
		// 这个地方只有 drag 能进
		if ctx.state.draggingPointer == 0 && !btn.JustPressed() {
			// The pointer is in drag operation
			if btn.JustReleased() {
				event = EventEndDrag
				log.Println("drag end real..", event)
				ctx.state.draggingPointer = -1
				ctx.state.draggingStart = mgl32.Vec2{}
			} else if btn.Down() {
				event = EventDragging
			} else {
				ctx.state.draggingPointer = -1
			}
		} else {
			if !checkDragOnly {
				if btn.JustPressed() {
					ctx.state.active = id
					event = EventDown
				}
				if btn.JustReleased() && ctx.state.active == id {
					event = EventUp
				}
			}

			// Check for drag events

			// 1. 此时未必是 drag，但是记下位置，接下来根据位移来得到是否是drag操作
			if btn.JustPressed() {
				ctx.state.draggingStart = p.MousePos
				log.Println("just pressed!!")
			}
			// 2. 如果接下来的移动超出了阈值，那么判断为 drag 操作
			// 超出初始按下位置周围20像素，认为是 drag 操作, 否则可能依然是点击（只是手滑了一下）
			if btn.Down() && bound.InRange(ctx.state.draggingStart) {
				startPosition := ctx.state.draggingStart
				dragThreshHold := float32(10)
				bb := Bound{0,startPosition[0]-dragThreshHold, startPosition[1]-dragThreshHold, 8, 8}
				//
				if !bb.InRange(p.MousePos) {
					event = EventStartDrag
					ctx.state.draggingStart = p.MousePos
					ctx.state.draggingPointer = 0

					log.Println("drag start real ..")
				}
			}

			if event > 0 {
				ctx.state.isLastEventPointerType = true
			}
		}

	}
	return event
}

func (ctx *Context) ImageButton(texId uint16, lyt LayoutManager, style *ImageButtonStyle) EventType{
	return EventNone
}

func (ctx *Context) Rect(w, h float32, style *RectStyle) (id int){
	bb := ctx.Layout.Bound

	x, y := Gui2Game(bb.X, bb.Y)

	var min, max mgl32.Vec2

	if ctx.Layout.Horizontal == Left2Right {
		min[0], max[0] = x, x+w
	} else {
		min[0], max[0] = x-w, x
	}

	if ctx.Layout.Vertical == Top2Bottom {
		min[1], max[1] = y-h, y
	} else {
		min[1], max[1] = y, y+h
	}

	bb.W, bb.H = w, h

	if style == nil {
		style = &ctx.Style.Rect
	} // todo
	if style.FillColor > 0 {
		ctx.DrawList.AddRectFilled(min, max, style.FillColor, style.Rounding, style.Corner)
	} else {
		ctx.DrawList.AddRect(min, max, style.StrokeColor, style.Rounding, style.Corner, style.Stroke)
	}
	return
}

// Render a rectangle shaped with optional rounding and borders(no border!) TODO
func (ctx *Context) renderFrame(x, y, w, h float32, fill uint32, rounding float32) {
	x, y = Gui2Game(x, y)
	min := mgl32.Vec2{x, y-h}
	max := mgl32.Vec2{x+w, y}

	// draw a filled rect
	ctx.DrawList.AddRectFilled(min, max, fill, rounding, FlagCornerAll)
	// border ? I don't think it's a good idea

	//log.Println("renderFrame:",min, max, fill, rounding)
}

// draw line rect
func (ctx *Context) RenderBorder(x, y, w, h float32, color uint32) {
	x, y = Gui2Game(x, y)
	min := mgl32.Vec2{x, y-h}
	max := mgl32.Vec2{x+w, y}

	ctx.DrawList.AddRect(min, max, color, 0, 0, 0.5)
}

// Slider 的绘制很简单，分别绘制滑动条和把手即可
// 难点在于跟踪把手的滑动距离
// Slider的风格，没有想好怎么控制，暂时使用两张图片
// 分别绘制Bar和Knob
// Slider 需要保存混动的结果，否则
func (ctx *Context) Slider(value float32, style *SliderStyle) (v1 float32){
	bb := ctx.Layout.Bound

	x, y := Gui2Game(bb.X, bb.Y)
	w, h := bb.W, bb.H

	//var bar, knob uint16
	//if style == nil {
	//	bar, knob = ctx.Style.Slider.Bar, ctx.Style.Slider.Knob
	//} else {
	//	bar, knob = style.Bar, style.Knob
	//}
	//
	//min, max := mgl32.Vec2{x, y}, mgl32.Vec2{x+w, y+h}
	//ctx.DrawList.AddImageNinePatch(bar, min, max, mgl32.Vec2{0, 0}, mgl32.Vec2{1, 1}, mgl32.Vec4{.5, .5, .5, .5}, 0xFFFFFFFF)
	//
	//
	//ctx.DrawList.AddImage()
	// Slider 需要设定一些自定义的属性，目前没有想好如何实现，先把逻辑实现了
	// 用两种颜色来绘制

	// 说明滑动了，那么应该使用最新的值，而不是传入的值
	if v := ctx.checkSlider(&bb.Bound); v != 0 {
		value = v
		v1 = v
	}

	min, max := mgl32.Vec2{x, y-h}, mgl32.Vec2{x+w, y}
	ctx.DrawList.AddRectFilled(min, max, 0xFFCDCDCD, 5, FlagCornerAll)

	centre := mgl32.Vec2{x+w*value, y-h/2}
	ctx.DrawList.AddCircleFilled(centre, 10, 0xFFABABAB, 12)

	return
}

// Scroll 效果的关键是使用裁切限制滚动区域，然后
// 通过计算拖拽，来得到争取的偏移
func (ctx *Context) StartScroll(size, offset mgl32.Vec2) {
	event := ctx.CheckEvent(nil, false)

	if event == EventStartDrag {
		ctx.capturePoint()
	} else if event == EventEndDrag {
		ctx.releasePointer()
	}
	// 好像算法也不是很难
}

func (ctx *Context) EndScroll() {
	//
}



// 这里的实现基于拖拽的实现，所以
// 只要正确的实现了拖拽，这里的就可以很容易的实现
func (ctx *Context) checkSlider(bound *Bound) (v float32){
	event := ctx.CheckEvent(bound, false)

	// log.Println("returrn evnt:", event == EventEndDrag)

	if event == EventStartDrag {
		ctx.state.pointerCapture = 123

		log.Println("start drag..")
	}

	if event == EventEndDrag {
		ctx.state.pointerCapture = -1

		log.Println("end drag..")
	}

	if ctx.state.isLastEventPointerType {
		// Update the knob position
		if (event == EventDragging) || (event == EventDown) {
			// default = Horizontal
			p1 := input.PointerPosition(0).MousePos[0]
			p0 := bound.X
			v = (p1 - p0)/bound.W

			if v > 1 {
				v = 1
			}
			if v < 0 {
				v = 0
			}
		}
	}
	return
}

func (ctx *Context) endSlider() {

}

func (ctx *Context) capturePoint() {

}

func (ctx *Context) releasePointer() {

}

func (ctx *Context) isLastEventPointerType() bool {
	return true
}

// Container: Window/PopupWindow/
func (ctx *Context) BeginGroup() {

}

func (ctx *Context) EndGroup() {

}

// Layout
// Layout 的时候清空 offset，这样避免把上一个 Layout 的 offset 代入到下一个布局
func (ctx *Context) BeginLayout(id ID, xtype LayoutType) {
	lm := &ctx.Layout
	lm.Bound.Offset.X = 0
	lm.Bound.Offset.Y = 0

	if bb, ok := lm.FindLayout(id); ok {
		var (
			x = lm.hGroup.X + ctx.Layout.Bound.X
			y = lm.hGroup.Y + ctx.Layout.Bound.Y
		)

		// debug draw - render group frame
		ctx.RenderBorder(x, y, bb.W, bb.H, 0xFF00FF00)
		log.Println("render border:", x, y, bb.W, bb.H)


		bb.X, bb.Y = lm.Bound.X, lm.Bound.Y
		lm.PushLayout(xtype, bb)
	} else {
		lm.NewLayout(id, xtype)
	}
}

func (ctx *Context) EndLayout() {
	ctx.Layout.EndLayout()

	//log.Println("end layout:", len(ctx.Layout.groupStack))
}

// Reference System: VirtualBounds
func (ctx *Context) PushVBounds(bounds mgl32.Vec4) {

}

// Clip:
func (ctx *Context) PushClipRect(minClip, maxClip mgl32.Vec2, intersectCurrent bool) {

}

// Theme:
func (ctx *Context) UseTheme(style *Style) {
	ctx.Style = style
}

func Gui2Game(x, y float32) (x1, y1 float32) {
	return x, screen.Height - y
}

func Game2Gui(x, y float32) (x1, y1 float32) {
	return x, screen.Height - y
}

var screen struct{
	Width, Height float32
}

func init() {
	screen.Width = 480
	screen.Height = 320
}