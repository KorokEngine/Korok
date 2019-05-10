// +build android ios

package res

import (
	"golang.org/x/mobile/asset"
)

// Open opens a named asset.
//
// Errors are of type *os.PathError.
//
// This must not be called from init when used in android apps.
func Open(name string) (File, error) {
	return asset.Open(name)
}
