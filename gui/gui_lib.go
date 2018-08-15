package gui

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
	"korok.io/korok/hid/input"
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

func (r *Rect) Offset(dx, dy float32) *Rect{
	r.X, r.Y = r.X+dx, r.Y+dy
	return r
}

func (r *Rect) Scale(skx, sky float32) *Rect{
	r.X, r.Y = r.X*skx, r.Y*sky
	r.W, r.H = r.W*skx, r.H*sky
	return r
}

func (r *Rect) InRange(p f32.Vec2) bool{
	if r.X < p[0] && p[0] < (r.X + r.W) && r.Y < p[1] && p[1] < (r.Y + r.H) {
		return true
	}
	return false
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
		var font = style.Font
		if font == nil {
			font = ctx.Theme.Font
		}
		sz := ctx.CalcTextSize(text, 0, font, style.Size)
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
	var (
		round = ctx.Theme.Button.Rounding
		offset f32.Vec2
		font font.Font
	)
	if style == nil {
		style = &ctx.Theme.Button
	}
	if style.Font != nil {
		font = style.Font
	} else {
		font = ctx.Theme.Font
	}

	if bb.W == 0 {
		textSize := ctx.CalcTextSize(text, 0, font, style.Size)
		extW := style.Padding.Left+style.Padding.Right
		extH := style.Padding.Top+style.Padding.Bottom
		bb.W = textSize[0] + extW
		bb.H = textSize[1] + extH
	} else {
		// if button has size, gravity will effect text's position
		var font = style.Font
		if font == nil {
			font = ctx.Theme.Font
		}
		textSize := ctx.CalcTextSize(text, 0, font, style.Size)
		g := style.Gravity
		offset[0] = (bb.W-textSize[0]-style.Padding.Left-style.Padding.Right) * g[0]
		offset[1] = (bb.H-textSize[1]-style.Padding.Top-style.Padding.Bottom) * g[1]
	}

	// Check Event
	event = ctx.ClickEvent(id, bb)

	// Render Frame
	if bg := style.Background; bg.Normal != (gfx.Color{}) {
		ctx.ColorBackground(event, bb, bg.Normal, bg.Pressed, round)
	} else {
		ctx.ColorBackground(event, bb, ThemeLight.Normal, ThemeLight.Pressed, round)
	}

	// Render Text
	bb.X +=  offset[0] + style.Padding.Left
	bb.Y +=  offset[1] + style.Padding.Top

	ctx.DrawText(bb, text, &style.TextStyle)
	return
}

func (ctx *Context) ImageBackground(eventType EventType) {

}

func (ctx *Context) ImageButton(id ID, normal, pressed gfx.Tex2D, bb *Rect, style *ImageButtonStyle) ( event EventType) {
	if style == nil {
		style = &ctx.Theme.ImageButton
	}
	event = ctx.ClickEvent(id, bb)
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

	ctx.DrawRect(bb, style.Bar, 5)
	ctx.DrawCircle(bb.X+bb.W*(*value), bb.Y+bb.H/2, 10, style.Knob)
	return
}

// Scroll 效果的关键是使用裁切限制滚动区域，然后
// 通过计算拖拽，来得到争取的偏移
func (ctx *Context) StartScroll(size, offset f32.Vec2) {

}

func (ctx *Context) EndScroll() {

}

func (ctx *Context) CheckSlider(id ID, bound *Rect) (v float32, e EventType) {
	event := ctx.DraggingEvent(id, bound)
	if (event & EventStartDrag) != 0 {
		ctx.state.pointerCapture = id
	}
	if (event & EventEndDrag) != 0 {
		ctx.state.pointerCapture = -1
	}
	// Update the knob position & default = Horizontal
	if (event & (EventDragging|EventWentDown)) != 0{
		p1 := (input.PointerPosition(0).MousePos[0])/screen.scaleX
		p0 := bound.X + ctx.Cursor.X
		v = (p1 - p0)/bound.W
		if v > 1 { v = 1 }
		if v < 0 { v = 0 }
	}
	e = event
	return
}

func (ctx *Context) ClickEvent(id ID, rect *Rect) EventType {
	var (
		event = EventNone
		c = ctx.Cursor
	)
	bb := Rect{(c.X+rect.X)*screen.scaleX, (c.Y+rect.Y)*screen.scaleY, rect.W*screen.scaleX,rect.H*screen.scaleY}
	if p  := input.PointerPosition(0); bb.InRange(p.MousePos) {
		btn := input.PointerButton(0)
		if btn.JustPressed() {
			ctx.state.active = id
			event = EventWentDown
		}
		if btn.JustReleased() && ctx.state.active == id {
			event = EventWentUp
			ctx.state.active = -1
		} else if btn.Down() && ctx.state.active == id {
			event |= EventDown
		}
	}
	return event
}

func (ctx *Context) DraggingEvent(id ID, bound *Rect) EventType {
	var (
		event = EventNone
		c = ctx.Cursor
	)

	bb := Rect{(c.X+bound.X)*screen.scaleX, (c.Y+bound.Y)*screen.scaleY, bound.W*screen.scaleX, bound.H*screen.scaleY}
	p  := input.PointerPosition(0)

	if bb.InRange(p.MousePos) || ctx.state.pointerCapture == id {
		// in-dragging, The pointer is in drag operation
		if btn := input.PointerButton(0); ctx.state.draggingPointer == id && !btn.JustPressed() {
			if btn.JustReleased() {
				event = EventEndDrag
				ctx.state.draggingPointer = -1
				ctx.state.draggingStart = f32.Vec2{}
				event |= EventWentUp
			} else if btn.Down() {
				event = EventDragging
			} else {
				ctx.state.draggingPointer = -1
			}
		} else {
			// Keep the click position, then use it to check a drag event
			if btn.JustPressed() {
				ctx.state.draggingStart = p.MousePos
			}
			// If the next movement out of thresh-hold, then it's a drag event
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
				}
			}
		}

	}
	return event
}

func (ctx *Context) DrawRect(bb *Rect, fill gfx.Color, round float32) {
	var (
		x = bb.X+ctx.Cursor.X
		y = bb.Y+ctx.Cursor.Y
	)
	x, y = Gui2Game(x, y)
	min := f32.Vec2{x * screen.scaleX, (y-bb.H) * screen.scaleY}
	max := f32.Vec2{(x+bb.W) * screen.scaleX, y * screen.scaleY}
	ctx.DrawList.AddRectFilled(min, max, fill.U32(), round * screen.scaleX, FlagCornerAll)
}

func (ctx *Context) DrawBorder(bb *Rect, color uint32, round, thick float32) {
	var (
		x = bb.X+ctx.Cursor.X
		y = bb.Y+ctx.Cursor.Y
	)
	x, y = Gui2Game(x, y)
	min := f32.Vec2{x * screen.scaleX, (y-bb.H) * screen.scaleY}
	max := f32.Vec2{(x+bb.W) * screen.scaleX, y * screen.scaleY}
	ctx.DrawList.AddRect(min, max, color, round * screen.scaleX, FlagCornerAll, thick)
}

func (ctx *Context) DrawDebugBorder(x, y, w, h float32, color uint32) {
	x, y = Gui2Game(x + ctx.Cursor.X, y + ctx.Cursor.Y)
	min := f32.Vec2{x * screen.scaleX, (y-h) * screen.scaleY}
	max := f32.Vec2{(x+w) * screen.scaleX, y * screen.scaleY}
	ctx.DrawList.AddRect(min, max, color, 0, FlagCornerNone, 1)
}

// default segment = 12 TODO, circle scale factor
func (ctx *Context) DrawCircle(x, y, radius float32, fill gfx.Color) {
	x, y = Gui2Game(x+ctx.Cursor.X, y+ctx.Cursor.Y)
	x = x * screen.scaleX
	y = y * screen.scaleY
	ctx.DrawList.AddCircleFilled(f32.Vec2{x, y}, radius * screen.scaleX, fill.U32(), 12)
}

func (ctx *Context) DrawImage(bound *Rect, tex gfx.Tex2D, style *ImageStyle) {
	min := f32.Vec2{bound.X+ctx.Cursor.X, bound.Y+ctx.Cursor.Y}
	if bound.W == 0 {
		sz := tex.Size()
		bound.W = sz.Width
		bound.H = sz.Height
	}
	max := min.Add(f32.Vec2{bound.W, bound.H})
	var color uint32
	if style != nil {
		color = style.Tint.U32()
	} else {
		color = ctx.Theme.Image.Tint.U32()
	}
	min[0], min[1] = Gui2Game(min[0], min[1])
	max[0], max[1] = Gui2Game(max[0], max[1])

	// scale
	min[0], min[1] = min[0] * screen.scaleX, min[1] * screen.scaleY
	max[0], max[1] = max[0] * screen.scaleX, max[1] * screen.scaleY

	rg := tex.Region()
	if rg.Rotated {
		ctx.DrawList.AddImageQuad(tex.Tex(),
			min, f32.Vec2{max[0], min[1]},  max, f32.Vec2{min[0], max[1]}, // xy
			f32.Vec2{rg.X2, rg.Y1}, f32.Vec2{rg.X2, rg.Y2}, f32.Vec2{rg.X1, rg.Y2}, f32.Vec2{rg.X1, rg.Y1},// uv
			color)
	} else {
		ctx.DrawList.AddImage(tex.Tex(), min, max, f32.Vec2{rg.X1, rg.Y1}, f32.Vec2{rg.X2, rg.Y2}, color)
	}
}

// 绘制元素, bb 存储相对于父容器的相对坐标..
func (ctx *Context) DrawText(bb *Rect, text string, style *TextStyle) (size f32.Vec2) {
	x, y := Gui2Game(bb.X+ctx.Cursor.X, bb.Y+ctx.Cursor.Y)
	var (
		font = style.Font
		fontSize = style.Size * screen.scaleX // TODO 字体缩放不能这么简单的考虑
		color = style.Color.U32()
		wrapWidth = (bb.W + 10) * screen.scaleX
		pos = f32.Vec2{x * screen.scaleX, y * screen.scaleY}
	)
	if font == nil {
		font = ctx.Theme.Font
	}
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
func (ctx *Context) ColorBackground(event EventType, bb *Rect, normal, pressed gfx.Color, round float32) {
	if (event & EventDown) != 0 {
		ctx.DrawRect(bb, pressed, round)
	} else {
		ctx.DrawRect(bb, normal, round)
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
	return x, screen.hintY - y
}

func Game2Gui(x, y float32) (x1, y1 float32) {
	return x, screen.hintY - y
}

type screenSize struct{
	rlWidth, rlHeight float32
	vtWidth, vtHeight float32
	// hint
	hintX, scaleX float32
	hintY, scaleY float32
}

func (sc *screenSize) SetRealSize(w, h float32) {
	sc.rlWidth, sc.rlHeight = w, h
	sc.updateHint()
}

func (sc *screenSize) SetVirtualSize(w, h float32) {
	sc.vtWidth, sc.vtHeight = w, h
	sc.updateHint()
}

func (sc *screenSize) updateHint() {
	if screen.rlWidth == 0 || screen.rlHeight == 0 {
		return
	}
	// update hint
	w := screen.vtWidth
	h := screen.vtHeight
	if w == 0 && h == 0 {
		screen.hintX = screen.rlWidth
		screen.hintY = screen.rlHeight
		screen.scaleX = 1
		screen.scaleY = 1
	} else if w == 0 {
		f := screen.rlHeight/h
		screen.scaleY = f
		screen.scaleX = f
		screen.hintY = h
		screen.hintX = screen.rlWidth/f
	} else if h == 0 {
		f := screen.rlWidth/w
		screen.scaleY = f
		screen.scaleX = f
		screen.hintX = w
		screen.hintY = screen.rlHeight/f
	} else {
		screen.scaleX = screen.rlWidth/w
		screen.scaleY = screen.rlHeight/h
		screen.hintX = w
		screen.hintY = h
	}
}

var screen screenSize