package font

import (
	"image"
	"fmt"
	"github.com/go-gl/gl/v3.2-core/gl"
)

// toRGBA translate the given image to RGBA format if necessary.
// Optionally scales it by the given amount
func toRGBA(src image.Image, scale int) *image.RGBA {
	if scale < 1 {
		scale = 1
	}

	dst, ok := src.(*image.RGBA)
	if ok && scale == 1 {
		return dst
	}

	// Scale image to match new size
	ib := src.Bounds()
	rect := image.Rect(0, 0, ib.Dx()*scale, ib.Dy()*scale)

	if !ok {
		// Image is not RGBA, so we create it
		dst = image.NewRGBA(rect)
	}

	for sy := 0; sy < ib.Dy(); sy++ {
		for sx := 0; sx < ib.Dx(); sx++ {
			dx := sx * scale
			dy := sy * scale
			pixel := src.At(sx, sy)

			for scy := 0; scy < scale; scy++ {
				for scx := 0; scx < scale; scx++ {
					dst.Set(dx+scx, dy+scy, pixel)
				}
 			}
		}
	}
	return dst
}

// checkGLError returns an opengl error if one exists
func checkGLError() error {
	errno := gl.GetError()
	if errno == gl.NO_ERROR {
		return nil
	}

	return fmt.Errorf("GL error: %d", errno)
}

