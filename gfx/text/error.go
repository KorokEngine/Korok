package text

import (
	"fmt"

	"github.com/go-gl/gl/v3.2-core/gl"
)

// checkGLError returns an opengl error if one exists
func checkGLError() error {
	errno := gl.GetError()
	if errno == gl.NO_ERROR {
		return nil
	}

	return fmt.Errorf("GL error: %d", errno)
}

