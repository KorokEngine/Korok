package gui

import (
	"korok.io/korok/gfx/font"
)

type Visibility uint8
const (
	Visible Visibility = iota
	InVisible
	Gone
)

type Style struct {
	Text TextStyle
	Button ButtonStyle
	Image ImageStyle
	ImageButton ImageButtonStyle
	Rect RectStyle
	Slider SliderStyle

	// global config..
	ColorNormal uint32
	ColorPressed uint32
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
	Color uint32
	Size float32
	Lines int
}

func (text *TextStyle) SetFont(f font.Font) *TextStyle {
	text.Font = f
	return text
}

func (text *TextStyle) SetColor(color uint32) *TextStyle {
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
	Color uint32
	Rounding float32
}

type ImageButtonStyle struct {
	ImageStyle
	Visibility
}

type ImageStyle struct {
	Visibility
	Padding
	TintColor uint32
}

type CheckBoxStyle struct {

}

type ProgressBarStyle struct {

}

type SliderStyle struct {
	Bar, Knob uint16
}

type RectStyle struct {
	Stroke float32
	FillColor uint32
	StrokeColor uint32
	Rounding float32
	Corner FlagCorner
}

//// 这样
func newLightTheme() *Style {
	return &Style{
		Text:TextStyle{
			Visibility:Visible,
			Color:0xFF000000,
			Size:12,
		},
		Button:ButtonStyle{
			TextStyle{Visibility:Visible, Color:0xFF000000, Size:12, Padding:Padding{10, 10, 10, 10}},
			Visible,
			0xFFCDCDCD,
			5,
		},
		Image:ImageStyle{
			Visible,
			Padding{0, 0, 0, 0},
			0xFFFFFFFF,
		},
		ImageButton:ImageButtonStyle{
			ImageStyle{
				Visibility:Visible,
				Padding: Padding{0, 0, 0, 0},
				TintColor:0xFFFFFFFF,
			},
			Visible,
		},
		Rect:RectStyle{
			2,
			0xFFCDCDCD,
			0xFF000000,
			5,
			FlagCornerNone,
		},
		Slider:SliderStyle{
			0, 0,
		},
		ColorNormal:0xFFCDCDCD,
		ColorPressed:0xFFABABAB,
		Spacing:4,
	}
}

func newDarkTheme() *Style {
	return &Style{}
}
