package gui

import "korok.io/korok/gfx"

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
	Rect RectStyle
}

type Padding struct {
	Left, Right, Top, Bottom float32
}

type TextStyle struct {
	Font gfx.FontSystem
	Visibility
	Color uint32
	Size float32
	Lines int
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
		TextStyle{
			Visibility:Visible,
			Color:0xFF000000,
			Size:12,
		},
		ButtonStyle{
			TextStyle{Visibility:Visible, Color:0xFF000000, Size:12},
			Visible,
			0xFFCDCDCD,
			5,
		},
		ImageStyle{
			Visible,
			Padding{0, 0, 0, 0},
			0xFFFFFFFF,
		},
		RectStyle{
			2,
			0xFFCDCDCD,
			0xFF000000,
			5,
			FlagCornerNone,
		},
	}
}

func newDarkTheme() *Style {
	return &Style{}
}
