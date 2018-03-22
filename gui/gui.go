package gui

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
)

/**
	Awesome GUI System

	实现优秀的布局系统是很困难的，尤其是在 Immediate-Mode，典型的 2-pass 布局也会变得不适，
	我觉得提供一个简单的移动光标的方法可能会更有用，比如：

	在程序中，记录一个结构：
	struct Bounds {
		left, top float32
		bottom, right float32
	}
	每个组件都包含一个类似的结构，这样就可以方便的计算组件的相对位置.
	这样便可以求出每个已经绘制的组件的位置。在 Group 开始前入栈，EndGroup 之后就清空位置记录

	可以实现的布局：

	垂直线性布局：
	gui.Button("btn.1")
	gui.BottomOf(""btn.1)
	gui.Button("btn.2")

	这样便可以将元素垂直绘制
	水平布局：
	gui.Button("btn.1")
	gui.ToRightOf(""btn.1)
	gui.Button("btn.2")

	这样便可以水平绘制元素

	相对布局：
	gui.BottomOf(""btn.1)
	gui.AlignLeftOf
	gui.ToRightOf()
	gui.ToLeftOf()
	gui.TopOf()  // 可能的问题是，这个元素还没有绘制，可以使用 XVBounds 提前输入这个位置

	这样便可以方便的实现相对布局，此时绘制的时候，将 Style {Align: Center}
	这样一元素中心对齐。

	表格布局：
	gui.Ratio("id", 1/5)
	这样变求出父类1/5位置

	上面方式虽然简单，但是无法父类的 padding 的情况，
	此时可以用 虚框，也即是往 Bounds 里面随便放入一个参照物

	gui.VirtualBounds("vb.1", bounds)

	这样计算表格的时候就更方便了，

	gui.Ratio("vb.1", 1/5)

	就来到了这个虚拟参照物的 1/5位置

	表格高度怎么算呢？

	其实很简单，提供了 Measure() 方法返回元素高度就可以了。

	以上没有考虑 Group 包含 Group 的情况！

	碉堡了...
*/

type LayoutType int

const (
	LinearVertical LayoutType = iota
	LinearHorizontal
	LinearOverLay
)

var layout bool

// Widgets: Text
func Text(id ID, text string, style *TextStyle) {
	if style == nil {
		style = &gContext.Style.Text
	}
	gContext.Text(id, text, style)
	return
}

func TextSizeColored(text string, color uint32, size float32) {

}

// Widgets: InputEditor
func InputText(hint string, style *InputStyle) {

}

// Widget: Image
func Image(id ID, texId uint16, uv f32.Vec4, style *ImageStyle) {
	gContext.Image(id, texId, uv, style)
}

// Widget: Button
func Button(id ID, text string, style *ButtonStyle) (event EventType) {
	return gContext.Button(id, text, style)
}

func ImageButton(id ID, normal, pressed gfx.Sprite, style *ImageButtonStyle) EventType{
	return gContext.ImageButton(id, normal, pressed, style)
}

func CheckBox(text string, style *CheckBoxStyle) bool {
	return false
}

// Widget: ProgressBar, Slider
func ProgressBar(fraction float32, style *ProgressBarStyle) {

}

func Slider(id ID, value *float32, style *SliderStyle) (v EventType){
	return gContext.Slider(id, value, style)
}

// Frame: Rect
func Rect(w, h float32, style *RectStyle) {
	if style == nil {
		style = &gContext.Style.Rect
	}
	gContext.DrawRect(&Bound{0, 0, w,h}, style.FillColor, style.Rounding)
}

// Widget: ListView TODO
func ListView() {

}

// 基于当前 Group 移动光标
func Offset(x, y float32) {
	gContext.Layout.Offset(x, y)
}

func Move(x, y float32) {
	gContext.Layout.Move(x, y)
}

func SetGravity(x, y float32) {
	gContext.Layout.SetGravity(x, y)
}

func SetPadding(top, left, right, bottom float32) {
	gContext.Layout.SetPadding(top, left, right, bottom)
}

func SetSize(w, h float32) {
	gContext.Layout.SetSize(w, h)
}

func Cursor() *cursor {
	return &gContext.Layout.Cursor
}

func BeginHorizontal(id ID) {
	gContext.BeginLayout(id, LinearHorizontal)
}

func EndHorizontal() {
	gContext.EndLayout()
}

func BeginVertical(id ID) {
	gContext.BeginLayout(id, LinearVertical)
}

func EndVertical() {
	gContext.EndLayout()
}

// 参数需要一个 Rect，暂时用 Cursor 代替
func BeginDock(id ID, w, h float32) {
	gContext.BeginLayout(id, LinearOverLay)
	gContext.Layout.SetSize(w, h)
}

func EndDock() {
	gContext.EndLayout()
}

// Theme:
func UseTheme(style *Style) {
	gContext.UseTheme(style)
}

func SetFont(font gfx.FontSystem) {
	gContext.Style.Text.Font = font
	gContext.Style.Button.Font = font
	texFont, _ := font.Tex()
	gContext.DrawList.PushTextureId(texFont)
}

// gui init, render and destroy
func Init() {

}

func Frame() {

}

func Destroy() {

}

func DefaultContext() *Context {
	return gContext
}

var ThemeLight *Style
var ThemeDark  *Style

////////// implementation
// 应该设计一种状态管理机制，用这套机制来维护状态
// 比如：PopupWindow/Animation/Toast/
// 这些控件的典型特点是状态是变化的，而且不方便用代码维护
// 如果用一个专门的系统，问题可以减轻很多

var gContext *Context

func init() {
	// default theme
	ThemeLight = newLightTheme()
	ThemeDark  = newDarkTheme()

	// default context
	gContext = NewContext(ThemeLight)
}
