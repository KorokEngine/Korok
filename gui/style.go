package gui

import (
	"korok.io/korok/gfx/font"
	"korok.io/korok/gfx"
	"korok.io/korok/math/f32"
)

type Visibility uint8
const (
	Visible Visibility = iota
	InVisible
	Gone
)

type Theme struct {
	Text        TextStyle
	Button      ButtonStyle
	Image       ImageStyle
	ImageButton ImageButtonStyle
	Slider      SliderStyle

	// global config..
	Normal gfx.Color
	Pressed gfx.Color
	// item 之间的空隙
	Spacing float32

	// default font
	font.Font
}

type Padding struct {
	Left, Right, Top, Bottom float32
}

type TextStyle struct {
	Font font.Font
	Color gfx.Color
	Size float32
	Lines int
	LineSpace float32
}

func (text *TextStyle) SetFont(f font.Font) *TextStyle {
	text.Font = f
	return text
}

func (text *TextStyle) SetColor(color gfx.Color) *TextStyle {
	text.Color = color
	return text
}

func (text *TextStyle) SetSize(size float32) *TextStyle {
	text.Size = size
	return text
}

type InputStyle struct {
	Visibility
	Color, HintColor uint32
	Size float32
}

type ButtonStyle struct {
	TextStyle
	Padding
	Background struct{
		Normal gfx.Color
		Pressed gfx.Color
	}
	Gravity f32.Vec2
	Rounding float32
}

type ImageButtonStyle struct {
	ImageStyle
	Padding
}

type ImageStyle struct {
	Tint gfx.Color
}

type CheckBoxStyle struct {

}

type ProgressBarStyle struct {

}

type SliderStyle struct {
	Bar, Knob gfx.Color
}

//// 这样
func newLightTheme() *Theme {
	return &Theme{
		Text:TextStyle{
			Color:gfx.Black,
			Size:12,
			LineSpace:6,
		},
		Button:ButtonStyle{
			TextStyle:TextStyle{Color:gfx.Black, Size:12},
			Padding:Padding{10, 10, 10, 10},
			Gravity:f32.Vec2{.5, .5},
			Rounding:5,
		},
		Image:ImageStyle{gfx.White},
		ImageButton:ImageButtonStyle{
			ImageStyle:ImageStyle{ Tint:gfx.White},
			Padding:Padding{0, 0,0,0},
		},
		Slider:SliderStyle{
			gfx.LTGray, gfx.Gray,
		},
		Normal: gfx.LTGray,
		Pressed: gfx.Gray,
		Spacing:4,
	}
}

func newDarkTheme() *Theme {
	return &Theme{}
}
