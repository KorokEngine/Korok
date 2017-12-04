package font

type Point struct {
	X, Y float32
}

// A Glyph describes metrics for a single font glyph.
// These indicate which area of a given image contains the
// glyph data and how the glyph should be spaced in a rendered string.
type Glyph struct {
	Id     rune `json:"id,string"` 	   // The id
	X      int `json:"x,string"`      // The x location of the glyph on a sprite sheet.
	Y      int `json:"y,string"`      // The y location of the glyph on a sprite sheet.
	Width  int `json:"width,string"`  // The width of the glyph on a sprite sheet.
	Height int `json:"height,string"` // The height of the glyph on a sprite sheet.

	// Advance determines the distance to the next glyph.
	// This is used to properly align non-monospaced fonts.
	Advance int `json:"advance,string"`
}

func (g *Glyph) GetTexturePosition(f *Font) (min, max Point) {
	min = Point{float32(g.X) / f.TexWidth, float32(g.Y) / f.TexHeight}
	max = Point{float32(g.X+g.Width)/ f.TexWidth, float32(g.Y+g.Height) / f.TexHeight}
	return
}

// A Charset represents a set of glyph descriptors for a font.
// Each glyph descriptor holds glyph metrics which are used to
// properly align the given glyph in the resulting rendered string.
type Charset []Glyph

// Scale scales all glyphs by the given factor and repositions them
// appropriately. A scale of 1 retains the original size. A scale of 2
// doubles the size of each glyph, etc.
//
// This is useful when the accompanying sprite sheet is scaled by the
// same factor. In this case, we want the glyph data to match up with the
// new image.
func (c Charset) Scale(factor int) {
	if factor <= 1 {
		// A factor of zero results in zero-sized glyphs and
		// is therefore not valid. A factor of 1 does not change
		// the glyphs, so we can ignore it.
		return
	}

	// Multiply each glyph field by the given factor
	// to scale them up to the new size.
	for i := range c {
		c[i].X *= factor
		c[i].Y *= factor
		c[i].Width *= factor
		c[i].Height *= factor
		c[i].Advance *= factor
	}
}
