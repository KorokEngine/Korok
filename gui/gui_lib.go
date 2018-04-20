package gui

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/hid/input"
	"korok.io/korok/gfx/dbg"
	"korok.io/korok/gfx/font"

	"log"
	"fmt"
)

type EventType uint8

const (
	EventWentDown  EventType = 1 << iota
	EventWentUp
	EventDown
	EventStartDrag
	EventEndDrag
	EventDragging
)

const EventNone = EventType(0)

func (et EventType) JustPressed() bool {
	return (et & EventWentDown) != 0
}

func (et EventType) JustReleased() bool {
	return (et & EventWentUp) != 0
}

func (et EventType) Down() bool {
	return (et & EventDown) != 0
}

func (et EventType) StartDrag() bool {
	return (et & EventStartDrag) != 0
}

func (et EventType) EndDrag() bool {
	return (et & EventEndDrag) != 0
}

func (et EventType) Dragging() bool {
	return (et & EventDragging) != 0
}

// UI绘制边界
type Rect struct {
	X, Y float32
	W, H float32
}

func (b *Rect) Offset(x, y float32) *Rect {
	b.X, b.Y = x, y
	return b
}

func (b *Rect) Size(w, h float32) {
	b.W, b.H = w, h
}

func (b *Rect) SizeAuto() {
	b.W, b.H = 0, 0
}

func (b *Rect) InRange(p f32.Vec2) bool{
	if p[0] < b.X || p[0] > (b.X + b.W) {
		return false
	}
	if p[1] < b.Y || p[1] > (b.Y + b.H) {
		return false
	}
	return true
}

type Cursor struct {
	X, Y float32
}

// 一个Context维护UI层逻辑:
// 1. 一个 DrawList，负责生成顶点
// 2. 共享的 Theme，指针引用，多个 Context 之间可以共享
// 3. State 维护，负责维护当前的 UI 状态，比如动画，按钮，列表...
// 窗口在 Context 之上维护，默认的Context可以随意的绘制
// 在窗口内绘制的UI会受到窗口的管理.
type Context struct {
	DrawList
	*Theme
	Cursor

	// ui global state
	state struct{
		hot, active ID
		mouseX, mouseY, mouseDown int

		// drag state
		draggingPointer ID
		draggingStart f32.Vec2

		isLastEventPointerType bool
		pointerCapture ID
	}

	// sqNum should be same for  layout and drawing
	sqNum int
}

func NewContext(style *Theme) *Context {
	c := &Context{
		Theme: style,
	}
	c.state.draggingPointer = -1
	c.state.isLastEventPointerType = false
	c.state.pointerCapture = -1
	c.DrawList.Initialize()
	return c
}

func (ctx *Context) Text(id ID, bb *Rect, text string, style *TextStyle) {
	if bb.W != 0 {
		ctx.DrawText(bb, text, style)
	} else {
		sz := ctx.CalcTextSize(text, 0, style.Font, style.Size)
		bb.W = sz[0]
		bb.H = sz[1]

		ctx.DrawText(bb, text, style)
	}
	return
}

// Widgets: InputEditor
func (ctx *Context) InputText(id ID, bb *Rect, hint string, style *InputStyle) {

}

// Widget: Image
func (ctx *Context) Image(id ID, bb *Rect, tex gfx.Tex2D, style *ImageStyle) {
	ctx.DrawImage(bb, tex, style)
}

// Widget: Button
func (ctx *Context) Button(id ID, bb *Rect, text string, style *ButtonStyle) (event EventType) {
	if style == nil {
		style = &ThemeLight.Button
	}

	var (
		round = ctx.Theme.Button.Rounding
	)

	if bb.W == 0 {
		textStyle := style
		textSize := ctx.CalcTextSize(text, 0, textStyle.Font, textStyle.Size)
		extW := style.Padding.Left+style.Padding.Right
		extH := style.Padding.Top+style.Padding.Bottom
		bb.W = textSize[0] + extW
		bb.H = textSize[1] + extH
	}

	// Check Event
	event = ctx.CheckEvent(id, bb, false)

	// Render Frame
	ctx.ColorBackground(event, bb, round)

	// Render Text
	ctx.DrawText(bb, text, &style.TextStyle)
	return
}


func (ctx *Context) renderTextClipped(text string, bb *Rect, style *TextStyle) {
	x, y := Gui2Game(bb.X, bb.Y)
	font := ctx.Theme.Text.Font
	if bb.W == 0 {
		ctx.DrawList.AddText(f32.Vec2{x, y}, text, font, 12, 0xFF000000, 0)
	} else {
		ctx.DrawList.AddText(f32.Vec2{x, y}, text, font, 12, 0xFF000000, bb.W)
	}
}

func (ctx *Context) ImageBackground(eventType EventType) {

}

func (ctx *Context) ImageButton(id ID, normal, pressed gfx.Tex2D, bb *Rect, style *ImageButtonStyle) ( event EventType) {
	if style == nil {
		style = &ctx.Theme.ImageButton
	}

	event = ctx.CheckEvent(id, bb, false)
	var tex gfx.Tex2D
	if event & EventDown != 0 {
		tex = pressed
	} else {
		tex = normal
	}
	ctx.DrawImage(bb, tex, &style.ImageStyle)
	return
}

func (ctx *Context) Slider(id ID, bb *Rect, value *float32, style *SliderStyle) (e EventType){
	if style == nil {
		style = &ctx.Theme.Slider
	}

	// 说明滑动了，那么应该使用最新的值，而不是传入的值
	if v, event := ctx.CheckSlider(id, bb); event & EventDragging != 0 {
		*value = v
		e = event
	}

	ctx.DrawRect(bb, 0xFFCDCDCD, 5)
	ctx.DrawCircle(bb.X+bb.W*(*value), bb.Y+bb.H/2, 10, 0xFFABABAB)

	return
}

// Scroll 效果的关键是使用裁切限制滚动区域，然后
// 通过计算拖拽，来得到争取的偏移
func (ctx *Context) StartScroll(size, offset f32.Vec2) {
	event := ctx.CheckEvent(123, nil, false)

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
func (ctx *Context) CheckSlider(id ID, bound *Rect) (v float32, e EventType) {
	event := ctx.CheckEvent(id, bound, false)
	if (event & EventStartDrag) != 0 {
		ctx.state.pointerCapture = id
		log.Println("start drag..")
	}

	if (event & EventEndDrag) != 0 {
		ctx.state.pointerCapture = -1
		log.Println("end drag..")
	}

	//
	//if ctx.state.isLastEventPointerType {
	//}
	// Update the knob position
	if (event & (EventDragging|EventWentDown)) != 0{
		// default = Horizontal
		p1 := input.PointerPosition(0).MousePos[0]
		p0 := bound.X + ctx.Cursor.X
		v = (p1 - p0)/bound.W

		dbg.DrawStrScaled(fmt.Sprintf("p0: %v, p1: %v", p0, p1), .6 )
		dbg.Return()

		if v > 1 {
			v = 1
		}
		if v < 0 {
			v = 0
		}
	}
	e = event
	return
}

func (ctx *Context) capturePoint() {

}

func (ctx *Context) releasePointer() {

}

func (ctx *Context) isLastEventPointerType() bool {
	return true
}

// Algorithm from FlatUI: http://google.github.io/flatui/
func (ctx *Context) CheckEvent(id ID, bound *Rect, checkDragOnly bool) EventType {
	var (
		event = EventNone
		cursor = ctx.Cursor
	)

	bb := Rect{cursor.X+bound.X, cursor.Y+bound.Y, bound.W, bound.H}
	p  := input.PointerPosition(0)


	if bb.InRange(p.MousePos) || ctx.state.pointerCapture == id {
		// in-dragging, The pointer is in drag operation
		if btn := input.PointerButton(0); ctx.state.draggingPointer == id && !btn.JustPressed() {
			if btn.JustReleased() {
				event = EventEndDrag
				log.Println("drag end real..", event)
				ctx.state.draggingPointer = -1
				ctx.state.draggingStart = f32.Vec2{}
				event |= EventWentUp
			} else if btn.Down() {
				event = EventDragging
			} else {
				ctx.state.draggingPointer = -1
			}
		} else {
			// Check event start, event as DragStart/Down/Up

			// 1. Regular pointer event handling
			if !checkDragOnly {
				if btn.JustPressed() {
					ctx.state.active = id
					event = EventWentDown
				}

				if btn.JustReleased() {
					event = EventWentUp
				} else if btn.Down() {
					event |= EventDown
				}
			}

			// 2. Check for drag events
			// 2.1 Keep the click position, then use it to check a drag event
			if btn.JustPressed() {
				ctx.state.draggingStart = p.MousePos
				log.Println("just pressed!!")
			}
			// 2.2 If the next movement out of thresh-hold, then it's a drag event
			if btn.Down() && bb.InRange(ctx.state.draggingStart) {
				var (
					startPosition  = ctx.state.draggingStart
					dragThreshHold = float32(8)
				)

				bb := Rect{startPosition[0]-dragThreshHold,
					startPosition[1]-dragThreshHold,
					dragThreshHold,
					dragThreshHold}

				// Start drag event
				if !bb.InRange(p.MousePos) {
					event |= EventStartDrag
					ctx.state.draggingStart = p.MousePos
					ctx.state.draggingPointer = id
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

func (ctx *Context) DrawRect(bb *Rect, color uint32, round float32) {
	var (
		x = bb.X + ctx.Cursor.X
		y = bb.Y + ctx.Cursor.Y
	)
	x, y = Gui2Game(x, y)
	min, max := f32.Vec2{x, y-bb.H}, f32.Vec2{x+bb.W, y}
	ctx.DrawList.AddRectFilled(min, max, color, round, FlagCornerAll)
}

func (ctx *Context) DrawBorder(bb *Rect, color uint32, round, thick float32) {
	var (
		x = bb.X + ctx.Cursor.X
		y = bb.Y + ctx.Cursor.Y
	)
	x, y = Gui2Game(x, y)
	min, max := f32.Vec2{x, y-bb.H}, f32.Vec2{x+bb.W, y}
	ctx.DrawList.AddRect(min, max, color, round, FlagCornerAll, thick)
}

func (ctx *Context) DrawDebugBorder(x, y, w, h float32, color uint32) {
	x, y = Gui2Game(x + ctx.Cursor.X, y + ctx.Cursor.Y)
	min, max := f32.Vec2{x, y-h}, f32.Vec2{x+w, y}
	ctx.DrawList.AddRect(min, max, color, 0, FlagCornerNone, 1)
}

// default segment = 12 TODO
func (ctx *Context) DrawCircle(x, y, radius float32, color uint32) {
	c := ctx.Cursor
	x, y = Gui2Game(x + c.X, y + c.Y)
	ctx.DrawList.AddCircleFilled(f32.Vec2{x, y}, radius, color, 12)
}

func (ctx *Context) DrawImage(bound *Rect, tex gfx.Tex2D, style *ImageStyle) {
	c := ctx.Cursor
	min := f32.Vec2{bound.X+c.X, bound.Y+c.Y}
	if bound.W == 0 {
		sz := tex.Size()
		bound.W = sz.Width
		bound.H = sz.Height
	}
	max := min.Add(f32.Vec2{bound.W, bound.H})
	var color uint32
	if style != nil {
		color = style.TintColor
	} else {
		color = ctx.Theme.Image.TintColor
	}
	min[0], min[1] = Gui2Game(min[0], min[1])
	max[0], max[1] = Gui2Game(max[0], max[1])
	rg := tex.Region()
	ctx.DrawList.AddImage(tex.Tex(), min, max, f32.Vec2{rg.X1, rg.Y1}, f32.Vec2{rg.X2, rg.Y2}, color)
}

// 绘制元素, bb 存储相对于父容器的相对坐标..
func (ctx *Context) DrawText(bb *Rect, text string, style *TextStyle) (size f32.Vec2) {
	// 1. 取出布局
	c := ctx.Cursor
	x, y := Gui2Game(bb.X+style.Left+c.X, bb.Y+style.Top+c.Y)

	// 2. 开始绘制
	var (
		font = style.Font
		fontSize = style.Size
		color = style.Color
		wrapWidth = bb.W + 10
		pos = f32.Vec2{x, y}
	)
	size = ctx.DrawList.AddText(pos, text, font, fontSize, color, wrapWidth)
	return
}

func (ctx *Context) CalcTextSize(text string, wrapWidth float32, fnt font.Font, fontSize float32) f32.Vec2 {
	return font.CalculateTextSize(text, fnt, fontSize)
}

// 偷师 flat-ui 中的设计，把空间的前景和背景分离，背景单独根据事件来变化..
// 在 Android 中，Widget的前景和背景都可以根据控件状态发生变化
// 但是在大部分UI中，比如 Text/Image 只会改变背景的状态
// 偷懒的自定义UI，不做任何状态的改变... 所以说呢, 我们也采用偷懒的做法呗。。
func (ctx *Context) ColorBackground(event EventType, bb *Rect, round float32) {
	if event == EventDown {
		ctx.DrawRect(bb, ThemeLight.ColorPressed, round)
	} else {
		ctx.DrawRect(bb, ThemeLight.ColorNormal, round)
	}
}

// Clip:
func (ctx *Context) PushClipRect(minClip, maxClip f32.Vec2, intersectCurrent bool) {

}

// Theme:
func (ctx *Context) UseTheme(style *Theme) {
	ctx.Theme = style
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