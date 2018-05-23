package gui

import (
	"korok.io/korok/gfx/font"
	"korok.io/korok/gfx"
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
}

type Padding struct {
	Left, Right, Top, Bottom float32
}

type TextStyle struct {
	Padding
	Font font.Font
	Visibility
	Color gfx.Color
	Size float32
	Lines int
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
	Visibility
	Background gfx.Color
	Rounding float32
}

type ImageButtonStyle struct {
	ImageStyle
	Visibility
}

type ImageStyle struct {
	Visibility
	Padding
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
			Visibility:Visible,
			Color:gfx.Blank,
			Size:12,
		},
		Button:ButtonStyle{
			TextStyle{Visibility:Visible, Color:gfx.Blank, Size:12, Padding:Padding{10, 10, 10, 10}},
			Visible,
			gfx.LTGray,
			5,
		},
		Image:ImageStyle{
			Visible,
			Padding{0, 0, 0, 0},
			gfx.While,
		},
		ImageButton:ImageButtonStyle{
			ImageStyle{
				Visibility:Visible,
				Padding: Padding{0, 0, 0, 0},
				Tint:gfx.While,
			},
			Visible,
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
