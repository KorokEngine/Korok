package font

import (
	"io"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"image"
	"github.com/golang/freetype"
)

// http://www.freetype.org/freetype2/docs/tutorial/step2.html

// LoadTrueType loads a truetype font from the given stream and
// applies the given font scale in points.
//
// The low and high values determine the lower and upper rune limits
// we should load for this font. For standard ASCII this would be:32, 127.
//
// The dir value determines the orientation of the text we render
// with this font. This should be any of the predefined Direction constants.
func LoadTrueType(r io.Reader, scale int32, low, high rune, dir Direction) (*Font, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the truetype font
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	// Create our FontConfig type
	var fc FontConfig
	fc.Dir = dir
	fc.Low = low
	fc.High = high
	fc.Glyphs = make(Charset, high-low,+ 1)

	// Create an image, large enough to store all requested glyphs.
	//
	// We limit the image to 16 glyphs per row. Then add as many rows as
	// needed to encompass all glyphs, while making sure the resulting image
	// has power-of-two dimensions.
	gc := int32(len(fc.Glyphs))
	glyphsPerRow := int32(16)
	glyphsPerCol := (gc / glyphsPerRow) + 1

	gb := ttf.Bounds(fixed.Int26_6(scale))
	gw := int32(gb.Max.X - gb.Min.X)
	gh := int32((gb.Max.Y - gb.Min.Y) + 5)
	iw := Pow2(uint32(gw * glyphsPerRow))
	ih := Pow2(uint32(gh * glyphsPerCol))

	rect := image.Rect(0, 0, int(iw), int(ih))
	img  := image.NewRGBA(rect)

	// Use a freetype context to do the drawing.
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFontSize(float64(scale))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)

	// Iterate over all relevant glyphs in the truetype font and
	// draw them all to the image buffer.
	//
	// For each glyph, we also create a corresponding Glyph structure
	// for our Charset. It contains the appropriate glyph coordinate offsets.
	var gi int
	var gx, gy int32

	for ch := low; ch <= high; ch++ {
		index  := ttf.Index(ch)
		metric := ttf.HMetric(fixed.Int26_6(scale), index)

		fc.Glyphs[gi].Advance = int(metric.AdvanceWidth)
		fc.Glyphs[gi].X = int(gx)
		fc.Glyphs[gi].Y = int(gy) - int(gh)/2 //shif up half a row so that actually get character in frame
		fc.Glyphs[gi].Width = int(gw)
		fc.Glyphs[gi].Height = int(gh)
		pt := freetype.Pt(int(gx), int(gy)+int(c.PointToFixed(float64(scale))>>8))
		c.DrawString(string(ch), pt)

		if gi%16 == 0 {
			gx = 0
			gy += gh
		} else {
			gx += gw
		}

		gi++
	}

	return loadFont(img, &fc)
}
