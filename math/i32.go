package math

// types for int/int32/int16...

func U32Clamp(v, low, high uint32) uint32 {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func U32Min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func U32Max(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

func U16Clamp(v, low, high uint16) uint16 {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func U16Min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func U16Max(a, b uint16) uint16 {
	if a < b {
		return b
	}
	return a
}

// Pow2 returns the first power-of-two value >= to n.
// This can be used to create suitable Texture dimensions.
func Pow2(x uint32) uint32 {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

// IsPow2 returns true if the given value if a power-of-two.
func IsPow2(x uint32) bool {
	return (x & (x-1)) == 0
}
