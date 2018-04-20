package font

import (
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx/bk"

	"image"
	"image/color"

	"unicode/utf8"
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
	Glyph(rune rune) (g Glyph, ok bool)
	Bounds() (gw, gh float32)
	Frame(rune rune) (x1, y1, x2, y2 float32)
}

type Disposer interface {
	Dispose()
}

// A fontAtlas allows rendering ofg text to an OpenGL context
type fontAtlas struct {
	id uint16

	texWidth float32
	texHeight float32

	gWidth  float32 // Largest glyph width.
	gHeight float32 // Largest glyph height.

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
	return f.gWidth, f.gHeight
}

func (f *fontAtlas) Frame(r rune) (u1, v1, u2, v2 float32) {
	g := f.glyphs[r]
	u1, v1 = float32(g.X)/ f.texWidth, float32(g.Y)/ f.texHeight
	u2, v2 = float32(g.X+g.Width)/ f.texWidth, float32(g.Y+g.Height)/ f.texHeight
	return
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

	f.texWidth = float32(ib.Dx())
	f.texHeight = float32(ib.Dy())

	if id, _ := bk.R.AllocTexture(img); id != bk.InvalidId {
		f.id = id
	}
	return checkGLError()
}

func Wrap(font Font, text string, wrap, fontSize float32) (n int, lines string) {
	size := len(text)
	line := make([]byte, 0, size)
	buff := make([]byte, 0, size*2)
	_, gh := font.Bounds()
	scale := fontSize/gh

	for i, w := 0, 0; i < size;{
		var lineSize float32
		var lastSpace = -1

		for j := 0 ; i < size; i, j = i+w, j+w {
			r, width := utf8.DecodeRuneInString(text[i:])
			w = width

			if r == '\n' {
				i, j = i+w, j+w
				goto NEW_LINE
			}

			line = append(line, text[i:i+w]...)
			if g, ok :=  font.Glyph(r); ok {
				lineSize += float32(g.Advance) * scale
			} else {
				// todo fallback ?
			}
			if r == ' ' || r == '\t' {
				lastSpace = j
			}

			if lineSize > wrap {
				i, j = i+w, j+w
				break
			}
		}

		// reach the end
		if lineSize < wrap {
			buff = append(buff, line...)
			n += 1
			break
		}

		if len(line) == 0 {

		}

		// if has space, break! or remove last char to fit line-width
		if lastSpace > 0 {
			i -= len(line) - lastSpace

			//log.Println("trim last space i:", i, "lastSpace:", lastSpace, "lineSize:", len(line))
			line = line[:lastSpace]
		} else {
			i -= w
			line = line[:len(line)-w]
		}

	NEW_LINE:
		line = append(line, '\n')
		buff = append(buff, line...)
		line = line[:0]
		n += 1
	}
	return n, string(buff)
}

func CalculateTextSize(text string, font Font, fontSize float32) f32.Vec2 {
	_, gh := font.Bounds()
	scale := fontSize/gh
	size := f32.Vec2{0, fontSize}

	for i, w := 0, 0; i < len(text); i += w {
		r, width := utf8.DecodeRuneInString(text[i:])
		w = width
		g, _ := font.Glyph(r)
		if r >= 32 {
			size[0] += float32(g.Advance) * scale
		}
	}
	return size
}