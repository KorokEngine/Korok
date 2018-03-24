package gui

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/hid/input"
	"korok.io/korok/gfx/dbg"

	"log"
	"fmt"
	"korok.io/korok/gfx/font"
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
func (ctx *Context) Text(id ID, text string, style *TextStyle)  *Element {
	var (
		elem, ready = ctx.BeginElement(id)
		size f32.Vec2
	)

	// draw text 最好返回最新的大小..
	if ready {
		size = ctx.DrawText(elem, text, style)
	} else {
		size = ctx.CalcTextSize(text, 0, style.Font, style.Size)
	}

	elem.Bound.W = size[0]
	elem.Bound.H = size[1]

	ctx.EndElement(elem)
	return nil
}

// Widgets: InputEditor
func (ctx *Context) InputText(hint string, lyt LayoutManager, style *InputStyle) {

}

// Widget: Image
func (ctx *Context) Image(id ID, texId uint16, uv f32.Vec4, style *ImageStyle) {
	var (
		elem, ready = ctx.BeginElement(id)
	)

	if ready {
		ctx.DrawImage(&elem.Bound, texId, uv, style)
	} else {
		size := ctx.Layout.Cursor.Bound
		elem.W = size.W
		elem.H = size.H
	}

	ctx.EndElement(elem)
}

// Widget: Button
func (ctx *Context) Button(id ID, text string, style *ButtonStyle) (event EventType) {
	if style == nil {
		style = &ThemeLight.Button
	}

	var (
		elem, ready = ctx.BeginElement(id)
		round = ctx.Style.Button.Rounding
	)

	if ready {
		bb := &elem.Bound
		// Check Event
		event = ctx.CheckEvent(id, bb, false)

		// Render Frame
		ctx.ColorBackground(event, bb, round)

		// Render Text
		ctx.DrawText(elem, text, &style.TextStyle)
	} else {
		textStyle := ctx.Style.Text
		textSize := ctx.CalcTextSize(text, 0, textStyle.Font, textStyle.Size)
		extW := style.Padding.Left+style.Padding.Right
		extH := style.Padding.Top+style.Padding.Bottom
		elem.W, elem.H = textSize[0]+extW, textSize[1]+extH
	}
	ctx.EndElement(elem)
	return
}


func (ctx *Context) renderTextClipped(text string, bb *Bound, style *TextStyle) {
	x, y := Gui2Game(bb.X, bb.Y)
	font := ctx.Style.Text.Font
	if bb.W == 0 {
		ctx.DrawList.AddText(f32.Vec2{x, y}, text, font, 12, 0xFF000000, 0)
	} else {
		ctx.DrawList.AddText(f32.Vec2{x, y}, text, font, 12, 0xFF000000, bb.W)
	}
}



func (ctx *Context) ImageBackground(eventType EventType) {

}

func (ctx *Context) ImageButton(id ID, normal, pressed gfx.Tex2D, style *ImageButtonStyle) ( event EventType) {
	if style == nil {
		style = &ctx.Style.ImageButton
	}
	var (
		elem, ready = ctx.BeginElement(id)
		bb = &elem.Bound
	)
	if ready {
		event = ctx.CheckEvent(id, bb, false)
		var tex gfx.Tex2D
		if event & EventDown != 0 {
			tex = pressed
		} else {
			tex = normal
		}
		rg := tex.Region()
		ctx.DrawImage(bb, tex.Tex(), f32.Vec4{rg.X1, rg.Y1, rg.X2, rg.Y2}, &style.ImageStyle)
	} else {
		size := ctx.Layout.Cursor.Bound
		elem.W = size.W
		elem.H = size.H
	}
	ctx.EndElement(elem)
	return
}

// Slider 的绘制很简单，分别绘制滑动条和把手即可
// 难点在于跟踪把手的滑动距离
// Slider的风格，没有想好怎么控制，暂时使用两张图片
// 分别绘制Bar和Knob
// Slider 需要保存混动的结果，否则
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
func (ctx *Context) Slider(id ID, value *float32, style *SliderStyle) (e EventType){
	if style == nil {
		style = &ctx.Style.Slider
	}

	var (
		elem, ready = ctx.BeginElement(id)
		bb = &elem.Bound
	)

	if ready {
		// 说明滑动了，那么应该使用最新的值，而不是传入的值
		if v, event := ctx.checkSlider(id, bb); event & EventDragging != 0 {
			*value = v
			e = event
		}

		ctx.DrawRect(bb, 0xFFCDCDCD, 5)
		ctx.DrawCircle(bb.X+bb.W*(*value), bb.Y+bb.H/2, 10, 0xFFABABAB)
	} else {
		// 设置默认的宽高
		if elem.W == 0 {
			elem.W = 100
		}
		if elem.H == 0 {
			elem.H = 10
		}
	}

	ctx.EndElement(elem)
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
func (ctx *Context) checkSlider(id ID, bound *Bound) (v float32, e EventType) {
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
		p0 := bound.X + ctx.Layout.hGroup.X
		v = (p1 - p0)/bound.W


		dbg.Move(10, 300)
		dbg.DrawStrScaled(fmt.Sprintf("p0: %v, p1: %v", p0, p1), .6 )

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
func (ctx *Context) CheckEvent(id ID, bound *Bound, checkDragOnly bool) EventType {
	var (
		event = EventNone
		g  = ctx.Layout.hGroup
	)

	bb := Bound{g.X+bound.X, g.Y + bound.Y, bound.W, bound.H}
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

				bb := Bound{startPosition[0]-dragThreshHold,
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

func (ctx *Context) DrawRect(bb *Bound, color uint32, round float32) {
	var (
		g = ctx.Layout.hGroup
		x = g.X + bb.X
		y = g.Y + bb.Y
	)
	x, y = Gui2Game(x, y)
	min, max := f32.Vec2{x, y-bb.H}, f32.Vec2{x+bb.W, y}
	ctx.DrawList.AddRectFilled(min, max, color, round, FlagCornerAll)
}

func (ctx *Context) DrawBorder(bb *Bound, color uint32, round, thick float32) {
	var (
		g = ctx.Layout.hGroup
		x = g.X + bb.X
		y = g.Y + bb.Y
	)
	x, y = Gui2Game(x, y)
	min, max := f32.Vec2{x, y-bb.H}, f32.Vec2{x+bb.W, y}
	ctx.DrawList.AddRect(min, max, color, round, FlagCornerAll, thick)
}

func (ctx *Context) DrawDebugBorder(x, y, w, h float32, color uint32) {
	x, y = Gui2Game(x, y)
	min, max := f32.Vec2{x, y-h}, f32.Vec2{x+w, y}
	ctx.DrawList.AddRect(min, max, color, 0, FlagCornerNone, 1)
}

// default segment = 12 TODO
func (ctx *Context) DrawCircle(x, y, radius float32, color uint32) {
	g := ctx.Layout.hGroup
	x, y = Gui2Game(g.X + x, g.Y + y)
	ctx.DrawList.AddCircleFilled(f32.Vec2{x, y}, radius, color, 12)
}

func (ctx *Context) DrawImage(bound *Bound, texId uint16, uv f32.Vec4, style *ImageStyle) {
	g := ctx.Layout.hGroup
	min := f32.Vec2{g.X+bound.X, g.Y+bound.Y}
	if bound.W == 0 {
		if ok, tex := bk.R.Texture(texId); ok {
			bound.W, bound.H = tex.Width, tex.Height
		}
	}
	max := min.Add(f32.Vec2{bound.W, bound.H})
	var color uint32
	if style != nil {
		color = style.TintColor
	} else {
		color = ctx.Style.Image.TintColor
	}
	min[0], min[1] = Gui2Game(min[0], min[1])
	max[0], max[1] = Gui2Game(max[0], max[1])
	ctx.DrawList.AddImage(texId, min, max, f32.Vec2{uv[0], uv[1]}, f32.Vec2{uv[2], uv[3]}, color)
}

// 绘制元素, bb 存储相对于父容器的相对坐标..
// Group 目前是绝对坐标
// Group + Offset = 当前绝对坐标..
func (ctx *Context) DrawText(bb *Element, text string, style *TextStyle) (size f32.Vec2) {
	// 1. 取出布局
	group := ctx.Layout.hGroup
	x, y := Gui2Game(group.X+bb.X+style.Left, group.Y+bb.Y+style.Top)

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

func (ctx *Context) CalcTextSize(text string, wrapWidth float32, font font.Font, fontSize float32) f32.Vec2 {
	fr := &FontRender{
		font: font,
		fontSize:fontSize,
	}
	return fr.CalculateTextSize1(text)
}

// 偷师 flat-ui 中的设计，把空间的前景和背景分离，背景单独根据事件来变化..
// 在 Android 中，Widget的前景和背景都可以根据控件状态发生变化
// 但是在大部分UI中，比如 Text/Image 只会改变背景的状态
// 偷懒的自定义UI，不做任何状态的改变... 所以说呢, 我们也采用偷懒的做法呗。。
func (ctx *Context) ColorBackground(event EventType, bb *Bound, round float32) {
	if event == EventDown {
		ctx.DrawRect(bb, ThemeLight.ColorPressed, round)
	} else {
		ctx.DrawRect(bb, ThemeLight.ColorNormal, round)
	}
}


// 计算单个UI元素
// 如果有大小则记录出偏移和Margin
// 否则只返回元素
func (ctx *Context) BeginElement(id ID) (elem *Element, ok bool){
	lm := &ctx.Layout
	if elem, ok = lm.Element(id); !ok {
		elem = lm.NewElement(id)
	} else {
		// 计算偏移
		elem.X = lm.Cursor.X + ctx.Layout.spacing
		elem.Y = lm.Cursor.Y + ctx.Layout.spacing

		// Each element's property
		if lm.Cursor.owner == id {
			// 计算 Margin 和 偏移
			if lm.Cursor.Flag & FlagMargin != 0 {
				elem.Margin = lm.Cursor.Margin
				elem.X += elem.Left
				elem.Y += elem.Top
			}

			// 计算大小
			if lm.Cursor.Flag & FlagSize != 0 {
				elem.Bound.W = lm.Cursor.W
				elem.Bound.H = lm.Cursor.H
			}

			// 清空标记
			lm.Cursor.owner = -1
			lm.Cursor.Flag = 0
		}

		// Gravity
		var (
			group = lm.hGroup
			gravity = group.Gravity
		)

		// Overlap group's gravity
		if lm.Cursor.owner == id && (lm.Cursor.Flag & FlagGravity != 0) {
			gravity = lm.Cursor.Gravity
		}
		switch group.LayoutType {
		case LinearHorizontal:
			elem.Y += (group.H - elem.H) * gravity.Y
		case LinearVertical:
			elem.X += (group.W - elem.W) * gravity.X
		case LinearOverLay:
			elem.Y += (group.H - elem.H) * gravity.Y
			elem.X += (group.W - elem.W) * gravity.X
		}
	}
	return
}

// 结束绘制, 每绘制完一个元素都要偏移一下光标
func (ctx *Context) EndElement(elem *Element) {
	ctx.Layout.Advance(elem)
	ctx.Layout.Extend(elem)
}

// Layout
func (ctx *Context) BeginLayout(id ID, xtype LayoutType) {
	var (
		lm = &ctx.Layout;
		ly, ok = lm.FindLayout(id)
	)

	if !ok {
		ly = lm.NewLayout(id, xtype)
	}

	// debug draw - render group frame
	if ok {
		var (
			x = lm.hGroup.X + ctx.Layout.Cursor.X
			y = lm.hGroup.Y + ctx.Layout.Cursor.Y
		)
		ctx.DrawDebugBorder(x, y, ly.W, ly.H, 0xFF00FF00)
	}

	lm.PushLayout(xtype, ly)
}

func (ctx *Context) EndLayout() {
	ctx.Layout.EndLayout()
}

// Reference System: VirtualBounds
func (ctx *Context) PushVBounds(bounds f32.Vec4) {

}

// Clip:
func (ctx *Context) PushClipRect(minClip, maxClip f32.Vec2, intersectCurrent bool) {

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