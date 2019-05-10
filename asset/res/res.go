package res

import (
	"io"
)

// File is an open asset.
type File interface {
	io.ReadCloser
	// io.Closer
}

type WFile interface {
	io.Writer
}
