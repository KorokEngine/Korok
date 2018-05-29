package gui

import (
	"korok.io/korok/gfx"
	"korok.io/korok/gfx/font"
)

//	Awesome GUI System

// Widgets: Text
func Text(id ID, bb Rect, text string, style *TextStyle) {
	if style == nil {
		style = &gContext.Theme.Text
	}
	gContext.Text(id, &bb, text, style)
	return
}

func TextSizeColored(id ID, bb Rect, text string, color gfx.Color, size float32) {
	sty := gContext.Theme.Text
	sty.Color = color
	sty.Size = size
	gContext.Text(id, &bb,text, &sty)
}

// Widgets: InputEditor
func InputText(id ID, bb Rect, hint string, style *InputStyle) {
	gContext.InputText(id, &bb, hint, style)
}

// Widget: Image
func Image(id ID, bb Rect, tex gfx.Tex2D, style *ImageStyle) {
	gContext.Image(id, &bb, tex, style)
}

// Widget: Button
func Button(id ID, bb Rect, text string, style *ButtonStyle) (event EventType) {
	return gContext.Button(id, &bb, text, style)
}

func ImageButton(id ID, bb Rect, normal, pressed gfx.Tex2D, style *ImageButtonStyle) EventType{
	return gContext.ImageButton(id, normal, pressed, &bb, style)
}

func CheckBox(id ID, bb Rect, text string, style *CheckBoxStyle) bool {
	return false
}

// Widget: ProgressBar, Slider
func ProgressBar(id ID, bb Rect, fraction float32, style *ProgressBarStyle) {

}

func Slider(id ID, bb Rect, value *float32, style *SliderStyle) (v EventType){
	return gContext.Slider(id, &bb, value, style)
}

// Widget: ColorRect
func ColorRect(bb Rect, fill gfx.Color, rounding float32) {
	gContext.DrawRect(&bb, fill, rounding)
}

// Widget: ListView TODO
func ListView() {

}

// Offset move the ui coordinate's origin by (dx, dy)
func Offset(dx, dy float32) {
	gContext.Cursor.X += dx
	gContext.Cursor.Y += dy
}

// Move sets the ui coordinate's origin to (x, y)
func Move(x, y float32) {
	gContext.Cursor.X = x
	gContext.Cursor.Y = y
}

// Theme:
func UseTheme(style *Theme) {
	gContext.UseTheme(style)
}

func SetFont(font font.Font) {
	gContext.Theme.Text.Font = font
	gContext.Theme.Button.Font = font
	texFont, _ := font.Tex2D()
	gContext.DrawList.PushTextureId(texFont)
}

// Set Z-Order for
func SetZOrder(z int16) (old int16){
	old = gContext.ZOrder
	gContext.DrawList.ZOrder = z
	return
}

// for internal usage, DO NOT call.
func SetScreenSize(w, h float32) {
	screen.SetRealSize(w, h)
}

// ScreenSize return the physical width&height of the screen.
func ScreenSize() (w, h float32) {
	return screen.rlWidth, screen.rlHeight
}

func VirtualSize() (w, h float32) {
	return screen.vtWidth, screen.vtHeight
}

// SetVirtualResolution set the virtual resolution.
func SetVirtualResolution(w, h float32) {
	screen.SetVirtualSize(w, h)
}

func DefaultContext() *Context {
	return gContext
}

var DebugDraw = false

var ThemeLight *Theme
var ThemeDark  *Theme

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
