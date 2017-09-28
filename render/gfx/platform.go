package gfx

import "unsafe"

/// Platform data
type PlatformData struct {
	ndt unsafe.Pointer			// Native display type
	nwh unsafe.Pointer			// Native window handle
	context unsafe.Pointer		// GL Context
	backBuffer unsafe.Pointer	// GL back-buffer
	backBufferDS unsafe.Pointer // Back-buffer depth/stencil
}

/// Set Platform data
func SetPlatformData(data *PlatformData) {

}

/// Internal data
type InternalData struct {
	caps *Caps 					// Renderer capabilities
	context unsafe.Pointer		// GL context
}

func GetInternalData() *InternalData {
	return nil
}



