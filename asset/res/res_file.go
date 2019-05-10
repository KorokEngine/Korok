// +build darwin linux windows

package res

import (
	"golang.org/x/mobile/asset"
)

// Open opens a named asset.
func Open(name string) (File, error) {
	return asset.Open(name)
}
