package text

import (
	"image"
	"github.com/go-gl/gl/v3.2-core/gl"
	"korok/gfx/bk"
)

// A Font allows rendering ofg text to an OpenGL context
type Font struct {
	config         *FontConfig // Character set for this font.
	Texture        uint32      // Holds the glyph Texture id.
	maxGlyphWidth  int         // Largest glyph width.
	maxGlyphHeight int         // Largest glyph height.
	TexWidth 	   float32
	TexHeight 	   float32
}

// loadFont loads the given font data. This does not deal with font scaling.
// Scaling should be handled by the independent Bitmap/TrueType loaders.
// We therefore expect the supplied image and charset to already be adjusted
// to the correct font scale
//
// The image should hold a sprite sheet, define the graphical layout for
// every glyph. The config describes font metadata
func loadFont(img *image.RGBA, config *FontConfig) (f *Font, err error) {
	f = new(Font)
	f.config = config

	// Resize image to enxt power-of-two.
	img = Pow2Image(img).(*image.RGBA)
	ib := img.Bounds()

	f.TexWidth = float32(ib.Dx())
	f.TexHeight = float32(ib.Dy())

	if id, _ := bk.R.AllocTexture(img); id != bk.InvalidId {
		f.Texture = uint32(id)
	}
	err = checkGLError()
	return
}

// Dir return the font's rendering orientation
func (f *Font) Dir() Direction {
	return f.config.Dir
}

// Low returns the font's lower rune bound
func (f *Font) Low() rune {
	return f.config.Low
}

// High returns the font['s upper rune bound
func (f *Font) High() rune {
	return f.config.High
}

// Glyphs returns the font's glyph descriptions
func (f *Font) Glyphs() Charset {
	return f.config.Glyphs
}

// Release release font resources.
// A font can no longer be used for rendering after this call completes
func (f *Font) Release() {
	gl.DeleteTextures(1, &f.Texture)
	// gl.DeleteList display list TODO
	f.config = nil
}

// Metrics returns the pixel width and height for the given string.
// This takes the scale and rendering direction of the font into account.
//
// Unknown runes will be counted as having the maximum glyph bounds as
// defined by Font.GlyphBounds().
func (f *Font) Metrics(text string) (int, int) {
	if len(text) == 0 {
		return 0, 0
	}

	gw, gh := f.GlyphBounds()

	if f.config.Dir == TopToBottom {
		return gw, f.advanceSize(text)
	}

	return f.advanceSize(text), gh
}

// advanceSize compute the pixel width or height for the given single-line
// input string. This iterate over all of its runes, finds the mathing
// Charset entry and adds up the Advance values
func (f *Font) advanceSize(line string) int {
	gw, gh := f.GlyphBounds()
	glyphs := f.config.Glyphs
	low    := f.config.Low
	indices:= []rune(line)

	var size int
	for _, r := range indices {
		r -= low

		if r >= 0 && int(r) < len(glyphs) {
			size += glyphs[r].Advance
			continue
		}

		if f.config.Dir == TopToBottom {
			size += gh
		} else {
			size += gw
		}
	}

	return size
}


// Printf draws the given string at the specified coordinates.
// It expects the string to be a single line. Line breaks are not
// handled as line breaks and are rendered as glyphs.
//
// In order to render multi-line text, it is up to the caller to split
// the text up into individual lines of adequate length and then call
// this method for each line seperately.
func (f *Font) Printf(x, y float32, fs string, argv ...interface{}) error {
	return nil
}

// GlyphBounds returns the largest width and height for any of the glyphs
// in the font. This constitutes the largest possible bounding box
// a single glyph will have
func (f *Font) GlyphBounds() (int, int) {
	return f.maxGlyphWidth, f.maxGlyphHeight
}

