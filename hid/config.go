package hid

import "korok.io/korok/math/f32"

type WindowOptions struct{
	Title string
	Width int
	Height int
	Clear f32.Vec4
	FullScreen bool
	NoVsync    bool
	NoTitleBar bool
	Resizable  bool
}
