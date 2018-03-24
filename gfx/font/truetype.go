package font

import (
	"github.com/golang/freetype/truetype"
	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"

	"io"
	"io/ioutil"
	"image"
	"log"
	"os"
	"image/png"
)

// http://www.freetype.org/freetype2/docs/tutorial/step2.html

// LoadTrueType loads a truetype fontAtlas from the given stream and
// applies the given fontAtlas scale in points.
//
// The low and high values determine the lower and upper rune limits
// we should load for this Font. For standard ASCII this would be:32, 127.
func LoadTrueType(r io.Reader, size int32, low, high rune, dir Direction) (*fontAtlas, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Read the TrueType Font
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	// Create an image(512*512) to store all requested glyphs.
	// TODO: image size
	iw, ih := 512, 512
	_, fg := image.Black, image.White
	rect := image.Rect(0, 0, iw, ih)
	img  := image.NewRGBA(rect)

	// Use a FreeType context to do the drawing.
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetFontSize(float64(size))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)

	// new font-atlas
	f := &fontAtlas{glyphs:make(map[rune]Glyph)}

	// Iterate over all relevant glyphs in the truetype fontAtlas and
	// draw them all to the image buffer.
	var (
		gb = ttf.Bounds(fixed.Int26_6(size))
		gw = (gb.Max.X - gb.Min.X)
		gh = (gb.Max.Y - gb.Min.Y)
		gx, gy int
	)

	for ch := low; ch <= high; ch++ {
		g := Glyph{Rune: ch}

		if gx + int(gw) > iw {
			gx, gy = 0, gy + int(gh)
		}

		index  := ttf.Index(ch)
		hm := ttf.HMetric(fixed.Int26_6(size), index)
		g.Advance = int(hm.AdvanceWidth)
		g.X = uint16(gx)
		g.Y = uint16(gy)
		g.Width = uint16(gw)
		g.Height = uint16(gh)

		f.addGlyphs(ch, g)
		pt := freetype.Pt(int(gx), int(gy)+int(c.PointToFixed(float64(size))>>6))
		c.DrawString(string(ch), pt)
		gx += int(gw)
	}

	// load image
	f.loadTex(img)

	// save baked fontAtlas-image
	// savePng(img)
	return f, nil
}

// debug only.
func savePng(img image.Image) {
	f, err := os.Create("ttf.png")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Println(err)
	}
}