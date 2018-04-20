package font

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx/bk"

	"image"
	"image/color"

)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

// Known directions.
const (
	LeftToRight Direction = iota // E.g.: Latin
	RightToLeft                  // E.g.: Arabic
	TopToBottom                  // E.g.: Old-Chinese
)

// Font provides all the information needed to render a Rune.
//
// Tex2D returns the low-level bk-texture.
// Glyph returns the matrix of rune.
// Bounds returns the largest width and height for any of the glyphs
// in the fontAtlas.
// SizeOf measure the text and returns the width and height.
type Font interface {
	Tex2D() (uint16, *bk.Texture2D)
	Glyph(rune rune) (Glyph, bool)
	Bounds() (float32, float32)
	SizeOf(text string, fontSize float32) f32.Vec2
}

type Disposer interface {
	Dispose()
}

// A fontAtlas allows rendering ofg text to an OpenGL context
type fontAtlas struct {
	id uint16
	width float32
	height float32

	glyphWidth  float32         // Largest glyph width.
	glyphHeight float32        // Largest glyph height.

	glyphs map[rune]Glyph
	w, h uint16

	regions []float32

	// fallback
}

func (f *fontAtlas) Tex2D() (id uint16, tex *bk.Texture2D) {
	id = f.id
	if ok, t := bk.R.Texture(id); ok {
		tex = t
	}
	return
}

// implement fontAtlas-system
func (f *fontAtlas) Glyph(rune rune) (g Glyph, ok bool) {
	g, ok = f.glyphs[rune]
	return
}

func (f *fontAtlas) Bounds() (float32, float32) {
	return f.glyphWidth, f.glyphHeight
}

// SizeOf returns the width and height for the given text. TODO:
func (f *fontAtlas) SizeOf(text string, fontSize float32) f32.Vec2 {
	return f32.Vec2{}
}

// Release release fontAtlas resources.
func (f *fontAtlas) Dispose() {
	bk.R.Free(f.id)
}

func (f *fontAtlas) addGlyphs(r rune, g Glyph) {
	f.glyphs[r] = g
}

func (f *fontAtlas) loadTex(img *image.RGBA) error {
	// Resize image to power-of-two.
	img = Pow2Image(img).(*image.RGBA)
	ib := img.Bounds()

	// add a white pixel at (0, 0)
	img.Set(0,0, color.White)

	f.width = float32(ib.Dx())
	f.height = float32(ib.Dy())

	if id, _ := bk.R.AllocTexture(img); id != bk.InvalidId {
		f.id = id
	}
	return checkGLError()
}

