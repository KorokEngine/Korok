package font

import (
	"korok.io/korok/math"
	"image"
	"fmt"
	"image/color"
)

// Pow2Image returns the given image, scaled to the smallest power-of-two
// dimensions larger or equal to the input dimensions.
// It preserves the image format and contents.
//
// This is useful if an image is to be used as an OpenGL Texture.
// These often require image data to have power-of-two dimensions.
func Pow2Image(src image.Image) image.Image {
	sb := src.Bounds()
	w, h := uint32(sb.Dx()), uint32(sb.Dy())

	if math.IsPow2(w) && math.IsPow2(h) {
		return src // Nothing to do.
	}

	rect := image.Rect(0, 0, int(math.Pow2(w)), int(math.Pow2(h)))

	switch src := src.(type) {
	case *image.Alpha:
		return copyImg(src, image.NewAlpha(rect))
	case *image.Alpha16:
		return copyImg(src, image.NewAlpha16(rect))
	case *image.Gray:
		return copyImg(src, image.NewGray(rect))
	case *image.Gray16:
		return copyImg(src, image.NewGray16(rect))
	case *image.NRGBA:
		return copyImg(src, image.NewNRGBA(rect))
	case *image.NRGBA64:
		return copyImg(src, image.NewNRGBA64(rect))
	case *image.Paletted:
		return copyImg(src, image.NewPaletted(rect, src.Palette))
	case *image.RGBA:
		return copyImg(src, image.NewRGBA(rect))
	case *image.RGBA64:
		return copyImg(src, image.NewRGBA64(rect))
	}

	panic(fmt.Sprintf("Unsupported image format: %T", src))
}

//
type copyable interface {
	image.Image
	Set(x, y int, clr color.Color)
}

func copyImg(src, dst copyable) image.Image {
	var x, y int
	sb := src.Bounds()

	for y = 0; y < sb.Dy(); y++ {
		for x = 0; x < sb.Dx(); x++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
	return dst
}


