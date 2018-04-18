package font

import(
	"io"

	_ "image/jpeg"
	_ "image/png"
	"image"
	"io/ioutil"
	"encoding/json"
)

// LoadBitmap loads a bitmap (raster) fontAtlas from the given
// sprite sheet and config files. It is optionally scaled by
// the given scale factor.
//
// A scale factor of 1 retains the original size. A factor of 2 doubles the
// fontAtlas size, etc. A scale factor of 0 is not valid and will default to 1.
//
// Supported image formats are 32-bit RGBA as PNG, JPEG.
func LoadBitmap(img, config io.Reader, scale int) (Font, error) {
	f := &fontAtlas{glyphs:make(map[rune]Glyph)}

	// load texture
	pix, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}
	f.loadTex(toRGBA(pix, scale))

	// load glyph data
	fc := &fontConfig{}
	data, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, fc)
	if err != nil {
		return nil, err
	}
	var gh, gw int
	// add glyphs
	for _, g := range fc.Glyphs {
		f.addGlyphs(g.Id, Glyph{
			Rune: g.Id,
			X: uint16(g.X), Y: uint16(g.Y),
			Width:uint16(g.Width),
			Height:uint16(g.Height),
			XOffset:uint16(g.XOffset),
			YOffset:uint16(g.YOffset),
			Advance:g.Advance,
		})
		if g.Width > gw {
			gw = g.Width
		}
		if g.Height > gh {
			gh = g.Height
		}
	}
	f.glyphWidth = gw
	f.glyphHeight = gh
	// log.Println("dump:", f)
	return f, nil
}

type fontConfig struct {
	Dir Direction `json:"direction"`

	// Lower rune boundary
	Low rune `json:"rune_low"`

	// Upper rune boundary.
	High rune `json:"rune_high"`

	// Glyphs holds a set of glyph descriptors, defining the location,
	// size and advance of each glyph in the sprite sheet.
	Glyphs []struct {
		Id     rune `json:"id,string"`
		X      int `json:"x,string"`
		Y      int `json:"y,string"`
		Width  int `json:"width,string"`
		Height int `json:"height,string"`
		XOffset int `json:"xoffset,string"`
		YOffset int `json:"yoffset,string"`
		Advance int `json:"advance,string"`
	} `json:"glyphs"`
}

