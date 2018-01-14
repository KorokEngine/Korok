package gui

import (
	"github.com/go-gl/mathgl/mgl32"
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
	LayoutLinearVertical LayoutType = iota
	LayoutLinearHorizontal
)

// Widgets: Text
func Text(text string, style *TextStyle) (id int){
	return gContext.Text(text, style)
}

func TextSizeColored(text string, color uint32, size float32) {

}

func Label(text string) {

}

// Widgets: InputEditor
func InputText(hint string, lyt Layout, style *InputStyle) {

}

// Widget: Image
func Image(texId uint16, uv mgl32.Vec4, style *ImageStyle) int{
	return gContext.Image(texId, uv, style)
}

// Widget: Button
func Button(text string, style *ButtonStyle) (id int, event EventType) {
	return gContext.Button(text, style)
}

func ImageButton(texId uint16, lyt Layout, style *ImageButtonStyle) EventType{
	return EventNone
}

func CheckBox(text string, lyt Layout, style *CheckBoxStyle) bool {
	return false
}

// Widget: ProgressBar
func ProgressBar(fraction float32, lyt Layout, style *ProgressBarStyle) {

}

// Frame: Rect
func Rect(w, h float32, style *RectStyle) (id int) {
	return gContext.Rect(w, h, style)
}

// Widget: ListView TODO
func ListView() {

}

// Container: Window/PopupWindow/
func OpenPopup(id string) {

}

func DismissPopup(id string) {

}

// 受管理的窗口对象,
// 窗口的状态会被系统记录!
func BeginWindow(name string) {

}

func EndWindow() {

}

// Layout: 方便坐标计算的布局系统
func BeginLayout(bb *Bound) (lyt *Layout) {
	lyt = &gContext.Layout
	lyt.Begin(bb)
	return
}

func EndLayout() {
	gContext.Layout.Reset()
}

// Reference System: VirtualBounds
func PushVBounds(bounds mgl32.Vec4) {

}

// Clip:
func PushClipRect(minClip, maxClip mgl32.Vec2, intersectCurrent bool) {

}

func Popup() {

}

// Theme:
func UseTheme(style *Style) {
	gContext.UseTheme(style)
}

func SetFont(font gfx.FontSystem) {
	gContext.Style.Text.Font = font
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





