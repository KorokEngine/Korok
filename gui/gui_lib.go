package gui

import (
	"github.com/go-gl/mathgl/mgl32"

	"korok.io/korok/gfx"
	"korok.io/korok/gfx/bk"
	"korok.io/korok/hid/input"
)

type EventType uint8
const (
	EventNone EventType = iota
	EventDown
	EventUp
	EventMove
)

// 一个Context维护UI层逻辑:
// 1. 一个 DrawList，负责生成顶点
// 2. 共享的 Style，指针引用，多个 Context 之间可以共享
// 3. State 维护，负责维护当前的 UI 状态，比如动画，按钮，列表...
// 窗口在 Context 之上维护，默认的Context可以随意的绘制
// 在窗口内绘制的UI会受到窗口的管理.
type Context struct {
	DrawList
	Layout
	*Style

	// ui global state
	state struct{
		hot, active int
		mouseX, mouseY, mouseDown int
	}
}

func NewContext(style *Style) *Context {
	c := &Context{
		Style:style,
	}
	c.DrawList.Initialize()
	return c
}

// Widgets: Text
func (ctx *Context) Text(text string, style *TextStyle) (id int){
	x,y := Gui2Game(ctx.Bound.X, ctx.Bound.Y)

	var font gfx.FontSystem
	var fontSize float32
	var color uint32
	var wrapWidth = ctx.Bound.W

	dft := &ctx.Style.Text
	if style == nil {
		font = dft.Font
		fontSize = dft.Size
		color = dft.Color
	} else {
		if style.Font == nil {
			font = dft.Font
		}
		if style.Size == 0 {
			fontSize = dft.Size
		}
		if style.Color == 0 {
			color = dft.Color
		}
		if style.Lines == 1 {
			wrapWidth = 0
		}
	}

	size := ctx.DrawList.AddText(mgl32.Vec2{x, y}, text, font, fontSize, color, wrapWidth)
	id = ctx.Layout.Push(&Bound{ctx.Bound.X, ctx.Bound.Y, size[0], size[1]})
	return
}

func (ctx *Context) Label(text string) {

}

func (ctx *Context) CalcTextSize(text string, wrapWidth float32, font gfx.FontSystem, fontSize float32) mgl32.Vec2 {
	fr := &FontRender{
		font: font,
		fontSize:fontSize,
	}
	return fr.CalculateTextSize1(text)
}

// Widgets: InputEditor
func (ctx *Context) InputText(hint string, lyt Layout, style *InputStyle) {

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
	id = ctx.Layout.Push(&bound)
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
	event = ctx.CheckEvent(&Bound{x, y, w, h})

	// Render Frame
	ctx.renderFrame(x, y, w, h, color, rounding)

	// Render Text
	x = bound.X + w/2 - textSize[0]/2
	y = bound.Y + h/2 - textSize[1]/2
	ctx.renderTextClipped(text, &Bound{x, y, 0, 0}, &style.TextStyle)

	// push bound
	bound.W, bound.H = w, h
	id = ctx.Layout.Push(&bound)
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
func (ctx *Context) CheckEvent(bound *Bound) EventType {
	event := EventNone
	if p := input.PointerPosition(0); bound.InRange(p.MousePos) {
		btn := input.PointerButton(0)
		id := int(0) // // todo 设计ID系统，记录每个按键的位置..
		if btn.JustPressed() {
			ctx.state.active = id
			event = EventDown
		}
		if btn.JustReleased() && ctx.state.active == id {
			event = EventUp
		}
	}
	return event
}

func (ctx *Context) ImageButton(texId uint16, lyt Layout, style *ImageButtonStyle) EventType{
	return EventNone
}

func (ctx *Context) CheckBox(text string, lyt Layout, style *CheckBoxStyle) bool {
	return false
}

// Widget: ProgressBar
func (ctx *Context) ProgressBar(fraction float32, lyt Layout, style *ProgressBarStyle) {

}

// Widget: ListView TODO
func (ctx *Context) ListView() {

}

func (ctx *Context) Rect(w, h float32, style *RectStyle) (id int){
	bb := ctx.Bound

	x, y := Gui2Game(bb.X, bb.Y)

	var min, max mgl32.Vec2

	if ctx.Horizontal == Left2Right {
		min[0], max[0] = x, x+w
	} else {
		min[0], max[0] = x-w, x
	}

	if ctx.Vertical == Top2Bottom {
		min[1], max[1] = y-h, y
	} else {
		min[1], max[1] = y, y+h
	}

	bb.W, bb.H = w, h
	id = ctx.Layout.Push(&bb)

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

// Container: Window/PopupWindow/
func (ctx *Context) OpenPopup(id string) {

}

func (ctx *Context) DismissPopup(id string) {

}

func (ctx *Context) BeginGroup() {

}

func (ctx *Context) EndGroup() {

}

// Layout: 方便坐标计算的布局系统
func (ctx *Context) BeginLayout() {

}

func (ctx *Context) EndLayout() {

}

// Reference System: VirtualBounds
func (ctx *Context) PushVBounds(bounds mgl32.Vec4) {

}

// Clip:
func (ctx *Context) PushClipRect(minClip, maxClip mgl32.Vec2, intersectCurrent bool) {

}

func (ctx *Context) Popup() {

}

// Theme:
func (ctx *Context) UseTheme(style *Style) {
	ctx.Style = style
}

type Window struct {

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

type ID int