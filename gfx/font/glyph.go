package font

type Point struct {
	X, Y float32
}

// A Glyph describes metrics for a single Font glyph.
// These indicate which area of a given image contains the
// glyph data and how the glyph should be spaced in a rendered string.
//
// Advance determines the distance to the next glyph.
// This is used to properly align non-monospaced fonts.
type Glyph struct {
	Rune     rune

	X, Y    uint16
	Width  uint16
	Height uint16

	XOffset uint16
	YOffset uint16

	Advance int
}
