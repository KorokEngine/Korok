// +build darwin linux windows

package res

import (
	"os"
)

// Open opens a named asset.
func Open(name string) (File, error) {
	return os.Open(name)
}
