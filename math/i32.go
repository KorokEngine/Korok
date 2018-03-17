package math

// types for int/int32/int16...

func UInt32Clamp(v, low, high uint32) uint32 {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func UInt32Min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func UInt32Max(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

func UInt16Clamp(v, low, high uint16) uint16 {
	if v < low {
		return low
	}
	if v > high {
		return high
	}
	return v
}

func UInt16Min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func UInt16Max(a, b uint16) uint16 {
	if a < b {
		return b
	}
	return a
}
