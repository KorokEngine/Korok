package text

import(
	"io"

	_ "image/jpeg"
	_ "image/png"
	"image"
)

// LoadBitmap loads a bitmap (raster) font from the given
// sprite sheet and config files. It is optionally scaled by
// the given scale factor.
//
// A scale factor of 1 retains the original size. A factor of 2 doubles the
// font size, etc. A scale factor of 0 is not valid and will default to 1.
//
// Supported image formats are 32-bit RGBA as PNG, JPEG.
func LoadBitmap(img, config io.Reader, scale int) (*Font, error) {
	pix, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}

	rgba := toRGBA(pix, scale)

	var fc FontConfig
	err = fc.Load(config)

	if err != nil {
		return nil, err
	}

	fc.Glyphs.Scale(scale)
	return loadFont(rgba, &fc)
}




